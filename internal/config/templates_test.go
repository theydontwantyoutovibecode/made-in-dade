package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultTemplates(t *testing.T) {
	got := DefaultTemplates()
	if len(got.Ordered) != 7 {
		t.Fatalf("expected 7 defaults, got %d", len(got.Ordered))
	}
	if got.Ordered[0].Name != "web-app" || got.Ordered[1].Name != "web-site" {
		t.Fatalf("unexpected ordering: %#v", got.Ordered)
	}
	if got.ByName["web-app"].URL == "" {
		t.Fatalf("expected web-app url")
	}
}

func TestLoadTemplatesMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.toml")
	got, err := LoadTemplates(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Ordered) != 7 {
		t.Fatalf("expected defaults when missing file")
	}
}

func TestPathsRespectXDGConfigHome(t *testing.T) {
	original := os.Getenv("XDG_CONFIG_HOME")
	t.Cleanup(func() {
		if original == "" {
			_ = os.Unsetenv("XDG_CONFIG_HOME")
			return
		}
		_ = os.Setenv("XDG_CONFIG_HOME", original)
	})
	tempDir := t.TempDir()
	if err := os.Setenv("XDG_CONFIG_HOME", tempDir); err != nil {
		t.Fatalf("set env: %v", err)
	}

	configDir, err := ConfigDir()
	if err != nil {
		t.Fatalf("config dir: %v", err)
	}
	if configDir != filepath.Join(tempDir, ConfigDirName) {
		t.Fatalf("unexpected config dir: %s", configDir)
	}

	projectsFile, err := ProjectsFile()
	if err != nil {
		t.Fatalf("projects file: %v", err)
	}
	if projectsFile != filepath.Join(configDir, ProjectsFileName) {
		t.Fatalf("unexpected projects file: %s", projectsFile)
	}

	templatesDir, err := TemplatesDir()
	if err != nil {
		t.Fatalf("templates dir: %v", err)
	}
	if templatesDir != filepath.Join(configDir, TemplatesDirName) {
		t.Fatalf("unexpected templates dir: %s", templatesDir)
	}

	caddyfilePath, err := CaddyfilePath()
	if err != nil {
		t.Fatalf("caddyfile path: %v", err)
	}
	if caddyfilePath != filepath.Join(configDir, CaddyfileName) {
		t.Fatalf("unexpected caddyfile path: %s", caddyfilePath)
	}

	configPath, err := ConfigFilePath()
	if err != nil {
		t.Fatalf("config file path: %v", err)
	}
	if configPath != filepath.Join(configDir, ConfigFileName) {
		t.Fatalf("unexpected config file path: %s", configPath)
	}

	proxyLogPath, err := ProxyLogPath()
	if err != nil {
		t.Fatalf("proxy log path: %v", err)
	}
	if proxyLogPath != filepath.Join(configDir, ProxyLogName) {
		t.Fatalf("unexpected proxy log path: %s", proxyLogPath)
	}

	proxyErrPath, err := ProxyErrPath()
	if err != nil {
		t.Fatalf("proxy err path: %v", err)
	}
	if proxyErrPath != filepath.Join(configDir, ProxyErrName) {
		t.Fatalf("unexpected proxy err path: %s", proxyErrPath)
	}

	plistPath, err := ProxyPlistPath()
	if err != nil {
		t.Fatalf("proxy plist path: %v", err)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("home dir: %v", err)
	}
	expectedPlist := filepath.Join(home, "Library", "LaunchAgents", ProxyLabel+".plist")
	if plistPath != expectedPlist {
		t.Fatalf("unexpected plist path: %s", plistPath)
	}
}

func TestLoadTemplatesOverridesAndAppend(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "templates.toml")
	content := strings.Join([]string{
		"[templates]",
		"web-app = \"https://override.example/repo.git\"",
		"custom-template = \"https://example.com/custom.git\"",
		"",
	}, "\n")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write: %v", err)
	}

	got, err := LoadTemplates(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ByName["web-app"].URL != "https://override.example/repo.git" {
		t.Fatalf("override not applied")
	}
	if got.ByName["custom-template"].URL != "https://example.com/custom.git" {
		t.Fatalf("custom template missing")
	}
	if got.Ordered[len(got.Ordered)-1].Name != "custom-template" {
		t.Fatalf("custom template should be appended last")
	}
}

func TestLoadTemplatesSkipsInvalidLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "templates.toml")
	content := strings.Join([]string{
		"# comment",
		"",
		"not a toml line",
		"key-only =",
		"= value-only",
		"[templates]",
		"valid = \"https://example.com/valid.git\"",
		"",
	}, "\n")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write: %v", err)
	}

	got, err := LoadTemplates(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ByName["valid"].URL != "https://example.com/valid.git" {
		t.Fatalf("expected valid template")
	}
}
