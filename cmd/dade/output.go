package main

import (
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type OutputSettings struct {
	Quiet   bool
	Verbose bool
	Styled  bool
	JSON    bool
}

func getOutputSettings(cmd *cobra.Command) OutputSettings {
	noColor, _ := cmd.Flags().GetBool("no-color")
	jsonOut, _ := cmd.Flags().GetBool("json")

	styled := term.IsTerminal(int(os.Stdout.Fd())) && !noColor

	return OutputSettings{
		Quiet:   flagQuiet,
		Verbose: flagVerbose,
		Styled:  styled,
		JSON:    jsonOut,
	}
}
