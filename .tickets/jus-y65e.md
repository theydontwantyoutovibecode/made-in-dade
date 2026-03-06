---
id: jus-y65e
status: closed
deps: []
links: []
created: 2026-03-02T02:59:33Z
type: task
priority: 1
assignee: Alex Cabrera
tags: [docs, ws-9, tui]
---
# Rewrite tui template README.md

Current README is a 3-line stub. Needs comprehensive documentation:
- What the template creates (Go TUI with Bubbletea v2, Bubbles, Lipgloss)
- Prerequisites (Go 1.22+)
- What setup.sh does
- Project structure (model.go, update.go, view.go, keys.go, styles.go, main.go)
- Elm Architecture explanation (Init, Model, View, Update)
- How dade dev works (go mod download, then go run . with DEBUG=1)
- How dade build works (go build with placeholders, cross-compilation, --release strips symbols)
- That dade start, share, tunnel, open do NOT apply (no port, no server)
- Bubbletea v2 key differences from v1
- .read-only manifest.txt contents (bubbletea, bubbles, lipgloss repos)
- AGENTS.md and .tickets workflow

