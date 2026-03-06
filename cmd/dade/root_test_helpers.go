package main

import (
	"testing"

	"github.com/spf13/pflag"
)

func resetRootFlags(t *testing.T) {
	t.Helper()
	flagQuiet = false
	flagVerbose = false
	flagNoColor = false
	flagJSON = false
	skipAutoSetup = true
	flags := rootCmd.PersistentFlags()
	resetBoolFlag(t, flags, "quiet")
	resetBoolFlag(t, flags, "verbose")
	resetBoolFlag(t, flags, "no-color")
	resetBoolFlag(t, flags, "json")
}

func resetBoolFlag(t *testing.T, flags *pflag.FlagSet, name string) {
	t.Helper()
	flag := flags.Lookup(name)
	if flag == nil {
		t.Fatalf("missing flag: %s", name)
	}
	if err := flag.Value.Set("false"); err != nil {
		t.Fatalf("reset %s: %v", name, err)
	}
	flag.Changed = false
}
