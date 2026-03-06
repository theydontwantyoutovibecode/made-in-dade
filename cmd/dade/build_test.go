package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/registry"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/ui"
)

func stubBuildCommand() buildCommand {
	return buildCommand{
		templatesDir: func() (string, error) { return "/tmp/nonexistent", nil },
		projectsFile: func() (string, error) { return "/tmp/nonexistent/projects.json", nil },
		readMarker:   registry.ReadMarker,
		readFile:     os.ReadFile,
		fileExists:   defaultFileExists,
		globMatch:    filepath.Glob,
		runCmd: func(ctx context.Context, dir, name string, args []string, env []string, so, se *strings.Builder) error {
			return nil
		},
		runShell: func(ctx context.Context, dir, cmdStr string, extraEnv []string) error {
			return nil
		},
	}
}

func TestBuildCmdProjectNotFound(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	cmd := defaultBuildCommand()
	code := cmd.run(context.Background(), []string{"nonexistent"}, ui.New(stdout, stderr, false), logger, buildOptions{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "not found") {
		t.Fatalf("expected not found error, got: %s", stderr.String())
	}
}

func TestBuildCmdAutoDetectGo(t *testing.T) {
	projectDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte("module example.com/myapp\n\ngo 1.22\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	var ranGoCmd bool
	var capturedArgs []string
	var capturedEnv []string

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)

	cmd := stubBuildCommand()
	cmd.runCmd = func(ctx context.Context, dir, name string, args []string, env []string, so, se *strings.Builder) error {
		if name == "go" {
			ranGoCmd = true
			capturedArgs = args
			capturedEnv = env
		}
		return nil
	}

	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()
	if err := os.Chdir(projectDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	code := cmd.run(context.Background(), []string{}, ui.New(stdout, stderr, false), logger, buildOptions{})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr.String())
	}
	if !ranGoCmd {
		t.Fatalf("expected go build to run")
	}
	if capturedArgs[0] != "build" {
		t.Fatalf("expected 'build' arg, got: %v", capturedArgs)
	}

	hasGOOS := false
	for _, e := range capturedEnv {
		if strings.HasPrefix(e, "GOOS=") {
			hasGOOS = true
		}
	}
	if !hasGOOS {
		t.Fatalf("expected GOOS in env")
	}
}

func TestBuildCmdAutoDetectGoRelease(t *testing.T) {
	projectDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte("module example.com/myapp\n\ngo 1.22\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	var capturedArgs []string

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)

	cmd := stubBuildCommand()
	cmd.runCmd = func(ctx context.Context, dir, name string, args []string, env []string, so, se *strings.Builder) error {
		capturedArgs = args
		return nil
	}

	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()
	_ = os.Chdir(projectDir)

	code := cmd.run(context.Background(), []string{}, ui.New(stdout, stderr, false), logger, buildOptions{release: true})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	found := false
	for _, a := range capturedArgs {
		if a == "-ldflags" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected -ldflags in release build args: %v", capturedArgs)
	}
}

func TestBuildCmdAutoDetectGoCrossCompile(t *testing.T) {
	projectDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte("module example.com/myapp\n\ngo 1.22\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	var capturedEnv []string

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)

	cmd := stubBuildCommand()
	cmd.runCmd = func(ctx context.Context, dir, name string, args []string, env []string, so, se *strings.Builder) error {
		capturedEnv = env
		return nil
	}

	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()
	_ = os.Chdir(projectDir)

	code := cmd.run(context.Background(), []string{}, ui.New(stdout, stderr, false), logger, buildOptions{
		targetOS:   "linux",
		targetArch: "amd64",
	})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	hasLinuxGOOS := false
	hasAmd64GOARCH := false
	for _, e := range capturedEnv {
		if e == "GOOS=linux" {
			hasLinuxGOOS = true
		}
		if e == "GOARCH=amd64" {
			hasAmd64GOARCH = true
		}
	}
	if !hasLinuxGOOS || !hasAmd64GOARCH {
		t.Fatalf("expected GOOS=linux and GOARCH=amd64 in env: %v", capturedEnv)
	}
}

func TestBuildCmdNoProjectDetected(t *testing.T) {
	projectDir := t.TempDir()

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)

	cmd := stubBuildCommand()

	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()
	_ = os.Chdir(projectDir)

	code := cmd.run(context.Background(), []string{}, ui.New(stdout, stderr, false), logger, buildOptions{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "Could not detect project type") {
		t.Fatalf("expected detection error, got: %s", stderr.String())
	}
}

