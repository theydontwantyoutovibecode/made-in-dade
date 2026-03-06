---
id: jus-jq4g
status: closed
deps: [jus-ocjp]
links: []
created: 2026-02-12T19:58:38Z
type: task
priority: 3
assignee: Alex Cabrera
parent: jus-6d0x
tags: [commands, serving, cloudflare]
---
# Implement 'dade tunnel' command

Implement the tunnel command for exposing projects via Cloudflare tunnels.

## Usage

```
dade tunnel             # Quick tunnel for current project
dade tunnel <name>      # Quick tunnel for named project
dade tunnel --persist   # Persistent tunnel (requires config)
```

## Quick Tunnel

Uses cloudflared's quick tunnel feature - no account needed:

```bash
cloudflared tunnel --url https://localhost:PORT
```

## Persistent Tunnel

For persistent tunnels with custom subdomains:

1. Requires ~/.config/dade/tunnel.toml:
```toml
[cloudflare]
tunnel_name = "my-tunnel"
domain = "mysite.com"
```

2. Creates tunnel config per project
3. Routes: <project>.<tunnel>.domain

## Implementation

```bash
cmd_tunnel() {
    local project_name=""
    local persist=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --persist) persist=true; shift ;;
            *) project_name="$1"; shift ;;
        esac
    done
    
    if [[ -z "$project_name" ]]; then
        if ! is_dade_project; then
            log_error "Not a dade project directory"
            exit 1
        fi
        project_name=$(get_current_project_name)
    fi
    
    if ! command -v cloudflared &>/dev/null; then
        log_error "cloudflared not installed"
        log_info "Install: brew install cloudflared"
        exit 1
    fi
    
    local port=$(get_project_port "$project_name")
    
    if ! is_project_running "$project_name"; then
        log_error "Project not running"
        log_info "Start first: dade start"
        exit 1
    fi
    
    if $persist; then
        start_persistent_tunnel "$project_name" "$port"
    else
        log_info "Starting quick tunnel for $project_name..."
        log_info "Press Ctrl+C to stop"
        cloudflared tunnel --url "https://localhost:$port"
    fi
}
```

## Dependencies

- cloudflared (optional, checked at runtime)
- For persistent: cloudflare account and tunnel setup

## Acceptance Criteria

- [ ] Quick tunnel works without config
- [ ] Shows public URL when tunnel starts
- [ ] --persist uses configured domain
- [ ] Error if cloudflared not installed
- [ ] Error if project not running

