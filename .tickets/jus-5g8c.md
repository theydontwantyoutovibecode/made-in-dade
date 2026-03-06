---
id: jus-5g8c
status: closed
deps: []
links: []
created: 2026-02-17T02:06:13Z
type: task
priority: 3
assignee: Alex Cabrera
parent: jus-gyx3
---
# Package/release guidance for Go binary

Define release packaging for Go: build instructions (darwin/arm64, amd64, linux), versioning tied to DADE_VERSION, binary naming, checksum generation, and optional Homebrew tap formula steps. Document in docs/release.md. No actual publishing required yet.

## Acceptance Criteria

- docs/release.md (or similar) describes build matrix and commands\n- Versioning guidance linked to source constant\n- Notes on checksums and optional brew tap steps\n- No publishing performed

