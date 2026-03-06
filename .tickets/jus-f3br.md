---
id: jus-f3br
status: closed
deps: []
links: []
created: 2026-02-12T19:54:08Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-cdr3
tags: [infrastructure, config]
---
# Add configuration directory initialization

Create the configuration directory structure and initialization logic.

## Directory Structure

```
~/.config/dade/
├── projects.json          # Project registry
├── templates/             # Installed template plugins
├── Caddyfile              # Auto-generated proxy config
├── config.toml            # User preferences (optional)
└── proxy.log              # Proxy output log
```

## Implementation

1. Add DADE_CONFIG_DIR variable (default: ~/.config/dade)
2. Create init_config() function that:
   - Creates directory structure
   - Initializes empty projects.json if missing
   - Creates templates/ directory
3. Call init_config() early in relevant commands

## Files to Create

- projects.json: `{}` (empty object)
- Caddyfile: minimal valid config

## Considerations

- Respect XDG_CONFIG_HOME if set
- Handle permission errors gracefully
- Migrate from any existing srv config (~/.config/srv/)

## Acceptance Criteria

- [ ] Config directory created on first run
- [ ] projects.json initialized as empty object
- [ ] templates/ directory created
- [ ] Existing srv config detected and offered migration

