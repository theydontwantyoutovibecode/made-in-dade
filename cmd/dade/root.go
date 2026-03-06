package main

import (
	"context"
	"fmt"
	"os"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/ui"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/version"
	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

var (
	flagQuiet     bool
	flagVerbose   bool
	flagNoColor   bool
	flagJSON      bool
	skipAutoSetup bool
)

var rootCmd = &cobra.Command{
	Use:                "dade",
	Short:              "CLI for scaffolding web application projects",
	Long:               "dade scaffolds web projects from curated templates and manages a local HTTPS proxy. Use it to create new projects, install template plugins, and run setup tasks. Global flags let you control output formatting for scripting and CI workflows.",
	Example:            "dade new myapp\ndade template list --json\ndade setup --check\ndade --help",
	Version:            version.Version,
	CompletionOptions:  cobra.CompletionOptions{DisableDefaultCmd: true},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		name := cmd.Name()
		if name == "setup" || name == "version" || name == "help" {
			return nil
		}
		if skipAutoSetup || config.IsConfigured() {
			return nil
		}
		output := getOutputSettings(cmd)
		console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
		logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
		logger.SetSilent(output.Quiet)
		logger.Info("First run detected, running setup...")
		fmt.Fprintln(cmd.OutOrStdout())
		impl := defaultSetupCommand()
		code := impl.run(cmd.Context(), nil, console, logger, output.Styled)
		if code != 0 {
			return fmt.Errorf("automatic setup failed, run 'dade setup' manually")
		}
		fmt.Fprintln(cmd.OutOrStdout())
		return nil
	},
}

func init() {
	rootCmd.AddGroup(
		&cobra.Group{ID: "dev", Title: "Development"},
		&cobra.Group{ID: "manage", Title: "Management"},
		&cobra.Group{ID: "system", Title: "System"},
	)
	rootCmd.SetHelpCommandGroupID("system")

	rootCmd.PersistentFlags().BoolVarP(&flagQuiet, "quiet", "q", false, "Suppress non-essential output. Default: false")
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "Enable verbose output. Default: false")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "Disable colored output. Default: false")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output in JSON format where supported. Default: false")
	rootCmd.MarkFlagsMutuallyExclusive("quiet", "verbose")
}

func Execute() {
	if err := fang.Execute(context.Background(), rootCmd, fang.WithVersion(version.Version)); err != nil {
		os.Exit(1)
	}
}
