---
id: jus-k78p
status: open
deps: [jus-2e3m]
links: []
created: 2026-02-12T19:53:53Z
type: epic
priority: 1
assignee: Alex Cabrera
parent: jus-nq0k
tags: [github, homebrew, distribution]
---
# GitHub & Distribution

Handle all GitHub repository management and distribution updates for the unified dade architecture.

## Repository Changes

### theydontwantyoutovibecode/dade
- Becomes the unified CLI (scaffolding + serving)
- Major version bump to 1.0.0
- New release with updated functionality

### theydontwantyoutovibecode/srv
- Archive repository
- Update README with deprecation notice and redirect to dade
- Set repository description to indicate archived status

### theydontwantyoutovibecode/homebrew-tap
- Update dade formula for new version
- Remove srv formula (or redirect)
- Update SHA256 checksums

## Release Process

1. Complete all code changes in dade
2. Tag release v1.0.0
3. Update Homebrew formula
4. Archive srv repository
5. Announce in README/releases

## Default Templates

Consider shipping default templates with dade:
- Option A: Bundle templates in the CLI repo
- Option B: Auto-install official templates on first run
- Option C: Just document how to install (current approach)

Recommend Option B: On 'dade install' or first 'dade new', offer to install official templates.

## Acceptance Criteria

- [ ] dade v1.0.0 released on GitHub
- [ ] Homebrew formula updated and working
- [ ] srv repository archived with redirect notice
- [ ] Installation works: brew install theydontwantyoutovibecode/tap/dade
- [ ] README documents the unified architecture

