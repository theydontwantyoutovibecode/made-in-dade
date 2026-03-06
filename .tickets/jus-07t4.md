---
id: jus-07t4
status: closed
deps: [jus-ikny]
links: []
created: 2026-02-12T20:30:21Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-ae77
tags: [commands, templates]
---
# Implement cmd_templates() command handler

Implement the templates command that lists installed templates.

## Function Implementation

```bash
cmd_templates() {
    local json_output=false
    [[ "${1:-}" == "--json" ]] && json_output=true
    
    init_config
    
    # Check if any templates installed
    if [[ ! -d "$DADE_TEMPLATES_DIR" ]] || \
       [[ -z "$(ls -A "$DADE_TEMPLATES_DIR" 2>/dev/null)" ]]; then
        log_info "No templates installed."
        echo ""
        log_info "Install with: dade install <git-url>"
        log_info "Or see official: dade install --list-official"
        return 0
    fi
    
    if $json_output; then
        echo "["
        local first=true
        for dir in "$DADE_TEMPLATES_DIR"/*/; do
            [[ ! -d "$dir" ]] && continue
            local name=$(basename "$dir")
            local manifest="$dir/dade.toml"
            local desc=$(parse_toml_value "$manifest" "template.description")
            local type=$(parse_toml_value "$manifest" "serve.type")
            local source=$(cat "$dir/.source" 2>/dev/null || echo "")
            
            $first || echo ","
            first=false
            echo "  {\"name\": \"$name\", \"description\": \"$desc\", \"type\": \"$type\", \"source\": \"$source\"}"
        done
        echo "]"
        return 0
    fi
    
    echo ""
    if has_gum; then
        gum style --bold "Installed Templates"
    else
        echo "Installed Templates"
        echo "==================="
    fi
    echo ""
    
    for dir in "$DADE_TEMPLATES_DIR"/*/; do
        [[ ! -d "$dir" ]] && continue
        local name=$(basename "$dir")
        local manifest="$dir/dade.toml"
        local desc=$(parse_toml_value "$manifest" "template.description")
        local type=$(parse_toml_value "$manifest" "serve.type")
        local source=$(cat "$dir/.source" 2>/dev/null || echo "local")
        
        if has_gum; then
            gum style --foreground 220 "  $name"
            gum style --foreground 240 "    $desc"
            gum style --foreground 240 "    Type: $type"
        else
            echo "  $name"
            echo "    $desc"
            echo "    Type: $type"
        fi
        echo ""
    done
}
```

## Acceptance Criteria

- [ ] Lists all installed templates
- [ ] Shows name, description, type
- [ ] Handles empty templates directory
- [ ] --json outputs valid JSON
- [ ] Pretty output with gum

