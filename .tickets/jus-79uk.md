---
id: jus-79uk
status: closed
deps: [jus-73u3]
links: []
created: 2026-02-12T20:28:48Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-wz9o
tags: [commands, setup]
---
# Implement dependency check functions

Implement functions to check and install required dependencies.

## Functions

```bash
check_dependency() {
    local name="$1"
    local brew_name="${2:-$1}"
    local required="${3:-true}"
    
    if command -v "$name" &>/dev/null; then
        log_success "$name installed"
        return 0
    fi
    
    if [[ "$required" == "true" ]]; then
        log_warn "$name not found"
        if confirm "Install $name via Homebrew?"; then
            spin "Installing $name" brew install "$brew_name"
            if command -v "$name" &>/dev/null; then
                log_success "$name installed"
                return 0
            fi
        fi
        log_error "Could not install $name"
        return 1
    else
        log_info "$name not found (optional)"
        return 0
    fi
}

check_all_dependencies() {
    local all_ok=true
    
    # Required
    check_dependency "jq" || all_ok=false
    check_dependency "caddy" || all_ok=false
    
    # Optional
    check_dependency "gum" "gum" false
    check_dependency "cloudflared" "cloudflared" false
    
    $all_ok
}
```

## Acceptance Criteria

- [ ] Checks for required dependencies (jq, caddy)
- [ ] Offers to install missing via Homebrew
- [ ] Handles optional dependencies (gum, cloudflared)
- [ ] Returns success/failure appropriately

