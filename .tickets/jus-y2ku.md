---
id: jus-y2ku
status: closed
deps: [jus-nmo3]
links: []
created: 2026-03-02T01:56:21Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-u0zj
tags: [cleanup, ws-1]
---
# Update internal references after rename

Update DefaultTemplates() in internal/config/templates.go to use new names and URLs (theydontwantyoutovibecode org). Update available_templates.go display names. Update any hardcoded template references in tests. Update AGENTS.md references in each template if they reference old paths. Ensure dade install --list-official shows correct info.

## Acceptance Criteria

1. DefaultTemplates() returns correct new names/URLs. 2. dade install --list-official shows all 6 templates with new URLs. 3. All tests pass.

