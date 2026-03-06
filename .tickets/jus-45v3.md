---
id: jus-45v3
status: closed
deps: [jus-q0si, jus-vt7v]
links: []
created: 2026-02-12T19:58:54Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-6d0x
tags: [commands, serving]
---
# Implement 'dade register' command

Implement the register command for adding existing directories to dade management.

## Usage

```
dade register [name]           # Register current directory
dade register --template <t>   # Specify template type
```

## Use Case

For projects created outside dade that want to use its serving features:
- Existing projects cloned from elsewhere
- Projects created manually
- Migrating from srv

## Implementation

```bash
cmd_register() {
    local name=""
    local template=""
    local project_dir=$(pwd)
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --template|-t) template="$2"; shift 2 ;;
            *) name="$1"; shift ;;
        esac
    done
    
    # Default name from directory
    if [[ -z "$name" ]]; then
        name=$(basename "$project_dir")
        if has_gum; then
            name=$(gum input --placeholder "project name" --value "$name")
        else
            read -p "Project name [$name]: " input
            name="${input:-$name}"
        fi
    fi
    
    # Validate name
    if [[ ! "$name" =~ ^[a-z0-9-]+$ ]]; then
        log_error "Invalid name: use lowercase, numbers, hyphens"
        exit 1
    fi
    
    # Check if already registered
    if [[ -f "$project_dir/.dade" ]]; then
        log_warn "Directory already registered"
        if ! confirm "Re-register?"; then
            exit 0
        fi
    fi
    
    # Check name conflict
    local existing_path=$(get_project_path "$name")
    if [[ -n "$existing_path" ]] && [[ "$existing_path" != "$project_dir" ]]; then
        log_error "Name '$name' already used by $existing_path"
        exit 1
    fi
    
    # Detect or prompt for template
    if [[ -z "$template" ]]; then
        template=$(detect_project_type "$project_dir")
        if [[ -z "$template" ]]; then
            template="static"  # Default fallback
        fi
        log_info "Detected template type: $template"
    fi
    
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
        # Could be node, next, etc
        echo "node"
    elif [[ -f "$dir/index.html" ]]; then
        echo "static"
    else
        echo ""
    fi
}
```

## Acceptance Criteria

- [ ] Current directory registered
- [ ] Name defaulted from directory
- [ ] Template detected or specified
- [ ] Port assigned
- [ ] .dade marker created
- [ ] Proxy updated

