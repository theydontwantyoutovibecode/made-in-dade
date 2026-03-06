package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	execx "github.com/theydontwantyoutovibecode/dade/internal/exec"
	"github.com/theydontwantyoutovibecode/dade/internal/logging"
	"github.com/theydontwantyoutovibecode/dade/internal/manifest"
	"github.com/theydontwantyoutovibecode/dade/internal/ui"
)

type installRunner struct {
	err error
}

func (r installRunner) Run(_ context.Context, _ string, _ ...string) error {
	return r.err
}

func (r installRunner) Output(_ context.Context, _ string, _ ...string) (string, error) {
	return "", r.err
}

func (r installRunner) LookPath(_ string) (string, error) {
	if r.err != nil {
		return "", r.err
	}
	return "/usr/bin/git", nil
}

func TestInstallCommandMissingManifest(t *testing.T) {
	logger := logging.New(&strings.Builder{}, &strings.Builder{}, false)
	cmd := defaultInstallCommand()
	cmd.runner = installRunner{}
	cmd.templatesDir = func() (string, error) { return t.TempDir(), nil }
	cmd.tempDir = func(_, _ string) (string, error) { return t.TempDir(), nil }
	cmd.readFile = func(string) ([]byte, error) { return nil, os.ErrNotExist }
	cmd.rename = func(string, string) error { return nil }
	cmd.removeAll = func(string) error { return nil }

	code := cmd.run(context.Background(), []string{"https://example.com/repo.git"}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
}

func TestInstallCommandInvalidManifest(t *testing.T) {
	logger := logging.New(&strings.Builder{}, &strings.Builder{}, false)
	cmd := defaultInstallCommand()
	cmd.runner = installRunner{}
	cmd.templatesDir = func() (string, error) { return t.TempDir(), nil }
	cmd.tempDir = func(_, _ string) (string, error) { return t.TempDir(), nil }
	cmd.readFile = func(string) ([]byte, error) { return []byte("[template]\nname=\"bad\"\n"), nil }
	cmd.rename = func(string, string) error { return nil }
	cmd.removeAll = func(string) error { return nil }

	code := cmd.run(context.Background(), []string{"https://example.com/repo.git"}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
}

func TestInstallCommandWritesSource(t *testing.T) {
	logger := logging.New(&strings.Builder{}, &strings.Builder{}, false)
	cmd := defaultInstallCommand()
	cmd.runner = installRunner{}
	cmd.templatesDir = func() (string, error) { return t.TempDir(), nil }
	cmd.tempDir = func(_, _ string) (string, error) { return t.TempDir(), nil }
	cmd.readFile = func(string) ([]byte, error) {
		data := strings.Join([]string{
			"[template]",
			"name = \"hypertext\"",
			"description = \"desc\"",
			"",
			"[serve]",
			"type = \"static\"",
			"",
		}, "\n")
		return []byte(data), nil
	}
	cmd.rename = func(oldPath, newPath string) error {
		return os.MkdirAll(newPath, 0755)
	}
	cmd.removeAll = func(string) error { return nil }
	cmd.writeFile = os.WriteFile

	code := cmd.run(context.Background(), []string{"https://example.com/repo.git"}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
}

func TestInstallCommandListOfficial(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	installCommandFactory = func() installCommand { return defaultInstallCommand() }
	defer func() { installCommandFactory = defaultInstallCommand }()

	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"template", "add", "--list-official"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
	if !strings.Contains(stdout.String(), "Available Templates") {
		t.Fatalf("expected official templates output")
	}
}

func TestInstallCmdMissingURL(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	installCommandFactory = func() installCommand {
		cmd := defaultInstallCommand()
		cmd.runner = installRunner{}
		cmd.templatesDir = func() (string, error) { return t.TempDir(), nil }
		cmd.tempDir = func(_, _ string) (string, error) { return t.TempDir(), nil }
		cmd.readFile = func(string) ([]byte, error) { return nil, os.ErrNotExist }
		cmd.rename = func(string, string) error { return nil }
		cmd.removeAll = func(string) error { return nil }
		return cmd
	}
	defer func() { installCommandFactory = defaultInstallCommand }()

	_ = templateAddCmd.Flags().Set("list-official", "false")
	_ = templateAddCmd.Flags().Set("name", "")
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"template", "add"})

	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected error when url missing")
	}
}

