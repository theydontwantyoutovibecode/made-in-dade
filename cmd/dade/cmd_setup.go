package main

import (
	"errors"
	"os"
	"strings"

	"github.com/theydontwantyoutovibecode/dade/internal/logging"
	"github.com/theydontwantyoutovibecode/dade/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var setupCmd = &cobra.Command{
	Use:     "setup",
	Short:   "First-time setup for dade",
	Long:    "Initialize dade configuration, install dependencies, set up the HTTPS proxy service, and optionally install official templates. This command is safe to re-run and will skip work that is already complete. Use --check for dependency-only verification in CI.",
	Example: "dade setup\ndade setup --check\ndade setup --yes\ndade setup --no-templates",
	GroupID: "system",
	RunE:    runSetupCmd,
}

func init() {
	rootCmd.AddCommand(setupCmd)

	setupCmd.Flags().Bool("check", false, "Only check dependencies, don't run setup. Default: false")
	setupCmd.Flags().BoolP("yes", "y", false, "Answer yes to all prompts. Default: false")
	setupCmd.Flags().Bool("install-deps", false, "Install missing dependencies via Homebrew without prompting. Default: false")
	setupCmd.Flags().Bool("skip-deps", false, "Skip dependency installation prompts (fail if missing). Default: false")

	setupCmd.Flags().Bool("migrate", false, "Migrate from srv without prompting. Default: false")
	setupCmd.Flags().Bool("no-migrate", false, "Skip srv migration without prompting. Default: false")

	setupCmd.Flags().Bool("trust-ca", false, "Trust Caddy CA without prompting. Requires sudo. Default: false")
	setupCmd.Flags().Bool("no-trust-ca", false, "Skip CA trust without prompting. Default: false")

	setupCmd.Flags().Bool("install-templates", false, "Install all official templates without prompting. Default: false")
	setupCmd.Flags().Bool("no-templates", false, "Skip template installation without prompting. Default: false")
	setupCmd.Flags().StringSlice("templates", nil, "Specific templates to install (e.g., --templates web-app,web-site). Default: none")
}

func runSetupCmd(cmd *cobra.Command, _ []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	flags, err := readSetupFlags(cmd)
	if err != nil {
		return err
	}

	cmdImpl := setupCommandFactory()
	cmdImpl.confirm = newSetupConfirmer(flags, logger)
	cmdImpl.spin = nil

	args := buildSetupArgs(flags)
	code := cmdImpl.run(cmd.Context(), args, console, logger, output.Styled)
	if code != 0 {
		return errors.New("setup command failed")
	}
	return nil
}

func buildSetupArgs(flags setupFlags) []string {
	args := []string{}
	if flags.check {
		args = append(args, "--check")
	}
	return args
}

func newSetupConfirmer(flags setupFlags, logger *logging.Logger) func(string) (bool, error) {
	interactive := term.IsTerminal(int(os.Stdin.Fd()))
	return func(question string) (bool, error) {
		if answer, ok := flags.answerFor(question); ok {
			return answer, nil
		}
		if !interactive {
			return false, nil
		}
		return defaultSetupConfirm(question, logger)
	}
}

type setupFlags struct {
	check           bool
	yes             bool
	installDeps     bool
	skipDeps        bool
	migrate         bool
	noMigrate       bool
	trustCA         bool
	noTrustCA       bool
	installTemplates bool
	noTemplates     bool
	templates       []string
}

func (f setupFlags) answerFor(question string) (bool, bool) {
	switch question {
	case "Install jq via Homebrew?", "Install caddy via Homebrew?":
		if f.installDeps || f.yes {
			return true, true
		}
		if f.skipDeps {
			return false, true
		}
	case "Migrate from srv?":
		if f.migrate || f.yes {
			return true, true
		}
		if f.noMigrate {
			return false, true
		}
	case "Trust Caddy CA? (requires sudo)":
		if f.trustCA || f.yes {
			return true, true
		}
		if f.noTrustCA {
			return false, true
		}
	case "Install official templates?":
		if f.installTemplates || f.yes || len(f.templates) > 0 {
			return true, true
		}
		if f.noTemplates {
			return false, true
		}
	default:
		if strings.HasPrefix(question, "Install ") && strings.Contains(question, "? (") {
			name := parseTemplateQuestion(question)
			if name != "" && len(f.templates) > 0 {
				for _, tpl := range f.templates {
					if tpl == name {
						return true, true
					}
				}
				return false, true
			}
			if f.installTemplates || f.yes {
				return true, true
			}
			if f.noTemplates {
				return false, true
			}
		}
	}
	return false, false
}

func parseTemplateQuestion(question string) string {
	prefix := "Install "
	if !strings.HasPrefix(question, prefix) {
		return ""
	}
	trimmed := strings.TrimPrefix(question, prefix)
	parts := strings.SplitN(trimmed, "? (", 2)
	if len(parts) < 1 {
		return ""
	}
	return strings.TrimSpace(parts[0])
}

func readSetupFlags(cmd *cobra.Command) (setupFlags, error) {
	flags := setupFlags{}
	var err error
	flags.check, err = cmd.Flags().GetBool("check")
	if err != nil {
		return flags, err
	}
	flags.yes, err = cmd.Flags().GetBool("yes")
	if err != nil {
		return flags, err
	}
	flags.installDeps, err = cmd.Flags().GetBool("install-deps")
	if err != nil {
		return flags, err
	}
	flags.skipDeps, err = cmd.Flags().GetBool("skip-deps")
	if err != nil {
		return flags, err
	}
	flags.migrate, err = cmd.Flags().GetBool("migrate")
	if err != nil {
		return flags, err
	}
	flags.noMigrate, err = cmd.Flags().GetBool("no-migrate")
	if err != nil {
		return flags, err
	}
	flags.trustCA, err = cmd.Flags().GetBool("trust-ca")
	if err != nil {
		return flags, err
	}
	flags.noTrustCA, err = cmd.Flags().GetBool("no-trust-ca")
	if err != nil {
		return flags, err
	}
	flags.installTemplates, err = cmd.Flags().GetBool("install-templates")
	if err != nil {
		return flags, err
	}
	flags.noTemplates, err = cmd.Flags().GetBool("no-templates")
	if err != nil {
		return flags, err
	}
	flags.templates, err = cmd.Flags().GetStringSlice("templates")
	if err != nil {
		return flags, err
	}
	if err := validateSetupFlags(flags); err != nil {
		return flags, err
	}
	if flags.yes {
		flags.installDeps = true
		flags.migrate = true
		flags.trustCA = true
		flags.installTemplates = true
	}
	return flags, nil
}

func validateSetupFlags(flags setupFlags) error {
	if flags.installDeps && flags.skipDeps {
		return errors.New("cannot use --install-deps and --skip-deps together")
	}
	if flags.migrate && flags.noMigrate {
		return errors.New("cannot use --migrate and --no-migrate together")
	}
	if flags.trustCA && flags.noTrustCA {
		return errors.New("cannot use --trust-ca and --no-trust-ca together")
	}
	if flags.installTemplates && flags.noTemplates {
		return errors.New("cannot use --install-templates and --no-templates together")
	}
	return nil
}

func defaultSetupConfirm(question string, logger *logging.Logger) (bool, error) {
	confirmed := false
	prompt := huh.NewConfirm().Title(question).Value(&confirmed)
	if err := prompt.Run(); err != nil {
		logger.Error("Failed to read confirmation")
		return false, err
	}
	return confirmed, nil
}
