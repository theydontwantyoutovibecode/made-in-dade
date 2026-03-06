package lifecycle

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// SetupRunner executes setup commands sequentially.
type SetupRunner struct {
	Dir    string
	Env    []string
	Stdout io.Writer
	Stderr io.Writer
}

// NewSetupRunner creates a SetupRunner for the given project directory.
func NewSetupRunner(dir string) *SetupRunner {
	return &SetupRunner{
		Dir:    dir,
		Env:    os.Environ(),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// WithEnv adds environment variables to the runner.
// Each entry should be in KEY=VALUE format.
func (r *SetupRunner) WithEnv(env []string) *SetupRunner {
	r.Env = append(r.Env, env...)
	return r
}

// WithOutput sets custom stdout/stderr writers.
func (r *SetupRunner) WithOutput(stdout, stderr io.Writer) *SetupRunner {
	r.Stdout = stdout
	r.Stderr = stderr
	return r
}

// RunCommands executes a list of shell commands sequentially.
// Stops and returns error on first failure.
func (r *SetupRunner) RunCommands(ctx context.Context, commands []string) error {
	for _, cmdStr := range commands {
		if err := r.runSingleCommand(ctx, cmdStr); err != nil {
			return fmt.Errorf("command %q failed: %w", cmdStr, err)
		}
	}
	return nil
}

// RunScript executes a shell script file.
func (r *SetupRunner) RunScript(ctx context.Context, scriptPath string) error {
	cmd := exec.CommandContext(ctx, "bash", scriptPath)
	cmd.Dir = r.Dir
	cmd.Env = r.Env
	cmd.Stdout = r.Stdout
	cmd.Stderr = r.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (r *SetupRunner) runSingleCommand(ctx context.Context, cmdStr string) error {
	cmdStr = strings.TrimSpace(cmdStr)
	if cmdStr == "" {
		return nil
	}

	cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)
	cmd.Dir = r.Dir
	cmd.Env = r.Env
	cmd.Stdout = r.Stdout
	cmd.Stderr = r.Stderr
	return cmd.Run()
}
