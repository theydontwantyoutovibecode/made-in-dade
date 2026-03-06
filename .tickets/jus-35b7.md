---
id: jus-35b7
status: closed
deps: [jus-6wq3]
links: []
created: 2026-03-02T01:56:52Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-uyus
tags: [templates, ws-2]
---
# Mark default vs user-installed templates in registry

Templates installed by ensureDefaultTemplates should be distinguishable from user-installed ones. Write a .default marker file alongside .source in the template directory. The dade new menu and dade templates list should show this distinction. This enables showing Default Templates and User-Installed Templates as separate sections.

## Acceptance Criteria

1. Default templates have .default marker. 2. User-installed templates do not. 3. loadPluginTemplates returns isDefault flag. 4. Tests cover both cases.

