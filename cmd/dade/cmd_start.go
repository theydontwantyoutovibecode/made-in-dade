package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/theydontwantyoutovibecode/dade/internal/config"
	"github.com/theydontwantyoutovibecode/dade/internal/lifecycle"
	"github.com/theydontwantyoutovibecode/dade/internal/logging"
	"github.com/theydontwantyoutovibecode/dade/internal/manifest"
	"github.com/theydontwantyoutovibecode/dade/internal/proxy"
	"github.com/theydontwantyoutovibecode/dade/internal/registry"
	"github.com/theydontwantyoutovibecode/dade/internal/serve"
	"github.com/theydontwantyoutovibecode/dade/internal/ui"
	"github.com/spf13/cobra"
)

type startCommand struct {
	templatesDir func() (string, error)
	projectsFile func() (string, error)
	readMarker   func(string) (registry.Marker, error)
	readFile     func(string) ([]byte, error)
	startStatic  func(ctx context.Context, runner serve.CommandRunner, port int, root string) (int, error)
	isPortInUse  func(int) bool
}

var startCommandFactory = defaultStartCommand

func defaultStartCommand() startCommand {
	return startCommand{
		templatesDir: config.TemplatesDir,
		projectsFile: config.ProjectsFile,
		readMarker:   registry.ReadMarker,
		readFile:     os.ReadFile,
		startStatic:  serve.StartStaticServer,
		isPortInUse:  isPortInUse,
	}
}

func runStartCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	portOverride, _ := cmd.Flags().GetInt("port")
	background, _ := cmd.Flags().GetBool("background")

	impl := startCommandFactory()
	code := impl.run(context.Background(), args, console, logger, portOverride, background)
	if code != 0 {
		return errors.New("start command failed")
	}
	return nil
}

func (c startCommand) run(ctx context.Context, args []string, console *ui.UI, logger *logging.Logger, portOverride int, background bool) int {
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

	serveType := "static"
	var mf manifest.Manifest
	if data, err := c.readFile(manifestPath); err == nil {
		if parsed, err := manifest.Parse(data); err == nil {
			mf = parsed
			if parsed.Serve.Type != "" {
				serveType = parsed.Serve.Type
			}
		}
	}

	needsProxy := manifest.NeedsProxy(mf)

	// Check if already running (only meaningful for projects with a port)
	if port > 0 && c.isPortInUse(port) {
		logger.Warn(fmt.Sprintf("Project already running on port %d", port))
		if needsProxy {
			logger.Info(fmt.Sprintf("URL: https://%s", config.ProjectDomain(projectName)))
		}
		return 0
	}

	portEnv := mf.Serve.PortEnv
	if portEnv == "" {
		portEnv = "PORT"
	}

	_ = console

	logger.Info(fmt.Sprintf("Starting %s in production mode...", projectName))

	switch serveType {
	case "static":
		staticRoot := projectDir
		if mf.Serve.Static.Root != "" {
			staticRoot = filepath.Join(projectDir, mf.Serve.Static.Root)
		}
		_, err := c.startStatic(ctx, nil, port, staticRoot)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to start static server: %v", err))
			return 1
		}
		logger.Success(fmt.Sprintf("Started: https://%s", config.ProjectDomain(projectName)))
		return 0

	case "command":
		// Use production command
		serveCmd := manifest.ServeCommand(mf, "prod")
		if serveCmd == "" {
			logger.Error("No serve.prod command defined in manifest")
			return 1
		}

		if background {
			return c.startBackground(ctx, projectDir, serveCmd, port, portEnv, projectName, mf, logger)
		}

		// Run in foreground with lifecycle management
		return c.startForeground(ctx, projectDir, serveCmd, port, portEnv, projectName, mf, logger)

	default:
		logger.Error(fmt.Sprintf("Unknown serve type: %s", serveType))
		return 1
	}
}

