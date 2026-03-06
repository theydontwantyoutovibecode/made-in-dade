---
id: jus-jcvx
status: closed
deps: []
links: []
created: 2026-02-12T20:30:50Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-q0si
tags: [templates, project]
---
# Implement write_project_marker() function

Implement the function that writes the .dade project marker file.

## Function Implementation

```bash
write_project_marker() {
    local project_dir="$1"
    local name="$2"
    local template="$3"
    local port="$4"
    local created=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    
    cat > "$project_dir/.dade" << EOF
{
  "name": "$name",
  "template": "$template",
  "port": $port,
  "created": "$created"
}
EOF
}
```

## File Format

```json
{
  "name": "myproject",
  "template": "django-hypermedia",
  "port": 3000,
  "created": "2024-01-15T10:30:00Z"
}
```

## Acceptance Criteria

- [ ] Creates .dade file
- [ ] Includes name, template, port, created
- [ ] Uses ISO 8601 date format
- [ ] Valid JSON output

