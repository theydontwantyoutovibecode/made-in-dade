---
id: jus-0vx2
status: closed
deps: []
links: []
created: 2026-02-17T02:05:41Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-gyx3
---
# Implement 'templates' command in Go with JSON/plain output

Port 'templates' command: list available templates from config/defaults + overrides. Support plain text output and --json flag. Plain output prints header and each template name, display name/description, and URL. For templates without display metadata, fall back to the template name. Always print guidance snippet for ~/.config/dade/templates.toml (mirrors Bash). JSON output returns array of objects with at least name, url; optionally description/type if available in future metadata. Handle empty templates gracefully with guidance. Provide tests covering default list, overrides, --json shape, empty state, and snippet output.

## Acceptance Criteria

- templates command prints list or JSON array with name/display/url\n- Includes default templates and merges overrides\n- --json outputs valid JSON with entries for each template\n- Output includes templates.toml guidance snippet\n- Empty state prints guidance\n- Tests cover defaults, overrides, empty, and snippet output

