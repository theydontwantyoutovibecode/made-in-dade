package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
	execx "github.com/theydontwantyoutovibecode/made-in-dade/internal/exec"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/lifecycle"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/manifest"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/proxy"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/readonly"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/registry"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/ui"
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev [name]",
	Short: "Start development server with full orchestration",
	Long: `Start a development server with setup commands, background processes, and 
the main server. Uses the [dev] section from the template manifest if available,
otherwise falls back to the serve.dev command.

The dev command:
1. Runs setup commands (dependencies, migrations)
2. Starts background processes (asset watchers)  
3. Starts the main dev server
4. Ensures HTTPS proxy is running`,
	Example: `dade dev              # Start current project in dev mode
dade dev myapp        # Start specific project
dade dev --skip-setup # Skip setup commands
dade dev --open       # Open in browser after starting`,
	GroupID: "dev",
	Args:    cobra.MaximumNArgs(1),
	RunE:    runDevCmd,
}

func init() {
	rootCmd.AddCommand(devCmd)
	devCmd.Flags().Bool("skip-setup", false, "Skip setup commands")
	devCmd.Flags().Bool("open", false, "Open project in browser after starting")
	devCmd.Flags().IntP("port", "p", 0, "Override port")
}

type devCommand struct {
	runner       execx.Runner
	templatesDir func() (string, error)
	projectsFile func() (string, error)
	readMarker   func(string) (registry.Marker, error)
	readFile     func(string) ([]byte, error)
	isPortInUse  func(int) bool
}

var devCommandFactory = defaultDevCommand

func defaultDevCommand() devCommand {
	return devCommand{
		runner:       execx.NewSystemRunner(),
		templatesDir: config.TemplatesDir,
		projectsFile: config.ProjectsFile,
		readMarker:   registry.ReadMarker,
		readFile:     os.ReadFile,
		isPortInUse:  isPortInUse,
	}
}

func runDevCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	skipSetup, _ := cmd.Flags().GetBool("skip-setup")
	openFlag, _ := cmd.Flags().GetBool("open")
	portOverride, _ := cmd.Flags().GetInt("port")

	impl := devCommandFactory()
	code := impl.run(context.Background(), args, console, logger, skipSetup, openFlag, portOverride)
	if code != 0 {
		return errors.New("dev command failed")
	}
	return nil
}

