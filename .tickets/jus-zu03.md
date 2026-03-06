---
id: jus-zu03
status: closed
deps: [jus-xcdp]
links: []
created: 2026-02-12T20:33:54Z
type: task
priority: 3
assignee: Alex Cabrera
parent: jus-ddig
tags: [commands, serving]
---
# Implement cmd_remove() command handler

Implement the remove command for unregistering projects.

## Function Implementation

```bash
cmd_remove() {
    local name=""
    local delete_files=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --files) delete_files=true; shift ;;
            *) name="$1"; shift ;;
        esac
    done
    
    if [[ -z "$name" ]]; then
        if is_dade_project; then
            name=$(get_current_project_name)
        else
            log_error "Project name required"
            exit 1
        fi
    fi
    
    local path=$(get_project_path "$name")
    
    if [[ -z "$path" ]]; then
        log_error "Project '$name' not found"
        exit 1
    fi
    
    # Stop if running
    if is_project_running "$name"; then
        stop_project "$name"
    fi
    
    # Handle file deletion
    if $delete_files; then
        echo ""
        log_warn "This will DELETE: $path"
        if confirm "Are you sure?"; then
            rm -rf "$path"
            log_success "Deleted $path"
        else
            delete_files=false
        fi
    fi
    
    # Unregister
    unregister_project "$name"
    
    # Clean up marker if files not deleted
    if ! $delete_files && [[ -d "$path" ]]; then
        rm -f "$path/.dade"
        rm -f "$path/.dade.pid"
    fi
    
    # Update proxy
    generate_caddyfile
    reload_proxy
    
    log_success "Removed $name from registry"
    if ! $delete_files && [[ -d "$path" ]]; then
        log_info "Files remain at: $path"
    fi
}
```

## Acceptance Criteria

- [ ] Unregisters project
- [ ] Stops server if running
- [ ] --files deletes project directory
- [ ] Requires confirmation for deletion
- [ ] Updates proxy config

