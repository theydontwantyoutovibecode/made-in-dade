---
id: jus-vf6p
status: open
deps: []
links: []
created: 2026-03-02T04:53:03Z
type: task
priority: 3
assignee: Alex Cabrera
tags: [cli, restructure]
---
# Keep dev, new, build, start, stop, open, share, setup as top-level commands

These are the daily-use commands and should remain at the top level for ergonomics. No changes needed other than ensuring they render correctly in the new help output with command groups.

Final top-level command surface:
  dade new [name]         Create a new project from a template
  dade dev [name]         Start development server
  dade build [name]       Build a compiled project
  dade start [name]       Start production server
  dade stop [name]        Stop a running project
  dade open [name]        Open project in browser
  dade share [name]       Share project via public tunnel
  dade setup              First-time setup
  dade project <cmd>      Manage project registry
  dade template <cmd>     Manage installed templates
  dade proxy <cmd>        Manage HTTPS proxy

That is 8 top-level verbs + 3 grouped nouns = 11 entries in help output, down from 21. Each group expands into 4-6 subcommands for power users.

