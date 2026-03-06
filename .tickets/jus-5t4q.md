---
id: jus-5t4q
status: closed
deps: [jus-wz9o]
links: []
created: 2026-02-12T20:01:41Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-k78p
tags: [migration, compatibility]
---
# Implement srv migration in dade setup

Implement migration logic for existing srv users in 'dade setup'.

## Detection

Check for existing srv installation:
- ~/.config/srv/ directory exists
- ~/Library/LaunchAgents/land.charm.srv.proxy.plist exists
- .srv files in projects

## Migration Steps

### 1. Detect srv

```bash
detect_srv_installation() {
    local found=false
    
    if [[ -d "$HOME/.config/srv" ]]; then
        found=true
    fi
    
    if [[ -f "$HOME/Library/LaunchAgents/land.charm.srv.proxy.plist" ]]; then
        found=true
    fi
    
    $found
}
```

### 2. Offer Migration

```bash
if detect_srv_installation; then
    log_info "Existing srv installation detected"
    if confirm "Migrate srv projects to dade?"; then
        migrate_from_srv
    fi
fi
```

### 3. Migration Process

```bash
migrate_from_srv() {
    log_info "Migrating from srv..."
    
    # Stop srv proxy
    if launchctl list | grep -q "land.charm.srv.proxy"; then
        launchctl bootout gui/$(id -u) "land.charm.srv.proxy.plist" 2>/dev/null || true
        log_success "Stopped srv proxy"
    fi
    
    # Copy projects.json
    if [[ -f "$HOME/.config/srv/projects.json" ]]; then
        local srv_projects="$HOME/.config/srv/projects.json"
        
        # For each srv project, create dade entry
        jq -r 'to_entries[] | "\(.key)|\(.value.port)|\(.value.path)"' "$srv_projects" | \
        while IFS='|' read -r name port path; do
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
                
                # Register in dade
                register_project "$name" "$port" "$path" "unknown"
                log_success "Migrated: $name"
            fi
        done
    fi
    
    # Archive old srv config (don't delete)
    if [[ -d "$HOME/.config/srv" ]]; then
        mv "$HOME/.config/srv" "$HOME/.config/srv.backup.$(date +%Y%m%d)"
        log_info "Old config backed up to ~/.config/srv.backup.*"
    fi
    
    # Remove old plist
    if [[ -f "$HOME/Library/LaunchAgents/land.charm.srv.proxy.plist" ]]; then
        rm -f "$HOME/Library/LaunchAgents/land.charm.srv.proxy.plist"
    fi
    
    log_success "Migration complete"
}
```

## Considerations

- Preserve port assignments to avoid breaking bookmarks
- Back up old config, don't delete
- Set template to "unknown" for migrated projects
- Remove old srv launchd service

## Acceptance Criteria

- [ ] srv installation detected
- [ ] User prompted before migration
- [ ] Projects.json entries migrated
- [ ] .srv files converted to .dade
- [ ] Port assignments preserved
- [ ] Old srv proxy stopped
- [ ] Old config backed up

