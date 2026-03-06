---
id: jus-drxv
status: closed
deps: [jus-a3c6]
links: []
created: 2026-02-19T18:58:59Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-dvrr
---
# Add find_available_port function to start.sh

Add a bash function 'find_available_port()' that starts from a given port and increments until an available port is found. Should have a maximum search range (e.g., 100 ports) to avoid infinite loops.

## Acceptance Criteria

- Function takes starting port as argument
- Returns first available port in range
- Exits with error if no port found in range
- Range limit is configurable (default 100 ports)

