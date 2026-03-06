---
id: jus-7z4w
status: open
deps: []
links: []
created: 2026-03-02T04:53:44Z
type: task
priority: 2
assignee: Alex Cabrera
tags: [cli, restructure, compat]
---
# Add hidden backward-compat aliases for renamed commands

After restructuring, keep the old command names as hidden aliases so existing scripts and muscle memory don't break immediately. These should not appear in help output.

Hidden aliases to add:
  dade install <url>     → routes to 'template add <url>'
  dade uninstall <name>  → routes to 'template remove <name>'  
  dade templates         → routes to 'template list'
  dade list              → routes to 'project list'
  dade register [name]   → routes to 'project add [name]'
  dade remove/rm [name]  → routes to 'project remove [name]'
  dade tunnel [name]     → routes to 'share --attach [name]'
  dade refresh           → routes to 'proxy reload'
  dade sync [path]       → routes to 'project sync [path]'
  dade port              → routes to 'project port'

Implementation: register hidden cobra commands (Hidden: true) that internally call the new command's RunE. This is a bridge — can be removed in a future major version.

