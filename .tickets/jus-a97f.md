---
id: jus-a97f
status: closed
deps: [jus-bx3x]
links: []
created: 2026-02-12T20:34:47Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-swo4
tags: [templates, documentation]
---
# Update django template AGENTS.md for dade

Update the Django template's AGENTS.md to document dade commands.

## Changes

Add section for dade development workflow:

```markdown
## Development with dade

If using dade CLI:

\`\`\`bash
# Start development server
dade start

# Open in browser
dade open

# Stop server
dade stop

# Share publicly
dade tunnel
\`\`\`

Your project is accessible at https://<project-name>.localhost

## Standalone Development

Without dade:

\`\`\`bash
./start.sh --dev
\`\`\`
```

## Also Update

- Development Commands section
- Any references to ./start.sh --dev

## Acceptance Criteria

- [ ] dade commands documented
- [ ] Standalone commands still documented
- [ ] URL format explained

