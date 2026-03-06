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
	"time"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/lifecycle"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/manifest"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/registry"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/tunnel"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/ui"
	"github.com/spf13/cobra"
)

var shareCmd = &cobra.Command{
	Use:   "share [name]",
	Short: "Start dev server and share via public tunnel",
	Long: `Start a development server and create a public Cloudflare tunnel to share it.
Combines the dev command workflow with cloudflared tunnel creation.

The share command:
1. Runs setup commands (same as dev)
2. Starts background processes (same as dev)
3. Starts the main dev server
4. Creates a Cloudflare tunnel (quick or named)
5. Displays the public URL

Use --attach to tunnel an already-running server without starting a new one.

Requires cloudflared to be installed: brew install cloudflared`,
	Example: `dade share              # Share current project
dade share myapp        # Share specific project
dade share --quick      # Force quick tunnel
dade share --attach     # Attach tunnel to running server`,
	GroupID: "dev",
	Args:    cobra.MaximumNArgs(1),
	RunE:    runShareCmd,
}

func init() {
	rootCmd.AddCommand(shareCmd)
	shareCmd.Flags().Bool("skip-setup", false, "Skip setup commands")
	shareCmd.Flags().Bool("quick", false, "Force quick tunnel (ignore named tunnel config)")
	shareCmd.Flags().Bool("attach", false, "Attach tunnel to already-running server (skip server startup)")
	shareCmd.Flags().IntP("port", "p", 0, "Override port")
}

type shareCommand struct {
	templatesDir func() (string, error)
	projectsFile func() (string, error)
	readMarker   func(string) (registry.Marker, error)
	readFile     func(string) ([]byte, error)
	isPortInUse  func(int) bool
}

var shareCommandFactory = defaultShareCommand

func defaultShareCommand() shareCommand {
	return shareCommand{
		templatesDir: config.TemplatesDir,
		projectsFile: config.ProjectsFile,
		readMarker:   registry.ReadMarker,
		readFile:     os.ReadFile,
		isPortInUse:  isPortInUse,
	}
}

func runShareCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	skipSetup, _ := cmd.Flags().GetBool("skip-setup")
	forceQuick, _ := cmd.Flags().GetBool("quick")
	attach, _ := cmd.Flags().GetBool("attach")
	portOverride, _ := cmd.Flags().GetInt("port")

	if attach {
		return runTunnelAttach(cmd, args)
	}

	impl := shareCommandFactory()
	code := impl.run(context.Background(), args, console, logger, skipSetup, forceQuick, portOverride)
	if code != 0 {
		return errors.New("share command failed")
	}
	return nil
}

func runTunnelAttach(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	if _, err := exec.LookPath("cloudflared"); err != nil {
		logger.Error("cloudflared not installed")
		logger.Info("Install: brew install cloudflared")
		return errors.New("share command failed")
	}

	var projectName string
	var port int

	portOverride, _ := cmd.Flags().GetInt("port")

	if len(args) > 0 {
		projectName = args[0]
		projectsPath, err := config.ProjectsFile()
		if err != nil {
			logger.Error("Failed to resolve projects file")
			return errors.New("share command failed")
		}
		project, ok, err := registry.Get(projectsPath, projectName)
		if err != nil {
			logger.Error("Failed to load project registry")
			return errors.New("share command failed")
		}
		if !ok {
			logger.Error(fmt.Sprintf("Project '%s' not found", projectName))
			return errors.New("share command failed")
		}
		port = project.Port
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			logger.Error("Failed to get current directory")
			return errors.New("share command failed")
		}
		if !registry.MarkerExists(cwd) {
			logger.Error("Not a dade project directory")
			logger.Info("Run 'dade new' or 'dade register' first")
			return errors.New("share command failed")
		}
		marker, err := registry.ReadMarker(cwd)
		if err != nil {
			logger.Error("Failed to read project marker")
			return errors.New("share command failed")
		}
		projectName = marker.Name
		port = marker.Port
	}

	if portOverride > 0 {
		port = portOverride
	}

	if !isPortInUse(port) {
		logger.Error("Project not running")
		logger.Info("Start first: dade dev")
		return errors.New("share command failed")
	}

	_ = console
	logger.Info(fmt.Sprintf("Starting quick tunnel for %s...", projectName))
	logger.Info("Press Ctrl+C to stop")

	localURL := fmt.Sprintf("http://localhost:%d", port)
	tunnelExec := exec.Command("cloudflared", "tunnel", "--url", localURL)
	tunnelExec.Stdout = os.Stdout
	tunnelExec.Stderr = os.Stderr
	tunnelExec.Stdin = os.Stdin
	tunnelExec.SysProcAttr = &syscall.SysProcAttr{Setpgid: false}

	if err := tunnelExec.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return nil
		}
		logger.Error(fmt.Sprintf("Tunnel failed: %v", err))
		return errors.New("share command failed")
	}
	return nil
}

