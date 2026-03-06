---
id: jus-jw2r
status: closed
deps: [jus-55ji]
links: []
created: 2026-02-17T15:12:23Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-fxaf
tags: [cli, fang, cobra, setup-command]
---
# Migrate 'setup' command to Cobra with headless flags

Migrate the 'dade setup' command from manual arg parsing to Cobra with comprehensive flags for headless operation.

## Background

dade is a CLI for scaffolding web projects. The 'setup' command performs first-time setup: dependency checks, config initialization, srv migration, Caddyfile generation, launchd proxy installation, CA trust, and official template installation.

## Current State

File: cmd/dade/setup.go

Current flags (manual parsing):
- --check: Only check dependencies, don't run setup

Interactive prompts (huh.NewConfirm) requiring headless alternatives:
1. "Install X via Homebrew?" - for missing jq, caddy (lines 278-317)
2. "Migrate from srv?" - when existing srv installation detected (lines 176-180)
3. "Trust Caddy CA? (requires sudo)" - CA trust confirmation (lines 235-246)
4. "Install official templates?" - template installation prompt (lines 248-258)
5. "Install <name>? (<description>)" - per-template confirmation (line 335)

## Implementation

1. Create cmd/dade/cmd_setup.go with Cobra command:

```go
package main

import (
    "github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
    Use:   "setup",
    Short: "First-time setup for dade",
    Long:  "Initialize dade configuration, install dependencies, set up the HTTPS proxy service, and optionally install official templates.",
    RunE:  runSetupCmd,
}

func init() {
    rootCmd.AddCommand(setupCmd)
    
    // Dependency handling
    setupCmd.Flags().Bool("check", false, "Only check dependencies, don't run setup")
    setupCmd.Flags().BoolP("yes", "y", false, "Answer yes to all prompts (auto-install deps, migrate, trust CA, install templates)")
    setupCmd.Flags().Bool("install-deps", false, "Install missing dependencies via Homebrew without prompting")
    setupCmd.Flags().Bool("skip-deps", false, "Skip dependency installation prompts (fail if missing)")
    
    // Migration handling
    setupCmd.Flags().Bool("migrate", false, "Migrate from srv without prompting")
    setupCmd.Flags().Bool("no-migrate", false, "Skip srv migration without prompting")
    
    // CA trust handling
    setupCmd.Flags().Bool("trust-ca", false, "Trust Caddy CA without prompting (requires sudo)")
    setupCmd.Flags().Bool("no-trust-ca", false, "Skip CA trust without prompting")
    
    // Template handling
    setupCmd.Flags().Bool("install-templates", false, "Install all official templates without prompting")
    setupCmd.Flags().Bool("no-templates", false, "Skip template installation without prompting")
    setupCmd.Flags().StringSlice("templates", nil, "Specific templates to install (e.g., --templates django-hypermedia,hypertext)")
}
```

2. Implement runSetupCmd that:
   - Checks for conflicting flags (e.g., --migrate and --no-migrate)
   - If --yes, sets all affirmative flags
   - Passes flag values to setupCommand via modified confirm function
   - Falls back to interactive prompts only if TTY and no flag provided

3. Modify setupCommand.confirm to accept flag overrides

## Flag Documentation (shown in --help)

```
Flags:
      --check              Only check dependencies, don't run setup
  -y, --yes                Answer yes to all prompts
      --install-deps       Install missing dependencies via Homebrew
      --skip-deps          Skip dependency installation (fail if missing)
      --migrate            Migrate from srv without prompting
      --no-migrate         Skip srv migration without prompting
      --trust-ca           Trust Caddy CA without prompting (requires sudo)
      --no-trust-ca        Skip CA trust without prompting
      --install-templates  Install all official templates
      --no-templates       Skip template installation
      --templates strings  Specific templates to install (comma-separated)
  -h, --help               help for setup
```

## Usage Examples (shown in --help)

```
Examples:
  dade setup                           # Interactive setup
  dade setup --check                   # Check dependencies only
  dade setup -y                        # Non-interactive, yes to all
  dade setup --skip-deps --no-migrate  # Headless, skip optional steps
  dade setup --templates django-hypermedia  # Install specific template
```

## Acceptance Criteria

- [ ] dade setup --help shows all flags with descriptions
- [ ] dade setup --check works as before
- [ ] dade setup -y runs without any prompts
- [ ] dade setup --install-deps installs missing deps
- [ ] dade setup --skip-deps fails if deps missing (no prompt)
- [ ] dade setup --migrate migrates without prompting
- [ ] dade setup --no-migrate skips migration without prompting
- [ ] dade setup --trust-ca trusts CA without prompting
- [ ] dade setup --no-trust-ca skips CA trust
- [ ] dade setup --install-templates installs all official
- [ ] dade setup --no-templates skips template installation
- [ ] dade setup --templates X installs only specified templates
- [ ] Conflicting flags (e.g., --migrate --no-migrate) produce error
- [ ] Existing tests pass or are updated appropriately

