package registry

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type srvMarker struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

func MigrateSrvMarker(projectDir string) (Marker, bool, error) {
	srvPath := filepath.Join(projectDir, ".srv")
	data, err := os.ReadFile(srvPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Marker{}, false, nil
		}
		return Marker{}, false, err
	}
	var srv srvMarker
	if err := json.Unmarshal(data, &srv); err != nil {
		return Marker{}, false, err
	}
	marker := Marker{
		Name:     srv.Name,
		Template: "unknown",
		Port:     srv.Port,
		Created:  time.Now().UTC().Format(time.RFC3339),
	}
	markerData, err := json.MarshalIndent(marker, "", "  ")
	if err != nil {
		return Marker{}, false, err
	}
	dadePath := filepath.Join(projectDir, ".dade")
	if err := os.WriteFile(dadePath, markerData, 0644); err != nil {
		return Marker{}, false, err
	}
	_ = os.Remove(srvPath)
	return marker, true, nil
}
