---
id: jus-du4p
status: open
deps: []
links: []
created: 2026-03-02T04:52:48Z
type: task
priority: 1
assignee: Alex Cabrera
tags: [cli, restructure]
---
# Absorb 'refresh' into 'proxy reload' subcommand

The 'refresh' command regenerates the Caddyfile and reloads the proxy. This is conceptually a proxy operation. Move it into the existing 'proxy' command group as 'proxy reload'.

Current → Proposed:
  dade refresh [--list]       → dade proxy reload [--list]

This eliminates one top-level command and puts the behavior where users would intuitively look for it.

Implementation:
- Add 'reload' subcommand to proxyCmd in cmd_proxy.go
- Move refresh logic from cmd_refresh.go into the new subcommand
- Remove cmd_refresh.go
- Update tests

