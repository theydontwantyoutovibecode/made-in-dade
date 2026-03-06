package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/theydontwantyoutovibecode/dade/internal/version"
)

func TestRootCommandHelp(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"--help"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected help to succeed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "dade") {
		t.Fatalf("expected help output to mention dade")
	}
}

func TestRootCommandVersion(t *testing.T) {
	if rootCmd.Version != version.Version {
		t.Fatalf("expected root version %s, got %s", version.Version, rootCmd.Version)
	}
}

func TestRootCommandHelpShowsGlobalFlags(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"--help"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected help to succeed: %v", err)
	}

	output := stdout.String()
	for _, flag := range []string{"--json", "--no-color", "--quiet", "--verbose"} {
		if !strings.Contains(output, flag) {
			t.Fatalf("missing %s in help output", flag)
		}
	}
	if !strings.Contains(output, "Examples:") {
		t.Fatalf("expected examples section")
	}
}

func TestRootCommandMutuallyExclusiveFlags(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"--quiet", "--verbose", "template", "list"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatalf("expected mutually exclusive flags to error")
	}
	if !strings.Contains(err.Error(), "quiet") || !strings.Contains(err.Error(), "verbose") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewHelpShowsGlobalFlags(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"new", "--help"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected help to succeed: %v", err)
	}

	output := stdout.String()
	for _, flag := range []string{"--json", "--no-color", "--quiet", "--verbose"} {
		if !strings.Contains(output, flag) {
			t.Fatalf("missing %s in help output", flag)
		}
	}
}

func TestNewHelpShowsExamples(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	resetRootFlags(t)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{"new", "--help"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected help to succeed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "Examples:") {
		t.Fatalf("expected examples section")
	}
	if !strings.Contains(output, "dade new") {
		t.Fatalf("expected example content")
	}
}
