package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestMigrateSrvMarker(t *testing.T) {
	root := t.TempDir()
	srv := srvMarker{Name: "legacy", Port: 4321}
	data, err := json.Marshal(srv)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".srv"), data, 0644); err != nil {
		t.Fatalf("write srv: %v", err)
	}

	marker, migrated, err := MigrateSrvMarker(root)
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if !migrated {
		t.Fatalf("expected migrated")
	}
	if marker.Name != "legacy" || marker.Port != 4321 || marker.Template != "unknown" {
		t.Fatalf("unexpected marker: %#v", marker)
	}
	if _, err := os.Stat(filepath.Join(root, ".srv")); err == nil {
		t.Fatalf("expected .srv to be removed")
	}
	if _, err := os.Stat(filepath.Join(root, ".dade")); err != nil {
		t.Fatalf("expected .dade to exist")
	}
}
