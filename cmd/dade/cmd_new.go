package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/theydontwantyoutovibecode/dade/internal/config"
	"github.com/theydontwantyoutovibecode/dade/internal/logging"
	"github.com/theydontwantyoutovibecode/dade/internal/ui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var newCmd = &cobra.Command{
	Use:     "new [name]",
	Short:   "Create a new project from a template",
	Long:    "Create a new project directory from a curated template. Templates are cloned from git repositories and initialized with a fresh git repo. If you omit the name, you can pass it via --name or be prompted when running interactively. For headless usage, provide --name or a positional name along with any flags.",
	Example: "dade new myapp\ndade new --template web-app myapp\ndade new --local ./templates/web-app --name myapp\ndade --json templates",
	GroupID: "dev",
	Args:    cobra.MaximumNArgs(1),
	RunE:    runNewCmd,
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringP("template", "t", "web-app", "Specify template name. Default: web-app")
	newCmd.Flags().String("local", "", "Use local template directory instead of cloning. Default: empty")
	newCmd.Flags().StringP("name", "n", "", "Project name (alternative to positional arg). Default: empty")
	newCmd.Flags().Bool("inspect", false, "Show details about a template without creating a project")
}

func runNewCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	if name == "" && len(args) > 0 {
		name = args[0]
	}

	templateName, err := cmd.Flags().GetString("template")
	if err != nil {
		return err
	}
	templateChanged := cmd.Flags().Changed("template")

	localPath, err := cmd.Flags().GetString("local")
	if err != nil {
		return err
	}

	inspect, err := cmd.Flags().GetBool("inspect")
	if err != nil {
		return err
	}

	if inspect {
		return runInspectTemplate(cmd, templateName, logger)
	}

	newArgs := make([]string, 0, 6)
	if name != "" {
		newArgs = append(newArgs, name)
	}
	if localPath != "" {
		newArgs = append(newArgs, "--local", localPath)
	}
	if templateChanged {
		newArgs = append(newArgs, "--template", templateName)
	}

	cmdImpl := newCommandFactory()
	interactive := term.IsTerminal(int(os.Stdin.Fd()))
	code := cmdImpl.run(context.Background(), newArgs, console, logger, interactive)
	if code != 0 {
		return errors.New("new command failed")
	}
	return nil
}

func runInspectTemplate(cmd *cobra.Command, templateName string, logger *logging.Logger) error {
	templatesDir, err := config.TemplatesDir()
	if err != nil {
		logger.Error("Failed to resolve templates directory")
		return errors.New("inspect failed")
	}
	installed, err := loadPluginTemplates(templatesDir, os.ReadDir, os.ReadFile)
	if err != nil {
		logger.Error("Failed to load templates")
		return errors.New("inspect failed")
	}
	match, ok := findPluginTemplate(installed, templateName)
	if !ok {
		logger.Error(fmt.Sprintf("Template '%s' not found", templateName))
		return errors.New("inspect failed")
	}

	tpl := match.Manifest.Template
	fmt.Fprintf(cmd.OutOrStdout(), "Name:        %s\n", tpl.Name)
	fmt.Fprintf(cmd.OutOrStdout(), "Description: %s\n", tpl.Description)
	if tpl.Version != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Version:     %s\n", tpl.Version)
	}
	if tpl.Author != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Author:      %s\n", tpl.Author)
	}
	if tpl.URL != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "URL:         %s\n", tpl.URL)
	}
	if len(tpl.Aliases) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "Aliases:     %s\n", strings.Join(tpl.Aliases, ", "))
	}

	scaffold := match.Manifest.Scaffold
	if scaffold.Setup != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Setup:       %s\n", scaffold.Setup)
	}

	serve := match.Manifest.Serve
	fmt.Fprintf(cmd.OutOrStdout(), "Serve Type:  %s\n", serve.Type)
	if serve.Dev != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Dev Command: %s\n", serve.Dev)
	}

	if isDefaultTemplate(match.Path) {
		fmt.Fprintln(cmd.OutOrStdout(), "Source:      default (bundled)")
	} else {
		sourcePath := filepath.Join(match.Path, ".source")
		if data, err := os.ReadFile(sourcePath); err == nil {
			fmt.Fprintf(cmd.OutOrStdout(), "Source:      %s\n", strings.TrimSpace(string(data)))
		}
	}

	return nil
}
