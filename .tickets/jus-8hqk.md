---
id: jus-8hqk
status: closed
deps: [jus-hrzu, jus-thwh, jus-ji80, jus-9dog]
links: []
created: 2026-02-12T20:37:08Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-1p4n
tags: [github, release]
---
# Create v1.0.0 git tag and release

Create the v1.0.0 git tag and GitHub release.

## Prerequisites

- All implementation tickets complete
- All tests passing
- README updated
- AGENTS.md created

## Steps

```bash
cd /Users/acabrera/Code/dade/dade

# Ensure version is set
grep -q 'DADE_VERSION="1.0.0"' dade

# Commit any final changes
git add -A
git commit -m "Prepare v1.0.0 release"

# Create annotated tag
git tag -a v1.0.0 -m "dade v1.0.0 - Unified CLI

Major release merging srv functionality:
- Template plugin system
- Automatic HTTPS via Caddy
- Framework-aware serving
- Cloudflare tunnel support
- Project registry and management
"

# Push tag
git push origin main --tags
```

## GitHub Release

Create release via GitHub UI with:
- Title: dade v1.0.0
- Body: Release notes (features, migration from srv, install instructions)
- Attach: nothing (Homebrew uses source tarball)

## Acceptance Criteria

- [ ] Version set to 1.0.0 in script
- [ ] Git tag v1.0.0 created
- [ ] Tag pushed to GitHub
- [ ] GitHub release created with notes

