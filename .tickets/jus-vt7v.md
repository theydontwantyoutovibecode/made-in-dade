---
id: jus-vt7v
status: closed
deps: [jus-f3br]
links: []
created: 2026-02-12T19:54:18Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-cdr3
tags: [infrastructure, registry]
---
# Implement project registry (projects.json)

Implement the project registry that tracks all dade-managed projects.

## Schema

```json
{
  "myproject": {
    "port": 3000,
    "path": "/Users/alex/Code/myproject",
    "template": "django-hypermedia",
    "created": "2024-01-15T10:30:00Z"
  }
}
```

## Functions to Implement

```bash
# Port management
next_port()              # Find next available port (start at 3000)
is_port_available()      # Check if port is free

# Project CRUD
register_project()       # Add project to registry
unregister_project()     # Remove project from registry
get_project()            # Get project by name
get_project_by_path()    # Get project by directory path
list_projects()          # List all projects

# Lookup helpers
get_project_port()       # Get port for project name
get_project_path()       # Get path for project name
get_project_template()   # Get template for project name
```

## Considerations

- Use jq for JSON manipulation (already a dependency)
- Atomic writes (write to temp, then mv)
- Handle concurrent access gracefully
- Validate JSON on read

## Acceptance Criteria

- [ ] Projects can be registered with name, port, path, template
- [ ] Port assignment finds next available port
- [ ] Projects can be looked up by name or path
- [ ] Registry survives invalid JSON (recovery mode)
- [ ] All CRUD operations work correctly

