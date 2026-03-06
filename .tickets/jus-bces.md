---
id: jus-bces
status: closed
deps: []
links: []
created: 2026-03-01T22:07:02Z
type: feature
priority: 1
assignee: Alex Cabrera
tags: [windows, setup, cross-platform]
---
# Add Windows support to dade setup command

Extend the setup command to detect Windows and use winget for package installation instead of Homebrew. Required for full Windows support.

## Acceptance Criteria

- Detect OS via runtime.GOOS
- Use winget on Windows for jq, caddy, gum, cloudflared
- Graceful fallback if winget unavailable
- Update checkDependency() to accept winget package IDs
- Test on Windows with PowerShell

