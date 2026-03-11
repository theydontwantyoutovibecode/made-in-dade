---
id: mid-xjd7
status: closed
deps: [mid-0jq8, mid-yz6r, mid-r64m]
links: []
created: 2026-03-11T18:13:21Z
type: task
priority: 2
assignee: Alex Cabrera
tags: [hot-reload, lifecycle, dev]
---
# Integrate hot-reload into dev command lifecycle

Update cmd_dev.go to properly integrate all hot-reload components: (1) Start SSE server when dev mode begins, (2) Start file watcher with SSE server callback, (3) Start Tailwind watch with coordination, (4) Ensure all processes shut down cleanly on Ctrl+C. The current implementation starts processes but doesn't coordinate them properly or ensure clean shutdown.

