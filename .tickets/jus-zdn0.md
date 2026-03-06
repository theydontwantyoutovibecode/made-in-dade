---
id: jus-zdn0
status: open
deps: []
links: []
created: 2026-03-02T04:57:27Z
type: task
priority: 1
assignee: Alex Cabrera
tags: [cli, restructure]
---
# Move start/stop under 'project' parent as subcommands

start and stop are production server lifecycle commands, not daily-use dev commands. They would typically be triggered by systemd/launchd, not typed interactively. Move them under the 'project' parent.

Current → Proposed:
  dade start [name]     → dade project start [name]
  dade stop [name]      → dade project stop [name]

This changes the final top-level surface to:

DEVELOPMENT
  new [name]           Create a new project from a template
  dev [name]           Start development server
  build [name]         Build a compiled project
  open [name]          Open project in browser
  share [name]         Share project via public tunnel

MANAGEMENT
  project <cmd>        list, register, remove, port, sync, start, stop
  template <cmd>       list, add, remove, update
  proxy <cmd>          start, stop, restart, status, logs, reload
  setup                First-time setup

That is 5 top-level verbs + 3 grouped nouns + setup = 9 entries total in help.

Add hidden compat aliases for 'start' and 'stop' at the top level (covered by jus-7z4w).

