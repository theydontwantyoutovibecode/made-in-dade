package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/registry"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/ui"
)

type uninstallStub struct {
	removed []string
}

func (s *uninstallStub) removeAll(path string) error {
	s.removed = append(s.removed, path)
	return nil
}

func TestUninstallMissingName(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	cmd := defaultUninstallCommand()
	cmd.templatesDir = func() (string, error) { return t.TempDir(), nil }
	cmd.confirm = func(string) (bool, error) { return false, nil }

	code := cmd.run(context.Background(), []string{}, ui.New(stdout, stderr, false), logging.New(stdout, stderr, false))
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "Template name is required") {
		t.Fatalf("expected missing name error")
	}
}

func TestUninstallMissingTemplateListsInstalled(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	templatesDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(templatesDir, "alpha"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cmd := defaultUninstallCommand()
	cmd.templatesDir = func() (string, error) { return templatesDir, nil }
	cmd.confirm = func(string) (bool, error) { return false, nil }

	code := cmd.run(context.Background(), []string{"missing"}, ui.New(stdout, stderr, false), logging.New(stdout, stderr, false))
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "Template 'missing' not found") {
		t.Fatalf("expected not found error")
	}
	if !strings.Contains(stdout.String(), "Installed templates:") {
		t.Fatalf("expected installed templates header")
	}
	if !strings.Contains(stdout.String(), "  - alpha") {
		t.Fatalf("expected installed template name")
	}
}

func TestUninstallWarnsAndAbortsWhenInUse(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	templatesDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(templatesDir, "alpha"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	projectsPath := filepath.Join(t.TempDir(), "projects.json")
	projects := map[string]registry.Project{
		"proj1": {Template: "alpha"},
	}
	if err := registry.Save(projectsPath, projects); err != nil {
		t.Fatalf("save: %v", err)
	}
	stub := &uninstallStub{}
	cmd := defaultUninstallCommand()
	cmd.templatesDir = func() (string, error) { return templatesDir, nil }
	cmd.projectsFile = func() (string, error) { return projectsPath, nil }
	cmd.loadProjects = registry.Load
	cmd.removeAll = stub.removeAll
	cmd.confirm = func(string) (bool, error) { return false, nil }

	code := cmd.run(context.Background(), []string{"alpha"}, ui.New(stdout, stderr, false), logging.New(stdout, stderr, false))
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if len(stub.removed) != 0 {
		t.Fatalf("expected no removal")
	}
	if !strings.Contains(stdout.String(), "Projects using this template") {
		t.Fatalf("expected warning about projects")
	}
	if !strings.Contains(stdout.String(), "proj1") {
		t.Fatalf("expected project name")
	}
}

func TestUninstallRemovesTemplate(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	templatesDir := t.TempDir()
	path := filepath.Join(templatesDir, "alpha")
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	projectsPath := filepath.Join(t.TempDir(), "projects.json")
	if err := registry.Save(projectsPath, map[string]registry.Project{}); err != nil {
		t.Fatalf("save: %v", err)
	}
	stub := &uninstallStub{}
	cmd := defaultUninstallCommand()
	cmd.templatesDir = func() (string, error) { return templatesDir, nil }
	cmd.projectsFile = func() (string, error) { return projectsPath, nil }
	cmd.loadProjects = registry.Load
	cmd.removeAll = stub.removeAll
	cmd.confirm = func(string) (bool, error) { return true, nil }

	code := cmd.run(context.Background(), []string{"alpha"}, ui.New(stdout, stderr, false), logging.New(stdout, stderr, false))
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if len(stub.removed) != 1 || stub.removed[0] != path {
		t.Fatalf("expected template removal")
	}
	if !strings.Contains(stdout.String(), "Removed template: alpha") {
		t.Fatalf("expected success message")
	}
}

func TestUninstallAllCancelled(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	stub := &uninstallStub{}
	cmd := defaultUninstallCommand()
	cmd.templatesDir = func() (string, error) { return t.TempDir(), nil }
	cmd.readDir = func(string) ([]os.DirEntry, error) { return nil, errors.New("unexpected") }
	cmd.removeAll = stub.removeAll
	cmd.confirm = func(string) (bool, error) { return false, nil }

	code := cmd.run(context.Background(), []string{"--all"}, ui.New(stdout, stderr, false), logging.New(stdout, stderr, false))
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if len(stub.removed) != 0 {
		t.Fatalf("expected no removal")
	}
}

func TestUninstallAllRemovesEntries(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	templatesDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(templatesDir, "alpha"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(templatesDir, "beta"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	stub := &uninstallStub{}
	cmd := defaultUninstallCommand()
	cmd.templatesDir = func() (string, error) { return templatesDir, nil }
	cmd.readDir = os.ReadDir
	cmd.removeAll = stub.removeAll
	cmd.confirm = func(string) (bool, error) { return true, nil }

	code := cmd.run(context.Background(), []string{"--all"}, ui.New(stdout, stderr, false), logging.New(stdout, stderr, false))
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if len(stub.removed) != 2 {
		t.Fatalf("expected all templates removed")
	}
	if !strings.Contains(stdout.String(), "All templates removed") {
		t.Fatalf("expected success message")
	}
}
