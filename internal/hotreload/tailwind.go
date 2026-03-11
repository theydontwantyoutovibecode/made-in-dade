package hotreload

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// TailwindWatcher watches for CSS changes and coordinates compilation with reload events
type TailwindWatcher struct {
	projectDir    string
	InputCSS      string
	OutputCSS     string
	reloadFunc    func()
	watcherCancel context.CancelFunc
	isCompiling  bool
	lastModTime   time.Time
}

// NewTailwindWatcher creates a new Tailwind watcher
func NewTailwindWatcher(projectDir, inputCSS, outputCSS string) *TailwindWatcher {
	return &TailwindWatcher{
		projectDir: projectDir,
		InputCSS:   inputCSS,
		OutputCSS:  outputCSS,
	}
}

// SetReloadFunc sets the function to call after successful compilation
func (tw *TailwindWatcher) SetReloadFunc(fn func()) {
	tw.reloadFunc = fn
}

// Start begins watching the input CSS file
func (tw *TailwindWatcher) Start(ctx context.Context) error {
	ctx, tw.watcherCancel = context.WithCancel(ctx)

	// Initialize lastModTime with current file time
	inputPath := filepath.Join(tw.projectDir, tw.InputCSS)
	if modTime, err := getModTime(inputPath); err == nil {
		tw.lastModTime = modTime
	}

	// Watch for input CSS file changes
	watchInterval := 500 * time.Millisecond

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(watchInterval):
				tw.checkAndCompile()
			}
		}
	}()

	return nil
}

// checkAndCompile checks if the input CSS file has changed and compiles if necessary
func (tw *TailwindWatcher) checkAndCompile() {
	inputPath := filepath.Join(tw.projectDir, tw.InputCSS)
	outputPath := filepath.Join(tw.projectDir, tw.OutputCSS)

	// Get current modification time
	currentModTime, err := getModTime(inputPath)
	if err != nil {
		return
	}

	// Check if file has been modified
	if currentModTime.After(tw.lastModTime) && !tw.isCompiling {
		tw.CompileAndReload(inputPath, outputPath)
		tw.lastModTime = currentModTime
	}
}

// CompileAndReload compiles Tailwind CSS and triggers a reload
func (tw *TailwindWatcher) CompileAndReload(inputPath, outputPath string) {
	tw.isCompiling = true
	defer func() { tw.isCompiling = false }()

	fmt.Printf("[Tailwind] Compiling %s...\n", inputPath)

	// Use project-local npm Tailwind v3 binary
	tailwindBin := filepath.Join(tw.projectDir, "node_modules", ".bin", "tailwindcss")

	// Run Tailwind compilation with content scanning
	cmd := exec.Command(tailwindBin, "-i", inputPath, "-o", outputPath, "--content", "index.html")
	cmd.Dir = tw.projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("[Tailwind] Compilation error: %v\n%s\n", err, string(output))
		return
	}

	fmt.Printf("[Tailwind] Compiled successfully to %s\n", outputPath)

	// Trigger reload after compilation
	if tw.reloadFunc != nil {
		tw.reloadFunc()
	}
}

// Stop stops the Tailwind watcher
func (tw *TailwindWatcher) Stop() {
	if tw.watcherCancel != nil {
		tw.watcherCancel()
	}
}

// getModTime returns the modification time of a file
func getModTime(path string) (time.Time, error) {
	var lastModTime time.Time

	// Find the file matching the pattern
	matches, err := filepath.Glob(path)
	if err != nil {
		return lastModTime, err
	}

	if len(matches) == 0 {
		return lastModTime, fmt.Errorf("file not found: %s", path)
	}

	// Get the modification time of the file
	info, err := os.Stat(matches[0])
	if err != nil {
		return lastModTime, err
	}

	return info.ModTime(), nil
}
