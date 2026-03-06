---
id: jus-ekuk
status: closed
deps: []
links: []
created: 2026-03-02T02:59:04Z
type: task
priority: 1
assignee: Alex Cabrera
tags: [docs, ws-9, web-app]
---
# Rewrite web-app template README.md

Current README is a 3-line stub. Needs comprehensive documentation:
- What the template creates (Django + HTMX + vanilla CSS)
- Prerequisites (Python 3.13+, uv)
- What setup.sh does
- Project structure
- How dade dev works (runs uv sync, migrations, starts runserver)
- How dade start works (gunicorn production mode)
- How dade share works (Cloudflare tunnel with ALLOWED_HOSTS)
- How dade build works (not applicable - explain why)
- .read-only manifest.txt contents and purpose
- Environment configuration (.env, settings modules)
- AGENTS.md and .tickets workflow

