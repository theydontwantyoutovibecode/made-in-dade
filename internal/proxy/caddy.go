package proxy

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/exec"
	execx "github.com/theydontwantyoutovibecode/made-in-dade/internal/exec"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/hosts"
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

	// Update hosts file if using .local TLD
	if config.DomainTLD() == ".local" {
		if err := updateHostsFileForEntries(entries); err != nil {
			// Don't fail the whole operation, just log the error
			fmt.Fprintf(os.Stderr, "Warning: failed to update hosts file: %v\n", err)
			fmt.Fprintf(os.Stderr, "  Note: .local domains require /etc/hosts entries for LAN access\n")
			fmt.Fprintf(os.Stderr, "  Run 'sudo dade proxy update-hosts' to fix\n")
		}
	}

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

func updateHostsFileForEntries(entries []registry.Entry) error {
	if !exec.CanUseSudo() {
		return fmt.Errorf("sudo not available, cannot update hosts file")
	}

	// Request sudo privileges
	if err := exec.RequestSudo(); err != nil {
		return fmt.Errorf("failed to get sudo access: %w", err)
	}

	// Get current domains in hosts file
	currentDomains, err := hosts.ListDomains()
	if err != nil {
		return fmt.Errorf("failed to read hosts file: %w", err)
	}

	// Build set of desired domains (all .local domains)
	desiredDomains := make(map[string]bool)
	for _, entry := range entries {
		if entry.Name == "" || entry.Project.Port <= 0 {
			continue
		}
		domain := fmt.Sprintf("%s.%s", entry.Name, config.LocalDomain())
		if strings.HasSuffix(domain, ".local") {
			desiredDomains[domain] = true
		}
	}

	// Remove old domains that are no longer needed
	currentSet := make(map[string]bool)
	for _, d := range currentDomains {
		currentSet[d] = true
	}

	for domain := range currentSet {
		if !desiredDomains[domain] && strings.HasSuffix(domain, ".local") {
			if removed, err := hosts.RemoveDomainEntry(domain); err != nil {
				return fmt.Errorf("failed to remove domain %s: %w", domain, err)
			} else if removed {
				fmt.Fprintf(os.Stderr, "Removed %s from hosts file\n", domain)
			}
		}
	}

	// Add new domains
	for domain := range desiredDomains {
		if !currentSet[domain] {
			// Use sudo to add the entry
			domainToAdd := domain
			if err := exec.RunWithSudo("bash", "-c",
				fmt.Sprintf("echo '127.0.0.1\t%s' >> %s", domainToAdd, hosts.HostsPath)); err != nil {
				return fmt.Errorf("failed to add domain %s: %w", domainToAdd, err)
			}
			fmt.Fprintf(os.Stderr, "Added %s to hosts file\n", domain)
		}
	}

	return nil
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
