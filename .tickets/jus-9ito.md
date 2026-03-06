---
id: jus-9ito
status: closed
deps: [jus-jjl8]
links: []
created: 2026-02-12T20:02:18Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-k78p
tags: [documentation, ai]
---
# Add dade AGENTS.md

Create AGENTS.md for AI agents working on or with dade.

## File Location

/Users/acabrera/Code/dade/dade/AGENTS.md

## Content

```markdown
# dade - AI Agent Guidelines

## Project Overview

dade is a bash CLI for scaffolding and serving web projects.
It combines project templating with local HTTPS development servers.

## Architecture

### Core Components

1. **Template Plugin System**
   - Templates installed to ~/.config/dade/templates/
   - Each template has dade.toml manifest
   - Scaffolding copies template to new project

2. **Project Registry**
   - ~/.config/dade/projects.json tracks projects
   - Each project has .dade marker file
   - Port assignment managed centrally

3. **Caddy Proxy**
   - Central reverse proxy for HTTPS
   - Routes *.localhost to project ports
   - Runs as launchd service

### Key Files

| File | Purpose |
|------|---------|
| dade | Main CLI script (~1500 lines bash) |
| ~/.config/dade/projects.json | Project registry |
| ~/.config/dade/templates/ | Installed templates |
| ~/.config/dade/Caddyfile | Generated proxy config |
| .dade | Per-project marker (JSON) |
| dade.toml | Template manifest (TOML) |

## Working with dade

### Adding Commands

Commands follow pattern:
\`\`\`bash
cmd_<name>() {
    # Parse arguments
    # Do work
    # Log result
}
\`\`\`

Add to main() case statement.

### Dependencies

- bash 4.0+ (associative arrays)
- jq (JSON processing)
- caddy (web server)
- gum (optional, UI)

### Testing Changes

\`\`\`bash
# Test directly
./dade <command>

# Test with new project
./dade new testproj --template hypertext
cd testproj
./dade start
./dade stop
\`\`\`

## Creating Templates

Templates must include dade.toml:

\`\`\`toml
[template]
name = "my-template"
description = "Description here"

[serve]
type = "static"  # or "command"
\`\`\`

For command-based templates:
\`\`\`toml
[serve]
type = "command"
dev = "./start.sh --dev"
port_env = "PORT"
\`\`\`

## Common Tasks

### Add new serve type
1. Update start_server() in dade
2. Add case for new type
3. Document in README

### Add new command
1. Create cmd_<name>() function
2. Add to main() case statement
3. Add to cmd_help() output
4. Update shell completions

### Debug proxy issues
\`\`\`bash
dade proxy status
cat ~/.config/dade/proxy.log
caddy validate --config ~/.config/dade/Caddyfile
\`\`\`
\`\`\`

## Acceptance Criteria

- [ ] AGENTS.md created in dade repo
- [ ] Architecture documented
- [ ] Key files listed
- [ ] Common tasks explained
- [ ] Template creation documented

