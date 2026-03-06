package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/registry"
)

func resetRemoveFlags(t *testing.T) {
	t.Helper()
	resetRootFlags(t)
	if f := projectRemoveCmd.Flags().Lookup("files"); f != nil {
		_ = f.Value.Set("false")
		f.Changed = false
	}
	if f := projectRemoveCmd.Flags().Lookup("yes"); f != nil {
		_ = f.Value.Set("false")
		f.Changed = false
	}
}

func TestRemoveCmdNotJustvibin(t *testing.T) {
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRemoveFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"project", "remove"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(stderr.String(), "Not a dade project") {
		t.Fatalf("expected not a project error")
	}
}

func TestRemoveCmdProjectNotFound(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRemoveFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"project", "remove", "nonexistent"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(stderr.String(), "not found") {
		t.Fatalf("expected not found error")
	}
}

func TestRemoveCmdSuccess(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	projectDir := t.TempDir()
	markerPath := filepath.Join(projectDir, ".dade")
	if err := os.WriteFile(markerPath, []byte(`{"name":"myapp","template":"hypertext","port":59999}`), 0644); err != nil {
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
	resetRemoveFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"project", "remove", "myapp"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
	if !strings.Contains(stdout.String(), "Removed: myapp") {
		t.Fatalf("expected success message")
	}

	loadedProjects, _ := registry.Load(projectsPath)
	if _, ok := loadedProjects["myapp"]; ok {
		t.Fatalf("expected project to be removed from registry")
	}

	if _, err := os.Stat(markerPath); !os.IsNotExist(err) {
		t.Fatalf("expected marker to be removed")
	}
}

func TestRemoveCmdAlias(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	projectDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(projectDir, ".dade"), []byte(`{"name":"myapp"}`), 0644); err != nil {
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
	resetRemoveFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"project", "rm", "myapp"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
	if !strings.Contains(stdout.String(), "Removed: myapp") {
		t.Fatalf("expected success message")
	}
}

func TestRemoveCmdWithFilesFlag(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	projectDir := t.TempDir()
	testFile := filepath.Join(projectDir, "testfile.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, ".dade"), []byte(`{"name":"myapp","port":59999}`), 0644); err != nil {
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
	resetRemoveFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"project", "remove", "myapp", "--files", "-y"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
	if !strings.Contains(stdout.String(), "Deleted") {
		t.Fatalf("expected deleted message, got: %s", stdout.String())
	}

	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		t.Fatalf("expected project directory to be deleted")
	}
}

func TestRemoveCmdWithFilesFlagCancelled(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	projectDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(projectDir, ".dade"), []byte(`{"name":"myapp","port":59999}`), 0644); err != nil {
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

	stdin := strings.NewReader("n\n")

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRemoveFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetIn(stdin)
	rootCmd.SetArgs([]string{"project", "remove", "myapp", "--files"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
	if !strings.Contains(stdout.String(), "Cancelled") {
		t.Fatalf("expected cancelled message, got: %s", stdout.String())
	}

	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Fatalf("expected project directory to remain")
	}
}
