---
id: jus-t35w
status: closed
deps: [jus-lvnb]
links: []
created: 2026-02-12T20:29:57Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-ovh7
tags: [commands, templates]
---
# Implement cmd_install() command handler

Implement the install command for adding template plugins.

## Function Implementation

```bash
cmd_install() {
    local url=""
    local name=""
    local list_official=false
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --name|-n) name="$2"; shift 2 ;;
            --list-official) list_official=true; shift ;;
            -*) log_error "Unknown option: $1"; exit 1 ;;
            *) url="$1"; shift ;;
        esac
    done
    
    if $list_official; then
        show_official_templates
        return 0
    fi
    
    if [[ -z "$url" ]]; then
        log_error "Usage: dade install <git-url>"
        exit 1
    fi
    
    init_config
    
    # Clone to temp directory
    local tmp_dir
    tmp_dir=$(mktemp -d)
    
    spin "Cloning template" git clone --depth 1 "$url" "$tmp_dir"
    
    # Validate manifest
    if [[ ! -f "$tmp_dir/dade.toml" ]]; then
        log_error "Template missing dade.toml manifest"
        rm -rf "$tmp_dir"
        exit 1
    fi
    
    if ! validate_manifest "$tmp_dir/dade.toml"; then
        rm -rf "$tmp_dir"
        exit 1
    fi
    
    # Get name from manifest if not provided
    if [[ -z "$name" ]]; then
        name=$(parse_toml_value "$tmp_dir/dade.toml" "template.name")
    fi
    
    local target="$DADE_TEMPLATES_DIR/$name"
    
    # Check if already installed
    if [[ -d "$target" ]]; then
        if confirm "Template '$name' already installed. Update?"; then
            rm -rf "$target"
        else
            rm -rf "$tmp_dir"
            return 0
        fi
    fi
    
    # Move to templates directory
    mv "$tmp_dir" "$target"
    
    # Store source URL
    echo "$url" > "$target/.source"
    
    log_success "Installed template: $name"
    
    local desc=$(parse_toml_value "$target/dade.toml" "template.description")
    log_info "$desc"
}
```

## Acceptance Criteria

- [ ] Clones template from git URL
- [ ] Validates manifest before installing
- [ ] Derives name from manifest
- [ ] Stores source URL for updates
- [ ] Handles already-installed templates
- [ ] --list-official shows available templates

