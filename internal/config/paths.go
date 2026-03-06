package config

import (
	"os"
	"path/filepath"
	"strings"
)

// Hostname returns the local hostname for use in domain names.
// Falls back to "localhost" if hostname cannot be determined.
func Hostname() string {
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		return "localhost"
	}
	// Remove any domain suffix (e.g., "macbook.local" -> "macbook")
	if idx := strings.Index(hostname, "."); idx != -1 {
		hostname = hostname[:idx]
	}
	return strings.ToLower(hostname)
}

// LocalDomain returns the base domain for local network sharing.
// Format: {hostname}.local (e.g., "macbook.local")
func LocalDomain() string {
	return Hostname() + ".local"
}

// ProjectDomain returns the full domain for a project.
// Format: {project}.{hostname}.local (e.g., "myapp.macbook.local")
func ProjectDomain(projectName string) string {
	return projectName + "." + LocalDomain()
}

const (
	ConfigDirName        = "dade"
	ProjectsFileName     = "projects.json"
	TemplatesDirName     = "templates"
	TemplatesFileName    = "templates.toml"
	CaddyfileName        = "Caddyfile"
	ConfigFileName       = "config.toml"
	ProxyLogName         = "proxy.log"
	ProxyErrName         = "proxy.err"
	ProxyLabel           = "land.charm.dade.proxy"
	BasePort             = 3000
	DefaultTemplatesPath = "~/.config/dade/templates.toml"
)

const launchAgentsDir = "Library/LaunchAgents"

func ConfigDir() (string, error) {
	base, err := baseConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, ConfigDirName), nil
}

func TemplatesPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, TemplatesFileName), nil
}

func TemplatesDir() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, TemplatesDirName), nil
}

func ProjectsFile() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ProjectsFileName), nil
}

func CaddyfilePath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, CaddyfileName), nil
}

func ConfigFilePath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ConfigFileName), nil
}

func ProxyLogPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ProxyLogName), nil
}

func ProxyErrPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ProxyErrName), nil
}

func ProxyPlistPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, launchAgentsDir, ProxyLabel+".plist"), nil
}

func baseConfigDir() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return xdg, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config"), nil
}
