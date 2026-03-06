---
id: jus-qvpa
status: closed
deps: [jus-55ji]
links: []
created: 2026-02-17T15:11:27Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-fxaf
tags: [cli, fang, cobra, templates-command]
---
# Migrate 'templates' command to Cobra with headless flags

Migrate the 'dade templates' command from manual arg parsing to Cobra subcommand.

## Background

dade is a CLI for scaffolding web projects. The 'templates' command lists installed template plugins from ~/.config/dade/templates/.

## Current State

Files:
- cmd/dade/templates_cmd.go: runTemplates function
- cmd/dade/templates.go: Helper functions (loadInstalledTemplates, templatesText, templatesJSON)

Current flags (manual parsing):
- --json: Output as JSON instead of styled text

No interactive prompts - this command is already headless-capable.

## Implementation

1. Create cmd/dade/cmd_templates.go with Cobra command:

```go
package main

import (
    "github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
    Use:   "templates",
    Short: "List installed templates",
    Long:  "List all template plugins installed in ~/.config/dade/templates/. Shows template name, description, serve type, and source URL.",
    RunE:  runTemplatesCmd,
}

func init() {
    rootCmd.AddCommand(templatesCmd)
    
    templatesCmd.Flags().Bool("json", false, "Output as JSON")
}
```

2. Implement runTemplatesCmd that:
   - Reads --json flag via cmd.Flags().GetBool("json")
   - Calls existing helper functions from templates.go
   - Uses console/logger from context or creates them

3. Remove old runTemplates from templates_cmd.go, keep helpers in templates.go

## Flag Documentation (shown in --help)

```
Flags:
      --json   Output as JSON
  -h, --help   help for templates
```

## Usage Examples (shown in --help)

```
Examples:
  dade templates           # List installed templates with styled output
  dade templates --json    # Output as JSON for scripting
```

## Acceptance Criteria

- [ ] dade templates --help shows all flags with descriptions
- [ ] dade templates lists installed templates
- [ ] dade templates --json outputs valid JSON array
- [ ] Empty state shows guidance message
- [ ] Existing tests pass or are updated appropriately

