---
id: jus-eeb9
status: closed
deps: [jus-drxv, jus-tk39, jus-earu]
links: []
created: 2026-02-19T18:59:18Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-dvrr
---
# Update start_development_server to use auto-detected port

Modify the start_development_server function to: (1) Check if assigned/default port is available, (2) If not, find next available port, (3) Call dade to update Caddy, (4) Start Django on the new port.

## Acceptance Criteria

- Function checks port availability before starting
- Auto-selects available port if needed
- Logs selected port clearly to user
- Caddy is updated before server starts
- Django receives correct PORT value

