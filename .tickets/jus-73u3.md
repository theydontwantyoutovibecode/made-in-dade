---
id: jus-73u3
status: closed
deps: [jus-rtbz]
links: []
created: 2026-02-12T20:28:39Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-aa04
tags: [commands, proxy]
---
# Implement cmd_proxy() command handler

Implement the proxy command handler that dispatches to service control functions.

## Function Implementation

```bash
cmd_proxy() {
    local action="${1:-status}"
    
    case "$action" in
        start)
            start_proxy_service
            ;;
        stop)
            stop_proxy_service
            ;;
        restart)
            restart_proxy_service
            ;;
        status)
            if is_proxy_running; then
                log_success "Proxy running"
                local count=$(jq 'length' "$DADE_PROJECTS_FILE" 2>/dev/null || echo 0)
                log_info "Serving $count project(s)"
            else
                log_warn "Proxy not running"
                log_info "Start with: dade proxy start"
            fi
            ;;
        logs)
            if [[ -f "$DADE_LOG" ]]; then
                tail -f "$DADE_LOG"
            else
                log_warn "No log file found"
            fi
            ;;
        *)
            log_error "Unknown action: $action"
            log_info "Usage: dade proxy [start|stop|restart|status|logs]"
            exit 1
            ;;
    esac
}
```

## Acceptance Criteria

- [ ] start/stop/restart/status/logs work
- [ ] status shows project count
- [ ] logs tails the log file
- [ ] Unknown action shows help

