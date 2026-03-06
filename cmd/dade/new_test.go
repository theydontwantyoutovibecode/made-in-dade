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

type fakeRunner struct {
	calls   []string
	err     error
	pathErr error
}

func (f *fakeRunner) Run(_ context.Context, name string, args ...string) error {
	f.calls = append(f.calls, strings.Join(append([]string{name}, args...), " "))
	return f.err
}

func (f *fakeRunner) Output(_ context.Context, name string, args ...string) (string, error) {
	f.calls = append(f.calls, strings.Join(append([]string{name}, args...), " "))
	return "", f.err
}

func (f *fakeRunner) LookPath(_ string) (string, error) {
	if f.pathErr != nil {
		return "", f.pathErr
	}
	return "/usr/bin/fake", nil
}

func TestNewCommandValidatesName(t *testing.T) {
	restore := withWorkDir(t)
	defer restore()
	logger := logging.New(&strings.Builder{}, &strings.Builder{}, false)
	cmd := defaultNewCommand()
	code := cmd.run(context.Background(), []string{"1bad"}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 1 {
		t.Fatalf("expected exit 1")
	}
}

func TestNewCommandUnknownTemplate(t *testing.T) {
	restore := withWorkDir(t)
	defer restore()
	logger := logging.New(&strings.Builder{}, &strings.Builder{}, false)
	templatesDir := t.TempDir()
	writePluginTemplate(t, templatesDir, "alpha", "")
	cmd := newTestCommand(t)
	cmd.templatesDir = func() (string, error) { return templatesDir, nil }
	code := cmd.run(context.Background(), []string{"proj", "--template", "unknown"}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 1 {
		t.Fatalf("expected exit 1")
	}
}

func TestNewCommandLocalCopy(t *testing.T) {
	restore := withWorkDir(t)
	defer restore()
	logger := logging.New(&strings.Builder{}, &strings.Builder{}, false)
	cmd := newTestCommand(t)
	local := t.TempDir()
	if err := os.WriteFile(filepath.Join(local, "hello.txt"), []byte("hello"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	code := cmd.run(context.Background(), []string{"proj", "--local", local}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 0 {
		t.Fatalf("expected exit 0")
	}
	if _, err := os.Stat(filepath.Join("proj", "hello.txt")); err != nil {
		t.Fatalf("expected copied file: %v", err)
	}
}

func TestNewCommandSetupSkippedWhenNonInteractive(t *testing.T) {
	restore := withWorkDir(t)
	defer restore()
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	templatesDir := t.TempDir()
	writePluginTemplate(t, templatesDir, "alpha", strings.Join([]string{
		"[scaffold]",
		"setup = \"echo hi\"",
		"setup_interactive = true",
		"",
	}, "\n"))
	cmd := newTestCommand(t)
	cmd.templatesDir = func() (string, error) { return templatesDir, nil }
	code := cmd.run(context.Background(), []string{"proj", "--template", "alpha"}, ui.New(stdout, stderr, false), logger, false)
	if code != 0 {
		t.Fatalf("expected exit 0")
	}
	if !strings.Contains(stdout.String(), "Skipping setup (requires a TTY)") {
		t.Fatalf("expected setup skip message")
	}
}

func TestNewCmdUsesNameFlag(t *testing.T) {
	restore := withWorkDir(t)
	defer restore()
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	local := t.TempDir()
	if err := os.WriteFile(filepath.Join(local, "hello.txt"), []byte("hello"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	newCommandFactory = func() newCommand {
		return newTestCommand(t)
	}
	defer func() { newCommandFactory = defaultNewCommand }()

	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"new", "--name", "proj", "--local", local})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
}

func TestNewCmdHonorsTemplateFlag(t *testing.T) {
	restore := withWorkDir(t)
	defer restore()
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	local := t.TempDir()
	if err := os.WriteFile(filepath.Join(local, "hello.txt"), []byte("hello"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	newCommandFactory = func() newCommand {
		return newTestCommand(t)
	}
	defer func() { newCommandFactory = defaultNewCommand }()

	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"new", "proj", "--template", "hypertext", "--local", local})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
}

func TestNewCmdMissingNameNonTTY(t *testing.T) {
	restore := withWorkDir(t)
	defer restore()
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	newCommandFactory = func() newCommand {
		return newTestCommand(t)
	}
	defer func() { newCommandFactory = defaultNewCommand }()

	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"new"})

	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected error for missing project name")
	}
}

func newTestCommand(t *testing.T) newCommand {
	t.Helper()
	cmd := defaultNewCommand()
	cmd.runner = &fakeRunner{}
	cmd.removeGitDir = func(string) error { return nil }
	cmd.nextPort = func(string) (int, error) { return 4000, nil }
	cmd.projectsFile = func() (string, error) { return filepath.Join(t.TempDir(), "projects.json"), nil }
	cmd.caddyfilePath = func() (string, error) { return filepath.Join(t.TempDir(), "Caddyfile"), nil }
	cmd.generateCaddy = func(context.Context, execx.Runner, string, string) error { return nil }
	cmd.reloadProxy = func(context.Context, execx.Runner, string) error { return nil }
	cmd.register = nil
	cmd.writeMarker = nil
	cmd.migrateSrv = nil
	cmd.spin = func(_ string, work func() error) error { return work() }
	return cmd
}

func writePluginTemplate(t *testing.T, templatesDir, name, extra string) string {
	t.Helper()
	templateDir := filepath.Join(templatesDir, name)
	if err := os.MkdirAll(templateDir, 0o755); err != nil {
		t.Fatalf("failed to create template dir: %v", err)
	}
	manifest := strings.Join([]string{
		"[template]",
		"name = \"" + name + "\"",
		"description = \"Test template\"",
		"",
		"[serve]",
		"type = \"static\"",
		"",
	}, "\n")
	if strings.TrimSpace(extra) != "" {
		manifest = strings.Join([]string{manifest, extra}, "\n")
	}
	if err := os.WriteFile(filepath.Join(templateDir, "dade.toml"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("failed to write manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templateDir, "hello.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatalf("failed to write template file: %v", err)
	}
	return templateDir
}

func withWorkDir(t *testing.T) func() {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	temp := t.TempDir()
	if err := os.Chdir(temp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	return func() {
		_ = os.Chdir(cwd)
	}
}

func TestNewCommandNoTemplatesInstalled(t *testing.T) {
	restore := withWorkDir(t)
	defer restore()
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	templatesDir := t.TempDir()
	cmd := newTestCommand(t)
	cmd.templatesDir = func() (string, error) { return templatesDir, nil }
	code := cmd.run(context.Background(), []string{"proj"}, ui.New(stdout, stderr, false), logger, false)
	if code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(stderr.String(), "No templates installed") {
		t.Fatalf("expected no templates error")
	}
	if !strings.Contains(stdout.String(), "dade install --list-official") {
		t.Fatalf("expected install guidance")
	}
}

func TestNewCommandExcludePatterns(t *testing.T) {
	restore := withWorkDir(t)
	defer restore()
	logger := logging.New(&strings.Builder{}, &strings.Builder{}, false)
	templatesDir := t.TempDir()
	writePluginTemplate(t, templatesDir, "alpha", strings.Join([]string{
		"[scaffold]",
		"exclude = [\"secret.txt\", \"build\"]",
		"",
	}, "\n"))
	templateDir := filepath.Join(templatesDir, "alpha")
	if err := os.WriteFile(filepath.Join(templateDir, "secret.txt"), []byte("secret"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(templateDir, "build"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templateDir, "build", "output.js"), []byte("output"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	cmd := newTestCommand(t)
	cmd.templatesDir = func() (string, error) { return templatesDir, nil }
	code := cmd.run(context.Background(), []string{"proj", "--template", "alpha"}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 0 {
		t.Fatalf("expected exit 0")
	}
	if _, err := os.Stat(filepath.Join("proj", "hello.txt")); err != nil {
		t.Fatalf("expected hello.txt copied: %v", err)
	}
	if _, err := os.Stat(filepath.Join("proj", "secret.txt")); !os.IsNotExist(err) {
		t.Fatalf("expected secret.txt excluded")
	}
	if _, err := os.Stat(filepath.Join("proj", "build")); !os.IsNotExist(err) {
		t.Fatalf("expected build dir excluded")
	}
}

func TestNewCommandRegistersProject(t *testing.T) {
	restore := withWorkDir(t)
	defer restore()
	logger := logging.New(&strings.Builder{}, &strings.Builder{}, false)
	templatesDir := t.TempDir()
	writePluginTemplate(t, templatesDir, "alpha", "")
	projectsPath := filepath.Join(t.TempDir(), "projects.json")
	var registeredName string
	var registeredPort int
	cmd := newTestCommand(t)
	cmd.templatesDir = func() (string, error) { return templatesDir, nil }
	cmd.projectsFile = func() (string, error) { return projectsPath, nil }
	cmd.register = func(path, name string, port int, projectPath, template string) (registry.Project, error) {
		registeredName = name
		registeredPort = port
		return registry.Project{Port: port, Path: projectPath, Template: template}, nil
	}
	code := cmd.run(context.Background(), []string{"proj", "--template", "alpha"}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 0 {
		t.Fatalf("expected exit 0")
	}
	if registeredName != "proj" {
		t.Fatalf("expected registered name 'proj', got %s", registeredName)
	}
	if registeredPort != 4000 {
		t.Fatalf("expected registered port 4000, got %d", registeredPort)
	}
}

func TestNewCommandWritesMarker(t *testing.T) {
	restore := withWorkDir(t)
	defer restore()
	logger := logging.New(&strings.Builder{}, &strings.Builder{}, false)
	templatesDir := t.TempDir()
	writePluginTemplate(t, templatesDir, "alpha", "")
	var markerDir, markerTemplate string
	var markerPort int
	cmd := newTestCommand(t)
	cmd.templatesDir = func() (string, error) { return templatesDir, nil }
	cmd.writeMarker = func(projectDir, name, template string, port int) (registry.Marker, error) {
		markerDir = projectDir
		markerTemplate = template
		markerPort = port
		return registry.Marker{Name: name, Template: template, Port: port}, nil
	}
	code := cmd.run(context.Background(), []string{"proj", "--template", "alpha"}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 0 {
		t.Fatalf("expected exit 0")
	}
	if markerDir != "proj" {
		t.Fatalf("expected marker dir 'proj', got %s", markerDir)
	}
	if markerTemplate != "alpha" {
		t.Fatalf("expected marker template 'alpha', got %s", markerTemplate)
	}
	if markerPort != 4000 {
		t.Fatalf("expected marker port 4000, got %d", markerPort)
	}
}

func TestNewCommandPickerShownForMultipleTemplates(t *testing.T) {
	restore := withWorkDir(t)
	defer restore()
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	templatesDir := t.TempDir()
	writePluginTemplate(t, templatesDir, "alpha", "")
	writePluginTemplate(t, templatesDir, "beta", "")
	cmd := newTestCommand(t)
	cmd.templatesDir = func() (string, error) { return templatesDir, nil }
	code := cmd.run(context.Background(), []string{"proj"}, ui.New(stdout, stderr, false), logger, false)
	if code != 1 {
		t.Fatalf("expected exit 1 (non-interactive picker fails)")
	}
	if !strings.Contains(stderr.String(), "template selection requires a TTY") {
		t.Fatalf("expected picker TTY error")
	}
}

var _ = execx.Runner(&fakeRunner{})
