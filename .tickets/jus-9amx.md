---
id: jus-9amx
status: closed
deps: [jus-nmo3]
links: []
created: 2026-03-02T01:59:01Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-g3fm
tags: [read-only, templates, ws-7]
---
# Create curated .read-only/manifest.txt for each template

Create a .read-only/manifest.txt in each template with relevant reference repos. web-site: htmx repo, example static sites. web-app: django repo, htmx repo, example django+htmx apps. ios-app: SwiftUI example apps, Apple sample code repos. android-app: Compose samples, Android architecture samples. cli: charm repos (fang, lipgloss, huh), example CLI apps. tui: charm repos (bubbletea, bubbles), example TUI apps.

## Acceptance Criteria

1. All 6 templates have manifest.txt. 2. Listed repos are public and cloneable. 3. manifest.txt has comments explaining each entry.

