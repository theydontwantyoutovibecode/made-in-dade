---
id: jus-7bvm
status: closed
deps: []
links: []
created: 2026-03-02T02:58:54Z
type: task
priority: 1
assignee: Alex Cabrera
tags: [docs, ws-9]
---
# Rewrite dade README.md

The current README has multiple inaccuracies:
- References old template names (django-hypermedia, hypertext) and old repo URLs (dade-with-*)
- Lists only 2 templates instead of 6
- Shows dade install commands that are unnecessary (templates are auto-installed)
- Implies all templates support share/tunnel/start (ios, android, cli, tui do not serve on ports)
- Says "Linux support planned" but dade is macOS-only
- The "Creating Templates" section shows command_windows which was removed
- The "How It Works" section oversimplifies and is partially wrong

Rewrite as comprehensive plain-language documentation covering:
1. What dade is (one paragraph, no marketing)
2. Installation (Homebrew + source)
3. Setup (what dade setup does)
4. All 6 templates with accurate descriptions of what each supports
5. Every command with flags, behavior, and which templates support it
6. Template system (dade.toml schema, creating custom templates)
7. HTTPS proxy system (how Caddy works, domains, LAN access)
8. .read-only reference libraries
9. .tickets and tk CLI
10. Project registry and configuration paths