func (c shareCommand) run(ctx context.Context, args []string, console *ui.UI, logger *logging.Logger, skipSetup, forceQuick bool, portOverride int) int {
	// Check cloudflared availability
	if !tunnel.IsAvailable() {
		logger.Error("cloudflared not installed")
		logger.Info("Install: brew install cloudflared")
		return 1
	}

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

	// Load template manifest
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

	// Determine serve command
	serveCmd := manifest.ServeCommand(mf, "dev")
	if serveCmd == "" {
		logger.Error("No serve command defined in manifest")
		return 1
	}

	portEnv := mf.Serve.PortEnv
	if portEnv == "" {
		portEnv = "PORT"
	}

	// Create lifecycle controller
	ctrl := lifecycle.NewController(projectDir)

	// Add environment variables from manifest
	if manifest.HasDevSection(mf) {
		ctrl.WithEnv(manifest.DevEnv(mf))
	}

	// Add share-specific environment variables
	if manifest.HasShareSection(mf) {
		ctrl.WithEnv(manifest.ShareEnv(mf))
	}

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var tun *tunnel.Tunnel
	go func() {
		<-sigChan
		logger.Info("Shutting down...")
		if tun != nil {
			_ = tun.Stop()
		}
		ctrl.Shutdown()
	}()

	_ = console

	logger.Info(fmt.Sprintf("Starting %s for sharing...", projectName))

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

	// Start server in background so we can start tunnel
	serverCmd, err := ctrl.StartServerBackground(ctx, serveCmd, port, portEnv)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to start server: %v", err))
		return 1
	}

	// Wait a moment for server to start
	time.Sleep(500 * time.Millisecond)

	// Start tunnel
	tunnelName := ""
	if manifest.HasShareSection(mf) && !forceQuick {
		tunnelName = manifest.ShareTunnelName(mf)
	}

	if tunnelName != "" && tunnel.HasNamedTunnel(ctx, tunnelName) {
		logger.Info(fmt.Sprintf("Starting named tunnel: %s", tunnelName))
		tun, err = tunnel.StartNamed(ctx, tunnelName, port)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to start named tunnel: %v", err))
			return 1
		}

		// Display custom domain if configured
		customDomain := manifest.ShareTunnelDomain(mf)
		if customDomain != "" {
			logger.Success(fmt.Sprintf("Sharing at: https://%s", customDomain))
		} else {
			logger.Success(fmt.Sprintf("Tunnel started: %s", tunnelName))
		}
	} else {
		logger.Info("Starting quick tunnel...")
		tun, err = tunnel.StartQuick(ctx, port)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to start tunnel: %v", err))
			return 1
		}

		// Wait for tunnel URL
		tunnelURL := tun.WaitForURL(20 * time.Second)
		if tunnelURL != "" {
			fmt.Println()
			logger.Success(fmt.Sprintf("Sharing at: %s", tunnelURL))
		} else {
			logger.Warn("Could not capture tunnel URL (tunnel may still be starting)")
		}
	}

	logger.Info("Press Ctrl+C to stop")
	fmt.Println()

	// Wait for server to exit
	_ = serverCmd.Wait()

	return 0
}
