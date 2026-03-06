package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/registry"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/serve"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/ui"
)

func TestStartCmdNotJustvibin(t *testing.T) {
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	cmd := defaultStartCommand()
	code := cmd.run(context.Background(), []string{}, ui.New(stdout, stderr, false), logger, 0, false)
	if code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(stderr.String(), "Not a dade project") {
		t.Fatalf("expected not a project error")
	}
}

func TestStartCmdProjectNotFound(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	cmd := defaultStartCommand()
	code := cmd.run(context.Background(), []string{"nonexistent"}, ui.New(stdout, stderr, false), logger, 0, false)
	if code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(stderr.String(), "not found") {
		t.Fatalf("expected not found error")
	}
}

func TestStartCmdAlreadyRunning(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	projectDir := t.TempDir()
	marker := registry.Marker{Name: "myapp", Template: "hypertext", Port: 59999}
	markerPath := filepath.Join(projectDir, ".dade")
	if err := os.WriteFile(markerPath, []byte(`{"name":"myapp","template":"hypertext","port":59999}`), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	projectsPath := filepath.Join(baseDir, "dade", "projects.json")
	if err := os.MkdirAll(filepath.Dir(projectsPath), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	projects := map[string]registry.Project{
		"myapp": {Port: marker.Port, Path: projectDir, Template: marker.Template},
	}
	if err := registry.Save(projectsPath, projects); err != nil {
		t.Fatalf("save: %v", err)
	}

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	cmd := defaultStartCommand()
	cmd.isPortInUse = func(int) bool { return true }
	code := cmd.run(context.Background(), []string{"myapp"}, ui.New(stdout, stderr, false), logger, 0, false)
	if code != 0 {
		t.Fatalf("expected exit 0")
	}
	if !strings.Contains(stdout.String(), "already running") {
		t.Fatalf("expected already running message")
	}
}

func TestStartCmdStaticServer(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	projectDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(projectDir, ".dade"), []byte(`{"name":"myapp","template":"hypertext","port":59999}`), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	templatesDir := filepath.Join(baseDir, "dade", "templates")
	templateDir := filepath.Join(templatesDir, "hypertext")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	manifestData := `[template]
name = "hypertext"
description = "Test"

[serve]
type = "static"
`
	if err := os.WriteFile(filepath.Join(templateDir, "dade.toml"), []byte(manifestData), 0644); err != nil {
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

	staticStarted := false
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	cmd := defaultStartCommand()
	cmd.isPortInUse = func(int) bool { return false }
	cmd.startStatic = func(ctx context.Context, runner serve.CommandRunner, port int, root string) (int, error) {
		staticStarted = true
		return 1234, nil
	}
	code := cmd.run(context.Background(), []string{"myapp"}, ui.New(stdout, stderr, false), logger, 0, false)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !staticStarted {
		t.Fatalf("expected static server to start")
	}
	expectedURL := "Started: https://" + config.ProjectDomain("myapp")
	if !strings.Contains(stdout.String(), expectedURL) {
		t.Fatalf("expected success message with %s, got: %s", expectedURL, stdout.String())
	}
}
