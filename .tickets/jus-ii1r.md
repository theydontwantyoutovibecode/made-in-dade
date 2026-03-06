---
id: jus-ii1r
status: closed
deps: [jus-wy9h]
links: []
created: 2026-02-12T20:27:41Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-vt7v
tags: [infrastructure, registry]
---
# Implement next_port() function

Implement the function that finds the next available port for a new project.

## Function Implementation

```bash
next_port() {
    local max_port=$DADE_BASE_PORT
    
    if [[ -f "$DADE_PROJECTS_FILE" ]]; then
        local ports
        ports=$(jq -r 'to_entries[].value.port // empty' "$DADE_PROJECTS_FILE" 2>/dev/null || echo "")
        
        while read -r p; do
            if [[ -n "$p" ]] && (( p >= max_port )); then
                max_port=$((p + 1))
            fi
        done <<< "$ports"
    fi
    
    echo "$max_port"
}
```

## Considerations

- Start from BASE_PORT (3000)
- Find maximum used port and add 1
- Handle empty projects file
- Handle malformed JSON gracefully

## Acceptance Criteria

- [ ] Returns 3000 for empty registry
- [ ] Returns max+1 for populated registry
- [ ] Handles empty/invalid JSON
- [ ] Works with gaps in port numbers

