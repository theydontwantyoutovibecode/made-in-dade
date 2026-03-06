---
id: jus-aa04
status: closed
deps: [jus-7fp1]
links: []
created: 2026-02-12T19:54:48Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-cdr3
tags: [commands, proxy]
---
# Implement 'dade proxy' command

Implement the proxy subcommand for managing the central Caddy proxy.

## Usage

```
dade proxy start     # Start proxy service
dade proxy stop      # Stop proxy service
dade proxy restart   # Restart proxy service
dade proxy status    # Show proxy status
dade proxy logs      # Tail proxy logs
```

## Implementation

```bash
cmd_proxy() {
    local action="${1:-status}"
    
    case "$action" in
        start)
            if is_proxy_running; then
                log_warn "Proxy already running"
            else
                install_proxy_service
                log_success "Proxy started"
            fi
            ;;
        stop)
            uninstall_proxy_service
            log_success "Proxy stopped"
            ;;
        restart)
            uninstall_proxy_service
            sleep 1
            install_proxy_service
            log_success "Proxy restarted"
            ;;
        status)
            if is_proxy_running; then
                log_success "Proxy running"
                # Show registered projects count
            else
                log_warn "Proxy not running"
            fi
            ;;
        logs)
            tail -f "$DADE_CONFIG_DIR/proxy.log"
            ;;
        *)
            log_error "Unknown action: $action"
            ;;
    esac
}
```

## Status Output

Show helpful information:
- Running/stopped status
- Number of registered projects
- Port range in use
- Caddyfile location

## Acceptance Criteria

- [ ] 'dade proxy start' starts the service
- [ ] 'dade proxy stop' stops the service
- [ ] 'dade proxy restart' restarts cleanly
- [ ] 'dade proxy status' shows accurate status
- [ ] 'dade proxy logs' tails log file

