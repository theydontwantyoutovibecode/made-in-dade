---
id: jus-tk39
status: closed
deps: []
links: []
created: 2026-02-19T18:59:04Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-dvrr
---
# Add --port/-p flag to start.sh for production mode

Add command-line argument parsing for --port/-p flag. In production mode, this flag is required if PORT env var is not set. In dev mode, this flag overrides auto-detection.

## Acceptance Criteria

- --port/-p accepts integer argument
- Production mode: uses PORT env var OR --port flag, errors if neither
- Dev mode: --port overrides auto-detection
- Help text updated to document the flag

