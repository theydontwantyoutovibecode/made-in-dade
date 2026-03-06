---
id: jus-c6us
status: closed
deps: [jus-f3wf]
links: []
created: 2026-02-12T20:35:12Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-zf02
tags: [templates, documentation]
---
# Update hypertext template AGENTS.md

Update the hypertext template's AGENTS.md to document dade commands.

## Changes

Update Development section:

```markdown
## Development

### With dade (recommended)

\`\`\`bash
dade start    # Start server
dade open     # Open in browser
dade stop     # Stop server
dade tunnel   # Share publicly
\`\`\`

Your site is at https://<project-name>.localhost

### Without dade

\`\`\`bash
python3 -m http.server 8000
# Visit http://localhost:8000
\`\`\`
```

## Acceptance Criteria

- [ ] dade commands documented
- [ ] Standalone fallback documented
- [ ] Old serve.sh references removed

