---
id: jus-mjxa
status: closed
deps: [jus-t35w]
links: []
created: 2026-02-12T20:30:43Z
type: task
priority: 3
assignee: Alex Cabrera
parent: jus-u8na
tags: [commands, templates]
---
# Implement cmd_update() command handler

Implement the update command for updating installed templates.

## Function Implementation

```bash
cmd_update() {
    local name="$1"
    
    if [[ "$name" == "--all" ]]; then
        local count=0
        for dir in "$DADE_TEMPLATES_DIR"/*/; do
            [[ ! -d "$dir" ]] && continue
            local tname=$(basename "$dir")
            if update_template "$tname"; then
                ((count++)) || true
            fi
        done
        log_success "Updated $count template(s)"
        return 0
    fi
    
    if [[ -z "$name" ]]; then
        log_error "Usage: dade update <template-name>"
        log_info "Or: dade update --all"
        exit 1
    fi
    
    update_template "$name"
}

update_template() {
    local name="$1"
    local template_dir="$DADE_TEMPLATES_DIR/$name"
    
    if [[ ! -d "$template_dir" ]]; then
        log_error "Template '$name' not found"
        return 1
    fi
    
    local source_url
    source_url=$(cat "$template_dir/.source" 2>/dev/null || true)
    
    if [[ -z "$source_url" ]]; then
        log_warn "No source URL for '$name' - cannot update"
        return 1
    fi
    
    # Clone to temp
    local tmp_dir
    tmp_dir=$(mktemp -d)
    
    if ! spin "Fetching $name" git clone --depth 1 "$source_url" "$tmp_dir"; then
        log_error "Failed to fetch $name"
        rm -rf "$tmp_dir"
        return 1
    fi
    
    # Preserve source file
    echo "$source_url" > "$tmp_dir/.source"
    
    # Replace
    rm -rf "$template_dir"
    mv "$tmp_dir" "$template_dir"
    
    log_success "Updated: $name"
    return 0
}
```

## Acceptance Criteria

- [ ] Updates single template by name
- [ ] --all updates all templates
- [ ] Preserves source URL
- [ ] Handles missing source URL
- [ ] Reports success/failure

