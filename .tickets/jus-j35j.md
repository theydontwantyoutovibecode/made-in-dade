---
id: jus-j35j
status: closed
deps: []
links: []
created: 2026-02-17T02:06:06Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-gyx3
---
# Implement tests for Go CLI flows

Add comprehensive tests for Go CLI: command parsing, help/version, error cases, new/templates flows using temp directories and fake runners. Include integration-style tests invoking main with args and capturing stdout/stderr/exit codes. Ensure coverage for validation errors, unknown template, existing dir, git failures, setup.sh execution paths (exec vs bash), empty templates list, JSON output.

## Acceptance Criteria

- Tests cover parsing and happy/error paths for commands\n- Tempdir-based tests for filesystem effects\n- Mocks/fakes for exec to simulate git failures/success\n- JSON output validated\n- CI-ready go test ./... passes

