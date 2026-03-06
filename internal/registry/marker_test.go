package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWriteMarker(t *testing.T) {
	projectDir := t.TempDir()
	marker, err := WriteMarker(projectDir, "demo", "hypertext", 3000)
	if err != nil {
		t.Fatalf("write marker: %v", err)
	}
	if marker.Name != "demo" || marker.Template != "hypertext" || marker.Port != 3000 {
		t.Fatalf("unexpected marker: %#v", marker)
	}
	if _, err := time.Parse(time.RFC3339, marker.Created); err != nil {
		t.Fatalf("expected valid created timestamp")
	}

	path := filepath.Join(projectDir, ".dade")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read marker: %v", err)
	}

	var decoded Marker
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if decoded.Name != "demo" || decoded.Template != "hypertext" || decoded.Port != 3000 {
		t.Fatalf("unexpected decoded marker: %#v", decoded)
	}

	readMarker, err := ReadMarker(projectDir)
	if err != nil {
		t.Fatalf("read marker: %v", err)
	}
	if readMarker.Name != "demo" || readMarker.Template != "hypertext" || readMarker.Port != 3000 {
		t.Fatalf("unexpected read marker: %#v", readMarker)
	}
	if !MarkerExists(projectDir) {
		t.Fatalf("expected marker to exist")
	}
}
