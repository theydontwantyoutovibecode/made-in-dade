---
id: jus-dvrr
status: open
deps: []
links: []
created: 2026-02-19T18:58:50Z
type: epic
priority: 2
assignee: Alex Cabrera
---
# Auto-detect available port in development mode

When running start.sh --dev, the script should automatically find an available port if the default (8000) or assigned port is in use. This allows multiple dade projects to run simultaneously. The system should: (1) Check if assigned port is available, (2) Auto-find next available port if not, (3) Update Caddy configuration to point to the new port, (4) Only apply in dev mode - production requires explicit --port flag.

## Acceptance Criteria

- start.sh --dev auto-selects available port when default is in use
- Caddy reverse proxy is updated to reflect the new port
- Production mode requires explicit port specification
- Port auto-selection is logged clearly to user
- System works across multiple simultaneous dade projects

