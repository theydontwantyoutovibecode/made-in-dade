package proxy

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
)

func TestBuildPlistIncludesPaths(t *testing.T) {
	content := buildPlist("/opt/homebrew/bin/caddy", "/tmp/Caddyfile", "/tmp/proxy.log", "/tmp/proxy.err")
	if !strings.Contains(content, config.ProxyLabel) {
		t.Fatalf("expected proxy label")
	}
	if !strings.Contains(content, "/opt/homebrew/bin/caddy") {
		t.Fatalf("expected caddy path")
	}
	if !strings.Contains(content, "/tmp/Caddyfile") {
		t.Fatalf("expected caddyfile path")
	}
	if !strings.Contains(content, "/tmp/proxy.log") {
		t.Fatalf("expected log path")
	}
	if !strings.Contains(content, "/tmp/proxy.err") {
		t.Fatalf("expected err path")
	}
}

func TestCreatePlistWritesFile(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	plistPath := filepath.Join(root, "proxy.plist")
	caddyfilePath := filepath.Join(root, "Caddyfile")
	logPath := filepath.Join(root, "proxy.log")
	errPath := filepath.Join(root, "proxy.err")

	runner := &fakeRunner{lookPath: func(name string) (string, error) {
		if name != "caddy" {
			return "", errors.New("unexpected binary")
		}
		return "/opt/homebrew/bin/caddy", nil
	}}

	if err := CreatePlist(ctx, runner, plistPath, caddyfilePath, logPath, errPath); err != nil {
		t.Fatalf("create plist: %v", err)
	}

	data, err := os.ReadFile(plistPath)
	if err != nil {
		t.Fatalf("read plist: %v", err)
	}
	if !strings.Contains(string(data), "/opt/homebrew/bin/caddy") {
		t.Fatalf("expected caddy path in plist")
	}
}

func TestInstallProxyServiceBootstrapsWhenStopped(t *testing.T) {
	ctx := context.Background()
	plistPath := "/tmp/dade.plist"
	uid := strconv.Itoa(os.Getuid())

	runner := &fakeRunner{run: func(name string, args ...string) error {
		if name == "launchctl" && len(args) > 0 && args[0] == "list" {
			return errors.New("not running")
		}
		return nil
	}}

	if err := InstallProxyService(ctx, runner, plistPath); err != nil {
		t.Fatalf("install: %v", err)
	}

	foundBootstrap := false
	for _, call := range runner.calls {
		if call.name == "launchctl" && len(call.args) >= 3 && call.args[0] == "bootstrap" {
			if call.args[1] != "gui/"+uid || call.args[2] != plistPath {
				t.Fatalf("unexpected bootstrap args")
			}
			foundBootstrap = true
		}
	}
	if !foundBootstrap {
		t.Fatalf("expected bootstrap call")
	}
}

func TestInstallProxyServiceSkipsWhenRunning(t *testing.T) {
	ctx := context.Background()
	runner := &fakeRunner{run: func(name string, args ...string) error {
		return nil
	}}

	if err := InstallProxyService(ctx, runner, "/tmp/dade.plist"); err != nil {
		t.Fatalf("install: %v", err)
	}
	for _, call := range runner.calls {
		if call.name == "launchctl" && len(call.args) > 0 && call.args[0] == "bootstrap" {
			t.Fatalf("unexpected bootstrap")
		}
	}
}

func TestUninstallProxyServiceBootoutWhenRunning(t *testing.T) {
	ctx := context.Background()
	plistPath := "/tmp/dade.plist"
	uid := strconv.Itoa(os.Getuid())

	runner := &fakeRunner{run: func(name string, args ...string) error {
		return nil
	}}

	if err := UninstallProxyService(ctx, runner, plistPath); err != nil {
		t.Fatalf("uninstall: %v", err)
	}

	foundBootout := false
	for _, call := range runner.calls {
		if call.name == "launchctl" && len(call.args) >= 3 && call.args[0] == "bootout" {
			if call.args[1] != "gui/"+uid || call.args[2] != plistPath {
				t.Fatalf("unexpected bootout args")
			}
			foundBootout = true
		}
	}
	if !foundBootout {
		t.Fatalf("expected bootout call")
	}
}

func TestUninstallProxyServiceSkipsWhenStopped(t *testing.T) {
	ctx := context.Background()
	runner := &fakeRunner{run: func(name string, args ...string) error {
		if name == "launchctl" && len(args) > 0 && args[0] == "list" {
			return errors.New("not running")
		}
		return nil
	}}

	if err := UninstallProxyService(ctx, runner, "/tmp/dade.plist"); err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	for _, call := range runner.calls {
		if call.name == "launchctl" && len(call.args) > 0 && call.args[0] == "bootout" {
			t.Fatalf("unexpected bootout")
		}
	}
}
