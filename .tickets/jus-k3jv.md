---
id: jus-k3jv
status: closed
deps: [jus-35b7]
links: []
created: 2026-03-02T01:57:11Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-d7n2
tags: [ux, ws-3]
---
# Redesign template selection menu with categories

Update choosePluginTemplate() in new.go to show two sections: Default Templates (with .default marker) and User-Installed Templates. Use Lipgloss styling to make sections visually distinct. Show template name, brief description, and stack summary for each. If no user-installed templates, only show defaults without the section header. Fallback numbered prompt should also show categories.

## Acceptance Criteria

1. Menu shows categorized templates. 2. Styled with Lipgloss. 3. Fallback prompt works. 4. Single-template case still works.

