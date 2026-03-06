---
id: jus-ju43
status: open
deps: [jus-1p4n]
links: []
created: 2026-02-12T20:00:50Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-k78p
tags: [github, deprecation]
---
# Archive srv repository

Archive the srv repository and add deprecation notice.

## Repository

theydontwantyoutovibecode/srv

## Steps

### 1. Update README

Replace entire README with deprecation notice:

```markdown
# ⚠️ This project has been merged into dade

**srv** functionality is now part of [dade](https://github.com/theydontwantyoutovibecode/dade).

## Migration

\`\`\`bash
# Uninstall srv (if installed via Homebrew)
brew uninstall srv

# Install dade
brew tap theydontwantyoutovibecode/tap
brew install dade

# Run setup (migrates existing projects)
dade setup
\`\`\`

## What Changed

- \`srv new\` → \`dade new\`
- \`srv start\` → \`dade start\`
- \`srv stop\` → \`dade stop\`
- \`srv list\` → \`dade list\`
- \`srv tunnel\` → \`dade tunnel\`
- \`srv proxy\` → \`dade proxy\`

dade includes all srv features plus:
- Template plugin system
- Framework-aware serving
- Better scaffolding

## Questions?

Open an issue on the [dade repository](https://github.com/theydontwantyoutovibecode/dade/issues).
\`\`\`

### 2. Update Repository Description

Set description to: "⚠️ DEPRECATED - merged into dade"

### 3. Archive Repository

Via GitHub UI or API:
- Settings → General → Archive this repository

### 4. Update Homebrew Tap

Remove or deprecate srv formula (separate ticket)

## Acceptance Criteria

- [ ] README replaced with deprecation notice
- [ ] Repository description updated
- [ ] Repository archived on GitHub
- [ ] Links to dade work

