package config

import (
	"errors"
	"os"
	"path/filepath"
)

const defaultCaddyfile = `{
	local_certs
}
`

func IsConfigured() bool {
	dir, err := ConfigDir()
	if err != nil {
		return false
	}
	projectsFile := filepath.Join(dir, ProjectsFileName)
	if _, err := os.Stat(projectsFile); err != nil {
		return false
	}
	return true
}

func InitConfig() (bool, error) {
	configDir, err := ConfigDir()
	if err != nil {
		return false, err
	}
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return false, err
	}

	templatesDir := filepath.Join(configDir, TemplatesDirName)
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return false, err
	}

	projectsFile := filepath.Join(configDir, ProjectsFileName)
	if err := ensureFile(projectsFile, []byte("{}")); err != nil {
		return false, err
	}

	caddyfile := filepath.Join(configDir, CaddyfileName)
	if err := ensureFile(caddyfile, []byte(defaultCaddyfile)); err != nil {
		return false, err
	}

	srvDetected, _ := DetectSrvConfig()
	return srvDetected, nil
}

func ensureFile(path string, content []byte) error {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return errors.New("path exists and is a directory")
		}
		return nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)
}
