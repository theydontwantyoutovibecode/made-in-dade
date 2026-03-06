---
id: jus-rtbz
status: closed
deps: [jus-qg1c]
links: []
created: 2026-02-12T20:28:31Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-7fp1
tags: [infrastructure, launchd]
---
# Implement proxy service control functions

Implement functions to start, stop, and check the proxy service.

## Functions

```bash
start_proxy_service() {
    if is_proxy_running; then
        log_warn "Proxy already running"
        return 0
    fi
    
    # Ensure plist exists
    create_plist
    
    # Bootstrap the service
    launchctl bootstrap gui/$(id -u) "$DADE_PLIST"
    
    sleep 1
    
    if is_proxy_running; then
        log_success "Proxy started"
    else
        log_error "Failed to start proxy"
        return 1
    fi
}

stop_proxy_service() {
    if ! is_proxy_running; then
        log_warn "Proxy not running"
        return 0
    fi
    
    launchctl bootout gui/$(id -u) "$DADE_PLIST" 2>/dev/null || true
    log_success "Proxy stopped"
}

restart_proxy_service() {
    stop_proxy_service
    sleep 1
    start_proxy_service
}

is_proxy_running() {
    launchctl list 2>/dev/null | grep -q "$DADE_PROXY_LABEL"
}
```

## Acceptance Criteria

- [ ] start_proxy_service bootstraps launchd job
- [ ] stop_proxy_service bootouts launchd job
- [ ] restart_proxy_service cycles the service
- [ ] is_proxy_running detects running state

