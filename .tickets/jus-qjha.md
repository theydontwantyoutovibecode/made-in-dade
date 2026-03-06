---
id: jus-qjha
status: closed
deps: [jus-b3ad]
links: []
created: 2026-02-12T20:33:27Z
type: task
priority: 3
assignee: Alex Cabrera
parent: jus-jq4g
tags: [commands, cloudflare]
---
# Implement cmd_tunnel() command handler

Implement the tunnel command for Cloudflare tunnels.

## Function Implementation

```bash
cmd_tunnel() {
    local name=""
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            *) name="$1"; shift ;;
        esac
    done
    
    if [[ -z "$name" ]]; then
        if is_dade_project; then
            name=$(get_current_project_name)
        else
            log_error "Not in a project directory"
            exit 1
        fi
    fi
    
    if ! command -v cloudflared &>/dev/null; then
        log_error "cloudflared not installed"
        log_info "Install: brew install cloudflared"
        exit 1
    fi
    
    local port
    if [[ -n "$(get_current_project_port)" ]]; then
        port=$(get_current_project_port)
    else
        port=$(get_project_port "$name")
    fi
    
    if [[ -z "$port" ]]; then
        log_error "Project '$name' not found"
        exit 1
    fi
    
    if ! is_project_running "$name"; then
        log_warn "Project not running locally"
        if confirm "Start $name first?"; then
            if is_dade_project; then
                cmd_start --bg
            else
                local path=$(get_project_path "$name")
                (cd "$path" && cmd_start --bg)
            fi
            sleep 2
        fi
    fi
    
    log_info "Starting tunnel for $name..."
    log_info "Local: https://$name.localhost"
    log_info "Press Ctrl+C to stop tunnel"
    echo ""
    
    cloudflared tunnel --url "https://localhost:$port"
}
```

## Acceptance Criteria

- [ ] Creates quick cloudflare tunnel
- [ ] Checks cloudflared is installed
- [ ] Offers to start project if not running
- [ ] Shows local and public URLs

