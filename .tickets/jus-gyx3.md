---
id: jus-gyx3
status: closed
deps: []
links: []
created: 2026-02-17T02:05:11Z
type: epic
priority: 1
assignee: Alex Cabrera
---
# Epic: Port dade Bash CLI to Go binary

Migrate the current Bash-based dade CLI (v0.1.0) into a Go binary while preserving feature parity and UX (gum-like styling/prompts/spinners via Charm libraries; plain fallback for non-tty). Deliver a Go module with reproducible builds, tests, and updated docs/install flows. Ensure behaviors: new/templates/help/version commands, config (~/.config/dade, templates.toml overrides), default template URLs/names, project name validation, local/remote template copying/cloning with git --depth=1, stripping template .git, running setup.sh if present/executable, and next-steps output. Include parity for logging/styling, spinner/interactive choices, and JSON output support where applicable. Include test coverage and release packaging guidance.

## Acceptance Criteria

- Go binary delivers all Bash CLI behaviors (new, templates, --help, --version)\n- Charm UI/logging libraries are used for styling, prompts, and spinners (no external gum dependency)\n- Module builds via go build ./...; tests cover core flows\n- Docs updated for install/usage\n- Tickets for each atomic migration task created and linked

