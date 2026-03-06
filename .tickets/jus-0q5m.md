---
id: jus-0q5m
status: closed
deps: [jus-jcvx]
links: []
created: 2026-02-12T20:30:56Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-q0si
tags: [templates, project]
---
# Implement project marker read functions

Implement helper functions for reading project marker files.

## Functions

```bash
is_dade_project() {
    [[ -f ".dade" ]]
}

get_current_project_name() {
    jq -r '.name // empty' .dade 2>/dev/null
}

get_current_project_port() {
    jq -r '.port // empty' .dade 2>/dev/null
}

get_current_project_template() {
    jq -r '.template // empty' .dade 2>/dev/null
}

read_project_marker() {
    local dir="${1:-.}"
    if [[ -f "$dir/.dade" ]]; then
        cat "$dir/.dade"
    else
        return 1
    fi
}
```

## Acceptance Criteria

- [ ] is_dade_project detects marker
- [ ] get_current_project_name returns name
- [ ] get_current_project_port returns port
- [ ] get_current_project_template returns template
- [ ] Functions work from project directory

