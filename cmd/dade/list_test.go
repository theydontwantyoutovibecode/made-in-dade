package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/registry"
)

func TestListCmdEmptyRegistry(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"project", "list"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
	if !strings.Contains(stdout.String(), "No projects registered") {
		t.Fatalf("expected empty state message")
	}
}

func TestListCmdShowsProjects(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	projectsPath := filepath.Join(baseDir, "dade", "projects.json")
	if err := os.MkdirAll(filepath.Dir(projectsPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	projects := map[string]registry.Project{
		"myapp": {Port: 4000, Path: "/tmp/myapp", Template: "hypertext"},
	}
	if err := registry.Save(projectsPath, projects); err != nil {
		t.Fatalf("save: %v", err)
	}

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"project", "list"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
	output := stdout.String()
	if !strings.Contains(output, "myapp") {
		t.Fatalf("expected project name in output")
	}
	expectedURL := "https://" + config.ProjectDomain("myapp")
	if !strings.Contains(output, expectedURL) {
		t.Fatalf("expected URL %s in output, got: %s", expectedURL, output)
	}
	if !strings.Contains(output, "hypertext") {
		t.Fatalf("expected template in output")
	}
}

func TestListCmdJSONOutput(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	projectsPath := filepath.Join(baseDir, "dade", "projects.json")
	if err := os.MkdirAll(filepath.Dir(projectsPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	projects := map[string]registry.Project{
		"app1": {Port: 4000, Path: "/tmp/app1", Template: "hypertext"},
		"app2": {Port: 4001, Path: "/tmp/app2", Template: "django"},
	}
	if err := registry.Save(projectsPath, projects); err != nil {
		t.Fatalf("save: %v", err)
	}

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"project", "list", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
	var payload []listProject
	if err := json.Unmarshal([]byte(stdout.String()), &payload); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(payload) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(payload))
	}
}

func TestListCmdRunningFilter(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	projectsPath := filepath.Join(baseDir, "dade", "projects.json")
	if err := os.MkdirAll(filepath.Dir(projectsPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	projects := map[string]registry.Project{
		"myapp": {Port: 59999, Path: "/tmp/myapp", Template: "hypertext"},
	}
	if err := registry.Save(projectsPath, projects); err != nil {
		t.Fatalf("save: %v", err)
	}

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"project", "list", "--running"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
	if !strings.Contains(stdout.String(), "No running projects") {
		t.Fatalf("expected no running projects message")
	}
}
