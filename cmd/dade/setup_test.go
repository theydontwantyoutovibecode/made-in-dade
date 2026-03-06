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

type setupRunner struct {
	calls    []string
	lookPath map[string]error
}

func (s *setupRunner) Run(_ context.Context, name string, args ...string) error {
	s.calls = append(s.calls, name)
	return nil
}

func (s *setupRunner) Output(_ context.Context, _ string, _ ...string) (string, error) {
	return "", nil
}

func (s *setupRunner) LookPath(name string) (string, error) {
	if err, ok := s.lookPath[name]; ok {
		return "", err
	}
	return "/usr/bin/" + name, nil
}

func TestSetupCheckOnlySuccess(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	console := ui.New(stdout, stderr, false)
	runner := &setupRunner{lookPath: map[string]error{}}
	cmd := defaultSetupCommand()
	cmd.runner = runner
	cmd.initConfig = func() (bool, error) { return false, nil }
	cmd.detectSrv = func() (bool, error) { return false, nil }
	cmd.confirm = func(string) (bool, error) { return false, nil }
	cmd.spin = func(string, func() error) error { return nil }

	code := cmd.run(context.Background(), []string{"--check"}, console, logger, false)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "All dependencies OK") {
		t.Fatalf("expected success message")
	}
}

func TestSetupCmdNoTemplatesSkipsInstall(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	setupCommandFactory = func() setupCommand {
		cmd := defaultSetupCommand()
		cmd.runner = &setupRunner{lookPath: map[string]error{}}
		cmd.initConfig = func() (bool, error) { return false, nil }
		cmd.detectSrv = func() (bool, error) { return false, nil }
		cmd.confirm = func(string) (bool, error) { return false, nil }
		cmd.spin = func(string, func() error) error { return nil }
		cmd.projectsFile = func() (string, error) { return "/tmp/projects.json", nil }
		cmd.caddyfilePath = func() (string, error) { return "/tmp/Caddyfile", nil }
		cmd.plistPath = func() (string, error) { return "/tmp/proxy.plist", nil }
		cmd.logPath = func() (string, error) { return "/tmp/proxy.log", nil }
		cmd.errPath = func() (string, error) { return "/tmp/proxy.err", nil }
		cmd.generateCaddy = func(context.Context, execx.Runner, string, string) error { return nil }
		cmd.createPlist = func(context.Context, execx.Runner, string, string, string, string) error { return nil }
		cmd.installProxy = func(context.Context, execx.Runner, string) error { return nil }
		cmd.trustCA = func(context.Context, execx.Runner) error { return nil }
		cmd.migrateSrv = func(context.Context, execx.Runner, string, *logging.Logger) error { return nil }
		cmd.installTemplate = func(context.Context, string, *ui.UI, *logging.Logger, bool) int { return 0 }
		return cmd
	}
	defer func() { setupCommandFactory = defaultSetupCommand }()

	_ = setupCmd.Flags().Set("check", "false")
	_ = setupCmd.Flags().Set("migrate", "false")
	_ = setupCmd.Flags().Set("no-migrate", "false")
	_ = setupCmd.Flags().Set("install-deps", "false")
	_ = setupCmd.Flags().Set("skip-deps", "false")
	_ = setupCmd.Flags().Set("trust-ca", "false")
	_ = setupCmd.Flags().Set("no-trust-ca", "false")
	_ = setupCmd.Flags().Set("install-templates", "false")
	_ = setupCmd.Flags().Set("no-templates", "false")
	_ = setupCmd.Flags().Set("yes", "false")
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"setup", "--no-templates"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
	if !strings.Contains(stdout.String(), "Setup complete") {
		t.Fatalf("expected setup completion")
	}
}

