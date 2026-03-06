---
id: jus-lz4t
status: closed
deps: [jus-bx3x]
links: []
created: 2026-02-12T20:34:52Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-swo4
tags: [templates, django]
---
# Update django template .gitignore

Update the Django template's .gitignore to include dade files.

## Add to .gitignore

```
# dade
.dade.pid
```

## Note

The .dade marker file should be committed (it defines project config).
Only the .dade.pid file should be ignored.

## Acceptance Criteria

- [ ] .dade.pid added to .gitignore
- [ ] .dade NOT ignored (should be committed)

