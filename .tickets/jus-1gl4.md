---
id: jus-1gl4
status: closed
deps: [jus-3mcd]
links: []
created: 2026-02-12T20:28:01Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-vt7v
tags: [infrastructure, registry]
---
# Implement registry lookup functions

Implement helper functions for looking up project information.

## Functions

```bash
get_project_port() {
    local name="$1"
    jq -r --arg n "$name" '.[$n].port // empty' "$DADE_PROJECTS_FILE"
}

get_project_path() {
    local name="$1"
    jq -r --arg n "$name" '.[$n].path // empty' "$DADE_PROJECTS_FILE"
}

get_project_template() {
    local name="$1"
    jq -r --arg n "$name" '.[$n].template // empty' "$DADE_PROJECTS_FILE"
}

project_exists() {
    local name="$1"
    [[ -n $(get_project_port "$name") ]]
}

list_project_names() {
    jq -r 'keys[]' "$DADE_PROJECTS_FILE" 2>/dev/null
}
```

## Acceptance Criteria

- [ ] get_project_port returns port or empty
- [ ] get_project_path returns path or empty
- [ ] get_project_template returns template or empty
- [ ] project_exists returns 0/1 correctly
- [ ] list_project_names returns all names

