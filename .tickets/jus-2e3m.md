---
id: jus-2e3m
status: closed
deps: [jus-6d0x]
links: []
created: 2026-02-12T19:53:42Z
type: epic
priority: 1
assignee: Alex Cabrera
parent: jus-nq0k
tags: [templates, migration]
---
# Template Migrations

Update existing template repositories to work with the new plugin architecture. Each template needs a dade.toml manifest and potentially updated scripts.

## Templates to Migrate

1. **dade-with-django-and-hypermedia**
   - Add dade.toml with command-based serving
   - Keep existing start.sh (works as-is)
   - Update AGENTS.md with new dade commands

2. **dade-with-hypertext**
   - Add dade.toml with static serving
   - Remove serve.sh (replaced by dade infrastructure)
   - Update AGENTS.md and README.md

## Manifest Requirements

Each template must have dade.toml in root with:
- Template metadata (name, description)
- Scaffold configuration (exclude patterns, setup script)
- Serve configuration (static vs command)

## Backward Compatibility

Templates should still work standalone:
- Django template: ./start.sh --dev still works
- Static template: can still use python -m http.server

The manifest just enables integration with dade's enhanced features (HTTPS, tunnels, etc.)

## Acceptance Criteria

- [ ] django-hypermedia template has dade.toml
- [ ] hypertext template has dade.toml
- [ ] Both templates work with new dade CLI
- [ ] Both templates still work standalone
- [ ] Documentation updated in both templates

