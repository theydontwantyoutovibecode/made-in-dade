---
id: jus-f3wf
status: closed
deps: [jus-psut]
links: []
created: 2026-02-12T20:35:06Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-zf02
tags: [templates, static]
---
# Remove serve.sh from hypertext template

Remove the serve.sh script from the hypertext template since dade handles serving.

## Action

Delete: /Users/acabrera/Code/dade/dade-with-hypertext/serve.sh

## Why

- dade now handles serving static files
- serve.sh is redundant
- Keeps template minimal

## Fallback Documentation

Update README.md to show fallback for users without dade:

```bash
# Without dade
python3 -m http.server 8000
```

## Acceptance Criteria

- [ ] serve.sh deleted
- [ ] README documents fallback

