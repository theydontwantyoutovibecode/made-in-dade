---
id: jus-ibzm
status: open
deps: []
links: []
created: 2026-03-02T04:58:21Z
type: task
priority: 1
assignee: Alex Cabrera
tags: [cli, restructure]
---
# Remove 'open' command, add --open flag to 'dev'

The 'open' command is only useful in the context of dev. Remove it as a top-level command and add an '--open' flag to 'dev' that opens the project URL in the browser after the dev server starts.

Current → Proposed:
  dade open [name]      → removed
  dade dev [name]       → dade dev [name] --open

Implementation:
- Add --open flag to devCmd in cmd_dev.go
- After dev server starts and proxy is confirmed, call openBrowserFunc if --open is set
- Remove cmd_open.go
- Keep open_test.go tests adapted for the flag
- Add hidden compat alias 'open' that prints a deprecation hint (covered by jus-7z4w)

Final top-level surface becomes:
  new, dev, build, share + project, template, proxy, setup = 8 entries

