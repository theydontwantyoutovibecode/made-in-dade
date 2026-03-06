---
id: jus-uyus
status: closed
deps: [jus-u0zj]
links: []
created: 2026-03-02T01:56:31Z
type: epic
priority: 0
assignee: Alex Cabrera
tags: [templates, bundling, ws-2]
---
# WS-2: Bundle Default Templates

Make dade ship with all 6 default templates pre-installed so users never need to run dade install for the defaults. On first run or dade setup, auto-clone all default templates into ~/.config/dade/templates/. The install command remains for user-added templates from arbitrary sources.

## Acceptance Criteria

1. Fresh install of dade → all 6 templates available immediately. 2. dade new shows all defaults without install. 3. dade install still works for custom templates. 4. dade update refreshes defaults.

