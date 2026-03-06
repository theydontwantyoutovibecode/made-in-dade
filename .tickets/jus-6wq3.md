---
id: jus-6wq3
status: closed
deps: [jus-y2ku]
links: []
created: 2026-03-02T01:56:39Z
type: task
priority: 0
assignee: Alex Cabrera
parent: jus-uyus
tags: [templates, setup, ws-2]
---
# Auto-install default templates on first run

Add a function ensureDefaultTemplates() that checks if each default template exists in ~/.config/dade/templates/. If missing, shallow-clone from the git URL. Call this from dade setup and as a pre-check in dade new. Should use the existing DefaultTemplates() registry. Show a spinner per template being cloned. Skip templates that are already installed. Caveats: git must be available (already required). Network errors should warn but not block — user can retry with dade setup.

## Acceptance Criteria

1. dade setup clones all missing default templates. 2. dade new calls ensureDefaultTemplates before showing menu. 3. Already-installed templates are skipped. 4. Network failure shows warning, does not crash. 5. Tests cover: all installed, some installed, none installed, git failure.

