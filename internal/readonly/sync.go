package readonly

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	execx "github.com/theydontwantyoutovibecode/dade/internal/exec"
	"github.com/theydontwantyoutovibecode/dade/internal/logging"
)

func SyncDeps(ctx context.Context, runner execx.Runner, projectDir string, logger *logging.Logger) error {
	manifestPath := filepath.Join(projectDir, ".read-only", "manifest.txt")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return nil
	}

	f, err := os.Open(manifestPath)
	if err != nil {
		return fmt.Errorf("open manifest: %w", err)
	}
	defer f.Close()

	readOnlyDir := filepath.Join(projectDir, ".read-only")
	if err := os.MkdirAll(readOnlyDir, 0755); err != nil {
		return fmt.Errorf("create .read-only dir: %w", err)
	}

	scanner := bufio.NewScanner(f)
	cloned := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		repoName := strings.TrimSuffix(filepath.Base(line), ".git")
		repoDir := filepath.Join(readOnlyDir, repoName)

		if _, err := os.Stat(repoDir); err == nil {
			continue
		}

		logger.Info(fmt.Sprintf("Syncing reference: %s", repoName))
		if err := runner.Run(ctx, "git", "clone", "--depth", "1", "--quiet", line, repoDir); err != nil {
			logger.Warn(fmt.Sprintf("Failed to sync %s: %v", repoName, err))
			continue
		}
		cloned++
	}

	if cloned > 0 {
		logger.Success(fmt.Sprintf("Synced %d reference libraries", cloned))
	}

	return scanner.Err()
}
