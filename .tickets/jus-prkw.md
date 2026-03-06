---
id: jus-prkw
status: closed
deps: [jus-b3ad]
links: []
created: 2026-02-12T20:33:17Z
type: task
priority: 3
assignee: Alex Cabrera
parent: jus-ghy8
tags: [commands, serving]
---
# Implement cmd_open() command handler

Implement the open command that opens project in browser.

## Function Implementation

```bash
cmd_open() {
    local name="$1"
    
    if [[ -z "$name" ]]; then
        if is_dade_project; then
            name=$(get_current_project_name)
        else
            log_error "Not in a project directory"
            log_info "Specify project: dade open <name>"
            exit 1
        fi
    fi
    
    if ! project_exists "$name"; then
        log_error "Project '$name' not found"
        exit 1
    fi
    
    local url="https://$name.localhost"
    
    # Check if running
    if ! is_project_running "$name"; then
        log_warn "Project not running"
        if confirm "Start $name first?"; then
            local path=$(get_project_path "$name")
            (cd "$path" && cmd_start --bg)
            sleep 1
        else
            log_info "Start with: dade start"
            return 0
        fi
    fi
    
    log_info "Opening $url"
    open_browser "$url"
}

open_browser() {
    local url="$1"
    
    # macOS
    if command -v open &>/dev/null; then
        open "$url"
    # Linux
    elif command -v xdg-open &>/dev/null; then
        xdg-open "$url"
    # WSL
    elif command -v wslview &>/dev/null; then
        wslview "$url"
    else
        log_warn "No browser opener found"
        log_info "Visit: $url"
    fi
}
```

## Acceptance Criteria

- [ ] Opens current project in browser
- [ ] Opens named project in browser
- [ ] Offers to start if not running
- [ ] Cross-platform browser opening

