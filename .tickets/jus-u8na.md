---
id: jus-u8na
status: closed
deps: [jus-ovh7]
links: []
created: 2026-02-12T19:56:32Z
type: task
priority: 3
assignee: Alex Cabrera
parent: jus-rg8c
tags: [commands, templates]
---
# Implement 'dade update' command

Implement the update command for updating installed template plugins.

## Usage

```
dade update <name>   # Update specific template
dade update --all    # Update all templates
```

## Implementation

```bash
cmd_update() {
    local name="$1"
    
    if [[ "$name" == "--all" ]]; then
        for dir in "$DADE_TEMPLATES_DIR"/*/; do
            local tname=$(basename "$dir")
            update_template "$tname"
        done
        return
    fi
    
    update_template "$name"
}

update_template() {
    local name="$1"
    local template_dir="$DADE_TEMPLATES_DIR/$name"
    
    if [[ ! -d "$template_dir" ]]; then
        log_error "Template '$name' not found"
        exit 1
    fi
    
    local source_url
    source_url=$(cat "$template_dir/.source" 2>/dev/null)
    
    if [[ -z "$source_url" ]]; then
        log_error "No source URL stored for '$name'"
        log_info "Re-install with: dade install <url> --name $name"
        exit 1
    fi
    
    # Clone to temp and replace
    local tmp_dir
    tmp_dir=$(mktemp -d)
    
    spin "Updating $name" git clone --depth 1 "$source_url" "$tmp_dir"
    
    # Preserve .source file
    cp "$template_dir/.source" "$tmp_dir/.source"
    
    # Replace
    rm -rf "$template_dir"
    mv "$tmp_dir" "$template_dir"
    
    log_success "Updated: $name"
}
```

## Version Tracking (Future Enhancement)

Could track installed version and show diff:
- Store template.version in cache
- Compare with new version after clone
- Show changelog if available

## Acceptance Criteria

- [ ] Single template updates via git pull
- [ ] --all updates all templates
- [ ] Source URL required (error if missing)
- [ ] .source file preserved after update

