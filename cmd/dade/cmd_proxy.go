package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/theydontwantyoutovibecode/dade/internal/logging"
	"github.com/theydontwantyoutovibecode/dade/internal/ui"
	"github.com/spf13/cobra"
)


var proxyCmd = &cobra.Command{
	Use:     "proxy",
	Short:   "Manage the local HTTPS proxy service",
	Long:    "Control the Caddy-based HTTPS proxy that provides local .localhost domains for your projects. Use subcommands to start, stop, restart, inspect status, or tail logs. Status supports JSON output for scripting via the global --json flag.",
	Example: "dade proxy start\ndade proxy status\ndade proxy status --json\ndade proxy logs --lines 200",
	GroupID: "manage",
}

var proxyStartCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start the proxy service",
	Long:    "Start the local HTTPS proxy service. If the proxy is already running, the command exits successfully without making changes.",
	Example: "dade proxy start",
	RunE:    runProxyStartCmd,
}

var proxyStopCmd = &cobra.Command{
	Use:     "stop",
	Short:   "Stop the proxy service",
	Long:    "Stop the local HTTPS proxy service. If the proxy is not running, the command exits successfully.",
	Example: "dade proxy stop",
	RunE:    runProxyStopCmd,
}

var proxyRestartCmd = &cobra.Command{
	Use:     "restart",
	Short:   "Restart the proxy service",
	Long:    "Restart the local HTTPS proxy service, reloading the configuration and registry.",
	Example: "dade proxy restart",
	RunE:    runProxyRestartCmd,
}

var proxyStatusCmd = &cobra.Command{
	Use:     "status",
	Short:   "Show proxy status",
	Long:    "Display whether the proxy is running, number of registered projects, and port range. Use --json for scripting-friendly output.",
	Example: "dade proxy status\ndade proxy status --json",
	RunE:    runProxyStatusCmd,
}

var proxyLogsCmd = &cobra.Command{
	Use:     "logs",
	Short:   "Tail proxy logs",
	Long:    "Stream the proxy service logs in real-time or show a fixed number of lines. Use --lines to print a finite set and --follow=false to stop streaming.",
	Example: "dade proxy logs\ndade proxy logs --lines 200\ndade proxy logs --follow=false",
	RunE:    runProxyLogsCmd,
}

var proxyReloadCmd = &cobra.Command{
	Use:     "reload",
	Short:   "Regenerate proxy configuration and reload",
	Long:    "Regenerate the Caddyfile from the project registry and reload the proxy. Use --list to show all project URLs after reloading.",
	Example: "dade proxy reload\ndade proxy reload --list",
	Args:    cobra.NoArgs,
	RunE:    runRefreshCmd,
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	proxyCmd.AddCommand(proxyStartCmd)
	proxyCmd.AddCommand(proxyStopCmd)
	proxyCmd.AddCommand(proxyRestartCmd)
	proxyCmd.AddCommand(proxyStatusCmd)
	proxyCmd.AddCommand(proxyLogsCmd)

	proxyReloadCmd.Flags().Bool("list", false, "List all project URLs after reload")
	proxyCmd.AddCommand(proxyReloadCmd)

	proxyLogsCmd.Flags().IntP("lines", "n", 0, "Number of lines to show. Default: 0 (follow)")
	proxyLogsCmd.Flags().BoolP("follow", "f", true, "Follow log output. Default: true")
}

func runProxyStartCmd(cmd *cobra.Command, _ []string) error {
	console, logger, _ := proxyIO(cmd)
	_ = console

	cmdImpl := proxyCommandFactory()
	code := cmdImpl.start(cmd.Context(), logger)
	if code != 0 {
		return errors.New("proxy start failed")
	}
	return nil
}

func runProxyStopCmd(cmd *cobra.Command, _ []string) error {
	console, logger, _ := proxyIO(cmd)
	_ = console

	cmdImpl := proxyCommandFactory()
	code := cmdImpl.stop(cmd.Context(), logger)
	if code != 0 {
		return errors.New("proxy stop failed")
	}
	return nil
}

func runProxyRestartCmd(cmd *cobra.Command, _ []string) error {
	console, logger, _ := proxyIO(cmd)
	_ = console

	cmdImpl := proxyCommandFactory()
	code := cmdImpl.restartService(cmd.Context(), logger)
	if code != 0 {
		return errors.New("proxy restart failed")
	}
	return nil
}

func runProxyStatusCmd(cmd *cobra.Command, _ []string) error {
	console, logger, styled := proxyIO(cmd)
	_ = console

	output := getOutputSettings(cmd)
	jsonOutput := output.JSON

	cmdImpl := proxyCommandFactory()
	if !jsonOutput {
		code := cmdImpl.status(cmd.Context(), logger)
		if code != 0 {
			return errors.New("proxy status failed")
		}
		return nil
	}

	running := cmdImpl.isRunning(cmd.Context(), cmdImpl.runner)
	projectsPath, err := cmdImpl.projectsFile()
	if err != nil {
		logger.Error("Failed to resolve projects registry")
		return errors.New("proxy status failed")
	}
	projects, err := cmdImpl.loadRegistry(projectsPath)
	if err != nil {
		logger.Error("Failed to load project registry")
		return errors.New("proxy status failed")
	}
	count, minPort, maxPort, ok := projectStats(projects)
	caddyfilePath, err := cmdImpl.caddyfilePath()
	if err != nil {
		logger.Error("Failed to resolve Caddyfile path")
		return errors.New("proxy status failed")
	}

	payload := map[string]interface{}{
		"running":        running,
		"projects":       count,
		"ports":          nil,
		"caddyfile_path": caddyfilePath,
	}
	if ok {
		payload["ports"] = map[string]int{"min": minPort, "max": maxPort}
	}

	encoded, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		logger.Error("Failed to render JSON")
		return errors.New("proxy status failed")
	}

	ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), styled).PrintHelp(string(encoded))
	return nil
}

func runProxyLogsCmd(cmd *cobra.Command, _ []string) error {
	console, logger, styled := proxyIO(cmd)
	_ = console
	_ = styled

	lines, err := cmd.Flags().GetInt("lines")
	if err != nil {
		return err
	}
	follow, err := cmd.Flags().GetBool("follow")
	if err != nil {
		return err
	}

	cmdImpl := proxyCommandFactory()
	logPath, err := cmdImpl.logPath()
	if err != nil {
		logger.Error("Failed to resolve proxy log path")
		return errors.New("proxy logs failed")
	}

	if cmdImpl.tail == nil {
		logger.Error("Log streaming unavailable")
		return errors.New("proxy logs failed")
	}

	if lines > 0 || !follow {
		return tailWithOptions(logPath, lines, follow)
	}

	code := cmdImpl.logs(logger)
	if code != 0 {
		return errors.New("proxy logs failed")
	}
	return nil
}

func proxyIO(cmd *cobra.Command) (*ui.UI, *logging.Logger, bool) {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)
	return console, logger, output.Styled
}

func tailWithOptions(path string, lines int, follow bool) error {
	if path == "" {
		return errors.New("path required")
	}
	args := []string{}
	if lines > 0 {
		args = append(args, "-n", fmt.Sprintf("%d", lines))
	}
	if follow {
		args = append(args, "-f")
	}
	args = append(args, path)
	cmd := exec.Command("tail", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
