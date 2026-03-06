package registry

import (
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestRegistryCRUD(t *testing.T) {
	path := filepath.Join(t.TempDir(), "projects.json")

	project, err := Register(path, "alpha", 3000, "/tmp/alpha", "hypertext")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if project.Port != 3000 || project.Path != "/tmp/alpha" || project.Template != "hypertext" {
		t.Fatalf("unexpected project: %#v", project)
	}

	loaded, ok, err := Get(path, "alpha")
	if err != nil || !ok {
		t.Fatalf("expected project")
	}
	if loaded.Path != "/tmp/alpha" {
		t.Fatalf("unexpected project path")
	}

	entry, ok, err := GetByPath(path, "/tmp/alpha")
	if err != nil || !ok {
		t.Fatalf("expected project by path")
	}
	if entry.Name != "alpha" {
		t.Fatalf("unexpected entry name")
	}

	entries, err := List(path)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(entries) != 1 || entries[0].Name != "alpha" {
		t.Fatalf("unexpected entries")
	}

	names, err := ListNames(path)
	if err != nil {
		t.Fatalf("list names: %v", err)
	}
	if len(names) != 1 || names[0] != "alpha" {
		t.Fatalf("unexpected names")
	}

	exists, err := Exists(path, "alpha")
	if err != nil || !exists {
		t.Fatalf("expected exists")
	}

	deleted, err := Unregister(path, "alpha")
	if err != nil || !deleted {
		t.Fatalf("expected delete")
	}

	exists, err = Exists(path, "alpha")
	if err != nil || exists {
		t.Fatalf("expected missing after delete")
	}
}

func TestRegistryInvalidJSONRecovery(t *testing.T) {
	path := filepath.Join(t.TempDir(), "projects.json")
	if err := os.WriteFile(path, []byte("{bad"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	projects, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(projects) != 0 {
		t.Fatalf("expected empty projects on invalid json")
	}
}

func TestNextPortSkipsUsed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "projects.json")

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()
	occupied := listener.Addr().(*net.TCPAddr).Port

	if _, err := Register(path, "occupied", occupied, "/tmp/occupied", "hypertext"); err != nil {
		t.Fatalf("register: %v", err)
	}

	port, err := NextPort(path)
	if err != nil {
		t.Fatalf("next port: %v", err)
	}
	if port != occupied+1 {
		t.Fatalf("expected next port to be max+1")
	}
}
