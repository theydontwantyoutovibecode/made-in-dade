---
id: jus-qjbm
status: closed
deps: [jus-sj8r]
links: []
created: 2026-03-02T15:41:20Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-k85j
tags: [cli, config]
---
# Register mac-app in DefaultTemplates and update tests

Add mac-app entry to config.DefaultTemplates() in internal/config/templates.go. Update default_templates_test.go if it checks count. Add aliases: mac, macos, appkit.


## Notes

**2026-03-02T15:43:30Z**

## Reviewed — test impact

internal/config/templates_test.go has TWO hardcoded assertions:
- Line 13: `len(got.Ordered) != 6` → must change to 7
- Line 29: `len(got.Ordered) != 6` → must change to 7

Aliases should be: mac, macos, desktop (NOT appkit — this is SwiftUI)
