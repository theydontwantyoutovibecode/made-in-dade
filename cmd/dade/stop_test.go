package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/registry"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/serve"
)

func TestStopCmdNotJustvibin(t *testing.T) {
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"project", "stop"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(stderr.String(), "Not a dade project") {
		t.Fatalf("expected not a project error")
	}
}

func TestStopCmdProjectNotFound(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"project", "stop", "nonexistent"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(stderr.String(), "not found") {
		t.Fatalf("expected not found error")
	}
}

func TestStopCmdNotRunning(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	projectDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(projectDir, ".dade"), []byte(`{"name":"myapp","template":"hypertext","port":59999}`), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	projectsPath := filepath.Join(baseDir, "dade", "projects.json")
	if err := os.MkdirAll(filepath.Dir(projectsPath), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	projects := map[string]registry.Project{
		"myapp": {Port: 59999, Path: projectDir, Template: "hypertext"},
	}
	if err := registry.Save(projectsPath, projects); err != nil {
		t.Fatalf("save: %v", err)
	}

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"project", "stop", "myapp"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
	if !strings.Contains(stdout.String(), "not running") {
		t.Fatalf("expected not running message")
	}
}

func TestStopCmdRemovesPIDFile(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	projectDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(projectDir, ".dade"), []byte(`{"name":"myapp","template":"hypertext","port":59999}`), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	pidFile := filepath.Join(projectDir, serve.DefaultPIDFile)
	if err := os.WriteFile(pidFile, []byte("99999999"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	projectsPath := filepath.Join(baseDir, "dade", "projects.json")
	if err := os.MkdirAll(filepath.Dir(projectsPath), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	projects := map[string]registry.Project{
		"myapp": {Port: 59999, Path: projectDir, Template: "hypertext"},
	}
	if err := registry.Save(projectsPath, projects); err != nil {
		t.Fatalf("save: %v", err)
	}

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"project", "stop", "myapp"})
	_ = rootCmd.Execute()

	if _, err := os.Stat(pidFile); !os.IsNotExist(err) {
		t.Fatalf("expected PID file to be removed")
	}
}
