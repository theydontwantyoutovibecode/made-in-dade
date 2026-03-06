package readonly

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
)

type stubRunner struct {
	clonedRepos []string
}

func (s *stubRunner) Run(_ context.Context, name string, args ...string) error {
	if name == "git" && len(args) > 0 && args[0] == "clone" {
		target := args[len(args)-1]
		s.clonedRepos = append(s.clonedRepos, target)
		return os.MkdirAll(target, 0755)
	}
	return nil
}

func (s *stubRunner) Output(_ context.Context, _ string, _ ...string) (string, error) {
	return "", nil
}

func (s *stubRunner) LookPath(name string) (string, error) {
	return "/usr/bin/" + name, nil
}

func TestSyncDepsNoManifest(t *testing.T) {
	dir := t.TempDir()
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)

	err := SyncDeps(context.Background(), &stubRunner{}, dir, logger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSyncDepsClones(t *testing.T) {
	dir := t.TempDir()
	readOnlyDir := filepath.Join(dir, ".read-only")
	if err := os.MkdirAll(readOnlyDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	manifest := "# reference repos\nhttps://github.com/example/repo1.git\nhttps://github.com/example/repo2.git\n"
	if err := os.WriteFile(filepath.Join(readOnlyDir, "manifest.txt"), []byte(manifest), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	runner := &stubRunner{}
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	logger := logging.New(stdout, stderr, false)

	if err := SyncDeps(context.Background(), runner, dir, logger); err != nil {
		t.Fatalf("sync: %v", err)
	}

	if len(runner.clonedRepos) != 2 {
		t.Fatalf("expected 2 clones, got %d", len(runner.clonedRepos))
	}

	runner.clonedRepos = nil
	if err := SyncDeps(context.Background(), runner, dir, logger); err != nil {
		t.Fatalf("second sync: %v", err)
	}
	if len(runner.clonedRepos) != 0 {
		t.Fatalf("expected 0 clones on re-run, got %d", len(runner.clonedRepos))
	}
}
