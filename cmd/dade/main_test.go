package main

import (
	"bytes"
	"strings"
	"testing"

)

func TestMainHelp(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"--help"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected help to succeed: %v", err)
	}
	if !strings.Contains(stdout.String(), "Usage") {
		t.Fatalf("expected help output")
	}
}

func TestMainUnknownCommand(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"nope"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatalf("expected unknown command error")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "unknown command") {
		t.Fatalf("unexpected error: %v", err)
	}
}
