---
id: jus-fgv0
status: open
deps: [jus-eeb9]
links: []
created: 2026-02-19T18:59:27Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-dvrr
---
# Add test for port auto-detection

Create manual or automated test to verify port auto-detection works correctly when multiple dade projects attempt to start on the same port.

## Acceptance Criteria

- Test starts two projects with same default port
- Second project auto-selects different port
- Both projects accessible via their .localhost URLs
- Caddy correctly routes to each backend

