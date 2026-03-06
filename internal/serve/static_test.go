package serve

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

type fakeRunner struct {
	cmd *exec.Cmd
	err error
}

func (f fakeRunner) Start(_ string, _ ...string) (*exec.Cmd, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.cmd, nil
}

func TestStartStaticServerWritesPID(t *testing.T) {
	root := t.TempDir()
	cmd := exec.Command("sleep", "5")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start sleep: %v", err)
	}
	defer func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	pid, err := StartStaticServer(ctx, fakeRunner{cmd: cmd}, 8080, root)
	if err != nil {
		t.Fatalf("start static server: %v", err)
	}
	if pid != cmd.Process.Pid {
		t.Fatalf("unexpected pid")
	}

	pidFile := filepath.Join(root, DefaultPIDFile)
	if _, err := os.Stat(pidFile); err != nil {
		t.Fatalf("expected pid file")
	}
}

func TestStartStaticServerRunnerError(t *testing.T) {
	ctx := context.Background()
	_, err := StartStaticServer(ctx, fakeRunner{err: errors.New("boom")}, 8080, t.TempDir())
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestStartStaticServerContextCancel(t *testing.T) {
	root := t.TempDir()
	cmd := exec.Command("sleep", "5")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start sleep: %v", err)
	}
	defer func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := StartStaticServer(ctx, fakeRunner{cmd: cmd}, 8080, root)
	if err == nil {
		t.Fatalf("expected error")
	}

	pidFile := filepath.Join(root, DefaultPIDFile)
	if _, err := os.Stat(pidFile); err == nil {
		t.Fatalf("expected pid file to be removed")
	}
}

func TestStartStaticServerDetectsExitedProcess(t *testing.T) {
	root := t.TempDir()
	cmd := exec.Command("true")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start true: %v", err)
	}
	_ = cmd.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := StartStaticServer(ctx, fakeRunner{cmd: cmd}, 8080, root)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, syscall.ESRCH) && err.Error() == "failed to start server" {
		return
	}
}
