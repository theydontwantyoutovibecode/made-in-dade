---
id: jus-4cer
status: closed
deps: [jus-gyx3]
links: []
created: 2026-02-17T02:05:00Z
type: task
priority: 1
assignee: Alex Cabrera
---
# Document Go migration requirements and parity for dade

Document the full behavioral spec of the current Bash CLI (dade v0.1.0) to guide the Go port. Include commands (new, templates, --help, --version), default template URLs (django-hypermedia -> https://github.com/theydontwantyoutovibecode/dade-with-django-and-hypermedia.git, hypertext -> https://github.com/theydontwantyoutovibecode/dade-with-hypertext.git), display names (Django + Hypermedia (HTMX, TailwindCSS), HTML + Hypertext), config dir (~/.config/dade, templates.toml overrides), interactions with gum (style, choose, input, spin) and their replacements in Go (lipgloss/log/huh/bubbles), project name validation regex ^[a-zA-Z][a-zA-Z0-9_-]*$, git init + template .git removal, local vs remote template copy/clone, setup.sh execution rules, headers/log messaging, and next steps output. Capture assumptions, gaps, and any new requirements for Go (binary name, installation flow). Deliver as markdown doc in repo (e.g., docs/go-migration-spec.md) with enough detail for downstream tickets.

## Acceptance Criteria

- Spec file added to repo documenting current Bash behavior and desired Go parity
- Includes defaults (template URLs/names, config dir), validation rules, and workflow steps for new/templates/help/version
- Notes gum vs fallback UX expectations and Charm replacements
- Identifies any open questions or decisions for Go implementation
- CI not required but doc must render in GitHub markdown

