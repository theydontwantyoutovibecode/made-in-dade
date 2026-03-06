---
id: jus-3ii0
status: closed
deps: []
links: []
created: 2026-02-12T20:31:59Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-ocjp
tags: [serving, static]
---
# Implement start_static_server() function

Implement the function that starts a Caddy file server for static projects.

## Function Implementation

```bash
start_static_server() {
    local port="$1"
    local root="${2:-.}"
    
    log_info "Starting static server on port $port..."
    
    # Start Caddy file-server in background
    caddy file-server --listen ":$port" --root "$root" &
    local pid=$!
    
    # Save PID
    echo "$pid" > .dade.pid
    
    # Wait briefly to check if it started
    sleep 0.5
    
    if kill -0 "$pid" 2>/dev/null; then
        return 0
    else
        log_error "Failed to start server"
        rm -f .dade.pid
        return 1
    fi
}
```

## Considerations

- Uses Caddy's file-server for consistency with proxy
- Runs in background
- PID saved for later stop
- Brief wait to verify startup

## Acceptance Criteria

- [ ] Starts Caddy file-server
- [ ] Runs in background
- [ ] Saves PID to .dade.pid
- [ ] Verifies process started
- [ ] Returns success/failure

