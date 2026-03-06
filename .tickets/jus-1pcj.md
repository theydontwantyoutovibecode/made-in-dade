---
id: jus-1pcj
status: closed
deps: [jus-wdib, jus-t6md]
links: []
created: 2026-02-23T16:55:09Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-3el5
tags: [template, django]
---
# Update django-hypermedia template for dev/share

Migrate django-hypermedia template to use new manifest sections.

Changes to dade.toml:
1. Add [dev] section with:
   - setup commands: uv sync, migrations
   - background: tailwindcss watcher
   - env: DJANGO_SETTINGS_MODULE

2. Add [share] section with:
   - env: ALLOWED_HOSTS, CSRF_TRUSTED_ORIGINS

3. Simplify serve commands:
   - dev: just the runserver command
   - prod: just the gunicorn command

Changes to start.sh:
1. Keep as thin wrapper or escape hatch
2. Move orchestration logic out
3. Target <100 lines

Template location:
- /Users/acabrera/Code/dade/dade-with-django-and-hypermedia/

Acceptance:
- Template works with dade dev
- Template works with dade share
- start.sh is dramatically simplified
- Backward compatible with existing projects

