package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/theydontwantyoutovibecode/dade/internal/config"
)

func TestTemplatesTextPlainIncludesGuidance(t *testing.T) {
	templates := config.DefaultTemplates()
	output := availableTemplatesText(templates, config.DefaultTemplatesPath, false)
	if !strings.Contains(output, "Available Templates") {
		t.Fatalf("missing templates header")
	}
	if !strings.Contains(output, "To add custom templates") {
		t.Fatalf("missing guidance snippet")
	}
	if !strings.Contains(output, "web-app") {
		t.Fatalf("missing default template")
	}
}

func TestTemplatesJSONOutput(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)
	writeInstalledTemplate(t, baseDir, "hypertext", "Vanilla HTML/CSS/JS + HTMX", "static", "https://example.com/hypertext.git")

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"template", "list", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
	var payload []installedTemplate
	if err := json.Unmarshal([]byte(stdout.String()), &payload); err != nil {
		t.Fatalf("invalid json output: %v", err)
	}
	if len(payload) != 1 {
		t.Fatalf("expected templates in json output")
	}
	if payload[0].Name != "hypertext" {
		t.Fatalf("expected template name in json output")
	}
	if payload[0].ServeType != "static" {
		t.Fatalf("expected serve type in json output")
	}
}

func TestTemplatesEmptyState(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", baseDir)

	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"template", "list"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
	output := stdout.String()
	if !strings.Contains(output, "No templates installed.") {
		t.Fatalf("missing empty state message")
	}
	if !strings.Contains(output, "dade template add <git-url>") {
		t.Fatalf("missing install guidance")
	}
	if !strings.Contains(output, "dade template add --list-official") {
		t.Fatalf("missing official guidance")
	}
}

func writeInstalledTemplate(t *testing.T, baseDir, name, description, serveType, source string) {
	t.Helper()
	templateDir := filepath.Join(baseDir, "dade", "templates", name)
	if err := os.MkdirAll(templateDir, 0o755); err != nil {
		t.Fatalf("failed to create template dir: %v", err)
	}
	manifest := strings.Join([]string{
		"[template]",
		"name = \"" + name + "\"",
		"description = \"" + description + "\"",
		"",
		"[serve]",
		"type = \"" + serveType + "\"",
		"",
	}, "\n")
	if err := os.WriteFile(filepath.Join(templateDir, "dade.toml"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("failed to write manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templateDir, ".source"), []byte(source), 0o644); err != nil {
		t.Fatalf("failed to write source: %v", err)
	}
}
