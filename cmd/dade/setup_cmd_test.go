package main

import (
	"context"
	"strings"
	"testing"

	execx "github.com/theydontwantyoutovibecode/made-in-dade/internal/exec"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/ui"
)

func TestSetupCmdYesSkipsPrompts(t *testing.T) {
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
	rootCmd.SetArgs([]string{"setup", "-y"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
	if !strings.Contains(stdout.String(), "Setup complete") {
		t.Fatalf("expected setup completion")
	}
}

func TestSetupCmdConflictingFlags(t *testing.T) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
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
	rootCmd.SetArgs([]string{"setup", "--migrate", "--no-migrate"})

	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected conflicting flags error")
	}
}
