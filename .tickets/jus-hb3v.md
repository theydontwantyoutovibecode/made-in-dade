---
id: jus-hb3v
status: closed
deps: [jus-vt7v]
links: []
created: 2026-02-12T19:54:28Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-cdr3
tags: [infrastructure, caddy]
---
# Implement Caddyfile generation

Implement automatic Caddyfile generation from the project registry.

## Caddyfile Structure

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

## Implementation

```bash
generate_caddyfile() {
    # 1. Write global options block with local_certs
    # 2. Iterate projects.json
    # 3. For each project, write reverse_proxy block
    # 4. Write to DADE_CADDYFILE
}
```

## When to Regenerate

- After register_project()
- After unregister_project()
- On 'dade sync' command
- On 'dade proxy restart'

## Reload Proxy

After generating new Caddyfile:
```bash
caddy reload --config "$DADE_CADDYFILE"
```

## Error Handling

- Validate Caddyfile syntax before applying: caddy validate --config file
- Keep backup of previous working config
- Rollback on validation failure

## Acceptance Criteria

- [ ] Caddyfile generated with all registered projects
- [ ] local_certs directive included for HTTPS
- [ ] Caddyfile validated before applying
- [ ] Proxy reloaded after regeneration
- [ ] Invalid configs don't break existing setup

