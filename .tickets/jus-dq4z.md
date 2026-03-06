---
id: jus-dq4z
status: open
deps: [jus-1pcj]
links: []
created: 2026-02-23T16:55:16Z
type: task
priority: 3
assignee: Alex Cabrera
parent: jus-3el5
tags: [migration, projects]
---
# Migrate marginalia and burnball to new dev/share system

Update existing projects to use new dev/share commands.

Projects:
- /Users/acabrera/Code/marginalia
- /Users/acabrera/Code/burnball

For each project:
1. Update .dade marker if needed
2. Add [dev] and [share] sections to local dade.toml (if supported)
3. Or rely on template defaults
4. Test with dade dev
5. Test with dade share
6. Verify start.sh still works as fallback

Validation:
- PostgreSQL setup still works
- Migrations run
- TailwindCSS watcher runs
- Tunnel creates public URL
- ALLOWED_HOSTS configured for tunnel

Acceptance:
- dade dev works for both projects
- dade share works for both projects
- No regression in functionality

