---
id: jus-3mcd
status: closed
deps: [jus-ii1r]
links: []
created: 2026-02-12T20:27:50Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-vt7v
tags: [infrastructure, registry]
---
# Implement register_project() function

Implement the function that adds a project to the registry.

## Function Implementation

```bash
register_project() {
    local name="$1"
    local port="$2"
    local path="$3"
    local template="${4:-unknown}"
    local created=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    
    local tmp
    tmp=$(mktemp)
    
    jq --arg n "$name" \
       --argjson p "$port" \
       --arg path "$path" \
       --arg t "$template" \
       --arg c "$created" \
       '.[$n] = {"port": $p, "path": $path, "template": $t, "created": $c}' \
       "$DADE_PROJECTS_FILE" > "$tmp"
    
    mv "$tmp" "$DADE_PROJECTS_FILE"
}
```

## Schema

```json
{
  "myproject": {
    "port": 3000,
    "path": "/Users/alex/Code/myproject",
    "template": "django-hypermedia",
    "created": "2024-01-15T10:30:00Z"
  }
}
```

## Acceptance Criteria

- [ ] Adds project to projects.json
- [ ] Includes port, path, template, created
- [ ] Uses atomic write (temp file + mv)
- [ ] Overwrites existing entry with same name

