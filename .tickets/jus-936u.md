---
id: jus-936u
status: open
deps: []
links: []
created: 2026-03-02T04:52:27Z
type: task
priority: 1
assignee: Alex Cabrera
tags: [cli, cleanup]
---
# Remove cobra's auto-generated completion command

Set CompletionOptions.DisableDefaultCmd = true on the root command. The completion command adds noise and we have no plans to support shell completions.

