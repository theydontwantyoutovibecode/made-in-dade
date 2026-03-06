---
id: jus-axr7
status: closed
deps: [jus-vt7v]
links: []
created: 2026-02-12T19:58:12Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-6d0x
tags: [commands, serving]
---
# Implement 'dade list' command

Implement the list command for showing all registered projects and their status.

## Usage

```
dade list           # List all projects
dade list --running # Only show running projects
dade list --json    # Output as JSON
```

## Output Format

With gum:
```
Projects

  ● myproject (https://myproject.localhost)
    Template: django-hypermedia | Port: 3000 | running
    Path: /Users/alex/Code/myproject

  ○ another (https://another.localhost)
    Template: hypertext | Port: 3001 | stopped
    Path: /Users/alex/Code/another

2 projects (1 running, 1 stopped)
```

## Implementation

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
    
    if $json_output; then
        list_projects_json "$running_only"
        return
    fi
    
    local total=0
    local running=0
    
    print_header "Projects"
    echo ""
    
    jq -r 'to_entries[] | "\(.key)|\(.value.port)|\(.value.path)|\(.value.template)"' \
        "$DADE_PROJECTS_FILE" | while IFS='|' read -r name port path template; do
        
        local status="stopped"
        local status_icon="○"
        local status_color="1"  # red
        
        if is_project_running "$name"; then
            status="running"
            status_icon="●"
            status_color="2"  # green
            ((running++)) || true
        fi
        
        if $running_only && [[ "$status" == "stopped" ]]; then
            continue
        fi
        
        ((total++)) || true
        
        # Display with gum or fallback
        if has_gum; then
            echo "  $(gum style --foreground "$status_color" "$status_icon") $name (https://$name.localhost)"
            gum style --foreground 240 "    Template: $template | Port: $port | $status"
            gum style --foreground 240 "    Path: $path"
        else
            echo "  $status_icon $name (https://$name.localhost)"
            echo "    Template: $template | Port: $port | $status"
            echo "    Path: $path"
        fi
        echo ""
    done
    
    # Summary
    if [[ $total -eq 0 ]]; then
        log_info "No projects registered."
        log_info "Create one: dade new <name>"
    else
        echo "$total project(s) ($running running)"
    fi
}

is_project_running() {
    local name="$1"
    local port=$(jq -r --arg n "$name" '.[$n].port' "$DADE_PROJECTS_FILE")
    lsof -i ":$port" &>/dev/null
}
```

## Acceptance Criteria

- [ ] All projects listed with status
- [ ] Running status detected correctly
- [ ] --running filters to running only
- [ ] --json outputs valid JSON
- [ ] Summary shows totals
- [ ] Empty state handled

