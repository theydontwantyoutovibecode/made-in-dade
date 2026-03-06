---
id: jus-55ji
status: closed
deps: [jus-xsvy]
links: []
created: 2026-02-17T15:10:46Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-fxaf
tags: [cli, fang, cobra]
---
# Create root Cobra command with Fang execution

Create the root Cobra command structure and integrate with Fang for styled help output, version handling, and shell completions.

## Background

dade is a CLI for scaffolding web projects. We are migrating from manual arg parsing to Cobra/Fang. Fang wraps Cobra to provide styled help pages, automatic --version, manpage generation, and shell completions.

## Current State

- Entry point: cmd/dade/main.go
- Current main() calls run(os.Args[1:]) which uses a switch statement for command dispatch
- Version is in internal/version/version.go (version.Version constant)
- Commands: new, templates, install, proxy, setup
- Global flags: --help/-h, --version/-v

## Implementation

1. Create new file: cmd/dade/root.go

2. Define root command:
```go
package main

import (
    "context"
    "os"
    
    "github.com/theydontwantyoutovibecode/dade/internal/version"
    "github.com/charmbracelet/fang"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:     "dade",
    Short:   "CLI for scaffolding web application projects",
    Long:    "dade is a command-line interface for creating new projects from curated templates, managing a local HTTPS proxy, and serving projects in development.",
    Version: version.Version,
}

func Execute() {
    if err := fang.Execute(context.Background(), rootCmd); err != nil {
        os.Exit(1)
    }
}
```

3. Update main.go to call Execute() instead of run()

4. Keep existing runWithIO function available for tests but mark for future removal

## File Structure After Change

- cmd/dade/root.go - Root command definition
- cmd/dade/main.go - Just calls Execute()
- cmd/dade/new.go - Will be updated in subsequent ticket
- (other command files unchanged initially)

## Acceptance Criteria

- [ ] cmd/dade/root.go exists with rootCmd definition
- [ ] main() function calls Execute()
- [ ] dade --help shows styled Fang help output
- [ ] dade --version shows version from internal/version
- [ ] dade (no args) shows help
- [ ] dade completion bash/zsh/fish/powershell generates completions
- [ ] Existing tests still pass (runWithIO preserved for now)

