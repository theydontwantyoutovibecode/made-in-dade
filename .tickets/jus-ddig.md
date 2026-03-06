---
id: jus-ddig
status: closed
deps: [jus-l4vb]
links: []
created: 2026-02-12T19:59:06Z
type: task
priority: 3
assignee: Alex Cabrera
parent: jus-6d0x
tags: [commands, serving]
---
# Implement 'dade remove' command

Implement the remove command for unregistering projects from dade.

## Usage

```
dade remove [name]      # Remove project from registry
dade remove --files     # Also delete project files (dangerous!)
```

## Implementation

```bash
cmd_remove() {
    local project_name=""
    local delete_files=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --files) delete_files=true; shift ;;
            *) project_name="$1"; shift ;;
        esac
    done
    
    if [[ -z "$project_name" ]]; then
        if is_dade_project; then
            project_name=$(get_current_project_name)
        else
            log_error "Project name required"
            exit 1
        fi
    fi
    
    local project_path=$(get_project_path "$project_name")
    
    if [[ -z "$project_path" ]]; then
        log_error "Project '$project_name' not found"
        exit 1
    fi
    
    # Stop if running
    if is_project_running "$project_name"; then
        stop_project "$project_name"
    fi
    
    if $delete_files; then
        if confirm "DELETE all files in $project_path?"; then
            rm -rf "$project_path"
            log_success "Deleted $project_path"
        else
            delete_files=false
        fi
    fi
    
    # Unregister
    unregister_project "$project_name"
    
    # Remove marker (if files not deleted)
    if ! $delete_files && [[ -d "$project_path" ]]; then
        rm -f "$project_path/.dade"
        rm -f "$project_path/.dade.pid"
    fi
    
    # Update proxy
    generate_caddyfile
    reload_proxy
    
    log_success "Removed $project_name from registry"
    if ! $delete_files; then
        log_info "Files remain at: $project_path"
    fi
}
```

## Safety

- --files requires explicit confirmation
- Default only removes from registry
- Stops running server first

## Acceptance Criteria

- [ ] Project unregistered from registry
- [ ] .dade marker removed
- [ ] Running server stopped first
- [ ] --files deletes with confirmation
- [ ] Proxy config updated