func TestInstallCommandRequiresGit(t *testing.T) {
	logger := logging.New(&strings.Builder{}, &strings.Builder{}, false)
	cmd := defaultInstallCommand()
	cmd.runner = installRunner{err: errors.New("missing")}
	cmd.templatesDir = func() (string, error) { return t.TempDir(), nil }
	cmd.tempDir = func(_, _ string) (string, error) { return t.TempDir(), nil }
	cmd.readFile = func(string) ([]byte, error) { return nil, os.ErrNotExist }
	cmd.rename = func(string, string) error { return nil }
	cmd.removeAll = func(string) error { return nil }

	code := cmd.run(context.Background(), []string{"https://example.com/repo.git"}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
}

func TestInstallCommandUsesTemplateNameFallback(t *testing.T) {
	logger := logging.New(&strings.Builder{}, &strings.Builder{}, false)
	cmd := defaultInstallCommand()
	cmd.runner = installRunner{}
	cmd.templatesDir = func() (string, error) { return t.TempDir(), nil }
	cmd.tempDir = func(_, _ string) (string, error) { return t.TempDir(), nil }
	cmd.readFile = func(string) ([]byte, error) {
		data := strings.Join([]string{
			"[template]",
			"name = \"fallback\"",
			"description = \"desc\"",
			"",
			"[serve]",
			"type = \"static\"",
			"",
		}, "\n")
		return []byte(data), nil
	}
	cmd.rename = func(oldPath, newPath string) error {
		return os.MkdirAll(newPath, 0755)
	}
	cmd.removeAll = func(string) error { return nil }
	cmd.writeFile = os.WriteFile

	code := cmd.run(context.Background(), []string{"https://example.com/repo.git"}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
}

func TestInstallCommandHonorsNameFlag(t *testing.T) {
	logger := logging.New(&strings.Builder{}, &strings.Builder{}, false)
	cmd := defaultInstallCommand()
	cmd.runner = installRunner{}
	cmd.templatesDir = func() (string, error) { return t.TempDir(), nil }
	cmd.tempDir = func(_, _ string) (string, error) { return t.TempDir(), nil }
	cmd.readFile = func(string) ([]byte, error) { return []byte("[template]\nname=\"orig\"\ndescription=\"desc\"\n[serve]\ntype=\"static\"\n"), nil }
	cmd.rename = func(oldPath, newPath string) error {
		if !strings.Contains(newPath, "custom") {
			return errors.New("unexpected target")
		}
		return os.MkdirAll(newPath, 0755)
	}
	cmd.removeAll = func(string) error { return nil }
	cmd.writeFile = os.WriteFile

	code := cmd.run(context.Background(), []string{"https://example.com/repo.git", "--name", "custom"}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
}

func TestInstallCommandRejectsExistingTemplate(t *testing.T) {
	logger := logging.New(&strings.Builder{}, &strings.Builder{}, false)
	cmd := defaultInstallCommand()
	cmd.runner = installRunner{}
	root := t.TempDir()
	cmd.templatesDir = func() (string, error) { return root, nil }
	cmd.tempDir = func(_, _ string) (string, error) { return t.TempDir(), nil }
	cmd.readFile = func(string) ([]byte, error) { return []byte("[template]\nname=\"dup\"\ndescription=\"desc\"\n[serve]\ntype=\"static\"\n"), nil }
	cmd.rename = func(oldPath, newPath string) error { return nil }
	cmd.removeAll = func(string) error { return nil }
	cmd.writeFile = os.WriteFile

	if err := os.MkdirAll(filepath.Join(root, "dup"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	code := cmd.run(context.Background(), []string{"https://example.com/repo.git"}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
}

func TestInstallCommandUsesManifestValidation(t *testing.T) {
	logger := logging.New(&strings.Builder{}, &strings.Builder{}, false)
	cmd := defaultInstallCommand()
	cmd.runner = installRunner{}
	cmd.templatesDir = func() (string, error) { return t.TempDir(), nil }
	cmd.tempDir = func(_, _ string) (string, error) { return t.TempDir(), nil }
	cmd.readFile = func(string) ([]byte, error) {
		return []byte("[template]\nname=\"valid\"\ndescription=\"desc\"\n[serve]\ntype=\"static\"\n"), nil
	}
	cmd.rename = func(oldPath, newPath string) error {
		return os.MkdirAll(newPath, 0755)
	}
	cmd.removeAll = func(string) error { return nil }
	cmd.writeFile = os.WriteFile

	code := cmd.run(context.Background(), []string{"https://example.com/repo.git"}, ui.New(&strings.Builder{}, &strings.Builder{}, false), logger, false)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
}

var _ = execx.Runner(installRunner{})
var _ = manifest.Manifest{}
