package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	execx "github.com/theydontwantyoutovibecode/made-in-dade/internal/exec"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/registry"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/ui"
)

func TestSyncCmdScansForProjects(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)
	t.Setenv("HOME", baseDir)

	projectDir := filepath.Join(baseDir, "myproject")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, ".dade"), []byte(`{"name":"myproject","port":4000,"template":"hypertext"}`), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	cmd := newTestSyncCommand(t)
	code := cmd.run(context.Background(), []string{baseDir}, ui.New(stdout, stderr, false), logger, false)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Found: myproject") {
		t.Fatalf("expected found message")
	}
	if !strings.Contains(stdout.String(), "Synced 1 project") {
		t.Fatalf("expected synced message")
	}
}

func TestSyncCmdCleanRemovesStale(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	projectsPath := filepath.Join(baseDir, "dade", "projects.json")
	if err := os.MkdirAll(filepath.Dir(projectsPath), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	projects := map[string]registry.Project{
		"stale": {Port: 4000, Path: "/nonexistent/path", Template: "hypertext"},
	}
	if err := registry.Save(projectsPath, projects); err != nil {
		t.Fatalf("save: %v", err)
	}

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	cmd := newTestSyncCommand(t)
	code := cmd.run(context.Background(), []string{}, ui.New(stdout, stderr, false), logger, true)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Removing stale: stale") {
		t.Fatalf("expected stale removal message")
	}
	if !strings.Contains(stdout.String(), "Removed 1 stale") {
		t.Fatalf("expected removed count")
	}
}

func TestSyncCmdEmptyDirectory(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	scanDir := t.TempDir()

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	cmd := newTestSyncCommand(t)
	code := cmd.run(context.Background(), []string{scanDir}, ui.New(stdout, stderr, false), logger, false)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Synced 0 project") {
		t.Fatalf("expected synced 0 message")
	}
}

func newTestSyncCommand(t *testing.T) syncCommand {
	t.Helper()
	cmd := defaultSyncCommand()
	cmd.generateCaddy = func(context.Context, execx.Runner, string, string) error { return nil }
	cmd.reloadProxy = func(context.Context, execx.Runner, string) error { return nil }
	return cmd
}
