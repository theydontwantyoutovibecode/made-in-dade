package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/theydontwantyoutovibecode/dade/internal/config"
	"github.com/theydontwantyoutovibecode/dade/internal/logging"
	"github.com/theydontwantyoutovibecode/dade/internal/registry"
	"github.com/theydontwantyoutovibecode/dade/internal/serve"
	"github.com/theydontwantyoutovibecode/dade/internal/ui"
	"github.com/spf13/cobra"
)

func runStopCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	var projectDir string
	var projectName string

	if len(args) > 0 {
		projectName = args[0]
		projectsPath, err := config.ProjectsFile()
		if err != nil {
			logger.Error("Failed to resolve projects file")
			return errors.New("stop command failed")
		}
		project, ok, err := registry.Get(projectsPath, projectName)
		if err != nil {
			logger.Error("Failed to load project registry")
			return errors.New("stop command failed")
		}
		if !ok {
			logger.Error(fmt.Sprintf("Project '%s' not found", projectName))
			return errors.New("stop command failed")
		}
		projectDir = project.Path
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			logger.Error("Failed to get current directory")
			return errors.New("stop command failed")
		}
		if !registry.MarkerExists(cwd) {
			logger.Error("Not a dade project directory")
			return errors.New("stop command failed")
		}
		projectDir = cwd
	}

	marker, err := registry.ReadMarker(projectDir)
	if err != nil {
		logger.Error("Failed to read project marker")
		return errors.New("stop command failed")
	}
	projectName = marker.Name

	pidFile := filepath.Join(projectDir, serve.DefaultPIDFile)
	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info(fmt.Sprintf("Project '%s' is not running", projectName))
			return nil
		}
		logger.Error("Failed to read PID file")
		return errors.New("stop command failed")
	}

	pid, err := strconv.Atoi(string(pidData))
	if err != nil {
		logger.Error("Invalid PID file")
		_ = os.Remove(pidFile)
		return errors.New("stop command failed")
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		logger.Info(fmt.Sprintf("Project '%s' is not running (process not found)", projectName))
		_ = os.Remove(pidFile)
		return nil
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		if err := process.Signal(syscall.SIGKILL); err != nil {
			logger.Warn(fmt.Sprintf("Failed to stop process %d", pid))
		}
	}

	_ = os.Remove(pidFile)
	_ = console
	logger.Success(fmt.Sprintf("Stopped: %s", projectName))
	return nil
}
