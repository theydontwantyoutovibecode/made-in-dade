---
id: mid-r64m
status: closed
deps: [mid-0jq8]
links: []
created: 2026-03-11T18:13:18Z
type: task
priority: 1
assignee: Alex Cabrera
tags: [hot-reload, caddy, html]
---
# Auto-inject client-side SSE script into HTML files

Modify the Caddyfile generation in cmd_dev.go to inject a small JavaScript snippet into all HTML files that: (1) Connects to the SSE endpoint, (2) Listens for 'reload' events, (3) Reloads the page when received. The script should be minimal and handle connection errors gracefully. This ensures the browser automatically refreshes when files change.

