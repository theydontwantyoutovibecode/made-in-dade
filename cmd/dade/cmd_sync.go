package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
	execx "github.com/theydontwantyoutovibecode/made-in-dade/internal/exec"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/proxy"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/registry"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/ui"
	"github.com/spf13/cobra"
)

type syncCommand struct {
	runner        execx.Runner
	projectsFile  func() (string, error)
	caddyfilePath func() (string, error)
	register      func(path, name string, port int, projectPath, template string) (registry.Project, error)
	generateCaddy func(context.Context, execx.Runner, string, string) error
	reloadProxy   func(context.Context, execx.Runner, string) error
	loadProjects  func(string) (map[string]registry.Project, error)
	unregister    func(string, string) (bool, error)
	saveProjects  func(string, map[string]registry.Project) error
}

var syncCommandFactory = defaultSyncCommand

func defaultSyncCommand() syncCommand {
	return syncCommand{
		runner:        execx.NewSystemRunner(),
		projectsFile:  config.ProjectsFile,
		caddyfilePath: config.CaddyfilePath,
		register:      registry.Register,
		generateCaddy: proxy.GenerateCaddyfile,
		reloadProxy:   proxy.ReloadProxy,
		loadProjects:  registry.Load,
		unregister:    registry.Unregister,
		saveProjects:  registry.Save,
	}
}

func runSyncCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	cleanMode, _ := cmd.Flags().GetBool("clean")

	impl := syncCommandFactory()
	code := impl.run(context.Background(), args, console, logger, cleanMode)
	if code != 0 {
		return errors.New("sync command failed")
	}
	return nil
}

func (c syncCommand) run(ctx context.Context, args []string, console *ui.UI, logger *logging.Logger, cleanMode bool) int {
	_ = console

	projectsPath, err := c.projectsFile()
	if err != nil {
		logger.Error("Failed to resolve projects file")
		return 1
	}

	caddyfilePath, err := c.caddyfilePath()
	if err != nil {
		logger.Error("Failed to resolve Caddyfile path")
		return 1
	}

	if cleanMode {
		return c.runClean(ctx, projectsPath, caddyfilePath, logger)
	}

	scanPath := os.Getenv("HOME")
	if len(args) > 0 {
		scanPath = args[0]
	}

	logger.Info(fmt.Sprintf("Scanning for .dade files in %s...", scanPath))

	if err := c.saveProjects(projectsPath, map[string]registry.Project{}); err != nil {
		logger.Error("Failed to clear registry")
		return 1
	}

	count := 0
	skipDirs := map[string]bool{
		".git":         true,
		"node_modules": true,
		".venv":        true,
		"venv":         true,
		"__pycache__":  true,
		".config":      true,
		"Library":      true,
		".Trash":       true,
	}

	err = filepath.WalkDir(scanPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() != ".dade" {
			return nil
		}
		if count >= 100 {
			return filepath.SkipAll
		}

		projectDir := filepath.Dir(path)
		marker, err := registry.ReadMarker(projectDir)
		if err != nil {
			return nil
		}
		if marker.Name == "" || marker.Port == 0 {
			return nil
		}

		if c.register != nil {
			if _, err := c.register(projectsPath, marker.Name, marker.Port, projectDir, marker.Template); err != nil {
				return nil
			}
		}
		logger.Success(fmt.Sprintf("Found: %s (%s)", marker.Name, projectDir))
		count++
		return nil
	})
	if err != nil && err != filepath.SkipAll {
		logger.Warn(fmt.Sprintf("Scan error: %v", err))
	}

	if c.generateCaddy != nil {
		_ = c.generateCaddy(ctx, c.runner, projectsPath, caddyfilePath)
	}
	if c.reloadProxy != nil {
		_ = c.reloadProxy(ctx, c.runner, caddyfilePath)
	}

	logger.Success(fmt.Sprintf("Synced %d project(s)", count))
	return 0
}

func (c syncCommand) runClean(ctx context.Context, projectsPath, caddyfilePath string, logger *logging.Logger) int {
	logger.Info("Cleaning stale registry entries...")

	projects, err := c.loadProjects(projectsPath)
	if err != nil {
		logger.Error("Failed to load projects")
		return 1
	}

	removed := 0
	for name, project := range projects {
		if _, err := os.Stat(project.Path); os.IsNotExist(err) {
			if c.unregister != nil {
				if _, err := c.unregister(projectsPath, name); err == nil {
					logger.Warn(fmt.Sprintf("Removing stale: %s (%s)", name, project.Path))
					removed++
				}
			}
			continue
		}
		if !registry.MarkerExists(project.Path) {
			if c.unregister != nil {
				if _, err := c.unregister(projectsPath, name); err == nil {
					logger.Warn(fmt.Sprintf("Removing stale: %s (%s)", name, project.Path))
					removed++
				}
			}
		}
	}

	if c.generateCaddy != nil {
		_ = c.generateCaddy(ctx, c.runner, projectsPath, caddyfilePath)
	}
	if c.reloadProxy != nil {
		_ = c.reloadProxy(ctx, c.runner, caddyfilePath)
	}

	logger.Success(fmt.Sprintf("Removed %d stale entries", removed))
	return 0
}
