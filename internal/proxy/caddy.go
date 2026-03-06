package proxy

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
	execx "github.com/theydontwantyoutovibecode/made-in-dade/internal/exec"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/registry"
)

func GenerateCaddyfile(ctx context.Context, runner execx.Runner, projectsPath, caddyfilePath string) error {
	if caddyfilePath == "" {
		return errors.New("caddyfile path is required")
	}
	if runner == nil {
		runner = execx.NewSystemRunner()
	}
	entries, err := registry.List(projectsPath)
	if err != nil {
		return err
	}
	content := buildCaddyfile(entries)

	if err := os.MkdirAll(filepath.Dir(caddyfilePath), 0755); err != nil {
		return err
	}

	tempFile, err := os.CreateTemp(filepath.Dir(caddyfilePath), "Caddyfile.*")
	if err != nil {
		return err
	}
	tempPath := tempFile.Name()
	defer func() {
		_ = os.Remove(tempPath)
	}()

	if _, err := tempFile.Write([]byte(content)); err != nil {
		_ = tempFile.Close()
		return err
	}
	if err := tempFile.Close(); err != nil {
		return err
	}

	if err := ValidateCaddyfile(ctx, runner, tempPath); err != nil {
		return err
	}

	if err := backupFile(caddyfilePath); err != nil {
		return err
	}

	if err := os.Rename(tempPath, caddyfilePath); err != nil {
		return err
	}
	return nil
}

func ValidateCaddyfile(ctx context.Context, runner execx.Runner, caddyfilePath string) error {
	if caddyfilePath == "" {
		return errors.New("caddyfile path is required")
	}
	if runner == nil {
		runner = execx.NewSystemRunner()
	}
	return runner.Run(ctx, "caddy", "validate", "--config", caddyfilePath)
}

func ReloadProxy(ctx context.Context, runner execx.Runner, caddyfilePath string) error {
	if caddyfilePath == "" {
		return errors.New("caddyfile path is required")
	}
	if runner == nil {
		runner = execx.NewSystemRunner()
	}
	if !IsProxyRunning(ctx, runner) {
		return nil
	}
	return runner.Run(ctx, "caddy", "reload", "--config", caddyfilePath)
}

func IsProxyRunning(ctx context.Context, runner execx.Runner) bool {
	if runner == nil {
		runner = execx.NewSystemRunner()
	}
	return runner.Run(ctx, "launchctl", "list", config.ProxyLabel) == nil
}

func buildCaddyfile(entries []registry.Entry) string {
	localDomain := config.LocalDomain()
	var builder strings.Builder
	builder.WriteString("{\n\tlocal_certs\n}\n\n")
	for _, entry := range entries {
		if entry.Name == "" || entry.Project.Port <= 0 {
			continue
		}
		builder.WriteString(fmt.Sprintf("https://%s.%s {\n\treverse_proxy localhost:%d\n}\n\n", entry.Name, localDomain, entry.Project.Port))
	}
	return builder.String()
}

func backupFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if info.IsDir() {
		return errors.New("caddyfile path is a directory")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path+".bak", data, info.Mode().Perm())
}
