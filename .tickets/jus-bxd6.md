---
id: jus-bxd6
status: closed
deps: []
links: []
created: 2026-02-12T20:32:19Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-ocjp
tags: [serving, status]
---
# Implement is_project_running() function

Implement the function that checks if a project's server is running.

## Function Implementation

```bash
is_project_running() {
    local name="${1:-}"
    local port
    
    if [[ -z "$name" ]]; then
        # Check current directory
        if [[ -f .dade.pid ]]; then
            local pid=$(cat .dade.pid)
            if kill -0 "$pid" 2>/dev/null; then
                return 0
            fi
        fi
        
        # Also check by port
        port=$(get_current_project_port)
    else
        port=$(get_project_port "$name")
    fi
    
    if [[ -n "$port" ]]; then
        lsof -i ":$port" &>/dev/null
        return $?
    fi
    
    return 1
}
```

## Checks

1. First check PID file if exists
2. Then check if port is in use
3. Both methods for robustness

## Acceptance Criteria

- [ ] Returns 0 if running, 1 if not
- [ ] Checks PID file first
- [ ] Falls back to port check
- [ ] Works with name or current directory

