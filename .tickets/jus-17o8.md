---
id: jus-17o8
status: closed
deps: [jus-vt7v]
links: []
created: 2026-02-12T19:59:21Z
type: task
priority: 3
assignee: Alex Cabrera
parent: jus-6d0x
tags: [commands, serving]
---
# Implement 'dade sync' command

Implement the sync command for rebuilding the project registry by scanning for .dade files.

## Usage

```
dade sync              # Scan home directory
dade sync <path>       # Scan specific path
dade sync --clean      # Remove entries for missing directories
```

## Use Case

- Registry got corrupted
- Moved projects around
- Migrating from another machine
- Cleaning up stale entries

## Implementation

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
    
    if $clean_mode; then
        sync_clean
        return
    fi
    
    log_info "Scanning for .dade files in $scan_path..."
    
    # Clear existing registry
    echo '{}' > "$DADE_PROJECTS_FILE"
    
    local count=0
    
    # Find all .dade files
    while IFS= read -r marker_file; do
        local project_dir=$(dirname "$marker_file")
        
        local name=$(jq -r '.name // empty' "$marker_file" 2>/dev/null)
        local port=$(jq -r '.port // empty' "$marker_file" 2>/dev/null)
        local template=$(jq -r '.template // empty' "$marker_file" 2>/dev/null)
        
        if [[ -n "$name" ]] && [[ -n "$port" ]]; then
            register_project "$name" "$port" "$project_dir" "$template"
            log_success "Found: $name ($project_dir)"
            ((count++)) || true
        fi
    done < <(find "$scan_path" -name ".dade" -type f 2>/dev/null | grep -v "/.config/" | head -100)
    
    # Regenerate Caddyfile
    generate_caddyfile
    reload_proxy
    
    log_success "Synced $count project(s)"
}

sync_clean() {
    log_info "Cleaning stale registry entries..."
    
    local removed=0
    local tmp=$(mktemp)
    
    jq -r 'to_entries[] | "\(.key)|\(.value.path)"' "$DADE_PROJECTS_FILE" | \
    while IFS='|' read -r name path; do
        if [[ ! -d "$path" ]] || [[ ! -f "$path/.dade" ]]; then
            log_warn "Removing stale: $name ($path)"
            unregister_project "$name"
            ((removed++)) || true
        fi
    done
    
    generate_caddyfile
    reload_proxy
    
    log_success "Removed $removed stale entries"
}
```

## Performance

- Limit scan depth or file count
- Skip common non-project directories
- Consider caching scan results

## Acceptance Criteria

- [ ] Scans for .dade files
- [ ] Rebuilds registry from found projects
- [ ] --clean removes entries for missing dirs
- [ ] Handles large home directories
- [ ] Proxy updated after sync

