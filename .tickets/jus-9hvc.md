---
id: jus-9hvc
status: closed
deps: [jus-1gl4]
links: []
created: 2026-02-12T20:28:10Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-hb3v
tags: [infrastructure, caddy]
---
# Implement generate_caddyfile() function

Implement the function that generates Caddyfile from project registry.

## Function Implementation

```bash
generate_caddyfile() {
    # Global options
    cat > "$DADE_CADDYFILE" << 'EOF'
{
    local_certs
}

EOF
    
    # Add entry for each project
    jq -r 'to_entries[] | "\(.key) \(.value.port)"' "$DADE_PROJECTS_FILE" 2>/dev/null | \
    while read -r name port; do
        if [[ -n "$name" ]] && [[ -n "$port" ]]; then
            cat >> "$DADE_CADDYFILE" << EOF
https://${name}.localhost {
    reverse_proxy localhost:${port}
}

EOF
        fi
    done
}
```

## Output Format

```
{
    local_certs
}

https://myproject.localhost {
    reverse_proxy localhost:3000
}

https://another.localhost {
    reverse_proxy localhost:3001
}
```

## Acceptance Criteria

- [ ] Generates valid Caddyfile
- [ ] Includes local_certs directive
- [ ] Adds reverse_proxy for each project
- [ ] Handles empty registry

