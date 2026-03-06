---
id: jus-qnt6
status: closed
deps: [jus-79uk]
links: []
created: 2026-02-12T20:28:59Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-wz9o
tags: [commands, setup]
---
# Implement cmd_setup() command handler

Implement the setup command that performs first-time installation.

## Function Implementation

```bash
cmd_setup() {
    local check_only=false
    [[ "${1:-}" == "--check" ]] && check_only=true
    
    print_header
    log_info "Setting up dade..."
    echo ""
    
    # Check dependencies
    if ! check_all_dependencies; then
        log_error "Missing required dependencies"
        exit 1
    fi
    
    if $check_only; then
        log_success "All dependencies OK"
        return 0
    fi
    
    # Initialize config
    init_config
    log_success "Configuration initialized"
    
    # Check for srv migration
    if detect_srv_installation; then
        if confirm "Migrate from srv?"; then
            migrate_from_srv
        fi
    fi
    
    # Generate Caddyfile
    generate_caddyfile
    log_success "Caddyfile generated"
    
    # Start proxy service
    start_proxy_service
    
    # Trust CA
    echo ""
    log_info "For HTTPS to work, Caddy's CA must be trusted."
    if confirm "Trust Caddy CA? (requires sudo)"; then
        sudo caddy trust
        log_success "CA trusted"
    fi
    
    # Offer to install templates
    echo ""
    log_info "Would you like to install official templates?"
    if confirm "Install official templates?"; then
        offer_official_templates
    fi
    
    echo ""
    log_success "Setup complete!"
    log_info "Create your first project: dade new myproject"
}
```

## Acceptance Criteria

- [ ] Checks all dependencies
- [ ] Initializes config directory
- [ ] Offers srv migration
- [ ] Starts proxy service
- [ ] Offers CA trust
- [ ] Offers official templates
- [ ] --check flag only checks deps

