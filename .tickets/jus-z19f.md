---
id: jus-z19f
status: closed
deps: []
links: []
created: 2026-02-23T16:54:21Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-3el5
tags: [manifest, parser]
---
# Extend manifest parser for [dev] and [share] sections

Extend the manifest package to parse new sections.

New [dev] section fields:
- setup: []string - Commands to run before server starts
- background: []string - Commands to run in background alongside server  
- env: []string - Environment variables for dev mode
- setup_script: string - Optional custom setup script path

New [share] section fields:
- env: []string - Additional env vars when sharing via tunnel
- tunnel_name: string - Named tunnel (if configured)
- tunnel_domain: string - Domain for named tunnel

Files to modify:
- internal/manifest/manifest.go

Acceptance:
- Parse new sections without breaking existing manifests
- Validate new fields appropriately
- Unit tests for new parsing logic

