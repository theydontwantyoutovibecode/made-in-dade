---
id: jus-i0c0
status: closed
deps: []
links: []
created: 2026-02-12T20:36:02Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-jjl8
tags: [commands, ux]
---
# Implement cmd_help() and main() dispatch

Implement the help command and main dispatch function.

## cmd_help() Implementation

```bash
cmd_help() {
    print_header
    
    if has_gum; then
        gum style --bold "Usage"
        echo ""
        echo "  dade <command> [args]"
        echo ""
        gum style --bold "Commands"
        echo ""
        echo "  setup           First-time setup (dependencies, proxy, CA)"
        echo "  install <url>   Install template from git repository"
        echo "  uninstall <n>   Remove installed template"
        echo "  templates       List installed templates"
        echo "  update <name>   Update template (or --all)"
        echo "  new [name]      Create new project from template"
        echo "  start           Start serving current project"
        echo "  stop            Stop serving current project"
        echo "  list            Show all projects and status"
        echo "  open            Open current project in browser"
        echo "  tunnel          Create public Cloudflare tunnel"
        echo "  register        Register existing directory"
        echo "  remove [name]   Unregister project"
        echo "  sync            Rebuild registry from .dade files"
        echo "  proxy <cmd>     Manage proxy (start|stop|status)"
        echo "  help            Show this help"
        echo ""
    else
        # Fallback without gum
        echo "dade - scaffold and serve web projects"
        echo ""
        echo "Usage: dade <command> [args]"
        echo ""
        echo "Commands:"
        echo "  setup, install, uninstall, templates, update"
        echo "  new, start, stop, list, open, tunnel"
        echo "  register, remove, sync, proxy, help"
        echo ""
    fi
}
```

## main() Implementation

```bash
main() {
    local cmd="${1:-help}"
    shift || true
    
    case "$cmd" in
        setup)      cmd_setup "$@" ;;
        install)    cmd_install "$@" ;;
        uninstall)  cmd_uninstall "$@" ;;
        templates)  cmd_templates "$@" ;;
        update)     cmd_update "$@" ;;
        new)        cmd_new "$@" ;;
        start)      cmd_start "$@" ;;
        stop)       cmd_stop "$@" ;;
        list|ls)    cmd_list "$@" ;;
        open)       cmd_open "$@" ;;
        tunnel)     cmd_tunnel "$@" ;;
        register)   cmd_register "$@" ;;
        remove|rm)  cmd_remove "$@" ;;
        sync)       cmd_sync "$@" ;;
        proxy)      cmd_proxy "$@" ;;
        help|--help|-h) cmd_help ;;
        --version|-v) echo "dade v$DADE_VERSION" ;;
        *)
            log_error "Unknown command: $cmd"
            echo ""
            cmd_help
            exit 1
            ;;
    esac
}

main "$@"
```

## Acceptance Criteria

- [ ] All commands listed in help
- [ ] main() dispatches to correct handlers
- [ ] --version shows version
- [ ] Unknown command shows help