func TestSetupCheckOnlyFailsWhenMissingRequired(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	console := ui.New(stdout, stderr, false)
	runner := &setupRunner{lookPath: map[string]error{"jq": errors.New("missing")}}
	cmd := defaultSetupCommand()
	cmd.runner = runner
	cmd.initConfig = func() (bool, error) { return false, nil }
	cmd.detectSrv = func() (bool, error) { return false, nil }
	cmd.confirm = func(string) (bool, error) { return false, nil }
	cmd.spin = func(string, func() error) error { return nil }

	code := cmd.run(context.Background(), []string{"--check"}, console, logger, false)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "Missing required dependencies") {
		t.Fatalf("expected missing deps error")
	}
}

func TestSetupRunsMigrationWhenConfirmed(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	console := ui.New(stdout, stderr, false)
	runner := &setupRunner{lookPath: map[string]error{}}
	migrated := false
	cmd := defaultSetupCommand()
	cmd.runner = runner
	cmd.initConfig = func() (bool, error) { return false, nil }
	cmd.detectSrv = func() (bool, error) { return true, nil }
	cmd.confirm = func(prompt string) (bool, error) {
		if strings.Contains(prompt, "Migrate") {
			return true, nil
		}
		return false, nil
	}
	cmd.spin = func(string, func() error) error { return nil }
	cmd.projectsFile = func() (string, error) { return "/tmp/projects.json", nil }
	cmd.caddyfilePath = func() (string, error) { return "/tmp/Caddyfile", nil }
	cmd.plistPath = func() (string, error) { return "/tmp/proxy.plist", nil }
	cmd.logPath = func() (string, error) { return "/tmp/proxy.log", nil }
	cmd.errPath = func() (string, error) { return "/tmp/proxy.err", nil }
	cmd.generateCaddy = func(context.Context, execx.Runner, string, string) error { return nil }
	cmd.createPlist = func(context.Context, execx.Runner, string, string, string, string) error { return nil }
	cmd.installProxy = func(context.Context, execx.Runner, string) error { return nil }
	cmd.trustCA = func(context.Context, execx.Runner) error { return nil }
	cmd.migrateSrv = func(context.Context, execx.Runner, string, *logging.Logger) error {
		migrated = true
		return nil
	}
	cmd.installTemplate = func(context.Context, string, *ui.UI, *logging.Logger, bool) int { return 0 }

	code := cmd.run(context.Background(), []string{}, console, logger, false)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !migrated {
		t.Fatalf("expected migration to run")
	}
}

func TestOfferOfficialTemplatesSkipsInstalled(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)
	console := ui.New(stdout, stderr, false)
	cmd := defaultSetupCommand()
	cmd.templatesDir = func() (string, error) { return "/tmp", nil }
	cmd.confirm = func(string) (bool, error) { return false, nil }
	cmd.installTemplate = func(context.Context, string, *ui.UI, *logging.Logger, bool) int { return 0 }

	cmdRunner := &setupRunner{lookPath: map[string]error{}}
	cmd.runner = cmdRunner
	cmd.initConfig = func() (bool, error) { return false, nil }
	cmd.detectSrv = func() (bool, error) { return false, nil }
	cmd.spin = func(string, func() error) error { return nil }
	cmd.projectsFile = func() (string, error) { return "/tmp/projects.json", nil }
	cmd.caddyfilePath = func() (string, error) { return "/tmp/Caddyfile", nil }
	cmd.plistPath = func() (string, error) { return "/tmp/proxy.plist", nil }
	cmd.logPath = func() (string, error) { return "/tmp/proxy.log", nil }
	cmd.errPath = func() (string, error) { return "/tmp/proxy.err", nil }
	cmd.generateCaddy = func(context.Context, execx.Runner, string, string) error { return nil }
	cmd.createPlist = func(context.Context, execx.Runner, string, string, string, string) error { return nil }
	cmd.installProxy = func(context.Context, execx.Runner, string) error { return nil }
	cmd.trustCA = func(context.Context, execx.Runner) error { return nil }
	cmd.migrateSrv = func(context.Context, execx.Runner, string, *logging.Logger) error { return nil }

	code := cmd.run(context.Background(), []string{}, console, logger, false)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Setup complete") {
		t.Fatalf("expected setup completion")
	}
}

var _ execx.Runner = &setupRunner{}
var _ registry.Project = registry.Project{}
