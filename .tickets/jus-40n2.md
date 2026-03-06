---
id: jus-40n2
status: closed
deps: [jus-3ii0]
links: []
created: 2026-02-12T20:32:12Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-ocjp
tags: [serving, command]
---
# Implement start_command_server() function

Implement the function that starts a command-based server from template config.

## Function Implementation

```bash
start_command_server() {
    local template="$1"
    local port="$2"
    local prod_mode="${3:-false}"
    
    local template_dir="$DADE_TEMPLATES_DIR/$template"
    local manifest="$template_dir/dade.toml"
    
    # Get serve command
    local cmd
    if [[ "$prod_mode" == "true" ]]; then
        cmd=$(parse_toml_value "$manifest" "serve.prod")
        # Fall back to dev if no prod command
        [[ -z "$cmd" ]] && cmd=$(parse_toml_value "$manifest" "serve.dev")
    else
        cmd=$(parse_toml_value "$manifest" "serve.dev")
        # Fall back to prod if no dev command
        [[ -z "$cmd" ]] && cmd=$(parse_toml_value "$manifest" "serve.prod")
    fi
    
    if [[ -z "$cmd" ]]; then
        log_error "No serve command defined in template"
        return 1
    fi
    
    # Get port environment variable name
    local port_env=$(parse_toml_value "$manifest" "serve.port_env")
    port_env="${port_env:-PORT}"
    
    log_info "Starting: $cmd"
    log_info "Port ($port_env): $port"
    
    # Export port and run command
    export "$port_env"="$port"
    
    # Run in foreground (user sees output)
    eval "$cmd"
}

start_command_server_bg() {
    local template="$1"
    local port="$2"
    local prod_mode="${3:-false}"
    
    # Same as above but background
    # ... (get cmd and port_env)
    
    export "$port_env"="$port"
    eval "$cmd" &
    local pid=$!
    echo "$pid" > .dade.pid
    
    sleep 1
    if kill -0 "$pid" 2>/dev/null; then
        return 0
    else
        return 1
    fi
}
```

## Considerations

- Foreground by default so user sees output
- Background variant for scripting
- PORT passed via configured env var
- Falls back between dev/prod commands

## Acceptance Criteria

- [ ] Runs template-defined command
- [ ] Passes port via configured env var
- [ ] Foreground shows command output
- [ ] Background variant saves PID
- [ ] Falls back between dev/prod

