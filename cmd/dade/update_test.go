package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/theydontwantyoutovibecode/dade/internal/logging"
	"github.com/theydontwantyoutovibecode/dade/internal/ui"
)

type updateRunner struct {
	err error
}

func (r updateRunner) Run(_ context.Context, _ string, _ ...string) error {
	return r.err
}

func (r updateRunner) Output(_ context.Context, _ string, _ ...string) (string, error) {
	return "", r.err
}

func (r updateRunner) LookPath(_ string) (string, error) {
	if r.err != nil {
		return "", r.err
	}
	return "/usr/bin/git", nil
}

func TestUpdateCommandMissingTemplate(t *testing.T) {
	logger := logging.New(&strings.Builder{}, &strings.Builder{}, false)
	cmd := defaultUpdateCommand()
	cmd.runner = updateRunner{}
	cmd.templatesDir = func() (string, error) { return t.TempDir(), nil }

	code := cmd.run(context.Background(), []string{"nonexistent"}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
}

func TestUpdateCommandMissingSourceFile(t *testing.T) {
	templatesDir := t.TempDir()
	templateDir := filepath.Join(templatesDir, "mytemplate")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	logger := logging.New(&strings.Builder{}, &strings.Builder{}, false)
	cmd := defaultUpdateCommand()
	cmd.runner = updateRunner{}
	cmd.templatesDir = func() (string, error) { return templatesDir, nil }

	code := cmd.run(context.Background(), []string{"mytemplate"}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
}

func TestUpdateCommandSuccess(t *testing.T) {
	templatesDir := t.TempDir()
	templateDir := filepath.Join(templatesDir, "mytemplate")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templateDir, ".source"), []byte("https://example.com/repo.git"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	stdout := &strings.Builder{}
	logger := logging.New(stdout, &strings.Builder{}, false)
	cmd := defaultUpdateCommand()
	cmd.runner = updateRunner{}
	cmd.templatesDir = func() (string, error) { return templatesDir, nil }
	cmd.tempDir = func(_, _ string) (string, error) { return t.TempDir(), nil }
	cmd.rename = func(_, newPath string) error {
		return os.MkdirAll(newPath, 0755)
	}
	cmd.removeAll = func(string) error { return nil }

	code := cmd.run(context.Background(), []string{"mytemplate"}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Updated: mytemplate") {
		t.Fatalf("expected success message, got: %s", stdout.String())
	}
}

func TestUpdateCommandAllNoTemplates(t *testing.T) {
	templatesDir := t.TempDir()

	stdout := &strings.Builder{}
	logger := logging.New(stdout, &strings.Builder{}, false)
	cmd := defaultUpdateCommand()
	cmd.runner = updateRunner{}
	cmd.templatesDir = func() (string, error) { return templatesDir, nil }

	code := cmd.run(context.Background(), []string{"--all"}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "No templates to update") {
		t.Fatalf("expected no templates message, got: %s", stdout.String())
	}
}

func TestUpdateCommandAllWithTemplates(t *testing.T) {
	templatesDir := t.TempDir()
	for _, name := range []string{"template1", "template2"} {
		dir := filepath.Join(templatesDir, name)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, ".source"), []byte("https://example.com/"+name+".git"), 0644); err != nil {
			t.Fatalf("write: %v", err)
		}
	}

	stdout := &strings.Builder{}
	logger := logging.New(stdout, &strings.Builder{}, false)
	cmd := defaultUpdateCommand()
	cmd.runner = updateRunner{}
	cmd.templatesDir = func() (string, error) { return templatesDir, nil }
	cmd.tempDir = func(_, _ string) (string, error) { return t.TempDir(), nil }
	cmd.rename = func(_, newPath string) error {
		return os.MkdirAll(newPath, 0755)
	}
	cmd.removeAll = func(string) error { return nil }

	code := cmd.run(context.Background(), []string{"--all"}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Updated 2 template(s)") {
		t.Fatalf("expected updated message, got: %s", stdout.String())
	}
}

func resetUpdateFlags(t *testing.T) {
	t.Helper()
	resetRootFlags(t)
	if f := templateUpdateCmd.Flags().Lookup("all"); f != nil {
		_ = f.Value.Set("false")
		f.Changed = false
	}
}

func TestUpdateCmdMissingName(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetUpdateFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"template", "update"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(stderr.String(), "Missing template name") {
		t.Fatalf("expected missing name error, got: %s", stderr.String())
	}
}
