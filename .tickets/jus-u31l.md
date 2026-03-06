---
id: jus-u31l
status: closed
deps: [jus-9hvc]
links: []
created: 2026-02-12T20:28:16Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-hb3v
tags: [infrastructure, caddy]
---
# Implement reload_proxy() function

Implement the function that reloads the Caddy proxy after config changes.

## Function Implementation

```bash
reload_proxy() {
    if is_proxy_running; then
        caddy reload --config "$DADE_CADDYFILE" 2>/dev/null || true
    fi
}

is_proxy_running() {
    launchctl list 2>/dev/null | grep -q "$DADE_PROXY_LABEL"
}
```

## When to Call

- After generate_caddyfile()
- After register_project()
- After unregister_project()
- On 'dade proxy restart'

## Acceptance Criteria

- [ ] Reloads Caddy config if running
- [ ] No error if proxy not running
- [ ] is_proxy_running works correctly

