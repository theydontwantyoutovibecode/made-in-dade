---
id: jus-wz9o
status: closed
deps: [jus-aa04]
links: []
created: 2026-02-12T19:55:00Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-cdr3
tags: [commands, setup]
---
# Implement 'dade setup' command

Implement the setup command for first-time installation and dependency checking.

## Usage

```
dade setup           # Interactive setup
dade setup --check   # Just check dependencies
```

## Setup Steps

1. **Check/install dependencies**
   - caddy (required)
   - jq (required)
   - gum (optional, but recommended)
   - cloudflared (optional, for tunnels)

2. **Initialize configuration**
   - Create ~/.config/dade/
   - Initialize projects.json
   - Create templates/ directory

3. **Set up proxy service**
   - Generate initial Caddyfile
   - Create launchd plist
   - Start proxy service

4. **Trust Caddy CA**
   - Prompt user for sudo
   - Run 'sudo caddy trust'

5. **Offer to install default templates**
   - Show list of official templates
   - Let user select which to install

## Dependency Installation

```bash
check_and_install_dep() {
    local name="$1"
    local brew_name="${2:-$1}"
    
    if ! command -v "$name" &>/dev/null; then
        if confirm "Install $name via Homebrew?"; then
            brew install "$brew_name"
        fi
    fi
}
```

## Migration from srv

If ~/.config/srv exists:
1. Detect existing srv installation
2. Offer to migrate projects
3. Copy projects.json entries
4. Update .srv files to .dade
5. Stop old srv proxy service

## Acceptance Criteria

- [ ] Dependencies checked and install offered
- [ ] Config directory initialized
- [ ] Proxy service configured and started
- [ ] Caddy CA trust offered
- [ ] Default templates offered
- [ ] srv migration works if applicable

