package proxy

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/registry"
)

type runCall struct {
	name string
	args []string
}

type fakeRunner struct {
	calls    []runCall
	run      func(name string, args ...string) error
	lookPath func(name string) (string, error)
}

func (f *fakeRunner) Run(_ context.Context, name string, args ...string) error {
	f.calls = append(f.calls, runCall{name: name, args: append([]string{}, args...)})
	if f.run != nil {
		return f.run(name, args...)
	}
	return nil
}

func (f *fakeRunner) Output(_ context.Context, _ string, _ ...string) (string, error) {
	return "", nil
}

func (f *fakeRunner) LookPath(name string) (string, error) {
	if f.lookPath != nil {
		return f.lookPath(name)
	}
	return "/usr/bin/true", nil
}

func TestBuildCaddyfile(t *testing.T) {
	entries := []registry.Entry{
		{Name: "alpha", Project: registry.Project{Port: 3000}},
		{Name: "", Project: registry.Project{Port: 3001}},
		{Name: "beta", Project: registry.Project{Port: 0}},
	}

	content := buildCaddyfile(entries)
	localDomain := config.LocalDomain()
	if !strings.Contains(content, "local_certs") {
		t.Fatalf("expected local_certs")
	}
	if !strings.Contains(content, "https://alpha."+localDomain) {
		t.Fatalf("expected alpha block, got: %s", content)
	}
	if strings.Contains(content, "beta."+localDomain) {
		t.Fatalf("unexpected beta block")
	}
}

func TestGenerateCaddyfileWritesAndBacksUp(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	projectsPath := filepath.Join(root, "projects.json")
	caddyfilePath := filepath.Join(root, "Caddyfile")

	projects := map[string]registry.Project{
		"alpha": {Port: 3000, Path: "/tmp/alpha", Template: "hypertext", Created: "now"},
	}
	if err := registry.Save(projectsPath, projects); err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := os.WriteFile(caddyfilePath, []byte("old"), 0644); err != nil {
		t.Fatalf("write old: %v", err)
	}

	runner := &fakeRunner{run: func(name string, args ...string) error {
		if name == "caddy" && len(args) > 0 && args[0] == "validate" {
			return nil
		}
		return nil
	}}

	if err := GenerateCaddyfile(ctx, runner, projectsPath, caddyfilePath); err != nil {
		t.Fatalf("generate: %v", err)
	}

	data, err := os.ReadFile(caddyfilePath)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	localDomain := config.LocalDomain()
	if !strings.Contains(string(data), "alpha."+localDomain) {
		t.Fatalf("expected alpha entry, got: %s", string(data))
	}

	backup, err := os.ReadFile(caddyfilePath + ".bak")
	if err != nil {
		t.Fatalf("backup: %v", err)
	}
	if string(backup) != "old" {
		t.Fatalf("unexpected backup contents")
	}
}

func TestGenerateCaddyfileValidationFailureKeepsExisting(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	projectsPath := filepath.Join(root, "projects.json")
	caddyfilePath := filepath.Join(root, "Caddyfile")

	projects := map[string]registry.Project{
		"alpha": {Port: 3000, Path: "/tmp/alpha", Template: "hypertext", Created: "now"},
	}
	if err := registry.Save(projectsPath, projects); err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := os.WriteFile(caddyfilePath, []byte("old"), 0644); err != nil {
		t.Fatalf("write old: %v", err)
	}

	runner := &fakeRunner{run: func(name string, args ...string) error {
		if name == "caddy" && len(args) > 0 && args[0] == "validate" {
			return errors.New("invalid")
		}
		return nil
	}}

	if err := GenerateCaddyfile(ctx, runner, projectsPath, caddyfilePath); err == nil {
		t.Fatalf("expected error")
	}

	data, err := os.ReadFile(caddyfilePath)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(data) != "old" {
		t.Fatalf("expected existing caddyfile to remain")
	}
}

func TestReloadProxySkipsWhenNotRunning(t *testing.T) {
	ctx := context.Background()
	runner := &fakeRunner{run: func(name string, args ...string) error {
		if name == "launchctl" {
			return errors.New("not running")
		}
		return nil
	}}

	if err := ReloadProxy(ctx, runner, "/tmp/Caddyfile"); err != nil {
		t.Fatalf("reload: %v", err)
	}
	for _, call := range runner.calls {
		if call.name == "caddy" && len(call.args) > 0 && call.args[0] == "reload" {
			t.Fatalf("unexpected reload")
		}
	}
}

func TestReloadProxyRunsWhenRunning(t *testing.T) {
	ctx := context.Background()
	runner := &fakeRunner{run: func(name string, args ...string) error {
		return nil
	}}

	if err := ReloadProxy(ctx, runner, "/tmp/Caddyfile"); err != nil {
		t.Fatalf("reload: %v", err)
	}

	foundLaunchctl := false
	foundReload := false
	for _, call := range runner.calls {
		if call.name == "launchctl" {
			foundLaunchctl = true
		}
		if call.name == "caddy" && len(call.args) > 0 && call.args[0] == "reload" {
			foundReload = true
		}
	}
	if !foundLaunchctl || !foundReload {
		t.Fatalf("expected launchctl and reload calls")
	}
}
