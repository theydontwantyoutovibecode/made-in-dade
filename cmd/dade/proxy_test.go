package main

import (
	"context"
	"errors"
	"strings"
	"testing"

	execx "github.com/theydontwantyoutovibecode/made-in-dade/internal/exec"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/registry"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/ui"
)

type proxyStub struct {
	createdPlist bool
	installed    bool
	uninstalled  bool
	restarted    bool
	tailed       bool
}

type stubRunner struct{}

func (stubRunner) Run(_ context.Context, _ string, _ ...string) error { return nil }
func (stubRunner) Output(_ context.Context, _ string, _ ...string) (string, error) {
	return "", nil
}
func (stubRunner) LookPath(_ string) (string, error) { return "/usr/bin/true", nil }

var _ execx.Runner = stubRunner{}

func TestProxyStatusRunning(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	console := ui.New(stdout, stderr, false)
	cmd := defaultProxyCommand()
	cmd.initConfig = func() (bool, error) { return false, nil }
	cmd.isRunning = func(context.Context, execx.Runner) bool { return true }
	cmd.loadRegistry = func(string) (map[string]registry.Project, error) {
		return map[string]registry.Project{
			"alpha": {Port: 3000},
			"beta":  {Port: 3001},
		}, nil
	}
	cmd.projectsFile = func() (string, error) { return "/tmp/projects.json", nil }
	cmd.caddyfilePath = func() (string, error) { return "/tmp/Caddyfile", nil }
	cmd.runner = nil

	code := cmd.run(context.Background(), []string{"status"}, console, logger)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	output := stdout.String()
	if !strings.Contains(output, "Proxy running") {
		t.Fatalf("expected running status")
	}
	if !strings.Contains(output, "Projects: 2") {
		t.Fatalf("expected project count")
	}
	if !strings.Contains(output, "Ports: 3000-3001") {
		t.Fatalf("expected ports")
	}
	if !strings.Contains(output, "Caddyfile: /tmp/Caddyfile") {
		t.Fatalf("expected caddyfile path")
	}
}

func TestProxyStatusJSON(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	proxyCommandFactory = func() proxyCommand {
		cmd := defaultProxyCommand()
		cmd.initConfig = func() (bool, error) { return false, nil }
		cmd.isRunning = func(context.Context, execx.Runner) bool { return true }
		cmd.projectsFile = func() (string, error) { return "projects.json", nil }
		cmd.loadRegistry = func(string) (map[string]registry.Project, error) {
			return map[string]registry.Project{"alpha": {Port: 3000}}, nil
		}
		cmd.caddyfilePath = func() (string, error) { return "Caddyfile", nil }
		return cmd
	}
	defer func() { proxyCommandFactory = defaultProxyCommand }()

	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"proxy", "status", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}

	if !strings.Contains(stdout.String(), "\"running\": true") {
		t.Fatalf("expected json output")
	}
}

func TestProxyLogsNoFollow(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	proxyCommandFactory = func() proxyCommand {
		cmd := defaultProxyCommand()
		cmd.logPath = func() (string, error) { return "test.log", nil }
		cmd.tail = nil
		return cmd
	}
	defer func() { proxyCommandFactory = defaultProxyCommand }()

	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"proxy", "logs", "--follow=false"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected error when tail unavailable")
	}
}

func TestProxyStartCreatesPlistAndInstalls(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	console := ui.New(stdout, stderr, false)
	stub := &proxyStub{}
	cmd := defaultProxyCommand()
	cmd.initConfig = func() (bool, error) { return false, nil }
	cmd.isRunning = func(context.Context, execx.Runner) bool { return false }
	cmd.createPlist = func(_ context.Context, _ execx.Runner, _, _, _, _ string) error {
		stub.createdPlist = true
		return nil
	}
	cmd.install = func(_ context.Context, _ execx.Runner, _ string) error {
		stub.installed = true
		return nil
	}
	cmd.plistPath = func() (string, error) { return "/tmp/proxy.plist", nil }
	cmd.caddyfilePath = func() (string, error) { return "/tmp/Caddyfile", nil }
	cmd.logPath = func() (string, error) { return "/tmp/proxy.log", nil }
	cmd.errPath = func() (string, error) { return "/tmp/proxy.err", nil }

	code := cmd.run(context.Background(), []string{"start"}, console, logger)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !stub.createdPlist || !stub.installed {
		t.Fatalf("expected plist creation and install")
	}
}

func TestProxyStopSkipsWhenNotRunning(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	console := ui.New(stdout, stderr, false)
	cmd := defaultProxyCommand()
	cmd.initConfig = func() (bool, error) { return false, nil }
	cmd.isRunning = func(context.Context, execx.Runner) bool { return false }

	code := cmd.run(context.Background(), []string{"stop"}, console, logger)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Proxy not running") {
		t.Fatalf("expected not running warning")
	}
}

func TestProxyRestartCallsRestart(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	console := ui.New(stdout, stderr, false)
	stub := &proxyStub{}
	cmd := defaultProxyCommand()
	cmd.initConfig = func() (bool, error) { return false, nil }
	cmd.restart = func(_ context.Context, _ execx.Runner, _ string) error {
		stub.restarted = true
		return nil
	}
	cmd.plistPath = func() (string, error) { return "/tmp/proxy.plist", nil }

	code := cmd.run(context.Background(), []string{"restart"}, console, logger)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !stub.restarted {
		t.Fatalf("expected restart call")
	}
}

func TestProxyLogsUsesTail(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	console := ui.New(stdout, stderr, false)
	stub := &proxyStub{}
	cmd := defaultProxyCommand()
	cmd.initConfig = func() (bool, error) { return false, nil }
	cmd.logPath = func() (string, error) { return "/tmp/proxy.log", nil }
	cmd.tail = func(path string) error {
		if path != "/tmp/proxy.log" {
			return errors.New("unexpected path")
		}
		stub.tailed = true
		return nil
	}

	code := cmd.run(context.Background(), []string{"logs"}, console, logger)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !stub.tailed {
		t.Fatalf("expected tail to be called")
	}
}

func TestProxyUnknownAction(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	console := ui.New(stdout, stderr, false)
	cmd := defaultProxyCommand()
	cmd.initConfig = func() (bool, error) { return false, nil }

	code := cmd.run(context.Background(), []string{"nope"}, console, logger)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "Unknown action") {
		t.Fatalf("expected error output")
	}
}

func TestProxyTooManyArgs(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	console := ui.New(stdout, stderr, false)
	cmd := defaultProxyCommand()
	cmd.initConfig = func() (bool, error) { return false, nil }

	code := cmd.run(context.Background(), []string{"status", "extra"}, console, logger)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "Too many arguments") {
		t.Fatalf("expected too many arguments error")
	}
}
