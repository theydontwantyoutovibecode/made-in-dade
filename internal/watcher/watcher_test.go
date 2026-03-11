package watcher

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewWatcher(t *testing.T) {
	w, err := New(DefaultPatterns())
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer w.Stop()

	if w.watcher == nil {
		t.Fatal("watcher is nil")
	}
}

func TestAddDirectory(t *testing.T) {
	w, err := New(DefaultPatterns())
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer w.Stop()

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := w.AddDirectory(tmpDir); err != nil {
		t.Fatalf("failed to add directory: %v", err)
	}
}

func TestFileChangeDetection(t *testing.T) {
	w, err := New(DefaultPatterns())
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer w.Stop()

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := w.AddDirectory(tmpDir); err != nil {
		t.Fatalf("failed to add directory: %v", err)
	}

	// Create a channel to receive change notifications
	changeChan := make(chan string, 10)
	w.OnChange(func(path string) {
		changeChan <- path
	})

	// Start the watcher
	ctx := context.Background()
	if err := w.Start(ctx); err != nil {
		t.Fatalf("failed to start watcher: %v", err)
	}

	// Give the watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.html")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Wait for the change notification
	select {
	case path := <-changeChan:
		if path != testFile {
			t.Errorf("expected path %s, got %s", testFile, path)
		}
	case <-time.After(2 * time.Second):
		t.Error("timed out waiting for file change")
	}
}

func TestPatternMatching(t *testing.T) {
	w, err := New(Pattern{
		Include: []string{"*.html", "*.css"},
		Exclude: []string{"*.swp", ".git/*"},
	})
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer w.Stop()

	tests := []struct {
		path     string
		expected bool
	}{
		{"test.html", true},       // Matches include pattern, not excluded
		{"test.css", true},        // Matches include pattern, not excluded
		{"test.js", false},        // Doesn't match include pattern
		{"test.swp", false},       // Matches exclude pattern
		{".git/test.html", false}, // Contains .git directory
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			// A file should be watched if it matches include pattern AND is not excluded
			shouldWatch := w.matchesPattern(tt.path) && !w.isExcluded(tt.path)
			if shouldWatch != tt.expected {
				t.Errorf("expected %v for %s, got %v (matches: %v, excluded: %v)",
					tt.expected, tt.path, shouldWatch, w.matchesPattern(tt.path), w.isExcluded(tt.path))
			}
		})
	}
}
