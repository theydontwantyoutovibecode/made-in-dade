---
id: jus-bz5o
status: closed
deps: [jus-3bvb, jus-dx0j, jus-0q5m]
links: []
created: 2026-02-12T20:31:34Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-keyx
tags: [commands, templates]
---
# Implement cmd_new() command handler

Implement the new command for creating projects from templates.

## Function Implementation

```bash
cmd_new() {
    local project_name=""
    local template=""
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --template|-t) template="$2"; shift 2 ;;
            -*) log_error "Unknown option: $1"; exit 1 ;;
            *) project_name="$1"; shift ;;
        esac
    done
    
    print_header
    init_config
    
    # Prompt for project name if not provided
    if [[ -z "$project_name" ]]; then
        if can_use_gum; then
            project_name=$(gum input --placeholder "my-project" --header "Project name:")
        else
            read -rp "Project name: " project_name
        fi
    fi
    
    # Validate name
    if [[ -z "$project_name" ]]; then
        log_error "Project name required"
        exit 1
    fi
    
    if [[ ! "$project_name" =~ ^[a-z][a-z0-9-]*$ ]]; then
        log_error "Invalid name: use lowercase letters, numbers, hyphens"
        exit 1
    fi
    
    # Check if directory exists
    if [[ -d "$project_name" ]]; then
        log_error "Directory '$project_name' already exists"
        exit 1
    fi
    
    # Choose template
    if [[ -z "$template" ]]; then
        template=$(choose_template)
    fi
    
    local template_dir="$DADE_TEMPLATES_DIR/$template"
    if [[ ! -d "$template_dir" ]]; then
        log_error "Template '$template' not found"
        exit 1
    fi
    
    log_info "Creating project: $project_name"
    log_info "Template: $template"
    echo ""
    
    # Copy template
    copy_template "$template_dir" "$project_name"
    log_success "Template copied"
    
    # Assign port and write marker
    local port=$(next_port)
    write_project_marker "$project_name" "$project_name" "$template" "$port"
    
    # Register project
    local full_path="$(pwd)/$project_name"
    register_project "$project_name" "$port" "$full_path" "$template"
    
    # Update proxy
    generate_caddyfile
    reload_proxy
    
    log_success "Project registered"
    
    # Run setup script if defined
    local manifest="$template_dir/dade.toml"
    local setup_cmd=$(parse_toml_value "$manifest" "scaffold.setup")
    
    if [[ -n "$setup_cmd" ]]; then
        echo ""
        log_info "Running setup..."
        (cd "$project_name" && eval "$setup_cmd")
    fi
    
    # Initialize git
    spin "Initializing git" git -C "$project_name" init
    log_success "Git initialized"
    
    echo ""
    log_success "Project '$project_name' created!"
    log_info "URL: https://$project_name.localhost"
    echo ""
    log_info "Next steps:"
    echo "  cd $project_name"
    echo "  dade start"
}
```

## Acceptance Criteria

- [ ] Creates project directory
- [ ] Copies template files
- [ ] Assigns port and writes marker
- [ ] Registers with project registry
- [ ] Updates proxy config
- [ ] Runs setup script if defined
- [ ] Initializes git repository
- [ ] Shows next steps

