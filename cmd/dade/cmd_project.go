package main

import (
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:     "project <command>",
	Short:   "Manage project registry",
	Long:    "List, register, remove, and manage dade projects. Projects are tracked in ~/.config/dade/projects.json.",
	GroupID: "manage",
}

var projectListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List registered projects",
	Long:    "List all projects registered with dade, showing their status, port, template, and path. Use --running to filter to running projects only, or --json for machine-readable output.",
	Example: "dade project list\ndade project list --running\ndade project list --json",
	RunE:    runListCmd,
}

var projectRegisterCmd = &cobra.Command{
	Use:   "register [name]",
	Short: "Register existing directory with dade",
	Long:  "Register an existing project directory with dade for serving and management. Automatically detects project type, assigns a port, creates .dade marker, and updates proxy configuration.",
	Example: `dade project register
dade project register myapp
dade project register -t static`,
	Args: cobra.MaximumNArgs(1),
	RunE: runRegisterCmd,
}

var projectRemoveCmd = &cobra.Command{
	Use:     "remove [name]",
	Aliases: []string{"rm"},
	Short:   "Remove a project from registry",
	Long:    "Unregister a project from dade management. This removes the project from the registry and deletes the .dade marker file, but does not delete project files unless --files is specified.",
	Example: `dade project remove myapp
dade project rm myapp
dade project remove myapp --files`,
	Args: cobra.MaximumNArgs(1),
	RunE: runRemoveCmd,
}

var projectPortCmd = &cobra.Command{
	Use:   "port",
	Short: "Show or update project port",
	Long:  "Display or update the port for the current dade project. When updating, the registry and Caddy configuration are updated automatically.",
	Example: `dade project port
dade project port --set 8001`,
	Args: cobra.NoArgs,
	RunE: runPortCmd,
}

var projectSyncCmd = &cobra.Command{
	Use:   "sync [path]",
	Short: "Rebuild project registry by scanning for .dade files",
	Long:  "Scan directories for .dade marker files and rebuild the project registry. Useful when projects have been moved or the registry is corrupted. Use --clean to remove stale entries instead of scanning.",
	Example: `dade project sync
dade project sync ~/Code
dade project sync --clean`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSyncCmd,
}

var projectStartCmd = &cobra.Command{
	Use:   "start [name]",
	Short: "Start production server",
	Long:  "Start a production server for a dade project. Runs the serve.prod command from the template manifest. For development, use 'dade dev' instead.",
	Example: `dade project start
dade project start myapp
dade project start --background`,
	Args: cobra.MaximumNArgs(1),
	RunE: runStartCmd,
}

var projectStopCmd = &cobra.Command{
	Use:   "stop [name]",
	Short: "Stop a running project",
	Long:  "Stop a running dade project by terminating its server process. Without arguments, stops the project in the current directory.",
	Example: `dade project stop
dade project stop myapp`,
	Args: cobra.MaximumNArgs(1),
	RunE: runStopCmd,
}

func init() {
	rootCmd.AddCommand(projectCmd)

	projectCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return projectListCmd.RunE(cmd, args)
	}

	projectListCmd.Flags().Bool("running", false, "Show only running projects")
	projectCmd.AddCommand(projectListCmd)

	projectRegisterCmd.Flags().StringP("template", "t", "", "Specify template type (static, django, node)")
	projectCmd.AddCommand(projectRegisterCmd)

	projectRemoveCmd.Flags().Bool("files", false, "Also delete project files (dangerous!)")
	projectRemoveCmd.Flags().BoolP("yes", "y", false, "Skip confirmation for --files")
	projectCmd.AddCommand(projectRemoveCmd)

	projectPortCmd.Flags().Int("set", 0, "Set the port to use")
	projectCmd.AddCommand(projectPortCmd)

	projectSyncCmd.Flags().Bool("clean", false, "Remove entries for missing directories instead of scanning")
	projectCmd.AddCommand(projectSyncCmd)

	projectStartCmd.Flags().IntP("port", "p", 0, "Override port")
	projectStartCmd.Flags().Bool("background", false, "Run in background (detach)")
	projectCmd.AddCommand(projectStartCmd)

	projectCmd.AddCommand(projectStopCmd)
}
