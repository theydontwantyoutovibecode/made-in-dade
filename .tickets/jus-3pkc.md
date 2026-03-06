---
id: jus-3pkc
status: closed
deps: [jus-07t4]
links: []
created: 2026-02-12T20:30:32Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-fqnp
tags: [commands, templates]
---
# Implement cmd_uninstall() command handler

Implement the uninstall command for removing template plugins.

## Function Implementation

```bash
cmd_uninstall() {
    local name="$1"
    
    if [[ -z "$name" ]]; then
        log_error "Usage: dade uninstall <template-name>"
        exit 1
    fi
    
    if [[ "$name" == "--all" ]]; then
        if confirm "Remove ALL installed templates?"; then
            rm -rf "$DADE_TEMPLATES_DIR"/*
            log_success "All templates removed"
        fi
        return 0
    fi
    
    local template_dir="$DADE_TEMPLATES_DIR/$name"
    
    if [[ ! -d "$template_dir" ]]; then
        log_error "Template '$name' not found"
        echo ""
        log_info "Installed templates:"
        for dir in "$DADE_TEMPLATES_DIR"/*/; do
            [[ -d "$dir" ]] && echo "  - $(basename "$dir")"
        done
        exit 1
    fi
    
    # Check if projects use this template
    local using_projects
    using_projects=$(jq -r --arg t "$name" \
        'to_entries[] | select(.value.template == $t) | .key' \
        "$DADE_PROJECTS_FILE" 2>/dev/null || true)
    
    if [[ -n "$using_projects" ]]; then
        log_warn "Projects using this template:"
        echo "$using_projects" | while read -r proj; do
            [[ -n "$proj" ]] && echo "  - $proj"
        done
        echo ""
        if ! confirm "Uninstall anyway? (projects will continue to work)"; then
            return 0
        fi
    fi
    
    rm -rf "$template_dir"
    log_success "Removed template: $name"
}
```

## Acceptance Criteria

- [ ] Removes template directory
- [ ] Warns about projects using template
- [ ] Shows error if not found
- [ ] --all removes all with confirmation

