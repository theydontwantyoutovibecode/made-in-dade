---
id: jus-l4vb
status: closed
deps: [jus-ocjp]
links: []
created: 2026-02-12T19:57:54Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-6d0x
tags: [commands, serving]
---
# Implement 'dade stop' command

Implement the stop command for stopping running project servers.

## Usage

```
dade stop           # Stop current project
dade stop <name>    # Stop specific project
dade stop --all     # Stop all running projects
```

## Implementation

```bash
cmd_stop() {
    local project_name="$1"
    
    if [[ "$project_name" == "--all" ]]; then
        stop_all_projects
        return
    fi
    
    if [[ -z "$project_name" ]]; then
        if ! is_dade_project; then
            log_error "Not a dade project directory"
            exit 1
        fi
        project_name=$(get_current_project_name)
    fi
    
    stop_project "$project_name"
}

stop_project() {
    local name="$1"
    local project_path=$(get_project_path "$name")
    
    if [[ -z "$project_path" ]]; then
        log_error "Project '$name' not found"
        exit 1
    fi
    
    # Try PID file first
    local pid_file="$project_path/.dade.pid"
    if [[ -f "$pid_file" ]]; then
        local pid=$(cat "$pid_file")
        if kill -0 "$pid" 2>/dev/null; then
            kill "$pid"
            log_success "Stopped $name (PID: $pid)"
        fi
        rm -f "$pid_file"
    else
        # Fallback: find process by port
        local port=$(get_project_port_from_registry "$name")
        local pid=$(lsof -ti ":$port" 2>/dev/null || true)
        if [[ -n "$pid" ]]; then
            kill "$pid" 2>/dev/null || true
            log_success "Stopped $name (port $port)"
        else
            log_warn "Project $name not running"
        fi
    fi
    
    # Update running registry
    untrack_running_process "$name"
}

stop_all_projects() {
    local count=0
    
    # Iterate running projects
    for project in $(list_running_projects); do
        stop_project "$project"
        ((count++)) || true
    done
    
    if [[ $count -eq 0 ]]; then
        log_info "No projects running"
    else
        log_success "Stopped $count project(s)"
    fi
}
```

## Graceful Shutdown

1. Send SIGTERM first
2. Wait briefly
3. Send SIGKILL if still running

```bash
graceful_kill() {
    local pid="$1"
    local timeout="${2:-5}"
    
    kill -TERM "$pid" 2>/dev/null || return 0
    
    for ((i=0; i<timeout; i++)); do
        if ! kill -0 "$pid" 2>/dev/null; then
            return 0
        fi
        sleep 1
    done
    
    kill -KILL "$pid" 2>/dev/null || true
}
```

## Acceptance Criteria

- [ ] Current project stopped
- [ ] Named project stopped
- [ ] --all stops all running projects
- [ ] PID file cleaned up
- [ ] Graceful shutdown attempted
- [ ] Port-based fallback works

