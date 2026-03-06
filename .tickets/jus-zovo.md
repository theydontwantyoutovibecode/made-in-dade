---
id: jus-zovo
status: closed
deps: [jus-qnt6]
links: []
created: 2026-02-12T20:36:30Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-5t4q
tags: [migration, compatibility]
---
# Implement migrate_from_srv() function

Implement the function that migrates from srv to dade.

## Function Implementation

```bash
detect_srv_installation() {
    [[ -d "$HOME/.config/srv" ]] || \
    [[ -f "$HOME/Library/LaunchAgents/land.charm.srv.proxy.plist" ]]
}

migrate_from_srv() {
    log_info "Migrating from srv..."
    
    # Stop srv proxy
    if launchctl list 2>/dev/null | grep -q "land.charm.srv.proxy"; then
        launchctl bootout gui/$(id -u) "$HOME/Library/LaunchAgents/land.charm.srv.proxy.plist" 2>/dev/null || true
        log_success "Stopped srv proxy"
    fi
    
    # Migrate projects
    local srv_projects="$HOME/.config/srv/projects.json"
    if [[ -f "$srv_projects" ]]; then
        local count=0
        
        while IFS='|' read -r name port path; do
            [[ -z "$name" ]] && continue
            
            if [[ -d "$path" ]]; then
                # Convert .srv to .dade
                if [[ -f "$path/.srv" ]]; then
                    cat > "$path/.dade" << EOF
{
  "name": "$name",
  "template": "unknown",
  "port": $port,
  "migrated_from": "srv"
}
EOF
                    rm -f "$path/.srv"
                fi
                
                register_project "$name" "$port" "$path" "unknown"
                log_success "Migrated: $name"
                ((count++)) || true
            fi
        done < <(jq -r 'to_entries[] | "\(.key)|\(.value.port)|\(.value.path)"' "$srv_projects" 2>/dev/null)
        
        log_success "Migrated $count project(s)"
    fi
    
    # Backup old config
    if [[ -d "$HOME/.config/srv" ]]; then
        local backup="$HOME/.config/srv.backup.$(date +%Y%m%d%H%M%S)"
        mv "$HOME/.config/srv" "$backup"
        log_info "Old config backed up to $backup"
    fi
    
    # Remove old plist
    rm -f "$HOME/Library/LaunchAgents/land.charm.srv.proxy.plist"
    
    log_success "Migration complete"
}
```

## Acceptance Criteria

- [ ] Detects srv installation
- [ ] Stops srv proxy
- [ ] Migrates projects with port preservation
- [ ] Converts .srv to .dade
- [ ] Backs up old config

