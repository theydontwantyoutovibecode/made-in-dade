---
id: jus-52py
status: closed
deps: [jus-4cer]
links: []
created: 2026-02-17T02:05:28Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-gyx3
---
# Implement config and template metadata handling in Go

Port configuration behavior from Bash: default config dir ~/.config/dade; load templates.toml overrides (TOML parsing) populating template URL map; maintain default templates: django-hypermedia -> https://github.com/theydontwantyoutovibecode/dade-with-django-and-hypermedia.git, hypertext -> https://github.com/theydontwantyoutovibecode/dade-with-hypertext.git; display names map for UX; expose struct(s)/functions to fetch template list with names/urls, merged with user overrides. Preserve insertion order for defaults (so help/templates output matches Bash ordering). Include detection of missing/unknown template with clear error. Provide unit tests covering parsing, override precedence, bad/missing file handling. Use go-toml or stdlib alternative (decide and note dependency).

## Acceptance Criteria

- Config loader reads ~/.config/dade/templates.toml when present\n- Defaults present when no overrides and ordering is deterministic\n- Overrides replace/add templates correctly\n- Unknown template errors with helpful message\n- Tests cover default, override, malformed file cases

