---
id: jus-bz78
status: closed
deps: []
links: []
created: 2026-02-17T02:05:51Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-gyx3
---
# Implement logging utilities and error handling in Go

Create centralized logging utilities matching Bash styles (log_info, log_success, log_warn, log_error) using charm/log plus lipgloss for color/styling when TTY; fall back to plain output when not. Ensure messages can be suppressed/adjusted for JSON outputs where needed. Provide consistent exit codes and error wrapping for user-facing messages. Include unit tests for formatting and exit codes.

## Acceptance Criteria

- Logging helpers available across commands with colored output fallback\n- Uses charm/log + lipgloss (no external gum dependency)\n- Error handling utilities for consistent exit codes/messages\n- Compatible with Charm UI layer but not dependent on it\n- Tests cover formatting and behavior

