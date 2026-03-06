---
id: jus-woz2
status: closed
deps: [jus-55ji]
links: []
created: 2026-02-17T15:11:39Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-fxaf
tags: [cli, fang, cobra, install-command]
---
# Migrate 'install' command to Cobra with headless flags

Migrate the 'dade install' command from manual arg parsing to Cobra subcommand.

## Background

dade is a CLI for scaffolding web projects. The 'install' command installs template plugins from git repositories into ~/.config/dade/templates/.

## Current State

File: cmd/dade/install.go

Current flags (manual parsing):
- --name NAME: Override template name (derived from manifest if not provided)
- --list-official: List official templates instead of installing
- Positional: <git-url> - Template repository URL

No interactive prompts - this command is already headless-capable.

## Implementation

1. Create cmd/dade/cmd_install.go with Cobra command:

```go
package main

import (
    "github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
    Use:   "install <git-url>",
    Short: "Install a template plugin from a git repository",
    Long:  "Clone a template repository and install it as a plugin. Templates must contain a dade.toml manifest file.",
    Args:  cobra.MaximumNArgs(1),
    RunE:  runInstallCmd,
}

func init() {
    rootCmd.AddCommand(installCmd)
    
    installCmd.Flags().StringP("name", "n", "", "Override template name (default: from manifest)")
    installCmd.Flags().Bool("list-official", false, "List official templates instead of installing")
}
```

2. Implement runInstallCmd that:
   - If --list-official, print official templates and return
   - Otherwise require git-url arg
   - Pass --name to installCommand if provided
   - Calls existing installCommand struct methods

3. Remove old arg parsing from install.go, keep business logic

## Flag Documentation (shown in --help)

```
Flags:
  -n, --name string     Override template name (default: from manifest)
      --list-official   List official templates instead of installing
  -h, --help            help for install
```

## Usage Examples (shown in --help)

```
Examples:
  dade install https://github.com/user/template.git    # Install from URL
  dade install https://github.com/user/repo.git -n my  # Install with custom name
  dade install --list-official                         # Show official templates
```

## Acceptance Criteria

- [ ] dade install --help shows all flags with descriptions
- [ ] dade install <url> installs template
- [ ] dade install <url> --name custom works
- [ ] dade install --list-official shows official templates
- [ ] dade install (no url, no --list-official) shows error
- [ ] Existing tests pass or are updated appropriately

