---
id: jus-dx0j
status: closed
deps: [jus-07t4]
links: []
created: 2026-02-12T20:31:17Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-keyx
tags: [commands, templates, ux]
---
# Implement choose_template() function

Implement the interactive template picker function.

## Function Implementation

```bash
choose_template() {
    local templates=()
    local template_names=()
    
    # Build list of installed templates
    for dir in "$DADE_TEMPLATES_DIR"/*/; do
        [[ ! -d "$dir" ]] && continue
        local name=$(basename "$dir")
        local manifest="$dir/dade.toml"
        local desc=$(parse_toml_value "$manifest" "template.description")
        
        templates+=("$name")
        template_names+=("$name - $desc")
    done
    
    if [[ ${#templates[@]} -eq 0 ]]; then
        log_error "No templates installed"
        log_info "Install with: dade install --list-official"
        exit 1
    fi
    
    if [[ ${#templates[@]} -eq 1 ]]; then
        echo "${templates[0]}"
        return
    fi
    
    # Show picker
    local choice
    if can_use_gum; then
        choice=$(gum choose --header "Select a template:" "${template_names[@]}")
    else
        echo "Select a template:" >&2
        select opt in "${template_names[@]}"; do
            if [[ -n "$opt" ]]; then
                choice="$opt"
                break
            fi
        done
    fi
    
    # Extract template name from choice
    local selected_name="${choice%% - *}"
    echo "$selected_name"
}
```

## Acceptance Criteria

- [ ] Shows picker with template names and descriptions
- [ ] Returns single template if only one installed
- [ ] Works with gum for pretty UI
- [ ] Falls back to bash select
- [ ] Exits with error if no templates

