---
id: jus-sibp
status: closed
deps: [jus-55ji]
links: []
created: 2026-02-17T15:11:10Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-fxaf
tags: [cli, fang, cobra, new-command]
---
# Migrate 'new' command to Cobra with headless flags

Migrate the 'dade new' command from manual arg parsing to Cobra, adding flags that enable fully headless operation without interactive prompts.

## Background

dade is a CLI for scaffolding web projects. The 'new' command creates a new project from a template. Currently it has interactive prompts that block headless/CI usage.

## Current State

File: cmd/dade/new.go

Current flags (manual parsing):
- --local PATH: Use local template directory
- --template NAME: Specify template (default: django-hypermedia)
- Positional: [name] - Project name

Interactive prompts requiring headless alternatives:
1. Template selection (huh.NewSelect) - triggered when multiple templates available and no --template flag
2. Project name input (huh.NewInput) - triggered when no positional arg and interactive TTY

Current behavior:
- If no project name and non-TTY: errors with "Project name is required"
- If no template specified and multiple available: prompts for selection (fails in non-TTY)

## Implementation

1. Create cmd/dade/cmd_new.go with Cobra command:

```go
package main

import (
    "github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
    Use:   "new [name]",
    Short: "Create a new project from a template",
    Long:  "Create a new project directory from a curated template. Templates are cloned from git repositories and initialized with a fresh git repo.",
    Args:  cobra.MaximumNArgs(1),
    RunE:  runNewCmd,
}

func init() {
    rootCmd.AddCommand(newCmd)
    
    newCmd.Flags().StringP("template", "t", "django-hypermedia", "Template name to use")
    newCmd.Flags().String("local", "", "Use local template directory instead of cloning")
    newCmd.Flags().StringP("name", "n", "", "Project name (alternative to positional arg)")
}
```

2. Implement runNewCmd that:
   - Reads flags via cmd.Flags()
   - Accepts name as positional arg OR --name flag
   - If --template provided, skips interactive selection
   - If name provided (either way), skips interactive prompt
   - Only falls back to interactive if TTY detected AND required value missing
   - Calls existing newCommand struct methods for actual work

3. Remove old arg parsing from new.go, keep business logic

## Flag Documentation (shown in --help)

```
Flags:
  -t, --template string   Template name to use (default "django-hypermedia")
      --local string      Use local template directory instead of cloning
  -n, --name string       Project name (alternative to positional arg)
  -h, --help              help for new
```

## Usage Examples (shown in --help)

```
Examples:
  dade new myproject                    # Create project with default template
  dade new myproject -t hypertext       # Create project with specific template
  dade new -n myproject -t hypertext    # Same using flags only
  dade new --local ./my-template myproj # Use local template directory
```

## Acceptance Criteria

- [ ] dade new --help shows all flags with descriptions
- [ ] dade new myproject creates project without prompts
- [ ] dade new myproject --template hypertext creates project with specified template
- [ ] dade new --name myproject works as alternative to positional
- [ ] dade new (no args, non-TTY) errors with clear message
- [ ] dade new (no args, TTY) falls back to interactive prompts
- [ ] --local flag works as before
- [ ] Existing tests pass or are updated appropriately

