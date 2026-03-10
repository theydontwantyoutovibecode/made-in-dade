---
id: mid-z748
status: open
deps: []
links: [mid-bjyk]
created: 2026-03-07T19:08:06Z
type: Hosts: Create sudo wrapper for privileged operations
priority: 1
assignee: Alex Cabrera
---
# Hosts: Create sudo wrapper for privileged operations

Create internal/exec/sudo.go with RunWithSudo() helper, add confirmation prompts for privileged operations, handle sudo password prompts gracefully

