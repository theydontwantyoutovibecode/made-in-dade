package lifecycle

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// Controller orchestrates the full dev/share lifecycle.
type Controller struct {
	ProjectDir string
	Env        []string
	Stdout     io.Writer
	Stderr     io.Writer

	setup      *SetupRunner
	background *BackgroundManager
	cleanup    *CleanupManager
	mainCmd    *exec.Cmd
}

// NewController creates a lifecycle controller for a project.
func NewController(projectDir string) *Controller {
	return &Controller{
		ProjectDir: projectDir,
		Env:        os.Environ(),
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		cleanup:    NewCleanupManager(),
	}
}

// WithEnv adds environment variables for all lifecycle phases.
func (c *Controller) WithEnv(env []string) *Controller {
	c.Env = append(c.Env, env...)
	return c
}

// WithOutput sets custom stdout/stderr for all processes.
func (c *Controller) WithOutput(stdout, stderr io.Writer) *Controller {
	c.Stdout = stdout
	c.Stderr = stderr
	return c
}

// RunSetup executes setup commands sequentially.
func (c *Controller) RunSetup(ctx context.Context, commands []string) error {
	if len(commands) == 0 {
		return nil
	}
	c.setup = NewSetupRunner(c.ProjectDir)
	c.setup.WithEnv(c.Env)
	c.setup.WithOutput(c.Stdout, c.Stderr)
	return c.setup.RunCommands(ctx, commands)
}

// RunSetupScript executes a setup script file.
func (c *Controller) RunSetupScript(ctx context.Context, scriptPath string) error {
	if scriptPath == "" {
		return nil
	}
	fullPath := scriptPath
	if !filepath.IsAbs(scriptPath) {
		fullPath = filepath.Join(c.ProjectDir, scriptPath)
	}
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("setup script not found: %s", scriptPath)
	}
	c.setup = NewSetupRunner(c.ProjectDir)
	c.setup.WithEnv(c.Env)
	c.setup.WithOutput(c.Stdout, c.Stderr)
	return c.setup.RunScript(ctx, fullPath)
}

// StartBackground launches background processes.
func (c *Controller) StartBackground(ctx context.Context, commands []string) error {
	if len(commands) == 0 {
		return nil
	}
	c.background = NewBackgroundManager(c.ProjectDir)
	c.background.WithEnv(c.Env)
	c.background.WithOutput(c.Stdout, c.Stderr)

	// Register cleanup for background processes
	c.cleanup.Register(func() {
		c.background.StopAll()
	})

	return c.background.StartAll(ctx, commands)
}

// StartServer launches the main server process.
// This blocks until the server exits or context is cancelled.
func (c *Controller) StartServer(ctx context.Context, cmdStr string, port int, portEnv string) error {
	if portEnv == "" {
		portEnv = "PORT"
	}

	cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)
	cmd.Dir = c.ProjectDir
	cmd.Env = append(c.Env, fmt.Sprintf("%s=%d", portEnv, port))
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	cmd.Stdin = os.Stdin

	// Set process group ID to ensure we can kill all child processes
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	c.mainCmd = cmd

	// Start the process
	if err := cmd.Start(); err != nil {
		return err
	}

	// Write PID file after starting (use the bash process PID)
	pidFile := filepath.Join(c.ProjectDir, ".dade.pid")
	pidStr := fmt.Sprintf("%d", cmd.Process.Pid)
	if err := os.WriteFile(pidFile, []byte(pidStr), 0644); err != nil {
		cmd.Process.Kill()
		return fmt.Errorf("failed to write PID file: %w", err)
	}

	// Register cleanup for main process and PID file
	cleanupPID := func() {
		_ = os.Remove(pidFile)
		if c.mainCmd != nil && c.mainCmd.Process != nil {
			// Kill entire process group (negative PID)
			_ = syscall.Kill(-c.mainCmd.Process.Pid, syscall.SIGTERM)
		}
	}
	c.cleanup.Register(cleanupPID)

	return cmd.Wait()
}

// StartServerBackground launches the main server in the background.
// Returns the process for monitoring.
func (c *Controller) StartServerBackground(ctx context.Context, cmdStr string, port int, portEnv string) (*exec.Cmd, error) {
	if portEnv == "" {
		portEnv = "PORT"
	}

	cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)
	cmd.Dir = c.ProjectDir
	cmd.Env = append(c.Env, fmt.Sprintf("%s=%d", portEnv, port))
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr

	c.mainCmd = cmd

	// Register cleanup for main process
	c.cleanup.Register(func() {
		if c.mainCmd != nil && c.mainCmd.Process != nil {
			_ = c.mainCmd.Process.Kill()
		}
	})

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

// HandleSignals starts signal handling for graceful shutdown.
// Call in a goroutine.
func (c *Controller) HandleSignals(ctx context.Context) {
	c.cleanup.HandleSignals(ctx)
}

// Shutdown triggers manual shutdown and cleanup.
func (c *Controller) Shutdown() {
	c.cleanup.Shutdown()
}

// RegisterCleanup adds a custom cleanup function.
func (c *Controller) RegisterCleanup(fn CleanupFunc) {
	c.cleanup.Register(fn)
}
