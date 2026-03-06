package config

import (
	"errors"
	"os"
	"path/filepath"
)

const srvProxyLabel = "land.charm.srv.proxy"

func DetectSrvConfig() (bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}

	srvConfigDir := filepath.Join(home, ".config", "srv")
	if _, err := os.Stat(srvConfigDir); err == nil {
		return true, nil
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return false, err
	}

	srvPlist := filepath.Join(home, "Library", "LaunchAgents", srvProxyLabel+".plist")
	if _, err := os.Stat(srvPlist); err == nil {
		return true, nil
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return false, err
	}

	return false, nil
}
