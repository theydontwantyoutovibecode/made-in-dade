---
id: jus-qg1c
status: closed
deps: [jus-u31l]
links: []
created: 2026-02-12T20:28:24Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-7fp1
tags: [infrastructure, launchd]
---
# Implement create_plist() function

Implement the function that creates the launchd plist for the proxy service.

## Function Implementation

```bash
create_plist() {
    local caddy_path
    caddy_path=$(which caddy)
    
    cat > "$DADE_PLIST" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>${DADE_PROXY_LABEL}</string>
    <key>ProgramArguments</key>
    <array>
        <string>${caddy_path}</string>
        <string>run</string>
        <string>--config</string>
        <string>${DADE_CADDYFILE}</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>${DADE_LOG}</string>
    <key>StandardErrorPath</key>
    <string>${DADE_ERR}</string>
</dict>
</plist>
EOF
}
```

## Considerations

- Detect Caddy path dynamically
- All paths must be absolute (no ~ or variables in plist)
- KeepAlive ensures restart on crash

## Acceptance Criteria

- [ ] Generates valid plist XML
- [ ] Uses absolute paths
- [ ] Finds Caddy binary location
- [ ] Logs to config directory

