---
id: jus-ghy8
status: closed
deps: [jus-ocjp]
links: []
created: 2026-02-12T19:58:25Z
type: task
priority: 3
assignee: Alex Cabrera
parent: jus-6d0x
tags: [commands, serving]
---
# Implement 'dade open' command

Implement the open command for opening projects in the browser.

## Usage

```
dade open           # Open current project
dade open <name>    # Open specific project
```

## Implementation

```bash
cmd_open() {
    local project_name="$1"
    
    if [[ -z "$project_name" ]]; then
        if ! is_dade_project; then
            log_error "Not a dade project directory"
            exit 1
        fi
        project_name=$(get_current_project_name)
    fi
    
    # Verify project exists
    if ! project_exists "$project_name"; then
        log_error "Project '$project_name' not found"
        exit 1
    fi
    
    local url="https://$project_name.localhost"
    
    # Check if running, offer to start if not
    if ! is_project_running "$project_name"; then
        log_warn "Project not running"
        if confirm "Start $project_name?"; then
            # Need to cd to project directory for start
            local path=$(get_project_path "$project_name")
            (cd "$path" && cmd_start)
            sleep 1
        else
            log_info "Start with: dade start"
            return
        fi
    fi
    
    log_info "Opening $url"
    
    # macOS
    if command -v open &>/dev/null; then
        open "$url"
    # Linux
    elif command -v xdg-open &>/dev/null; then
        xdg-open "$url"
    else
        log_error "No browser opener found"
        log_info "Visit: $url"
    fi
}
```

## Cross-Platform

- macOS: open
- Linux: xdg-open
- WSL: could detect and use cmd.exe /c start

## Acceptance Criteria

- [ ] Opens browser to project URL
- [ ] Works from project directory
- [ ] Works with project name argument
- [ ] Offers to start if not running
- [ ] Cross-platform browser opening

