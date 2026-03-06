---
id: jus-yf8t
status: closed
deps: []
links: []
created: 2026-03-02T01:56:07Z
type: task
priority: 0
assignee: Alex Cabrera
parent: jus-u0zj
tags: [cleanup, windows, ws-1]
---
# Drop Windows support from all templates

Remove all Windows-specific files and configuration from every template and from the dade binary. dade is macOS-only. Files to remove: setup.ps1 from cli, tui, android, mobile templates. Fields to remove from dade.toml: setup_windows, command_windows, dev_windows, setup_windows in prod section. Fields to remove from cmd_build.go: Windows-specific cross-compilation logic, command_windows handling in buildFromManifest. Close tickets: jus-bces, jus-wisx, jus-dnxe.

## Acceptance Criteria

1. No setup.ps1 files exist in any template. 2. No *_windows fields in any dade.toml. 3. cmd_build.go has no Windows-specific code paths. 4. manifest.go has no Windows-specific fields. 5. Windows tickets closed. 6. All tests pass.

