package proxy

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
	execx "github.com/theydontwantyoutovibecode/made-in-dade/internal/exec"
)

func CreatePlist(ctx context.Context, runner execx.Runner, plistPath, caddyfilePath, logPath, errPath string) error {
	if runner == nil {
		runner = execx.NewSystemRunner()
	}
	if plistPath == "" || caddyfilePath == "" || logPath == "" || errPath == "" {
		return errors.New("plist paths are required")
	}
	caddyPath, err := runner.LookPath("caddy")
	if err != nil {
		return err
	}
	content := buildPlist(caddyPath, caddyfilePath, logPath, errPath)
	if err := os.MkdirAll(filepath.Dir(plistPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(plistPath, []byte(content), 0644)
}

func InstallProxyService(ctx context.Context, runner execx.Runner, plistPath string) error {
	if runner == nil {
		runner = execx.NewSystemRunner()
	}
	if plistPath == "" {
		return errors.New("plist path is required")
	}
	if IsProxyRunning(ctx, runner) {
		return nil
	}
	uid, err := currentUID()
	if err != nil {
		return err
	}
	return runner.Run(ctx, "launchctl", "bootstrap", "gui/"+uid, plistPath)
}

func UninstallProxyService(ctx context.Context, runner execx.Runner, plistPath string) error {
	if runner == nil {
		runner = execx.NewSystemRunner()
	}
	if plistPath == "" {
		return errors.New("plist path is required")
	}
	if !IsProxyRunning(ctx, runner) {
		return nil
	}
	uid, err := currentUID()
	if err != nil {
		return err
	}
	if err := runner.Run(ctx, "launchctl", "bootout", "gui/"+uid, plistPath); err != nil {
		return err
	}
	return nil
}

func RestartProxyService(ctx context.Context, runner execx.Runner, plistPath string) error {
	if err := UninstallProxyService(ctx, runner, plistPath); err != nil {
		return err
	}
	return InstallProxyService(ctx, runner, plistPath)
}

func buildPlist(caddyPath, caddyfilePath, logPath, errPath string) string {
	entries := []string{
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>",
		"<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">",
		"<plist version=\"1.0\">",
		"<dict>",
		"    <key>Label</key>",
		fmt.Sprintf("    <string>%s</string>", config.ProxyLabel),
		"    <key>ProgramArguments</key>",
		"    <array>",
		fmt.Sprintf("        <string>%s</string>", caddyPath),
		"        <string>run</string>",
		"        <string>--config</string>",
		fmt.Sprintf("        <string>%s</string>", caddyfilePath),
		"    </array>",
		"    <key>RunAtLoad</key>",
		"    <true/>",
		"    <key>KeepAlive</key>",
		"    <true/>",
		"    <key>StandardOutPath</key>",
		fmt.Sprintf("    <string>%s</string>", logPath),
		"    <key>StandardErrorPath</key>",
		fmt.Sprintf("    <string>%s</string>", errPath),
		"</dict>",
		"</plist>",
	}
	return strings.Join(entries, "\n")
}

func currentUID() (string, error) {
	uid := os.Getuid()
	if uid <= 0 {
		return "", errors.New("invalid user id")
	}
	return strconv.Itoa(uid), nil
}
