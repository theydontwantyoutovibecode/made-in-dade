---
id: jus-g3fm
status: closed
deps: [jus-u0zj]
links: []
created: 2026-03-02T01:58:45Z
type: epic
priority: 1
assignee: Alex Cabrera
tags: [read-only, agents, ws-7]
---
# WS-7: .read-only Integration

Every template ships with a curated .read-only/manifest.txt listing preferred libraries and example repos for AI agent context. Move the sync logic from start.sh into the dade binary so it works across all templates. dade dev calls syncReadOnlyDeps() before starting the dev server.

## Acceptance Criteria

1. Every template has .read-only/manifest.txt. 2. dade dev syncs .read-only deps. 3. Django template start.sh sync removed (handled by binary). 4. Tests cover sync logic.

