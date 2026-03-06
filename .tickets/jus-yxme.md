---
id: jus-yxme
status: closed
deps: []
links: []
created: 2026-03-02T06:14:24Z
type: epic
priority: 1
assignee: Alex Cabrera
tags: [manifest, dev, templates]
---
# Template-defined dev messages

Templates should be able to register custom messages shown during `dade dev`. Non-web templates (iOS, Android, CLI, TUI) shouldn't show HTTPS proxy URLs or run proxy setup. Each template controls what the user sees.

## Design

Add a `[dev.messages]` section to dade.toml that lets templates define: `ready` (shown when dev server starts), `running` (shown while running). A new `serve.proxy` boolean (default true for web, false when default_port=0) controls whether proxy setup runs. The dev/start commands read these fields and show template-specific output instead of hardcoded HTTPS URLs.

## Acceptance Criteria

1. iOS/Android/CLI/TUI templates show no HTTPS URL or proxy output during dev. 2. Web templates continue showing HTTPS URL as before. 3. Templates can define custom ready messages. 4. All existing tests pass. 5. Changes pushed to all repos.

