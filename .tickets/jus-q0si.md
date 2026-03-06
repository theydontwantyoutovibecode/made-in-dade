---
id: jus-q0si
status: closed
deps: [jus-jngy]
links: []
created: 2026-02-12T19:57:07Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-rg8c
tags: [templates, project]
---
# Implement .dade project marker file

Implement the .dade project marker file that identifies dade-managed projects.

## File Format

JSON format for easy parsing:

```json
{
  "name": "myproject",
  "template": "django-hypermedia",
  "port": 3000,
  "created": "2024-01-15T10:30:00Z"
}
```

## Functions

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

read_project_marker() {
    local project_dir="${1:-.}"
    local marker="$project_dir/.dade"
    
    if [[ ! -f "$marker" ]]; then
        return 1
    fi
    
    cat "$marker"
}

get_current_project_name() {
    jq -r '.name' .dade 2>/dev/null
}

get_current_project_port() {
    jq -r '.port' .dade 2>/dev/null
}

get_current_project_template() {
    jq -r '.template' .dade 2>/dev/null
}

is_dade_project() {
    [[ -f ".dade" ]]
}
```

## Migration from .srv

If .srv file exists (from old srv tool):
1. Read name and port from .srv
2. Create .dade with template="unknown"
3. Optionally remove .srv

## .gitignore

Template should include .dade in .gitignore? Or commit it?

Arguments for committing:
- Project settings are part of project
- Team members get same config

Arguments for ignoring:
- Port might differ per machine
- Personal preference

Recommend: Commit it, but allow port override via environment.

## Acceptance Criteria

- [ ] .dade marker created on 'new'
- [ ] Contains name, template, port, created
- [ ] Helper functions read marker values
- [ ] is_dade_project() detects projects
- [ ] Migration from .srv works

