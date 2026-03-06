---
id: jus-xw6m
status: open
deps: []
links: []
created: 2026-03-02T04:53:18Z
type: task
priority: 3
assignee: Alex Cabrera
tags: [cli, docs]
---
# Update README and documentation for new command structure

After the command restructure, update:
- README.md command reference section
- All template READMEs that reference commands
- AGENTS.md if it references specific commands
- Help text on all commands (Short, Long, Example fields)

Ensure the documentation matches the new surface:
  dade template list/add/remove/update
  dade project list/add/remove/port/sync  
  dade proxy start/stop/restart/status/logs/reload
  dade share --attach (replacing tunnel)

