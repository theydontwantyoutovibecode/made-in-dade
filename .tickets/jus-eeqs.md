---
id: jus-eeqs
status: closed
deps: [jus-1gl4]
links: []
created: 2026-02-12T20:34:06Z
type: task
priority: 3
assignee: Alex Cabrera
parent: jus-17o8
tags: [commands, serving]
---
# Implement cmd_sync() command handler

Implement the sync command for rebuilding the registry.

## Function Implementation

```bash
cmd_sync() {
    local scan_path="$HOME"
    local clean_mode=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --clean) clean_mode=true; shift ;;
            *) scan_path="$1"; shift ;;
        esac
    done
    
    init_config
    
    if $clean_mode; then
        log_info "Cleaning stale registry entries..."
        local removed=0
        
        while IFS='|' read -r name path; do
            [[ -z "$name" ]] && continue
            if [[ ! -d "$path" ]] || [[ ! -f "$path/.dade" ]]; then
                log_warn "Removing stale: $name"
                unregister_project "$name"
                ((removed++)) || true
            fi
        done < <(jq -r 'to_entries[] | "\(.key)|\(.value.path)"' "$DADE_PROJECTS_FILE" 2>/dev/null)
        
        generate_caddyfile
        reload_proxy
        
        log_success "Removed $removed stale entries"
        return 0
    fi
    
    log_info "Scanning for projects in $scan_path..."
    
    # Clear registry
    echo '{}' > "$DADE_PROJECTS_FILE"
    
    local count=0
    
    while IFS= read -r marker; do
        [[ -z "$marker" ]] && continue
        local dir=$(dirname "$marker")
        
        local name=$(jq -r '.name // empty' "$marker" 2>/dev/null)
        local port=$(jq -r '.port // empty' "$marker" 2>/dev/null)
        local template=$(jq -r '.template // empty' "$marker" 2>/dev/null)
        
        if [[ -n "$name" ]] && [[ -n "$port" ]]; then
            register_project "$name" "$port" "$dir" "$template"
            log_success "Found: $name"
            ((count++)) || true
        fi
    done < <(find "$scan_path" -name ".dade" -type f 2>/dev/null | grep -v "/.config/" | head -100)
    
    generate_caddyfile
    reload_proxy
    
    log_success "Synced $count project(s)"
}
```

## Acceptance Criteria

- [ ] Scans for .dade files
- [ ] Rebuilds registry from found projects
- [ ] --clean removes stale entries
- [ ] Updates proxy config

