---
id: jus-t6md
status: closed
deps: [jus-wdib, jus-6033]
links: []
created: 2026-02-23T16:54:52Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-3el5
tags: [commands, share, tunnel]
---
# Implement dade share command

New command to share project via public tunnel.

Usage:
  dade share [name]       # Start dev server + create tunnel

Workflow:
1. Run full dev setup (same as dade dev)
2. Start dev server
3. Set additional env from [share.env]
4. Check for named tunnel in [share.tunnel_name]
5. Start cloudflared tunnel (named or quick)
6. Capture and display public URL
7. Handle signals for clean shutdown

Named tunnel support:
- If tunnel_name set, use: cloudflared tunnel run <name>
- If tunnel_domain set, show custom domain URL
- Otherwise use quick tunnel

Flags:
- --quick: Force quick tunnel even if named configured
- --port: Override port

Files to modify/create:
- cmd/dade/cmd_share.go (new)
- cmd/dade/share.go (new)

Acceptance:
- Starts server and tunnel with one command
- Uses named tunnel if configured
- Falls back to quick tunnel
- Shows public URL
- Clean shutdown on Ctrl+C

