package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Marker struct {
	Name     string `json:"name"`
	Template string `json:"template"`
	Port     int    `json:"port"`
	Created  string `json:"created"`
}

func WriteMarker(projectDir, name, template string, port int) (Marker, error) {
	marker := Marker{
		Name:     name,
		Template: template,
		Port:     port,
		Created:  time.Now().UTC().Format(time.RFC3339),
	}
	data, err := json.MarshalIndent(marker, "", "  ")
	if err != nil {
		return Marker{}, err
	}
	path := filepath.Join(projectDir, ".dade")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return Marker{}, err
	}
	return marker, nil
}

func ReadMarker(projectDir string) (Marker, error) {
	path := filepath.Join(projectDir, ".dade")
	data, err := os.ReadFile(path)
	if err != nil {
		return Marker{}, err
	}
	var marker Marker
	if err := json.Unmarshal(data, &marker); err != nil {
		return Marker{}, err
	}
	return marker, nil
}

func MarkerExists(projectDir string) bool {
	path := filepath.Join(projectDir, ".dade")
	if info, err := os.Stat(path); err == nil {
		return !info.IsDir()
	}
	return false
}

// UpdateMarkerPort updates the port in an existing marker file, preserving other fields.
func UpdateMarkerPort(projectDir string, port int) (Marker, error) {
	marker, err := ReadMarker(projectDir)
	if err != nil {
		return Marker{}, err
	}
	marker.Port = port
	data, err := json.MarshalIndent(marker, "", "  ")
	if err != nil {
		return Marker{}, err
	}
	path := filepath.Join(projectDir, ".dade")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return Marker{}, err
	}
	return marker, nil
}
