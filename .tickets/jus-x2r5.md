---
id: jus-x2r5
status: open
deps: []
links: []
created: 2026-03-02T04:52:41Z
type: task
priority: 1
assignee: Alex Cabrera
tags: [cli, restructure]
---
# Group project registry commands under 'dade project' parent

Create a 'project' parent command with subcommands: list, add, remove, port, sync. This replaces five top-level commands (list, register, remove, port, sync) with one parent.

Current → Proposed:
  dade list                   → dade project list
  dade register [name]        → dade project add [name]
  dade remove [name]          → dade project remove [name]
  dade port                   → dade project port
  dade sync [path]            → dade project sync [path]

The verb-noun pattern reads like sentences: 'project list', 'project add myapp', 'project remove old-one'.

The parent command 'dade project' with no subcommand should default to 'project list' for convenience.

The 'rm' alias on remove should be preserved as an alias on 'project remove'.

Implementation:
- Create cmd_project.go with the parent command
- Move list logic into 'project list' subcommand
- Move register logic into 'project add' subcommand
- Move remove logic into 'project remove' subcommand
- Move port logic into 'project port' subcommand
- Move sync logic into 'project sync' subcommand
- Remove old cmd_list.go, cmd_register.go, cmd_remove.go, cmd_port.go, cmd_sync.go
- Update tests accordingly

