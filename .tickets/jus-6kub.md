---
id: jus-6kub
status: closed
deps: [jus-yaim]
links: []
created: 2026-03-02T01:59:42Z
type: task
priority: 0
assignee: Alex Cabrera
parent: jus-odye
tags: [checkpoint, testing, ws-6]
---
# CHECKPOINT: Full integration test before publish

Before pushing to the new org, run a full integration test: 1) go test ./... 2) Build dade binary. 3) dade setup (verify deps install). 4) dade new test-web-site (verify template scaffolding). 5) dade dev in the test project (verify it starts). 6) dade new test-ios (verify iOS scaffold). 7) Clean up test projects. Verify all 6 templates work end-to-end.

## Acceptance Criteria

1. All unit tests pass. 2. Binary builds. 3. At least web-site and cli templates scaffold and start correctly. 4. Template selection menu shows all 6 defaults.

