package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/registry"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/ui"
)

type uninstallCommand struct {
	removeAll    func(string) error
	readDir      func(string) ([]os.DirEntry, error)
	templatesDir func() (string, error)
	projectsFile func() (string, error)
	loadProjects func(string) (map[string]registry.Project, error)
	confirm      func(string) (bool, error)
}

func defaultUninstallCommand() uninstallCommand {
	return uninstallCommand{
		removeAll:    os.RemoveAll,
		readDir:      os.ReadDir,
		templatesDir: config.TemplatesDir,
		projectsFile: config.ProjectsFile,
		loadProjects: registry.Load,
	}
}

var uninstallCommandFactory = defaultUninstallCommand

func runUninstall(args []string, console *ui.UI, logger *logging.Logger) int {
	cmd := uninstallCommandFactory()
	return cmd.run(context.Background(), args, console, logger)
}

func (c uninstallCommand) run(_ context.Context, args []string, console *ui.UI, logger *logging.Logger) int {
	if _, err := config.InitConfig(); err != nil {
		logger.Error("Failed to initialize config")
		return 1
	}

	all := false
	name := ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--all":
			all = true
		default:
			if strings.HasPrefix(args[i], "-") {
				logger.Error(fmt.Sprintf("Unknown option: %s", args[i]))
				return 1
			}
			if name == "" {
				name = args[i]
			}
		}
	}

	templatesDir, err := c.templatesDir()
	if err != nil {
		logger.Error("Failed to resolve templates directory")
		return 1
	}

	if all {
		confirmed, err := c.confirm("Remove ALL installed templates?")
		if err != nil {
			logger.Error("Failed to read confirmation")
			return 1
		}
		if !confirmed {
			return 0
		}

		entries, err := c.readDir(templatesDir)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				logger.Info("No templates installed.")
				return 0
			}
			logger.Error("Failed to read templates directory")
			return 1
		}
		if len(entries) == 0 {
			logger.Info("No templates installed.")
			return 0
		}
		for _, entry := range entries {
			if err := c.removeAll(filepath.Join(templatesDir, entry.Name())); err != nil {
				logger.Error("Failed to remove templates")
				return 1
			}
		}
		logger.Success("All templates removed")
		return 0
	}

	if name == "" {
		logger.Error("Template name is required")
		return 1
	}

	templatePath := filepath.Join(templatesDir, name)
	info, err := os.Stat(templatePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logger.Error(fmt.Sprintf("Template '%s' not found", name))
			printInstalledTemplates(templatesDir, console, logger)
			return 1
		}
		logger.Error("Failed to check template")
		return 1
	}
	if info.IsDir() == false {
		logger.Error(fmt.Sprintf("Template '%s' not found", name))
		printInstalledTemplates(templatesDir, console, logger)
		return 1
	}

	projectsPath, err := c.projectsFile()
	if err != nil {
		logger.Error("Failed to resolve project registry")
		return 1
	}
	projects, err := c.loadProjects(projectsPath)
	if err != nil {
		logger.Error("Failed to load project registry")
		return 1
	}

	using := projectsUsingTemplate(projects, name)
	if len(using) > 0 {
		logger.Warn("Projects using this template:")
		for _, project := range using {
			logger.Info(fmt.Sprintf("  - %s", project))
		}
		confirmed, err := c.confirm("Uninstall anyway? (projects will continue to work)")
		if err != nil {
			logger.Error("Failed to read confirmation")
			return 1
		}
		if !confirmed {
			return 0
		}
	}

	if err := c.removeAll(templatePath); err != nil {
		logger.Error("Failed to remove template")
		return 1
	}
	logger.Success(fmt.Sprintf("Removed template: %s", name))
	return 0
}

func printInstalledTemplates(templatesDir string, console *ui.UI, logger *logging.Logger) {
	installed, err := loadInstalledTemplates(templatesDir)
	if err != nil || len(installed) == 0 {
		return
	}
	console.PrintHelp("Installed templates:")
	for _, tpl := range installed {
		console.PrintHelp(fmt.Sprintf("  - %s", tpl.Name))
	}
}

func projectsUsingTemplate(projects map[string]registry.Project, template string) []string {
	if len(projects) == 0 {
		return nil
	}
	using := make([]string, 0)
	for name, project := range projects {
		if project.Template == template {
			using = append(using, name)
		}
	}
	sort.Strings(using)
	return using
}
