---
id: jus-odp6
status: open
deps: []
links: []
created: 2026-03-02T05:03:59Z
type: task
priority: 2
assignee: Alex Cabrera
tags: [cli, restructure, docs]
---
# Audit all template repos for old command references

After the command restructure, audit every template repository for references to old command names and update them.

Files to check in each of the 6 templates:
- README.md — usage examples, command references
- AGENTS.md — AI agent instructions referencing commands
- setup.sh — any dade CLI calls
- start.sh — any dade CLI calls
- dade.toml — help text or comments referencing commands

Old → New mappings to search for:
  dade install        → dade template add
  dade uninstall      → dade template remove
  dade templates      → dade template list
  dade update         → dade template update
  dade list           → dade project list
  dade register       → dade project register
  dade remove / rm    → dade project remove
  dade port           → dade project port
  dade sync           → dade project sync
  dade start          → dade project start
  dade stop           → dade project stop
  dade open           → dade dev --open
  dade tunnel         → dade share --attach
  dade refresh        → dade proxy reload

Also check the main dade README.md (already covered by jus-xw6m but this ticket is specifically about the 6 template repos).

Template repos:
- web-app-made-in-dade
- web-site-made-in-dade
- ios-made-in-dade
- android-made-in-dade
- cli-made-in-dade
- tui-made-in-dade

