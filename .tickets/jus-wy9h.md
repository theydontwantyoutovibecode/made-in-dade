---
id: jus-wy9h
status: closed
deps: [jus-3cof]
links: []
created: 2026-02-12T20:27:33Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-f3br
tags: [infrastructure, setup]
---
# Implement init_config() function

Implement the init_config function that creates the configuration directory structure.

## Function Implementation

```bash
init_config() {
    # Create main config directory
    mkdir -p "$DADE_CONFIG_DIR"
    
    # Create templates directory
    mkdir -p "$DADE_TEMPLATES_DIR"
    
    # Initialize empty projects.json if missing
    if [[ ! -f "$DADE_PROJECTS_FILE" ]]; then
        echo '{}' > "$DADE_PROJECTS_FILE"
    fi
    
    # Create minimal Caddyfile if missing
    if [[ ! -f "$DADE_CADDYFILE" ]]; then
        cat > "$DADE_CADDYFILE" << 'EOF'
{
    local_certs
}
EOF
    fi
}
```

## When to Call

- At start of setup command
- At start of new command
- At start of register command
- At start of install command
- Basically any command that writes to config

## Acceptance Criteria

- [ ] Creates ~/.config/dade/
- [ ] Creates templates/ subdirectory
- [ ] Initializes empty projects.json
- [ ] Creates minimal Caddyfile
- [ ] Idempotent (safe to call multiple times)

