package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDirAndRemoveGit(t *testing.T) {
	src := t.TempDir()
	sub := filepath.Join(src, "sub")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sub, "file.txt"), []byte("hello"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	gitDir := filepath.Join(src, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir git: %v", err)
	}

	dst := filepath.Join(t.TempDir(), "dst")
	if err := CopyDir(src, dst); err != nil {
		t.Fatalf("copy: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dst, "sub", "file.txt")); err != nil {
		t.Fatalf("expected file copied: %v", err)
	}
	if err := RemoveGitDir(dst); err != nil {
		t.Fatalf("remove git: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dst, ".git")); !os.IsNotExist(err) {
		t.Fatalf("expected .git removed")
	}
}

func TestIsExecutable(t *testing.T) {
	file := filepath.Join(t.TempDir(), "exec.sh")
	if err := os.WriteFile(file, []byte("echo hi"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	exec, err := IsExecutable(file)
	if err != nil {
		t.Fatalf("is exec: %v", err)
	}
	if exec {
		t.Fatalf("expected not executable")
	}
	if err := os.Chmod(file, 0755); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	exec, err = IsExecutable(file)
	if err != nil {
		t.Fatalf("is exec: %v", err)
	}
	if !exec {
		t.Fatalf("expected executable")
	}
}
