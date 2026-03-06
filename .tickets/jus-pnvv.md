---
id: jus-pnvv
status: closed
deps: [jus-3mcd]
links: []
created: 2026-02-12T20:27:55Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-vt7v
tags: [infrastructure, registry]
---
# Implement unregister_project() function

Implement the function that removes a project from the registry.

## Function Implementation

```bash
unregister_project() {
    local name="$1"
    
    local tmp
    tmp=$(mktemp)
    
    jq --arg n "$name" 'del(.[$n])' "$DADE_PROJECTS_FILE" > "$tmp"
    mv "$tmp" "$DADE_PROJECTS_FILE"
}
```

## Acceptance Criteria

- [ ] Removes project from projects.json
- [ ] Uses atomic write
- [ ] No error if project doesn't exist

