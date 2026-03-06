---
id: jus-fqnp
status: closed
deps: [jus-ae77]
links: []
created: 2026-02-12T19:56:21Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-rg8c
tags: [commands, templates]
---
# Implement 'dade uninstall' command

Implement the uninstall command for removing installed template plugins.

## Usage

```
dade uninstall <name>   # Remove template
dade uninstall --all    # Remove all templates (with confirmation)
```

## Implementation

```bash
cmd_uninstall() {
    local name="$1"
    
    if [[ "$name" == "--all" ]]; then
        if confirm "Remove ALL installed templates?"; then
            rm -rf "$DADE_TEMPLATES_DIR"/*
            log_success "All templates removed"
        fi
        return
    fi
    
    local template_dir="$DADE_TEMPLATES_DIR/$name"
    
    if [[ ! -d "$template_dir" ]]; then
        log_error "Template '$name' not found"
        log_info "Installed templates:"
        cmd_templates
        exit 1
    fi
    
    # Check if any projects use this template
    local using_projects
    using_projects=$(jq -r --arg t "$name" \
        'to_entries[] | select(.value.template == $t) | .key' \
        "$DADE_PROJECTS_FILE" 2>/dev/null)
    
    if [[ -n "$using_projects" ]]; then
        log_warn "The following projects use this template:"
        echo "$using_projects" | while read -r proj; do
            echo "  - $proj"
        done
        if ! confirm "Uninstall anyway? (projects will still work)"; then
            exit 0
        fi
    fi
    
    rm -rf "$template_dir"
    log_success "Removed template: $name"
}
```

## Considerations

- Warn if projects depend on the template
- Projects continue working (they have the code, just lose update capability)
- Tab completion for template names

## Acceptance Criteria

- [ ] Template directory removed
- [ ] Warning shown if projects use template
- [ ] Error if template not found
- [ ] --all removes all with confirmation

