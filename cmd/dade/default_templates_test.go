package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
)

func TestEnsureDefaultTemplatesCreatesMarker(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	templatesDir := filepath.Join(baseDir, "dade", "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	tplDir := filepath.Join(templatesDir, "web-app")
	if err := os.MkdirAll(tplDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tplDir, ".default"), []byte(""), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	if !isDefaultTemplate(tplDir) {
		t.Fatalf("expected isDefaultTemplate to return true")
	}

	noMarkerDir := filepath.Join(templatesDir, "custom")
	if err := os.MkdirAll(noMarkerDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if isDefaultTemplate(noMarkerDir) {
		t.Fatalf("expected isDefaultTemplate to return false for user template")
	}
}

func TestEnsureDefaultTemplatesSkipsExisting(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	templatesDir := filepath.Join(baseDir, "dade", "templates")

	var clonedURLs []string
	cmd := defaultTemplatesCommand{
		templatesDir: func() (string, error) { return templatesDir, nil },
		writeFile:    os.WriteFile,
		runner: &fakeRunner{
			err: nil,
		},
	}

	origRun := cmd.runner
	cloneRunner := &cloneFakeRunner{inner: origRun.(*fakeRunner), clonedURLs: &clonedURLs, templatesDir: templatesDir}
	cmd.runner = cloneRunner

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)

	if err := cmd.run(context.Background(), logger); err != nil {
		t.Fatalf("first run: %v", err)
	}
	firstCount := len(clonedURLs)
	if firstCount == 0 {
		t.Fatalf("expected at least one clone")
	}

	clonedURLs = nil
	cloneRunner.clonedURLs = &clonedURLs
	if err := cmd.run(context.Background(), logger); err != nil {
		t.Fatalf("second run: %v", err)
	}
	if len(clonedURLs) != 0 {
		t.Fatalf("expected no clones on second run, got %d", len(clonedURLs))
	}
}

type cloneFakeRunner struct {
	inner        *fakeRunner
	clonedURLs   *[]string
	templatesDir string
}

func (c *cloneFakeRunner) Run(ctx context.Context, name string, args ...string) error {
	if name == "git" && len(args) > 0 && args[0] == "clone" {
		url := args[len(args)-2]
		target := args[len(args)-1]
		*c.clonedURLs = append(*c.clonedURLs, url)
		if err := os.MkdirAll(target, 0755); err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(target, "dade.toml"), []byte("[template]\nname = \"test\"\n"), 0644)
	}
	return c.inner.Run(ctx, name, args...)
}

func (c *cloneFakeRunner) Output(ctx context.Context, name string, args ...string) (string, error) {
	return c.inner.Output(ctx, name, args...)
}

func (c *cloneFakeRunner) LookPath(name string) (string, error) {
	return "/usr/bin/" + name, nil
}
