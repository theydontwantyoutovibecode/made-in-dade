package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
	execx "github.com/theydontwantyoutovibecode/made-in-dade/internal/exec"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/manifest"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/ui"
)

type installCommand struct {
	runner      execx.Runner
	readFile    func(string) ([]byte, error)
	writeFile   func(string, []byte, os.FileMode) error
	removeAll   func(string) error
	tempDir     func(string, string) (string, error)
	rename      func(string, string) error
	templatesDir func() (string, error)
	spin        func(message string, work func() error) error
}

func defaultInstallCommand() installCommand {
	return installCommand{
		runner:       execx.NewSystemRunner(),
		readFile:     os.ReadFile,
		writeFile:    os.WriteFile,
		removeAll:    os.RemoveAll,
		tempDir:      os.MkdirTemp,
		rename:       os.Rename,
		templatesDir: config.TemplatesDir,
	}
}

var installCommandFactory = defaultInstallCommand

func runInstall(args []string, console *ui.UI, logger *logging.Logger, styled bool) int {
	cmd := installCommandFactory()
	return cmd.run(context.Background(), args, console, logger, styled)
}

func (c installCommand) run(ctx context.Context, args []string, console *ui.UI, logger *logging.Logger, styled bool) int {
	if _, err := config.InitConfig(); err != nil {
		logger.Error("Failed to initialize config")
		return 1
	}

	name := ""
	url := ""
	listOfficial := false
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--name":
			if i+1 >= len(args) {
				logger.Error("Missing value for --name")
				return 1
			}
			name = args[i+1]
			i++
		case "--list-official":
			listOfficial = true
		default:
			if strings.HasPrefix(args[i], "-") {
				logger.Error(fmt.Sprintf("Unknown option: %s", args[i]))
				return 1
			}
			if url == "" {
				url = args[i]
			}
		}
	}

	if listOfficial {
		console.PrintHelp(officialTemplatesText(styled))
		return 0
	}
	if url == "" {
		logger.Error("Missing template git URL")
		return 1
	}
	if !execx.CommandAvailable(c.runner, "git") {
		logger.Error("git is required to install templates")
		return 1
	}

	templatesDir, err := c.templatesDir()
	if err != nil {
		logger.Error("Failed to resolve templates directory")
		return 1
	}

	tmpDir, err := c.tempDir("", "dade-template-*")
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

	if err := spin("Cloning template", func() error {
		return c.runner.Run(ctx, "git", "clone", "--depth", "1", url, tmpDir)
	}); err != nil {
		logger.Error(err.Error())
		return 1
	}

	manifestPath := filepath.Join(tmpDir, "dade.toml")
	manifestData, err := c.readFile(manifestPath)
	if err != nil {
		logger.Error("Template missing dade.toml manifest")
		return 1
	}
	parsedManifest, err := manifest.Parse(manifestData)
	if err != nil {
		logger.Error(err.Error())
		return 1
	}
	if err := manifest.Validate(parsedManifest); err != nil {
		logger.Error(err.Error())
		return 1
	}
	if name == "" {
		name = manifest.TemplateName(parsedManifest)
	}
	if name == "" {
		logger.Error("Template name is required")
		return 1
	}

	target := filepath.Join(templatesDir, name)
	if _, err := os.Stat(target); err == nil {
		logger.Error(fmt.Sprintf("Template '%s' already installed", name))
		return 1
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Error("Failed to check existing template")
		return 1
	}

	if err := c.rename(tmpDir, target); err != nil {
		logger.Error("Failed to install template")
		return 1
	}
	tmpDir = ""

	sourcePath := filepath.Join(target, ".source")
	if err := c.writeFile(sourcePath, []byte(url), 0644); err != nil {
		logger.Error("Failed to write template source")
		return 1
	}

	logger.Success(fmt.Sprintf("Installed template: %s", name))
	return 0
}

func officialTemplatesText(styled bool) string {
	return availableTemplatesText(config.DefaultTemplates(), config.DefaultTemplatesPath, styled)
}
