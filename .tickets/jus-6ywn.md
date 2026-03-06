---
id: jus-6ywn
status: closed
deps: [jus-k3jv]
links: []
created: 2026-03-02T01:57:19Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-d7n2
tags: [ux, ws-3]
---
# Add template inspect mode to dade new

When the user highlights a template in the selection menu, allow them to press i or ? to see full details: stack components, what files/directories get created, build/dev commands, and any setup requirements. Show this in a styled panel. Press esc or backspace to return to selection. This requires reading additional fields from the manifest or a dedicated [template.details] section. Consider adding a stack field to [template] in dade.toml.

## Acceptance Criteria

1. Pressing i/? shows template details. 2. Details include stack, structure, commands. 3. Esc returns to selection. 4. Works in non-interactive mode (--template flag skips).

