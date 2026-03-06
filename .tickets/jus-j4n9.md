---
id: jus-j4n9
status: closed
deps: []
links: []
created: 2026-02-17T02:05:47Z
type: task
priority: 3
assignee: Alex Cabrera
parent: jus-gyx3
---
# Add Charm UI layer in Go

Implement a Charm-based UI layer mirroring Bash UX: use lipgloss for headers/styles, charm/log for colored logs, huh for input/select prompts, and bubbles/spinner (with bubbletea runtime) for spinners. No external gum dependency. Design an interface with fallback to plain output when non-tty. Apply to commands: header, log_info/success/warn/error, spinner wrapper, interactive template selection (if multiple templates and no --template provided). Ensure prompts/spinners are gated by TTY detection so non-interactive environments behave without hangs. Provide tests where feasible via abstraction/mocking; otherwise, structure so UI dependencies are optional and don't break when unavailable.

## Acceptance Criteria

- Charm UI components integrated (lipgloss, log, huh, bubbles/bubbletea)\n- Header/log helpers mirror Bash look-and-feel when TTY supports styling\n- Interactive flows use huh select/input when TTY, else fallback/error\n- Spinners render via bubbles when TTY; fallback prints “msg... done/failed”\n- Non-interactive mode never prompts and never hangs\n- Tests or mock coverage for fallback logic

