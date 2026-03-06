package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/theydontwantyoutovibecode/dade/internal/config"
	execx "github.com/theydontwantyoutovibecode/dade/internal/exec"
	"github.com/theydontwantyoutovibecode/dade/internal/logging"
)

type defaultTemplatesCommand struct {
	runner       execx.Runner
	templatesDir func() (string, error)
	writeFile    func(string, []byte, os.FileMode) error
}

func defaultDefaultTemplatesCommand() defaultTemplatesCommand {
	return defaultTemplatesCommand{
		runner:       execx.NewSystemRunner(),
		templatesDir: config.TemplatesDir,
		writeFile:    os.WriteFile,
	}
}

func ensureDefaultTemplates(ctx context.Context, logger *logging.Logger) error {
	cmd := defaultDefaultTemplatesCommand()
	return cmd.run(ctx, logger)
}

func (c defaultTemplatesCommand) run(ctx context.Context, logger *logging.Logger) error {
	templatesDir, err := c.templatesDir()
	if err != nil {
		return fmt.Errorf("resolve templates dir: %w", err)
	}
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return fmt.Errorf("create templates dir: %w", err)
	}

	if !execx.CommandAvailable(c.runner, "git") {
		return errors.New("git is required to install default templates")
	}

	defaults := config.DefaultTemplates()
	for _, tpl := range defaults.Ordered {
		target := filepath.Join(templatesDir, tpl.Name)
		manifestPath := filepath.Join(target, "dade.toml")
		if _, err := os.Stat(manifestPath); err == nil {
			continue
		}

		_ = os.RemoveAll(target)

		logger.Info(fmt.Sprintf("Installing default template: %s", tpl.Name))
		if err := c.runner.Run(ctx, "git", "clone", "--depth", "1", tpl.URL, target); err != nil {
			logger.Warn(fmt.Sprintf("Failed to install %s: %v", tpl.Name, err))
			continue
		}

		sourcePath := filepath.Join(target, ".source")
		_ = c.writeFile(sourcePath, []byte(tpl.URL), 0644)

		defaultMarker := filepath.Join(target, ".default")
		_ = c.writeFile(defaultMarker, []byte(""), 0644)

		logger.Success(fmt.Sprintf("Installed: %s", tpl.Name))
	}

	return nil
}

func isDefaultTemplate(templatePath string) bool {
	_, err := os.Stat(filepath.Join(templatePath, ".default"))
	return err == nil
}
