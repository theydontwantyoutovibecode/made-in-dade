---
id: jus-7hfd
status: closed
deps: [jus-nmo3]
links: []
created: 2026-03-02T01:59:12Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-yi7f
tags: [tk, setup, ws-8]
---
# Install tk via Homebrew in dade setup

Add tk (wedow/ticket) to the list of dependencies installed by dade setup. Currently setup checks for jq, caddy, gum (optional), cloudflared (optional). Add tk as a required dependency: brew install wedow/tap/tk. Update the setup command in cmd_setup.go and setup.go.

## Acceptance Criteria

1. dade setup installs tk. 2. tk available after setup. 3. Tests updated.

