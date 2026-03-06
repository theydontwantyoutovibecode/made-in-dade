---
id: jus-a3c6
status: closed
deps: []
links: []
created: 2026-02-19T18:58:55Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-dvrr
---
# Add port availability check function to start.sh

Add a bash function 'is_port_available()' that checks if a given port is in use by attempting to bind to it. Use netcat or /dev/tcp probe as fallback.

## Acceptance Criteria

- Function returns 0 (true) if port is available
- Function returns 1 (false) if port is in use
- Works on macOS and Linux
- Handles edge cases (port 0, invalid port numbers)