func (c startCommand) startBackground(ctx context.Context, projectDir, serveCmd string, port int, portEnv, projectName string, mf manifest.Manifest, logger *logging.Logger) int {
	cmd := exec.CommandContext(ctx, "bash", "-c", serveCmd)
	cmd.Dir = projectDir
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%d", portEnv, port))

	if err := cmd.Start(); err != nil {
		logger.Error(fmt.Sprintf("Failed to start server: %v", err))
		return 1
	}

	pid := cmd.Process.Pid
	pidFile := filepath.Join(projectDir, serve.DefaultPIDFile)
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		_ = cmd.Process.Kill()
		logger.Error(fmt.Sprintf("Failed to write PID file: %v", err))
		return 1
	}

	if manifest.NeedsProxy(mf) {
		c.ensureProxy(ctx, logger)
		logger.Success(fmt.Sprintf("Started (PID %d): https://%s", pid, config.ProjectDomain(projectName)))
	} else {
		logger.Success(fmt.Sprintf("Started (PID %d)", pid))
	}
	return 0
}

func (c startCommand) startForeground(ctx context.Context, projectDir, serveCmd string, port int, portEnv, projectName string, mf manifest.Manifest, logger *logging.Logger) int {
	// Create lifecycle controller
	ctrl := lifecycle.NewController(projectDir)

	// Add production environment variables
	if manifest.HasProdSection(mf) {
		ctrl.WithEnv(manifest.ProdEnv(mf))
	}

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Shutting down...")
		ctrl.Shutdown()
	}()

	// Run production setup if defined
	if manifest.HasProdSection(mf) {
		setupScript := manifest.ProdSetupScript(mf)
		setupCmds := manifest.ProdSetupCommands(mf)

		if setupScript != "" {
			logger.Info("Running production setup script...")
			if err := ctrl.RunSetupScript(ctx, setupScript); err != nil {
				logger.Error(fmt.Sprintf("Setup script failed: %v", err))
				return 1
			}
		} else if len(setupCmds) > 0 {
			logger.Info("Running production setup commands...")
			if err := ctrl.RunSetup(ctx, setupCmds); err != nil {
				logger.Error(fmt.Sprintf("Setup failed: %v", err))
				return 1
			}
		}
	}

	// Ensure HTTPS proxy is running (only for web projects)
	if manifest.NeedsProxy(mf) {
		c.ensureProxy(ctx, logger)
	}

	// Display ready message
	if manifest.NeedsProxy(mf) {
		logger.Success(fmt.Sprintf("Ready: https://%s", config.ProjectDomain(projectName)))
	} else {
		logger.Success(fmt.Sprintf("Starting %s...", projectName))
	}
	logger.Info("Press Ctrl+C to stop")
	fmt.Println()

	// Start main server (blocking)
	if err := ctrl.StartServer(ctx, serveCmd, port, portEnv); err != nil {
		if ctx.Err() != nil {
			return 0
		}
		logger.Error(fmt.Sprintf("Server exited: %v", err))
		return 1
	}

	return 0
}

func (c startCommand) ensureProxy(ctx context.Context, logger *logging.Logger) {
	projectsPath, err := c.projectsFile()
	if err != nil {
		return
	}
	caddyfilePath, err := config.CaddyfilePath()
	if err != nil {
		return
	}
	if err := proxy.GenerateCaddyfile(ctx, nil, projectsPath, caddyfilePath); err != nil {
		logger.Warn(fmt.Sprintf("Failed to regenerate Caddyfile: %v", err))
		return
	}

	if !proxy.IsProxyRunning(ctx, nil) {
		plistPath, err := config.ProxyPlistPath()
		if err != nil {
			return
		}
		logPath, _ := config.ProxyLogPath()
		errPath, _ := config.ProxyErrPath()
		if err := proxy.CreatePlist(ctx, nil, plistPath, caddyfilePath, logPath, errPath); err != nil {
			logger.Warn(fmt.Sprintf("Failed to create proxy plist: %v", err))
		} else if err := proxy.InstallProxyService(ctx, nil, plistPath); err != nil {
			logger.Warn(fmt.Sprintf("Failed to start proxy: %v", err))
		} else {
			logger.Info("Started HTTPS proxy")
		}
	} else if err := proxy.ReloadProxy(ctx, nil, caddyfilePath); err != nil {
		logger.Warn(fmt.Sprintf("Failed to reload proxy: %v", err))
	}
}
