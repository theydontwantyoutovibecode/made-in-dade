---
id: jus-pvzg
status: closed
deps: []
links: []
created: 2026-02-17T02:06:02Z
type: task
priority: 3
assignee: Alex Cabrera
parent: jus-gyx3
---
# Add interactive selection support for templates in Go

Implement interactive template selection: when user runs 'dade new' without --template and multiple templates exist, prompt to choose (huh select when TTY, fallback to numbered select). Handle non-interactive detection to avoid hanging (return error if stdin not tty). Integrate with template loader from config. Tests for fallback selection logic (non-interactive error, selection mapping).

## Acceptance Criteria

- Choice prompt shown when template not specified and multiple templates available\n- Uses huh select when TTY, else numbered menu\n- Non-interactive mode errors gracefully without blocking\n- Tests cover fallback logic and template mapping

