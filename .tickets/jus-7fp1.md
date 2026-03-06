---
id: jus-7fp1
status: closed
deps: [jus-hb3v]
links: []
created: 2026-02-12T19:54:39Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-cdr3
tags: [infrastructure, launchd, macos]
---
# Implement launchd proxy service

Implement the launchd service that keeps the Caddy proxy running.

## Plist Location

~/Library/LaunchAgents/land.charm.dade.proxy.plist

## Plist Content

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>land.charm.dade.proxy</string>
    <key>ProgramArguments</key>
    <array>
        <string>/opt/homebrew/bin/caddy</string>
        <string>run</string>
        <string>--config</string>
        <string>/Users/USER/.config/dade/Caddyfile</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/Users/USER/.config/dade/proxy.log</string>
    <key>StandardErrorPath</key>
    <string>/Users/USER/.config/dade/proxy.err</string>
</dict>
</plist>
```

## Functions

```bash
create_plist()           # Generate plist file
install_proxy_service()  # bootstrap launchd service
uninstall_proxy_service() # bootout launchd service
is_proxy_running()       # Check if service is loaded
```

## Commands

- dade proxy start: bootstrap if not running
- dade proxy stop: bootout service
- dade proxy restart: bootout + bootstrap
- dade proxy status: check if running

## Considerations

- Detect Caddy path dynamically (which caddy)
- Expand ~ to actual home directory in plist
- Handle existing srv plist (offer migration)

## Acceptance Criteria

- [ ] Plist generated with correct paths
- [ ] Service starts on 'dade proxy start'
- [ ] Service stops on 'dade proxy stop'
- [ ] Service survives logout/restart (KeepAlive)
- [ ] Logs written to config directory

