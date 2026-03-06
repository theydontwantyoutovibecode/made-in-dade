---
id: jus-xmk3
status: closed
deps: [jus-6wq3]
links: []
created: 2026-03-02T01:56:46Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-uyus
tags: [templates, ux, ws-2]
---
# Add template alias support

Allow dade new <alias> to resolve template aliases. E.g. dade new website → resolves to web-site template. dade new ios → resolves to ios-app. Add aliases field to manifest [template] section. Add alias resolution to the template lookup in new.go. Aliases: web-site=[website,site], web-app=[webapp], ios-app=[ios], android-app=[android], cli=[], tui=[]. Update manifest parser for the new aliases field.

## Acceptance Criteria

1. dade new website works. 2. dade new ios works. 3. dade new webapp works. 4. Aliases shown in template list. 5. Manifest parser handles aliases field. 6. Tests for alias resolution.

