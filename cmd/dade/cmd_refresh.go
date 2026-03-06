package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
	execx "github.com/theydontwantyoutovibecode/made-in-dade/internal/exec"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/proxy"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/registry"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/ui"
	"github.com/spf13/cobra"
)

type refreshCommand struct {
	runner        execx.Runner
	projectsFile  func() (string, error)
	caddyfilePath func() (string, error)
	listProjects  func(string) ([]registry.Entry, error)
	generateCaddy func(context.Context, execx.Runner, string, string) error
	reloadProxy   func(context.Context, execx.Runner, string) error
	localDomain   func() string
}

var refreshCommandFactory = defaultRefreshCommand

func defaultRefreshCommand() refreshCommand {
	return refreshCommand{
		runner:        execx.NewSystemRunner(),
		projectsFile:  config.ProjectsFile,
		caddyfilePath: config.CaddyfilePath,
		listProjects:  registry.List,
		generateCaddy: proxy.GenerateCaddyfile,
		reloadProxy:   proxy.ReloadProxy,
		localDomain:   config.LocalDomain,
	}
}

func runRefreshCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	showList, _ := cmd.Flags().GetBool("list")

	impl := refreshCommandFactory()
	code := impl.run(context.Background(), console, logger, showList)
	if code != 0 {
		return errors.New("refresh command failed")
	}
	return nil
}

func (c refreshCommand) run(ctx context.Context, console *ui.UI, logger *logging.Logger, showList bool) int {
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

	entries, err := c.listProjects(projectsPath)
	if err != nil {
		logger.Error("Failed to list projects")
		return 1
	}

	if len(entries) == 0 {
		logger.Info("No projects registered")
		return 0
	}

	if c.generateCaddy != nil {
		if err := c.generateCaddy(ctx, c.runner, projectsPath, caddyfilePath); err != nil {
			logger.Error("Failed to generate Caddyfile")
			return 1
		}
	}

	if c.reloadProxy != nil {
		if err := c.reloadProxy(ctx, c.runner, caddyfilePath); err != nil {
			logger.Error("Failed to reload proxy")
			return 1
		}
	}

	localDomain := c.localDomain()
	logger.Success(fmt.Sprintf("Refreshed %d project(s)", len(entries)))
	logger.Info(fmt.Sprintf("Domain: *.%s", localDomain))

	if showList {
		logger.Info("")
		logger.Info("Project URLs:")
		for _, entry := range entries {
			logger.Info(fmt.Sprintf("  https://%s.%s", entry.Name, localDomain))
		}
	}

	return 0
}
