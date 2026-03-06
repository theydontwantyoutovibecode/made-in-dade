---
id: jus-ri29
status: open
deps: []
links: []
created: 2026-03-02T04:52:35Z
type: task
priority: 1
assignee: Alex Cabrera
tags: [cli, restructure]
---
# Group template commands under 'dade template' parent

Create a 'template' parent command with subcommands: list, add, remove, update. This replaces four top-level commands (templates, install, uninstall, update) with one parent.

Current → Proposed:
  dade templates              → dade template list
  dade install <url>          → dade template add <url>
  dade uninstall <name>       → dade template remove <name>
  dade update <name>          → dade template update <name>

The verb-noun pattern reads naturally: 'template list', 'template add django', 'template remove old-one'.

The parent command 'dade template' with no subcommand should default to 'template list' for convenience.

Implementation:
- Create cmd_template.go with the parent command
- Move install.go logic into a 'template add' subcommand
- Move uninstall.go logic into a 'template remove' subcommand  
- Move update.go logic into a 'template update' subcommand
- Move templates.go logic into a 'template list' subcommand
- Remove old cmd_install.go, cmd_uninstall.go, cmd_update.go, cmd_templates.go
- Update tests accordingly

