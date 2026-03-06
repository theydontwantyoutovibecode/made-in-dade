---
id: jus-njz0
status: closed
deps: [jus-w9pp]
links: []
created: 2026-02-17T15:13:34Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-fxaf
tags: [cli, fang, cobra, documentation]
---
# Add comprehensive help documentation to all commands

Ensure all commands have comprehensive help documentation including long descriptions, examples, and flag descriptions.

## Background

dade is a CLI for scaffolding web projects. With Cobra/Fang, help is auto-generated from command metadata. This ticket ensures all help text is thorough and useful.

## Current State

After Cobra migration, each command has:
- Use: command syntax
- Short: one-line description
- Long: detailed description (may be minimal)
- Example: usage examples (may be missing)
- Flags with descriptions

## Implementation

1. For each command, ensure Long description covers:
   - What the command does
   - When to use it
   - Key behaviors and defaults
   - Related commands

2. For each command, add Example field with 3-5 examples:
   - Basic usage
   - Common flags
   - Advanced usage
   - CI/scripting usage

3. For each flag, ensure description covers:
   - What it does
   - Default value (if not obvious)
   - Valid values (if constrained)
   - When to use it

4. Review Fang's styled output to ensure formatting looks good

## Commands to Document

### new
Long: Detailed explanation of project creation, template system, setup.sh execution
Examples: basic, with template, with local, headless CI

### templates  
Long: Where templates are stored, what info is shown
Examples: basic, JSON for scripting

### install
Long: How template installation works, manifest requirements
Examples: from URL, with custom name, list official

### proxy
Long: What the proxy does, how it works with projects
Examples: start/stop/restart/status/logs, JSON status

### setup
Long: What setup does, order of operations, idempotency
Examples: interactive, headless yes-to-all, check only

## Flag Description Standards

Format: "<action>. <default if any>. <notes if any>"

Examples:
- "Specify template name. Default: django-hypermedia"
- "Output in JSON format for scripting"
- "Skip all prompts with 'no' answer"
- "Trust Caddy CA certificate. Requires sudo"

## Acceptance Criteria

- [ ] All commands have Long description (3+ sentences)
- [ ] All commands have Example field (3+ examples)
- [ ] All flags have clear descriptions
- [ ] dade --help output is useful and complete
- [ ] dade <cmd> --help output is useful and complete
- [ ] Help renders well in terminal (test with Fang styling)

