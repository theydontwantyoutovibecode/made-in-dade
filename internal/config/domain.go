package config

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	defaultDomainTLD  = ".localhost"
	legacyDomainTLD   = ".local"
	configFileName     = "config.toml"
)

var (
	domainTLD     string
	domainTLDOnce sync.Once
)

// DomainTLD returns the configured top-level domain for project domains.
// Defaults to ".localhost" for new installations, ".local" for legacy.
// Can be overridden via DADE_DOMAIN_TLD environment variable or config file.
func DomainTLD() string {
	domainTLDOnce.Do(func() {
		// Check environment variable first
		if env := os.Getenv("DADE_DOMAIN_TLD"); env != "" {
			domainTLD = normalizeDomainTLD(env)
			return
		}

		// Check config file
		if tld, err := readDomainTLDFromConfig(); err == nil && tld != "" {
			domainTLD = normalizeDomainTLD(tld)
			return
		}

		// Check for legacy installation (projects.json exists)
		if isLegacyInstallation() {
			domainTLD = legacyDomainTLD
			return
		}

		// Default for new installations
		domainTLD = defaultDomainTLD
	})
	return domainTLD
}

// readDomainTLDFromConfig reads the domain_tld from config.toml
func readDomainTLDFromConfig() (string, error) {
	configDir, err := ConfigDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(configDir, configFileName)
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	// Simple parser for domain_tld setting
	// Format: domain_tld = ".localhost"
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' || line[0] == ';' {
			continue
		}
		if strings.Contains(line, "domain_tld") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				return strings.TrimSpace(trimQuotes(strings.TrimSpace(parts[1]))), nil
			}
		}
	}

	return "", nil
}

// isLegacyInstallation checks if this is an existing dade installation
// that was using .local domains before the migration
func isLegacyInstallation() bool {
	configDir, err := ConfigDir()
	if err != nil {
		return false
	}
	projectsFile := filepath.Join(configDir, ProjectsFileName)
	if _, err := os.Stat(projectsFile); err != nil {
		return false
	}

	// If config file exists with domain_tld, it's already migrated
	configPath := filepath.Join(configDir, configFileName)
	if _, err := os.Stat(configPath); err == nil {
		return false
	}

	return true
}

// normalizeDomainTLD ensures the TLD starts with a dot
func normalizeDomainTLD(tld string) string {
	tld = strings.TrimSpace(tld)
	if len(tld) == 0 {
		return defaultDomainTLD
	}
	if tld[0] != '.' {
		tld = "." + tld
	}
	return tld
}

// SetDomainTLD updates the configured domain TLD in the config file
func SetDomainTLD(tld string) error {
	configDir, err := ConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, configFileName)
	tld = normalizeDomainTLD(tld)

	// Read existing config or create new
	var content []byte
	if data, err := os.ReadFile(configPath); err == nil {
		content = data
	} else if !os.IsNotExist(err) {
		return err
	}

	// Update or add domain_tld setting
	updated := updateDomainTLDInConfig(content, tld)

	// Write back
	if err := os.WriteFile(configPath, updated, 0644); err != nil {
		return err
	}

	// Reset the cached value
	domainTLD = ""
	domainTLDOnce = sync.Once{}

	return nil
}

// updateDomainTLDInConfig updates the domain_tld setting in config content
func updateDomainTLDInConfig(content []byte, tld string) []byte {
	lines := strings.Split(string(content), "\n")
	found := false
	var result []string

	for _, line := range lines {
		lineStr := strings.TrimSpace(line)
		if strings.Contains(lineStr, "domain_tld") {
			result = append(result, "domain_tld = \""+tld+"\"")
			found = true
		} else {
			result = append(result, line)
		}
	}

	if !found && len(result) > 0 {
		// Add at the end
		result = append(result, "domain_tld = \""+tld+"\"")
	}

	return []byte(strings.Join(result, "\n"))
}

// trimQuotes removes surrounding quotes from a string
func trimQuotes(s string) string {
	if len(s) >= 2 && (s[0] == '"' || s[0] == '\'') && s[len(s)-1] == s[0] {
		return s[1 : len(s)-1]
	}
	return s
}