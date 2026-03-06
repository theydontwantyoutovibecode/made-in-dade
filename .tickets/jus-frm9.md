---
id: jus-frm9
status: closed
deps: []
links: []
created: 2026-03-02T02:59:41Z
type: task
priority: 1
assignee: Alex Cabrera
tags: [docs, ws-9]
---
# Delete obsolete docs/ directory

The docs/ directory contains 3 files that are outdated and no longer useful:
- dev-share-architecture.md - internal design doc from before implementation, references old template names, describes planned features that are now built
- go-migration-spec.md - spec for porting bash to Go, completed long ago, references old template names and bash-era behavior
- release.md - 5-line stub with a build matrix that includes Linux (not supported)

These should be deleted. The README.md is the single source of documentation. Internal architecture details are not needed as standalone docs.

