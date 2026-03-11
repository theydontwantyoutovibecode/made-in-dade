package watcher

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// Pattern represents a file pattern to watch
type Pattern struct {
	Include []string // File patterns to include (e.g., "*.html", "*.css")
	Exclude []string // File patterns to exclude (e.g., "*.swp", ".git/*")
}

// DefaultPatterns returns common patterns for web development
func DefaultPatterns() Pattern {
	return Pattern{
		Include: []string{"*.html", "*.css", "*.js", "*.htm"},
		Exclude: []string{"*.swp", "*.tmp", ".git/*", "node_modules/*", ".dade*", "css/output.css"},
	}
}

// ChangeFunc is the callback function called when a file changes
type ChangeFunc func(path string)

// Watcher monitors files for changes and triggers callbacks
type Watcher struct {
	watcher   *fsnotify.Watcher
	patterns  Pattern
	callbacks map[string][]ChangeFunc
	ctx       context.Context
	cancel    context.CancelFunc
}

// New creates a new file watcher
func New(patterns Pattern) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &Watcher{
		watcher:   fsWatcher,
		patterns:  patterns,
		callbacks: make(map[string][]ChangeFunc),
	}, nil
}

// AddDirectory adds a directory to watch recursively
func (w *Watcher) AddDirectory(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors and continue
		}

		// Skip excluded paths
		if w.isExcluded(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Watch directories, not files
		if info.IsDir() {
			if err := w.watcher.Add(path); err != nil {
				// Some directories can't be watched, just log and continue
				fmt.Printf("Warning: failed to watch %s: %v\n", path, err)
			}
		}

		return nil
	})
}

// OnChange registers a callback for file changes
func (w *Watcher) OnChange(callback ChangeFunc) {
	// Use a generic key for all callbacks
	w.callbacks["*"] = append(w.callbacks["*"], callback)
}

// Start begins watching for file changes
func (w *Watcher) Start(ctx context.Context) error {
	w.ctx, w.cancel = context.WithCancel(ctx)

	go func() {
		defer w.watcher.Close()

		for {
			select {
			case <-w.ctx.Done():
				return
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}

				// Only process create and write events
				if event.Op&(fsnotify.Create|fsnotify.Write) == 0 {
					continue
				}

				// Skip if path is excluded
				if w.isExcluded(event.Name) {
					continue
				}

				// Skip if path doesn't match include patterns
				if !w.matchesPattern(event.Name) {
					continue
				}

				// Trigger all callbacks
				for _, cb := range w.callbacks["*"] {
					cb(event.Name)
				}

			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("Watcher error: %v\n", err)
			}
		}
	}()

	return nil
}

// Stop stops watching for file changes
func (w *Watcher) Stop() {
	if w.cancel != nil {
		w.cancel()
	}
}

// isExcluded checks if a path matches any exclude pattern
func (w *Watcher) isExcluded(path string) bool {
	for _, pattern := range w.patterns.Exclude {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
		// Check for directory exclusion
		if strings.Contains(path, strings.TrimSuffix(pattern, "/*")) {
			return true
		}
	}
	return false
}

// matchesPattern checks if a path matches any include pattern
func (w *Watcher) matchesPattern(path string) bool {
	filename := filepath.Base(path)
	for _, pattern := range w.patterns.Include {
		if matched, _ := filepath.Match(pattern, filename); matched {
			return true
		}
	}
	return len(w.patterns.Include) == 0 // Include all if no patterns specified
}
