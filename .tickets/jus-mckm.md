---
id: jus-mckm
status: closed
deps: [jus-bxd6]
links: []
created: 2026-02-12T20:32:46Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-l4vb
tags: [serving, process]
---
# Implement stop_project() function

Implement the function that stops a project's server.

## Function Implementation

```bash
stop_project() {
    local name="${1:-}"
    local project_dir
    
    if [[ -z "$name" ]]; then
        project_dir="."
        name=$(get_current_project_name)
    else
        project_dir=$(get_project_path "$name")
    fi
    
    if [[ -z "$name" ]]; then
        log_error "Project not found"
        return 1
    fi
    
    local stopped=false
    
    # Try PID file first
    local pid_file="$project_dir/.dade.pid"
    if [[ -f "$pid_file" ]]; then
        local pid=$(cat "$pid_file")
        if kill -0 "$pid" 2>/dev/null; then
            graceful_kill "$pid"
            stopped=true
        fi
        rm -f "$pid_file"
    fi
    
    # Also try by port
    if ! $stopped; then
        local port
        if [[ -z "${1:-}" ]]; then
            port=$(get_current_project_port)
        else
            port=$(get_project_port "$name")
        fi
        
        if [[ -n "$port" ]]; then
            local pid=$(lsof -ti ":$port" 2>/dev/null || true)
            if [[ -n "$pid" ]]; then
                graceful_kill "$pid"
                stopped=true
            fi
        fi
    fi
    
    if $stopped; then
        log_success "Stopped $name"
    else
        log_warn "$name was not running"
    fi
}

graceful_kill() {
    local pid="$1"
    local timeout="${2:-3}"
    
    # Send SIGTERM
    kill -TERM "$pid" 2>/dev/null || return 0
    
    # Wait for graceful shutdown
    for ((i=0; i<timeout; i++)); do
        if ! kill -0 "$pid" 2>/dev/null; then
            return 0
        fi
        sleep 1
    done
    
    # Force kill if still running
    kill -KILL "$pid" 2>/dev/null || true
}
```

## Acceptance Criteria

- [ ] Stops process by PID file
- [ ] Falls back to port-based lookup
- [ ] Graceful shutdown (SIGTERM first)
- [ ] Force kill after timeout
- [ ] Cleans up PID file

