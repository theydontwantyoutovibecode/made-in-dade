---
id: jus-mkgk
status: closed
deps: []
links: []
created: 2026-03-02T02:59:27Z
type: task
priority: 1
assignee: Alex Cabrera
tags: [docs, ws-9, cli]
---
# Rewrite cli template README.md

Current README is a 3-line stub. Needs comprehensive documentation:
- What the template creates (Go CLI with Cobra, Fang, Lipgloss, Huh)
- Prerequisites (Go 1.22+)
- What setup.sh does
- Project structure (cmd/, internal/config/, main.go)
- How dade dev works (go mod download, then go run .)
- How dade build works (go build with placeholders, cross-compilation, --release strips symbols)
- That dade start, share, tunnel, open do NOT apply (no port, no server)
- Cobra command structure
- Fang configuration
- .read-only manifest.txt contents (charm repos)
- AGENTS.md and .tickets workflow

