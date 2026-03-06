---
id: jus-xcdp
status: closed
deps: [jus-mckm]
links: []
created: 2026-02-12T20:32:54Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-l4vb
tags: [commands, serving]
---
# Implement cmd_stop() command handler

Implement the stop command handler.

## Function Implementation

```bash
cmd_stop() {
    local name=""
    local stop_all=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --all|-a) stop_all=true; shift ;;
            *) name="$1"; shift ;;
        esac
    done
    
    if $stop_all; then
        local count=0
        for project in $(list_project_names); do
            if is_project_running "$project"; then
                stop_project "$project"
                ((count++)) || true
            fi
        done
        
        if [[ $count -eq 0 ]]; then
            log_info "No projects running"
        else
            log_success "Stopped $count project(s)"
        fi
        return 0
    fi
    
    if [[ -n "$name" ]]; then
        stop_project "$name"
    elif is_dade_project; then
        stop_project
    else
        log_error "Not in a project directory"
        log_info "Specify project: dade stop <name>"
        log_info "Or stop all: dade stop --all"
        exit 1
    fi
}
```

## Acceptance Criteria

- [ ] Stops current project
- [ ] Stops named project
- [ ] --all stops all running projects
- [ ] Reports count stopped

