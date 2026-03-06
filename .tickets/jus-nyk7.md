---
id: jus-nyk7
status: closed
deps: [jus-0q5m]
links: []
created: 2026-02-12T20:33:42Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-45v3
tags: [commands, serving]
---
# Implement cmd_register() command handler

Implement the register command for existing directories.

## Function Implementation

```bash
cmd_register() {
    local name=""
    local template=""
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --template|-t) template="$2"; shift 2 ;;
            *) name="$1"; shift ;;
        esac
    done
    
    local project_dir=$(pwd)
    
    # Check if already registered
    if [[ -f "$project_dir/.dade" ]]; then
        log_warn "Directory already registered"
        if ! confirm "Re-register?"; then
            return 0
        fi
    fi
    
    # Default name from directory
    if [[ -z "$name" ]]; then
        local default_name=$(basename "$project_dir" | tr '[:upper:]' '[:lower:]' | tr ' ' '-')
        if can_use_gum; then
            name=$(gum input --placeholder "project name" --value "$default_name")
        else
            read -p "Project name [$default_name]: " name
            name="${name:-$default_name}"
        fi
    fi
    
    # Validate name
    if [[ ! "$name" =~ ^[a-z][a-z0-9-]*$ ]]; then
        log_error "Invalid name: use lowercase letters, numbers, hyphens"
        exit 1
    fi
    
    # Check name conflict
    local existing=$(get_project_path "$name")
    if [[ -n "$existing" ]] && [[ "$existing" != "$project_dir" ]]; then
        log_error "Name '$name' already used by $existing"
        exit 1
    fi
    
    # Detect or set template
    if [[ -z "$template" ]]; then
        template=$(detect_project_type "$project_dir")
        template="${template:-static}"
        log_info "Detected type: $template"
    fi
    
    init_config
    
    # Assign port
    local port=$(next_port)
    
    # Write marker
    write_project_marker "$project_dir" "$name" "$template" "$port"
    
    # Register
    register_project "$name" "$port" "$project_dir" "$template"
    
    # Update proxy
    generate_caddyfile
    reload_proxy
    
    log_success "Registered: $name"
    log_info "URL: https://$name.localhost"
    log_info "Start: dade start"
}

detect_project_type() {
    local dir="$1"
    
    if [[ -f "$dir/manage.py" ]]; then
        echo "django-hypermedia"
    elif [[ -f "$dir/package.json" ]]; then
        if grep -q "next" "$dir/package.json" 2>/dev/null; then
            echo "nextjs"
        else
            echo "node"
        fi
    elif [[ -f "$dir/index.html" ]]; then
        echo "static"
    else
        echo ""
    fi
}
```

## Acceptance Criteria

- [ ] Registers current directory
- [ ] Prompts for name if not given
- [ ] Detects project type
- [ ] Assigns port
- [ ] Updates proxy config