func TestBuildCmdFromManifest(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	projectDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(projectDir, ".dade"), []byte(`{"name":"myapp","template":"cli","port":59999}`), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	templatesDir := filepath.Join(baseDir, "dade", "templates")
	templateDir := filepath.Join(templatesDir, "cli")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	manifestData := `[template]
name = "cli"
description = "CLI application"

[serve]
type = "command"
dev = "go run ."

[build]
command = "go build -o {{output}}/{{name}} ."
output = "./bin"
release_flags = "-ldflags '-s -w'"
pre = ["go mod tidy"]
post = ["echo done"]
`
	if err := os.WriteFile(filepath.Join(templateDir, "dade.toml"), []byte(manifestData), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	projectsPath := filepath.Join(baseDir, "dade", "projects.json")
	if err := os.MkdirAll(filepath.Dir(projectsPath), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	projects := map[string]registry.Project{
		"myapp": {Port: 59999, Path: projectDir, Template: "cli"},
	}
	if err := registry.Save(projectsPath, projects); err != nil {
		t.Fatalf("save: %v", err)
	}

	var shellCmds []string

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)

	cmd := defaultBuildCommand()
	cmd.runShell = func(ctx context.Context, dir, cmdStr string, extraEnv []string) error {
		shellCmds = append(shellCmds, cmdStr)
		return nil
	}

	code := cmd.run(context.Background(), []string{"myapp"}, ui.New(stdout, stderr, false), logger, buildOptions{})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s; stdout: %s", code, stderr.String(), stdout.String())
	}

	if !strings.Contains(stdout.String(), "Built for") {
		t.Fatalf("expected success message, got: %s", stdout.String())
	}

	if len(shellCmds) < 3 {
		t.Fatalf("expected at least 3 shell commands (pre + build + post), got %d: %v", len(shellCmds), shellCmds)
	}

	if shellCmds[0] != "go mod tidy" {
		t.Fatalf("expected pre-build 'go mod tidy', got: %s", shellCmds[0])
	}

	if shellCmds[len(shellCmds)-1] != "echo done" {
		t.Fatalf("expected post-build 'echo done', got: %s", shellCmds[len(shellCmds)-1])
	}
}

func TestBuildCmdCustomOutputDir(t *testing.T) {
	projectDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte("module example.com/myapp\n\ngo 1.22\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	var capturedArgs []string

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)

	cmd := stubBuildCommand()
	cmd.runCmd = func(ctx context.Context, dir, name string, args []string, env []string, so, se *strings.Builder) error {
		capturedArgs = args
		return nil
	}

	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()
	_ = os.Chdir(projectDir)

	code := cmd.run(context.Background(), []string{}, ui.New(stdout, stderr, false), logger, buildOptions{outputDir: "./dist"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	hasDistPath := false
	for _, a := range capturedArgs {
		if strings.Contains(a, "dist") {
			hasDistPath = true
		}
	}
	if !hasDistPath {
		t.Fatalf("expected custom output dir in args: %v", capturedArgs)
	}
}

func TestBuildCmdAutoDetectXcode(t *testing.T) {
	projectDir := t.TempDir()

	xcodeDir := filepath.Join(projectDir, "MyApp.xcodeproj")
	if err := os.MkdirAll(xcodeDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	var ranXcodebuild bool
	var capturedArgs []string

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)

	cmd := stubBuildCommand()
	cmd.runCmd = func(ctx context.Context, dir, name string, args []string, env []string, so, se *strings.Builder) error {
		if name == "xcodebuild" {
			ranXcodebuild = true
			capturedArgs = args
		}
		return nil
	}

	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()
	_ = os.Chdir(projectDir)

	code := cmd.run(context.Background(), []string{}, ui.New(stdout, stderr, false), logger, buildOptions{})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr.String())
	}
	if !ranXcodebuild {
		t.Fatalf("expected xcodebuild to run")
	}

	hasScheme := false
	for i, a := range capturedArgs {
		if a == "-scheme" && i+1 < len(capturedArgs) && capturedArgs[i+1] == "MyApp" {
			hasScheme = true
		}
	}
	if !hasScheme {
		t.Fatalf("expected -scheme MyApp in args: %v", capturedArgs)
	}
}

func TestBuildCmdAutoDetectGradle(t *testing.T) {
	projectDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(projectDir, "build.gradle.kts"), []byte(""), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	var ranGradle bool

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)

	cmd := stubBuildCommand()
	cmd.runCmd = func(ctx context.Context, dir, name string, args []string, env []string, so, se *strings.Builder) error {
		if name == "gradle" || strings.HasSuffix(name, "gradlew") {
			ranGradle = true
		}
		return nil
	}

	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()
	_ = os.Chdir(projectDir)

	code := cmd.run(context.Background(), []string{}, ui.New(stdout, stderr, false), logger, buildOptions{})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr.String())
	}
	if !ranGradle {
		t.Fatalf("expected gradle to run")
	}
}


