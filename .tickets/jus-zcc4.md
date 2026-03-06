---
id: jus-zcc4
status: closed
deps: [jus-bx3x]
links: []
created: 2026-02-12T20:34:39Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-swo4
tags: [templates, django]
---
# Update django template start.sh for PORT env var

Verify and update the Django template's start.sh to properly respect PORT environment variable.

## Current Usage

The script already uses PORT:
```bash
PORT="${PORT:-8000}"
```

## Verification

1. Check that PORT is used consistently
2. Ensure development server binds to PORT
3. Ensure Gunicorn binds to PORT

## Changes If Needed

The script should:
1. Default PORT to 8000 if not set
2. Use $PORT in Django runserver command
3. Use $PORT in Gunicorn command

## Acceptance Criteria

- [ ] PORT env var is read
- [ ] Default is 8000 if not set
- [ ] Dev server uses PORT
- [ ] Prod server uses PORT

