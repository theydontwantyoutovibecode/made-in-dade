---
id: jus-ysvu
status: closed
deps: [jus-u0zj]
links: []
created: 2026-03-02T01:58:01Z
type: epic
priority: 1
assignee: Alex Cabrera
tags: [agents, documentation, ws-5]
---
# WS-5: Rewrite AGENTS.md for all templates

Completely rewrite every template's AGENTS.md. Remove all references to novice users. Focus on: 1) Clear stack documentation. 2) How to use .read-only/ for reference libraries. 3) Enforce tk CLI for ticket-driven development. 4) Every request must produce granular tickets with: background context, verbose description, atomic subtasks, definition of done, caveats/edge cases, required test coverage. 5) Agent analyzes plan holistically, links tickets, creates rollup/checkpoint tickets. 6) One ticket = one commit, never work multiple tickets. 7) Before implementation: surface open questions, validate plan. 8) Agent fills gaps creatively — user is vibecoding and may never see the code.

## Acceptance Criteria

1. All 6 templates have rewritten AGENTS.md. 2. No mention of novice/beginner. 3. tk CLI workflow documented. 4. .read-only documented. 5. Ticket format specification included.

