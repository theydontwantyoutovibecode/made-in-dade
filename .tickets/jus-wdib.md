---
id: jus-wdib
status: closed
deps: [jus-z19f, jus-syr2]
links: []
created: 2026-02-23T16:54:45Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-3el5
tags: [commands, dev]
---
# Implement dade dev command

New command for development server orchestration.

Usage:
  dade dev [name]         # Start dev server
  dade dev --background   # Detach and run in background

Workflow:
1. Resolve project (by name or current directory)
2. Read template manifest
3. If [dev.setup] exists: run setup commands sequentially
4. If [dev.background] exists: start background processes
5. Set environment from [dev.env]
6. Start main dev server (serve.dev)
7. Ensure Caddy proxy is running
8. Display HTTPS URL
9. Handle signals for clean shutdown

Flags:
- --background/-b: Detach and run in background (write PID)
- --skip-setup: Skip setup commands
- --port: Override port

Fallback behavior:
- If no [dev] section, run serve.dev directly (current start behavior)

Files to modify/create:
- cmd/dade/cmd_dev.go (new)
- cmd/dade/dev.go (new)

Acceptance:
- Works with existing templates (backward compatible)
- Works with new [dev] section
- Clean output with project URL
- Ctrl+C cleanly stops server and background processes

