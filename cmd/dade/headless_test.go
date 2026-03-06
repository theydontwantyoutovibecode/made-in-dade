package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/theydontwantyoutovibecode/dade/internal/registry"
)

func TestHeadlessNewRequiresArgs(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"new"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatalf("new without args should fail in headless mode")
	}
	if !strings.Contains(stderr.String(), "name") || !strings.Contains(stderr.String(), "required") {
		if !strings.Contains(stderr.String(), "template") && !strings.Contains(stderr.String(), "required") {
			t.Logf("stderr: %s", stderr.String())
		}
	}
}

func TestHeadlessTemplatesWithJSON(t *testing.T) {
	templatesDir := t.TempDir()
	templateDir := filepath.Join(templatesDir, "mytemplate")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	manifest := `[template]
name = "mytemplate"
description = "Test template"

[serve]
type = "static"
`
	if err := os.WriteFile(filepath.Join(templateDir, "dade.toml"), []byte(manifest), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	configDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configDir)
	if err := os.Symlink(templateDir, filepath.Join(configDir, "dade", "templates", "mytemplate")); err != nil {
		if err := os.MkdirAll(filepath.Join(configDir, "dade", "templates"), 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		src, _ := os.ReadFile(filepath.Join(templateDir, "dade.toml"))
		destDir := filepath.Join(configDir, "dade", "templates", "mytemplate")
		if err := os.MkdirAll(destDir, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(destDir, "dade.toml"), src, 0644); err != nil {
			t.Fatalf("write: %v", err)
		}
	}

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"template", "list", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("templates --json failed: %v", err)
	}
	if !strings.Contains(stdout.String(), "mytemplate") {
		t.Fatalf("expected template in output, got: %s", stdout.String())
	}
}

func TestHeadlessInstallRequiresURL(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"template", "add"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("install without args should fail")
	}
	if !strings.Contains(stderr.String(), "URL") && !strings.Contains(stderr.String(), "url") {
		t.Logf("expected URL error, got: %s", stderr.String())
	}
}

func TestHeadlessInstallListOfficial(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"template", "add", "--list-official"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("install --list-official failed: %v", err)
	}
}

func TestHeadlessListWithJSON(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	projectsPath := filepath.Join(baseDir, "dade", "projects.json")
	if err := os.MkdirAll(filepath.Dir(projectsPath), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	projects := map[string]registry.Project{
		"testproj": {Port: 59999, Path: "/tmp/testproj", Template: "hypertext"},
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
		t.Fatalf("list --json failed: %v", err)
	}
	if !strings.Contains(stdout.String(), "testproj") {
		t.Fatalf("expected project in output, got: %s", stdout.String())
	}
}

func TestHeadlessProxyStatus(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"proxy", "status"})
	_ = rootCmd.Execute()
}

func TestHeadlessSetupCheck(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"setup", "--check"})
	_ = rootCmd.Execute()
}

func TestHeadlessUpdateRequiresName(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetUpdateFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"template", "update"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("update without args should fail")
	}
	if !strings.Contains(stderr.String(), "Missing template name") {
		t.Logf("expected missing name error, got: %s", stderr.String())
	}
}

func TestHeadlessRemoveRequiresName(t *testing.T) {
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
		t.Fatalf("remove without args should fail")
	}
}

func TestHeadlessQuietFlag(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"--quiet", "template", "list"})
	_ = rootCmd.Execute()
}

func TestHeadlessNoColorFlag(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"--no-color", "template", "list"})
	_ = rootCmd.Execute()
}
