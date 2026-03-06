---
id: jus-4l44
status: closed
deps: [jus-bxd6]
links: []
created: 2026-02-12T20:33:06Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-axr7
tags: [commands, serving]
---
# Implement cmd_list() command handler

Implement the list command that shows all projects.

## Function Implementation

```bash
cmd_list() {
    local running_only=false
    local json_output=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --running) running_only=true; shift ;;
            --json) json_output=true; shift ;;
            *) shift ;;
        esac
    done
    
    init_config
    
    if $json_output; then
        list_projects_json "$running_only"
        return 0
    fi
    
    local total=0
    local running=0
    
    echo ""
    if has_gum; then
        gum style --bold "Projects"
    else
        echo "Projects"
        echo "========"
    fi
    echo ""
    
    while IFS='|' read -r name port path template; do
        [[ -z "$name" ]] && continue
        
        local status="stopped"
        local icon="○"
        
        if is_project_running "$name"; then
            status="running"
            icon="●"
            ((running++)) || true
        fi
        
        if $running_only && [[ "$status" == "stopped" ]]; then
            continue
        fi
        
        ((total++)) || true
        
        if has_gum; then
            if [[ "$status" == "running" ]]; then
                echo "  $(gum style --foreground 2 "$icon") $name (https://$name.localhost)"
            else
                echo "  $(gum style --foreground 1 "$icon") $name (https://$name.localhost)"
            fi
            gum style --foreground 240 "    Template: $template | Port: $port | $status"
        else
            echo "  $icon $name (https://$name.localhost)"
            echo "    Template: $template | Port: $port | $status"
        fi
        echo ""
    done < <(jq -r 'to_entries[] | "\(.key)|\(.value.port)|\(.value.path)|\(.value.template)"' "$DADE_PROJECTS_FILE" 2>/dev/null)
    
    if [[ $total -eq 0 ]]; then
        log_info "No projects registered."
        log_info "Create one: dade new <name>"
    else
        echo "$total project(s) ($running running)"
    fi
}
```

## Acceptance Criteria

- [ ] Lists all projects
- [ ] Shows running status
- [ ] --running shows only running
- [ ] --json outputs JSON
- [ ] Shows URL, template, port

