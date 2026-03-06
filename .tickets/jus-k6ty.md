---
id: jus-k6ty
status: open
deps: [jus-8hqk]
links: []
created: 2026-02-12T20:36:52Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-ju43
tags: [github, deprecation]
---
# Update srv repository with deprecation notice

Update srv repository README with deprecation notice.

## New README.md

```markdown
# \u26a0\ufe0f srv has been merged into dade

This project's functionality is now part of [dade](https://github.com/theydontwantyoutovibecode/dade).

## Migration

\`\`\`bash
# Uninstall srv
brew uninstall srv

# Install dade
brew tap theydontwantyoutovibecode/tap
brew install dade

# Run setup (migrates your projects)
dade setup
\`\`\`

## Command Mapping

| srv | dade |
|-----|-----------|
| srv new | dade new |
| srv start | dade start |
| srv stop | dade stop |
| srv list | dade list |
| srv open | dade open |
| srv tunnel | dade tunnel |
| srv proxy | dade proxy |
| srv register | dade register |
| srv remove | dade remove |
| srv sync | dade sync |

## What's New in dade

- Template plugin system (install templates from git)
- Framework-aware serving (Django, Node, etc.)
- Better scaffolding

## Questions?

[Open an issue on dade](https://github.com/theydontwantyoutovibecode/dade/issues)
\`\`\`
```

## Acceptance Criteria

- [ ] README replaced with deprecation notice
- [ ] Migration steps documented
- [ ] Command mapping shown
- [ ] Links to dade work

