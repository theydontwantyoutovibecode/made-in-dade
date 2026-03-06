package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/theydontwantyoutovibecode/dade/internal/config"
	"github.com/theydontwantyoutovibecode/dade/internal/logging"
	"github.com/theydontwantyoutovibecode/dade/internal/registry"
	"github.com/theydontwantyoutovibecode/dade/internal/ui"
	"github.com/spf13/cobra"
)

func init() {
	for _, alias := range backwardCompatAliases() {
		rootCmd.AddCommand(alias)
	}
}

func backwardCompatAliases() []*cobra.Command {
	hidden := func(use string, target func(cmd *cobra.Command, args []string) error) *cobra.Command {
		return &cobra.Command{
			Use:    use,
			Hidden: true,
			RunE:   target,
		}
	}

	hiddenArgs := func(use string, n int, target func(cmd *cobra.Command, args []string) error) *cobra.Command {
		c := hidden(use, target)
		c.Args = cobra.MaximumNArgs(n)
		return c
	}

	aliases := []*cobra.Command{
		hiddenArgs("templates", 0, runTemplatesCmd),
		hiddenArgs("install", 1, runInstallCmd),
		hiddenArgs("uninstall", 1, runUninstallCmd),
		hiddenArgs("update", 1, runUpdateCmd),
		hiddenArgs("list", 0, runListCmd),
		hiddenArgs("register", 1, runRegisterCmd),
		hiddenArgs("remove", 1, runRemoveCmd),
		hiddenArgs("rm", 1, runRemoveCmd),
		hidden("port", runPortCmd),
		hiddenArgs("sync", 1, runSyncCmd),
		hiddenArgs("start", 1, runStartCmd),
		hiddenArgs("stop", 1, runStopCmd),
		hidden("refresh", runRefreshCmd),
		hiddenArgs("tunnel", 1, runTunnelAttach),
		hiddenArgs("open", 1, runOpenCompat),
	}

	return aliases
}

func runOpenCompat(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	_ = console

	var projectName string

	if len(args) > 0 {
		projectName = args[0]
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			logger.Error("Failed to get current directory")
			return errors.New("open command failed")
		}
		if !registry.MarkerExists(cwd) {
			logger.Error("Not a dade project directory")
			logger.Info("Run 'dade new' or 'dade register' first")
			return errors.New("open command failed")
		}
		marker, err := registry.ReadMarker(cwd)
		if err != nil {
			logger.Error("Failed to read project marker")
			return errors.New("open command failed")
		}
		projectName = marker.Name
	}

	projectsPath, err := config.ProjectsFile()
	if err != nil {
		logger.Error("Failed to resolve projects file")
		return errors.New("open command failed")
	}

	_, ok, err := registry.Get(projectsPath, projectName)
	if err != nil {
		logger.Error("Failed to load project registry")
		return errors.New("open command failed")
	}
	if !ok {
		logger.Error(fmt.Sprintf("Project '%s' not found", projectName))
		return errors.New("open command failed")
	}

	url := fmt.Sprintf("https://%s", config.ProjectDomain(projectName))
	if err := openBrowserFunc(url); err != nil {
		logger.Error(fmt.Sprintf("Failed to open browser: %v", err))
		return errors.New("open command failed")
	}

	logger.Success(fmt.Sprintf("Opened: %s", url))
	return nil
}
