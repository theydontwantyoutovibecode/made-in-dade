---
id: jus-lrn6
status: closed
deps: []
links: []
created: 2026-02-17T02:05:54Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-gyx3
---
# Implement version/help commands in Go

Port --version/-v and --help/-h commands: version prints dade v0.1.0 and exits 0; help prints usage matching Bash (commands new/templates/help/version, options for new), includes Available Templates and Configuration sections, with lipgloss styling when TTY and plain fallback otherwise. Ensure unknown command path shows error then help and exits 1. Provide tests for outputs and exit codes.

## Acceptance Criteria

- Version/help output matches Bash content (updated for Go binary wording)\n- Exit codes: version/help 0, unknown command 1\n- Help includes usage lines, options, available templates, and config path\n- Styled output uses lipgloss when TTY, plain otherwise\n- Tests cover version/help/unknown command paths

