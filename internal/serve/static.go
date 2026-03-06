package serve

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

type CommandRunner interface {
	Start(name string, args ...string) (*exec.Cmd, error)
}

type SystemRunner struct{}

func (SystemRunner) Start(name string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(name, args...)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

const DefaultPIDFile = ".dade.pid"

func StartStaticServer(ctx context.Context, runner CommandRunner, port int, root string) (int, error) {
	if runner == nil {
		runner = SystemRunner{}
	}
	if root == "" {
		root = "."
	}

	cmd, err := runner.Start("caddy", "file-server", "--listen", fmt.Sprintf(":%d", port), "--root", root)
	if err != nil {
		return 0, err
	}
	pid := cmd.Process.Pid
	pidFile := filepath.Join(root, DefaultPIDFile)
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		_ = cmd.Process.Kill()
		return 0, err
	}

	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		_ = cmd.Process.Kill()
		_ = os.Remove(pidFile)
		return 0, ctx.Err()
	case <-timer.C:
	}

	if err := cmd.Process.Signal(syscall.Signal(0)); err != nil {
		_ = os.Remove(pidFile)
		return 0, errors.New("failed to start server")
	}
	return pid, nil
}