func (c devCommand) run(ctx context.Context, args []string, console *ui.UI, logger *logging.Logger, skipSetup bool, openBrowser bool, portOverride int) int {
	var projectDir string
	var projectName string

	// Resolve project
	if len(args) > 0 {
		projectName = args[0]
		projectsPath, err := c.projectsFile()
		if err != nil {
			logger.Error("Failed to resolve projects file")
			return 1
		}
		project, ok, err := registry.Get(projectsPath, projectName)
		if err != nil {
			logger.Error("Failed to load project registry")
			return 1
		}
		if !ok {
			logger.Error(fmt.Sprintf("Project '%s' not found", projectName))
			return 1
		}
		projectDir = project.Path
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			logger.Error("Failed to get current directory")
			return 1
		}
		if !registry.MarkerExists(cwd) {
			logger.Error("Not a dade project directory")
			logger.Info("Run 'dade new' or 'dade register' first")
			return 1
		}
		projectDir = cwd
	}

	// Read project marker
	marker, err := c.readMarker(projectDir)
	if err != nil {
		logger.Error("Failed to read project marker")
		return 1
	}
	projectName = marker.Name
	port := marker.Port
	if portOverride > 0 {
		port = portOverride
	}
	templateName := marker.Template

	// Load template manifest (before port check so we know if proxy is needed)
	templatesDir, err := c.templatesDir()
	if err != nil {
		logger.Error("Failed to resolve templates directory")
		return 1
	}

	templateDir := filepath.Join(templatesDir, templateName)
	manifestPath := filepath.Join(templateDir, "dade.toml")

	var mf manifest.Manifest
	if data, err := c.readFile(manifestPath); err == nil {
		if parsed, err := manifest.Parse(data); err == nil {
			mf = parsed
		}
	}

	needsProxy := manifest.NeedsProxy(mf)

	// Check if already running (only meaningful for projects with a port)
	if port > 0 && c.isPortInUse(port) {
		pidFile := filepath.Join(projectDir, ".dade.pid")
		pidData, err := os.ReadFile(pidFile)

		if err == nil {
			pidStr := strings.TrimSpace(string(pidData))
			pid, _ := strconv.Atoi(pidStr)

			if pid > 0 {
				process, err := os.FindProcess(pid)
				if err == nil {
					if err := process.Signal(syscall.Signal(0)); err != nil {
						// Process is not running, stale PID file
						logger.Info("Removing stale PID file")
						_ = os.Remove(pidFile)
					} else {
						// Process is running
						logger.Warn(fmt.Sprintf("Project already running on port %d", port))
						if needsProxy {
							logger.Info(fmt.Sprintf("URL: https://%s", config.ProjectDomain(projectName)))
						}
						return 0
					}
				}
			}
		}

		if os.IsNotExist(err) {
			// No PID file but port is in use - orphaned process
			logger.Warn(fmt.Sprintf("Found orphaned process on port %d", port))
			if err := killProcessOnPort(port); err != nil {
				logger.Warn(fmt.Sprintf("Failed to kill orphaned process: %v", err))
				logger.Info(fmt.Sprintf("URL: https://%s", config.ProjectDomain(projectName)))
				return 0
			}
			logger.Info("Cleaned up orphaned process")
		}
	}

	// Determine serve type and command
	serveType := manifest.ServeType(mf)
	serveCmd := manifest.ServeCommand(mf, "dev")

	portEnv := mf.Serve.PortEnv
	if portEnv == "" {
		portEnv = "PORT"
	}

	switch serveType {
	case "static":
		// Handle static server separately - no command needed
		if serveCmd != "" {
			logger.Warn("serve.dev command ignored for static templates")
		}
	case "command":
		if serveCmd == "" {
			logger.Error("No serve.dev command defined in manifest")
			return 1
		}
	default:
		logger.Error(fmt.Sprintf("Unknown serve type: %s", serveType))
		return 1
	}

	// Create lifecycle controller
	ctrl := lifecycle.NewController(projectDir)

	// Add environment variables from manifest
	if manifest.HasDevSection(mf) {
		ctrl.WithEnv(manifest.DevEnv(mf))
	}

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Shutting down...")
		ctrl.Shutdown()
	}()

	_ = console

	logger.Info(fmt.Sprintf("Starting %s in development mode...", projectName))

	// Sync .read-only reference dependencies
	if err := readonly.SyncDeps(ctx, c.runner, projectDir, logger); err != nil {
		logger.Warn(fmt.Sprintf("Failed to sync read-only deps: %v", err))
	}

	// Run setup commands
	if !skipSetup && manifest.HasDevSection(mf) {
		setupScript := manifest.DevSetupScript(mf)
		setupCmds := manifest.DevSetupCommands(mf)

		if setupScript != "" {
			logger.Info("Running setup script...")
			if err := ctrl.RunSetupScript(ctx, setupScript); err != nil {
				logger.Error(fmt.Sprintf("Setup script failed: %v", err))
				return 1
			}
		} else if len(setupCmds) > 0 {
			logger.Info("Running setup commands...")
			if err := ctrl.RunSetup(ctx, setupCmds); err != nil {
				logger.Error(fmt.Sprintf("Setup failed: %v", err))
				return 1
			}
		}
		logger.Success("Setup complete")
	}

	// Start background processes
	if manifest.HasDevSection(mf) {
		bgCmds := manifest.DevBackgroundCommands(mf)
		if len(bgCmds) > 0 {
			logger.Info(fmt.Sprintf("Starting %d background process(es)...", len(bgCmds)))
			if err := ctrl.StartBackground(ctx, bgCmds); err != nil {
				logger.Error(fmt.Sprintf("Failed to start background processes: %v", err))
				return 1
			}
		}
	}

	// Ensure HTTPS proxy is running (only for web projects)
	if needsProxy {
		projectsPath, err := c.projectsFile()
		if err == nil {
			caddyfilePath, err := config.CaddyfilePath()
			if err == nil {
				if err := proxy.GenerateCaddyfile(ctx, nil, projectsPath, caddyfilePath); err != nil {
					logger.Warn(fmt.Sprintf("Failed to regenerate Caddyfile: %v", err))
				} else {
					if !proxy.IsProxyRunning(ctx, nil) {
						plistPath, err := config.ProxyPlistPath()
						if err == nil {
							logPath, _ := config.ProxyLogPath()
							errPath, _ := config.ProxyErrPath()
							if err := proxy.CreatePlist(ctx, nil, plistPath, caddyfilePath, logPath, errPath); err != nil {
								logger.Warn(fmt.Sprintf("Failed to create proxy plist: %v", err))
							} else if err := proxy.InstallProxyService(ctx, nil, plistPath); err != nil {
								logger.Warn(fmt.Sprintf("Failed to start proxy: %v", err))
							} else {
								logger.Info("Started HTTPS proxy")
							}
						}
					} else if err := proxy.ReloadProxy(ctx, nil, caddyfilePath); err != nil {
						logger.Warn(fmt.Sprintf("Failed to reload proxy: %v", err))
					}
				}
			}
		}
	}

	// Display domain TLD info
	domainTLD := config.DomainTLD()
	if domainTLD == ".local" {
		logger.Info("Note: .local domains work on this machine only. For LAN access, update /etc/hosts manually")
	} else if domainTLD != ".localhost" {
		logger.Info(fmt.Sprintf("Using custom domain TLD: %s", domainTLD))
	}

	// Display ready message
	if readyMsg := manifest.DevReadyMessage(mf); readyMsg != "" {
		logger.Success(readyMsg)
	} else if needsProxy {
		projectURL := fmt.Sprintf("https://%s", config.ProjectDomain(projectName))
		logger.Success(fmt.Sprintf("Ready: %s", projectURL))
		if openBrowser {
			if err := openBrowserFunc(projectURL); err != nil {
				logger.Warn(fmt.Sprintf("Failed to open browser: %v", err))
			}
		}
	} else {
		logger.Success(fmt.Sprintf("Starting %s...", projectName))
	}

	if runningMsg := manifest.DevRunningMessage(mf); runningMsg != "" {
		logger.Info(runningMsg)
	} else {
		logger.Info("Press Ctrl+C to stop")
	}
	fmt.Println()

	// Start main server (blocking)
	switch serveType {
	case "static":
		staticRoot := projectDir
		if mf.Serve.Static.Root != "" {
			staticRoot = filepath.Join(projectDir, mf.Serve.Static.Root)
		}
		// Build caddy command for static serving (suppress logs for cleaner output)
		serveCmd := fmt.Sprintf("caddy file-server --listen :%[1]d --root %[2]s 2>/dev/null", port, staticRoot)
		if err := ctrl.StartServer(ctx, serveCmd, port, portEnv); err != nil {
			// Check if it was a signal-based shutdown
			if ctx.Err() != nil {
				return 0
			}
			logger.Error(fmt.Sprintf("Server exited: %v", err))
			return 1
		}
	case "command":
		if err := ctrl.StartServer(ctx, serveCmd, port, portEnv); err != nil {
			// Check if it was a signal-based shutdown
			if ctx.Err() != nil {
				return 0
			}
			logger.Error(fmt.Sprintf("Server exited: %v", err))
			return 1
		}
	}

	return 0
}

func killProcessOnPort(port int) error {
	// Use lsof on macOS and Linux
	cmd := exec.Command("lsof", "-ti", fmt.Sprintf(":%d", port))
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to find process on port %d: %w", port, err)
	}

	pidStr := strings.TrimSpace(string(output))
	if pidStr == "" {
		return fmt.Errorf("no process found on port %d", port)
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return fmt.Errorf("invalid PID: %w", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process %d: %w", pid, err)
	}

	// Try SIGTERM first
	if err := process.Signal(syscall.SIGTERM); err != nil {
		// Fall back to SIGKILL
		if err := process.Signal(syscall.SIGKILL); err != nil {
			return fmt.Errorf("failed to kill process %d: %w", pid, err)
		}
	}

	return nil
}
