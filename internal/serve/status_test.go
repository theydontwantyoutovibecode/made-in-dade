package serve

import (
	"os"
	"path/filepath"
	"testing"
)

type stubLookup struct {
	current int
	project int
	currentErr error
	projectErr error
}

func (s stubLookup) CurrentPort(_ string) (int, error) {
	return s.current, s.currentErr
}

func (s stubLookup) ProjectPort(_ string) (int, error) {
	return s.project, s.projectErr
}

func TestIsProjectRunningWithPIDFile(t *testing.T) {
	root := t.TempDir()
	pidFile := filepath.Join(root, DefaultPIDFile)
	if err := os.WriteFile(pidFile, []byte("99999"), 0644); err != nil {
		t.Fatalf("write pid: %v", err)
	}

	running, err := IsProjectRunning(root, "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if running {
		t.Fatalf("expected not running for invalid pid")
	}
}

func TestIsProjectRunningWithPortLookup(t *testing.T) {
	root := t.TempDir()
	lookup := stubLookup{current: 0, project: 0}
	if running, err := IsProjectRunning(root, "", lookup); err != nil || running {
		t.Fatalf("expected not running with empty port")
	}
}
