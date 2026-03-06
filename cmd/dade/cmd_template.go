package main

import (
	"errors"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:     "template <command>",
	Short:   "Manage installed templates",
	Long:    "List, add, remove, and update template plugins. Templates are stored in ~/.config/dade/templates/.",
	GroupID: "manage",
}

var templateListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List installed templates",
	Long:    "List template plugins installed in ~/.config/dade/templates/. Shows template name, description, serve type, and source URL for each entry. Use --json for machine-readable output.",
	Example: "dade template list\ndade template list --json",
	RunE:    runTemplatesCmd,
}

var templateAddCmd = &cobra.Command{
	Use:   "add <git-url>",
	Short: "Add a template from a git repository",
	Long:  "Clone a template repository and install it as a plugin. Templates must contain a dade.toml manifest file. Use --list-official to browse curated templates without installing, or --name to override the installed template name.",
	Example: `dade template add https://github.com/acme/my-template.git
dade template add --name custom-name https://github.com/acme/my-template.git
dade template add --list-official`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInstallCmd,
}

var templateRemoveCmd = &cobra.Command{
	Use:     "remove <name>",
	Short:   "Remove an installed template",
	Long:    "Remove an installed template plugin by name. Use --all to remove all installed templates with confirmation.",
	Example: "dade template remove web-site\ndade template remove --all",
	Args:    cobra.MaximumNArgs(1),
	RunE:    runUninstallCmd,
}

var templateUpdateCmd = &cobra.Command{
	Use:   "update <name>",
	Short: "Update templates from their source repositories",
	Long:  "Re-fetch a template from its original git source URL. Use --all to update all templates at once.",
	Example: `dade template update web-site
dade template update --all`,
	Args: cobra.MaximumNArgs(1),
	RunE: runUpdateCmd,
}

func init() {
	rootCmd.AddCommand(templateCmd)

	templateCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return templateListCmd.RunE(cmd, args)
	}

	templateCmd.AddCommand(templateListCmd)

	templateAddCmd.Flags().StringP("name", "n", "", "Override template name (default: from manifest)")
	templateAddCmd.Flags().Bool("list-official", false, "List official templates instead of installing")
	templateCmd.AddCommand(templateAddCmd)

	templateRemoveCmd.Flags().Bool("all", false, "Remove all installed templates with confirmation")
	templateCmd.AddCommand(templateRemoveCmd)

	templateUpdateCmd.Flags().Bool("all", false, "Update all installed templates")
	templateCmd.AddCommand(templateUpdateCmd)
}

func runTemplatesCmd(cmd *cobra.Command, _ []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	jsonOutput := output.JSON

	templatesDir, err := config.TemplatesDir()
	if err != nil {
		logger.Error("Failed to resolve templates directory")
		return errors.New("template list failed")
	}
	installed, err := loadInstalledTemplates(templatesDir)
	if err != nil {
		logger.Error("Failed to load installed templates")
		return errors.New("template list failed")
	}

	if len(installed) == 0 {
		logger.Info("No templates installed.")
		logger.Info("Add one: dade template add <git-url>")
		logger.Info("Or see official: dade template add --list-official")
		return nil
	}

	if jsonOutput {
		payload, err := templatesJSON(installed)
		if err != nil {
			logger.Error("Failed to render JSON")
			return errors.New("template list failed")
		}
		console.PrintHelp(payload)
		return nil
	}

	console.PrintHelp(templatesText(installed, output.Styled))
	return nil
}

func runInstallCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	listOfficial, err := cmd.Flags().GetBool("list-official")
	if err != nil {
		return err
	}

	if !listOfficial && len(args) == 0 {
		logger.Error("Missing template git URL")
		return errors.New("template add failed")
	}

	cmdImpl := installCommandFactory()

	argsToRun := make([]string, 0, 4)
	if listOfficial {
		argsToRun = append(argsToRun, "--list-official")
	}
	if name != "" {
		argsToRun = append(argsToRun, "--name", name)
	}
	if len(args) > 0 {
		argsToRun = append(argsToRun, args[0])
	}

	code := cmdImpl.run(cmd.Context(), argsToRun, console, logger, output.Styled)
	if code != 0 {
		return errors.New("template add failed")
	}
	return nil
}

func runUninstallCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		return err
	}

	cmdImpl := uninstallCommandFactory()
	cmdImpl.confirm = func(question string) (bool, error) {
		if output.Quiet {
			return false, nil
		}
		confirmed := false
		prompt := huh.NewConfirm().Title(question).Value(&confirmed)
		if err := prompt.Run(); err != nil {
			logger.Error("Failed to read confirmation")
			return false, err
		}
		return confirmed, nil
	}

	argsToRun := make([]string, 0, 2)
	if all {
		argsToRun = append(argsToRun, "--all")
	}
	if len(args) > 0 {
		argsToRun = append(argsToRun, args[0])
	}

	code := cmdImpl.run(cmd.Context(), argsToRun, console, logger)
	if code != 0 {
		return errors.New("template remove failed")
	}
	return nil
}

func runUpdateCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	updateAll, _ := cmd.Flags().GetBool("all")

	if !updateAll && len(args) == 0 {
		logger.Error("Missing template name")
		logger.Info("Usage: dade template update <name>")
		logger.Info("Or: dade template update --all")
		return errors.New("template update failed")
	}

	cmdImpl := updateCommandFactory()

	argsToRun := make([]string, 0, 2)
	if updateAll {
		argsToRun = append(argsToRun, "--all")
	}
	if len(args) > 0 {
		argsToRun = append(argsToRun, args[0])
	}

	code := cmdImpl.run(cmd.Context(), argsToRun, console, logger, output.Styled)
	if code != 0 {
		return errors.New("template update failed")
	}
	return nil
}
