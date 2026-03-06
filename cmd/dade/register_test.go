package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	execx "github.com/theydontwantyoutovibecode/dade/internal/exec"
	"github.com/theydontwantyoutovibecode/dade/internal/logging"
	"github.com/theydontwantyoutovibecode/dade/internal/registry"
	"github.com/theydontwantyoutovibecode/dade/internal/ui"
)

func TestRegisterCmdInvalidName(t *testing.T) {
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
	cmd := newTestRegisterCommand(t)
	code := cmd.run(context.Background(), []string{"1bad"}, ui.New(stdout, stderr, false), logger, "")
	if code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(stderr.String(), "Invalid name") {
		t.Fatalf("expected invalid name error")
	}
}

func TestRegisterCmdAlreadyRegistered(t *testing.T) {
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, ".dade"), []byte(`{"name":"myapp"}`), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	cmd := newTestRegisterCommand(t)
	code := cmd.run(context.Background(), []string{"myapp"}, ui.New(stdout, stderr, false), logger, "")
	if code != 0 {
		t.Fatalf("expected exit 0")
	}
	if !strings.Contains(stdout.String(), "already registered") {
		t.Fatalf("expected already registered message")
	}
}

func TestRegisterCmdNameConflict(t *testing.T) {
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()

	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	projectsPath := filepath.Join(baseDir, "dade", "projects.json")
	if err := os.MkdirAll(filepath.Dir(projectsPath), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	projects := map[string]registry.Project{
		"myapp": {Port: 4000, Path: "/other/path", Template: "static"},
	}
	if err := registry.Save(projectsPath, projects); err != nil {
		t.Fatalf("save: %v", err)
	}

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	cmd := newTestRegisterCommand(t)
	code := cmd.run(context.Background(), []string{"myapp"}, ui.New(stdout, stderr, false), logger, "")
	if code != 1 {
		t.Fatalf("expected exit 1")
	}
	if !strings.Contains(stderr.String(), "already used") {
		t.Fatalf("expected name conflict error")
	}
}

func TestRegisterCmdSuccess(t *testing.T) {
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()

	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "index.html"), []byte("<h1>hi</h1>"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	var registeredName string
	var markerWritten bool

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	cmd := newTestRegisterCommand(t)
	cmd.register = func(path, name string, port int, projectPath, template string) (registry.Project, error) {
		registeredName = name
		return registry.Project{Port: port, Path: projectPath, Template: template}, nil
	}
	cmd.writeMarker = func(projectDir, name, template string, port int) (registry.Marker, error) {
		markerWritten = true
		return registry.Marker{Name: name, Template: template, Port: port}, nil
	}
	code := cmd.run(context.Background(), []string{"myapp"}, ui.New(stdout, stderr, false), logger, "")
	if code != 0 {
		t.Fatalf("expected exit 0")
	}
	if registeredName != "myapp" {
		t.Fatalf("expected registered name 'myapp', got %s", registeredName)
	}
	if !markerWritten {
		t.Fatalf("expected marker to be written")
	}
	if !strings.Contains(stdout.String(), "Registered: myapp") {
		t.Fatalf("expected success message")
	}
	if !strings.Contains(stdout.String(), "Detected template type: static") {
		t.Fatalf("expected detected template message")
	}
}

func TestRegisterCmdDetectsDjango(t *testing.T) {
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()

	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "manage.py"), []byte("#!/usr/bin/env python"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	var registeredTemplate string

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	cmd := newTestRegisterCommand(t)
	cmd.register = func(path, name string, port int, projectPath, template string) (registry.Project, error) {
		registeredTemplate = template
		return registry.Project{Port: port, Path: projectPath, Template: template}, nil
	}
	code := cmd.run(context.Background(), []string{"myapp"}, ui.New(stdout, stderr, false), logger, "")
	if code != 0 {
		t.Fatalf("expected exit 0")
	}
	if registeredTemplate != "web-app" {
		t.Fatalf("expected web-app, got %s", registeredTemplate)
	}
}

func newTestRegisterCommand(t *testing.T) registerCommand {
	t.Helper()
	cmd := defaultRegisterCommand()
	cmd.generateCaddy = func(context.Context, execx.Runner, string, string) error { return nil }
	cmd.reloadProxy = func(context.Context, execx.Runner, string) error { return nil }
	cmd.nextPort = func(string) (int, error) { return 4000, nil }
	return cmd
}
