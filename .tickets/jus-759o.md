---
id: jus-759o
status: closed
deps: [jus-7hfd]
links: []
created: 2026-03-02T01:59:17Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-yi7f
tags: [tk, scaffold, ws-8]
---
# Initialize .tickets/ in dade new scaffold

After copying the template and running git init, create a .tickets/ directory in the new project. This ensures tk has a place to store tickets immediately. Could also run tk init if that command exists, or simply mkdir .tickets.

## Acceptance Criteria

1. dade new creates .tickets/ dir. 2. tk commands work in new project immediately.

