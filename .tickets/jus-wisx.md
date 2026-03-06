---
id: jus-wisx
status: closed
deps: []
links: []
created: 2026-03-01T22:07:08Z
type: feature
priority: 1
assignee: Alex Cabrera
tags: [windows, proxy, cross-platform]
---
# Add Windows service management for Caddy proxy

Create Windows equivalent of launchd.go for managing the Caddy proxy service. Options include Windows Task Scheduler or NSSM (Non-Sucking Service Manager).

## Acceptance Criteria

- Create internal/proxy/windows.go with build tags
- Implement InstallProxyService, UninstallProxyService, RestartProxyService for Windows
- Use NSSM or schtasks for service management
- Test proxy start/stop on Windows

