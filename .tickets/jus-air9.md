---
id: jus-air9
status: open
deps: []
links: []
created: 2026-03-02T04:52:55Z
type: task
priority: 2
assignee: Alex Cabrera
tags: [cli, restructure]
---
# Merge 'tunnel' into 'share' with a flag

The 'tunnel' command is a lightweight version of 'share' — it creates a tunnel to an already-running server. Having both as top-level commands is confusing. Merge 'tunnel' into 'share --attach' (or 'share --no-server').

Current → Proposed:
  dade share [name]           → dade share [name]        (unchanged, starts server + tunnel)
  dade tunnel [name]          → dade share --attach [name] (tunnel only, server must be running)

The --attach flag communicates intent clearly: you're attaching a tunnel to an existing process.

Implementation:
- Add --attach flag to shareCmd
- When --attach is set, skip server startup and just create the tunnel (current tunnel behavior)
- Remove cmd_tunnel.go
- Update tests

