package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitConfigCreatesStructure(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	if _, err := InitConfig(); err != nil {
		t.Fatalf("init config: %v", err)
	}

	configDir, err := ConfigDir()
	if err != nil {
		t.Fatalf("config dir: %v", err)
	}
	if info, err := os.Stat(configDir); err != nil || !info.IsDir() {
		t.Fatalf("expected config dir to exist")
	}

	templatesDir := filepath.Join(configDir, TemplatesDirName)
	if info, err := os.Stat(templatesDir); err != nil || !info.IsDir() {
		t.Fatalf("expected templates dir to exist")
	}

	projectsFile := filepath.Join(configDir, ProjectsFileName)
	data, err := os.ReadFile(projectsFile)
	if err != nil {
		t.Fatalf("read projects file: %v", err)
	}
	if string(data) != "{}" {
		t.Fatalf("unexpected projects file contents: %s", string(data))
	}

	caddyfile := filepath.Join(configDir, CaddyfileName)
	caddyData, err := os.ReadFile(caddyfile)
	if err != nil {
		t.Fatalf("read caddyfile: %v", err)
	}
	if string(caddyData) != defaultCaddyfile {
		t.Fatalf("unexpected caddyfile contents: %s", string(caddyData))
	}
}

func TestInitConfigIsIdempotent(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	if _, err := InitConfig(); err != nil {
		t.Fatalf("init config: %v", err)
	}
	if _, err := InitConfig(); err != nil {
		t.Fatalf("init config again: %v", err)
	}
}
