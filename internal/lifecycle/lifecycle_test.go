package lifecycle

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSetupRunnerRunCommands(t *testing.T) {
	dir := t.TempDir()

	var stdout, stderr bytes.Buffer
	runner := NewSetupRunner(dir)
	runner.WithOutput(&stdout, &stderr)

	commands := []string{
		"echo hello",
		"echo world",
	}

	ctx := context.Background()
	if err := runner.RunCommands(ctx, commands); err != nil {
		t.Fatalf("RunCommands failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "hello") || !strings.Contains(output, "world") {
		t.Fatalf("expected output to contain hello and world, got: %s", output)
	}
}

func TestSetupRunnerRunCommandsFailure(t *testing.T) {
	dir := t.TempDir()

	var stdout, stderr bytes.Buffer
	runner := NewSetupRunner(dir)
	runner.WithOutput(&stdout, &stderr)

	commands := []string{
		"echo first",
		"exit 1",
		"echo should-not-run",
	}

	ctx := context.Background()
	err := runner.RunCommands(ctx, commands)
	if err == nil {
		t.Fatal("expected error from failing command")
	}
	if !strings.Contains(err.Error(), "exit 1") {
		t.Fatalf("expected error to mention failing command, got: %v", err)
	}
}

func TestSetupRunnerRunScript(t *testing.T) {
	dir := t.TempDir()
	scriptPath := filepath.Join(dir, "setup.sh")
	if err := os.WriteFile(scriptPath, []byte("#!/bin/bash\necho script-ran\n"), 0755); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	runner := NewSetupRunner(dir)
	runner.WithOutput(&stdout, &stderr)

	ctx := context.Background()
	if err := runner.RunScript(ctx, scriptPath); err != nil {
		t.Fatalf("RunScript failed: %v", err)
	}

	if !strings.Contains(stdout.String(), "script-ran") {
		t.Fatalf("expected output to contain script-ran, got: %s", stdout.String())
	}
}

func TestSetupRunnerWithEnv(t *testing.T) {
	dir := t.TempDir()

	var stdout, stderr bytes.Buffer
	runner := NewSetupRunner(dir)
	runner.WithEnv([]string{"MY_VAR=test-value"})
	runner.WithOutput(&stdout, &stderr)

	ctx := context.Background()
	if err := runner.RunCommands(ctx, []string{"echo $MY_VAR"}); err != nil {
		t.Fatalf("RunCommands failed: %v", err)
	}

	if !strings.Contains(stdout.String(), "test-value") {
		t.Fatalf("expected env var in output, got: %s", stdout.String())
	}
}

func TestBackgroundManagerStartAndStop(t *testing.T) {
	dir := t.TempDir()

	var stdout bytes.Buffer
	mgr := NewBackgroundManager(dir)
	mgr.WithOutput(&stdout, &stdout)

	ctx := context.Background()
	proc, err := mgr.Start(ctx, "sleep 10")
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	if proc.Cmd.Process == nil {
		t.Fatal("expected process to be running")
	}

	procs := mgr.Processes()
	if len(procs) != 1 {
		t.Fatalf("expected 1 process, got %d", len(procs))
	}

	mgr.StopAll()

	// Give it a moment to clean up
	time.Sleep(100 * time.Millisecond)
}

func TestBackgroundManagerStartAll(t *testing.T) {
	dir := t.TempDir()

	var stdout bytes.Buffer
	mgr := NewBackgroundManager(dir)
	mgr.WithOutput(&stdout, &stdout)

	ctx := context.Background()
	err := mgr.StartAll(ctx, []string{"sleep 10", "sleep 10"})
	if err != nil {
		t.Fatalf("StartAll failed: %v", err)
	}

	procs := mgr.Processes()
	if len(procs) != 2 {
		t.Fatalf("expected 2 processes, got %d", len(procs))
	}

	mgr.StopAll()
}

func TestCleanupManagerRunsInReverseOrder(t *testing.T) {
	mgr := NewCleanupManager()

	var order []int
	mgr.Register(func() { order = append(order, 1) })
	mgr.Register(func() { order = append(order, 2) })
	mgr.Register(func() { order = append(order, 3) })

	mgr.Shutdown()

	if len(order) != 3 {
		t.Fatalf("expected 3 cleanups, got %d", len(order))
	}
	if order[0] != 3 || order[1] != 2 || order[2] != 1 {
		t.Fatalf("expected LIFO order [3,2,1], got %v", order)
	}
}

func TestCleanupManagerOnlyRunsOnce(t *testing.T) {
	mgr := NewCleanupManager()

	var count int
	mgr.Register(func() { count++ })

	mgr.Shutdown()
	mgr.Shutdown()

	if count != 1 {
		t.Fatalf("expected cleanup to run once, ran %d times", count)
	}
}

func TestControllerIntegration(t *testing.T) {
	dir := t.TempDir()

	var stdout, stderr bytes.Buffer
	ctrl := NewController(dir)
	ctrl.WithOutput(&stdout, &stderr)
	ctrl.WithEnv([]string{"TEST_VAR=integration"})

	ctx := context.Background()

	// Run setup
	err := ctrl.RunSetup(ctx, []string{"echo setup-ran"})
	if err != nil {
		t.Fatalf("RunSetup failed: %v", err)
	}

	if !strings.Contains(stdout.String(), "setup-ran") {
		t.Fatalf("expected setup output, got: %s", stdout.String())
	}
}
