---
id: jus-3el5
status: closed
deps: []
links: []
created: 2026-02-23T16:54:14Z
type: epic
priority: 0
assignee: Alex Cabrera
tags: [core, commands, dx]
---
# Epic: dade dev/share commands

Move development server orchestration from template start.sh scripts into dade core.

Current state: Templates like django-hypermedia have 1000+ line start.sh scripts handling dev servers, port detection, dependency installation, asset watchers, tunnel setup, etc.

Goal: dade dev and dade share commands that orchestrate development workflow while templates define hooks for template-specific behavior.

See docs/dev-share-architecture.md for full design.

