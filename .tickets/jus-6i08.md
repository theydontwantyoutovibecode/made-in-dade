---
id: jus-6i08
status: open
deps: []
links: []
created: 2026-03-02T04:53:12Z
type: task
priority: 2
assignee: Alex Cabrera
tags: [cli, restructure]
---
# Add cobra command groups to organize help output

Use cobra's AddGroup and GroupID to visually organize the help output into sections:

  DEVELOPMENT
    new        Create a new project from a template
    dev        Start development server  
    build      Build a compiled project
    start      Start production server
    stop       Stop a running project
    open       Open project in browser
    share      Share project via public tunnel

  MANAGEMENT
    project    Manage project registry
    template   Manage installed templates
    proxy      Manage HTTPS proxy
    setup      First-time setup

This makes the help output scannable. Development commands (used daily) are visually separated from management commands (used occasionally).

Implementation:
- Call rootCmd.AddGroup with two groups in root.go
- Set GroupID on each command's cobra.Command definition
- Customize the help template if needed to render groups cleanly

