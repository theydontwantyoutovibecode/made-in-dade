---
id: jus-b3ad
status: closed
deps: [jus-40n2, jus-bxd6]
links: []
created: 2026-02-12T20:32:34Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-ocjp
tags: [commands, serving]
---
# Implement cmd_start() command handler

Implement the start command handler that serves projects.

## Function Implementation

```bash
cmd_start() {
    local prod_mode=false
    local background=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --prod) prod_mode=true; shift ;;
            --bg|--background) background=true; shift ;;
            *) shift ;;
        esac
    done
    
    # Must be in a project directory
    if ! is_dade_project; then
        log_error "Not a dade project directory"
        log_info "Create a project: dade new <name>"
        log_info "Or register this directory: dade register"
        exit 1
    fi
    
    local name=$(get_current_project_name)
    local port=$(get_current_project_port)
    local template=$(get_current_project_template)
    
    # Check if already running
    if is_project_running; then
        log_warn "Project already running"
        log_info "URL: https://$name.localhost"
        log_info "Stop with: dade stop"
        return 0
    fi
    
    # Get serve type from template
    local serve_type="static"
    local template_dir="$DADE_TEMPLATES_DIR/$template"
    if [[ -f "$template_dir/dade.toml" ]]; then
        serve_type=$(parse_toml_value "$template_dir/dade.toml" "serve.type")
    fi
    
    serve_type="${serve_type:-static}"
    
    log_info "Starting $name..."
    
    case "$serve_type" in
        static)
            if $background; then
                start_static_server "$port"
                log_success "Server running in background"
            else
                log_info "Serving static files on port $port"
                log_success "URL: https://$name.localhost"
                log_info "Press Ctrl+C to stop"
                echo ""
                # Foreground - don't background the process
                caddy file-server --listen ":$port" --root "."
            fi
            ;;
        command)
            if $background; then
                start_command_server_bg "$template" "$port" "$prod_mode"
                log_success "Server running in background"
            else
                log_success "URL: https://$name.localhost"
                echo ""
                start_command_server "$template" "$port" "$prod_mode"
            fi
            ;;
        *)
            log_error "Unknown serve type: $serve_type"
            exit 1
            ;;
    esac
}
```

## Acceptance Criteria

- [ ] Detects project from .dade
- [ ] Routes to static or command server
- [ ] --prod flag for production mode
- [ ] --bg flag for background
- [ ] Shows URL on start
- [ ] Handles already-running case

