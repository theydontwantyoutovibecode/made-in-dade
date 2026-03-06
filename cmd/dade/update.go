package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
	execx "github.com/theydontwantyoutovibecode/made-in-dade/internal/exec"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/ui"
)

type updateCommand struct {
	runner       execx.Runner
	readFile     func(string) ([]byte, error)
	writeFile    func(string, []byte, os.FileMode) error
	removeAll    func(string) error
	tempDir      func(string, string) (string, error)
	rename       func(string, string) error
	templatesDir func() (string, error)
	readDir      func(string) ([]os.DirEntry, error)
	spin         func(message string, work func() error) error
}

func defaultUpdateCommand() updateCommand {
	return updateCommand{
		runner:       execx.NewSystemRunner(),
		readFile:     os.ReadFile,
		writeFile:    os.WriteFile,
		removeAll:    os.RemoveAll,
		tempDir:      os.MkdirTemp,
		rename:       os.Rename,
		templatesDir: config.TemplatesDir,
		readDir:      os.ReadDir,
	}
}

var updateCommandFactory = defaultUpdateCommand

func (c updateCommand) run(ctx context.Context, args []string, console *ui.UI, logger *logging.Logger, styled bool) int {
	if _, err := config.InitConfig(); err != nil {
		logger.Error("Failed to initialize config")
		return 1
	}

	updateAll := false
	name := ""
	for _, arg := range args {
		if arg == "--all" {
			updateAll = true
		} else {
			name = arg
		}
	}

	if !execx.CommandAvailable(c.runner, "git") {
		logger.Error("git is required to update templates")
		return 1
	}

	templatesDir, err := c.templatesDir()
	if err != nil {
		logger.Error("Failed to resolve templates directory")
		return 1
	}

	if updateAll {
		return c.updateAll(ctx, templatesDir, logger, styled)
	}

	if name == "" {
		logger.Error("Missing template name")
		return 1
	}

	return c.updateTemplate(ctx, templatesDir, name, logger, styled)
}

func (c updateCommand) updateAll(ctx context.Context, templatesDir string, logger *logging.Logger, styled bool) int {
	entries, err := c.readDir(templatesDir)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info("No templates installed")
			return 0
		}
		logger.Error("Failed to read templates directory")
		return 1
	}

	updated := 0
	failed := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if code := c.updateTemplate(ctx, templatesDir, name, logger, styled); code == 0 {
			updated++
		} else {
			failed++
		}
	}

	if updated > 0 {
		logger.Success(fmt.Sprintf("Updated %d template(s)", updated))
	}
	if failed > 0 {
		logger.Warn(fmt.Sprintf("Failed to update %d template(s)", failed))
	}
	if updated == 0 && failed == 0 {
		logger.Info("No templates to update")
	}

	return 0
}

func (c updateCommand) updateTemplate(ctx context.Context, templatesDir, name string, logger *logging.Logger, styled bool) int {
	templateDir := filepath.Join(templatesDir, name)

	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		logger.Error(fmt.Sprintf("Template '%s' not found", name))
		return 1
	}

	sourcePath := filepath.Join(templateDir, ".source")
	sourceData, err := c.readFile(sourcePath)
	if err != nil {
		logger.Warn(fmt.Sprintf("No source URL for '%s' - cannot update", name))
		return 1
	}
	sourceURL := string(sourceData)
	if sourceURL == "" {
		logger.Warn(fmt.Sprintf("Empty source URL for '%s' - cannot update", name))
		return 1
	}

	tmpDir, err := c.tempDir("", "dade-update-*")
	if err != nil {
		logger.Error("Failed to create temp directory")
		return 1
	}
	defer func() {
		if tmpDir != "" {
			_ = c.removeAll(tmpDir)
		}
	}()

	spin := c.spin
	if spin == nil {
		spin = func(_ string, work func() error) error { return work() }
	}

	if err := spin(fmt.Sprintf("Fetching %s", name), func() error {
		return c.runner.Run(ctx, "git", "clone", "--depth", "1", sourceURL, tmpDir)
	}); err != nil {
		logger.Error(fmt.Sprintf("Failed to fetch %s: %v", name, err))
		return 1
	}

	if err := c.writeFile(filepath.Join(tmpDir, ".source"), sourceData, 0644); err != nil {
		logger.Error("Failed to preserve source URL")
		return 1
	}

	if err := c.removeAll(templateDir); err != nil {
		logger.Error(fmt.Sprintf("Failed to remove old template: %v", err))
		return 1
	}

	if err := c.rename(tmpDir, templateDir); err != nil {
		logger.Error(fmt.Sprintf("Failed to install updated template: %v", err))
		return 1
	}
	tmpDir = ""

	logger.Success(fmt.Sprintf("Updated: %s", name))
	return 0
}
