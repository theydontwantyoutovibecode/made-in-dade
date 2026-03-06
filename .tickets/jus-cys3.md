---
id: jus-cys3
status: closed
deps: []
links: []
created: 2026-02-12T19:55:30Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-rg8c
tags: [templates, schema, design]
---
# Define dade.toml manifest schema

Define and document the dade.toml manifest schema that templates must include.

## Full Schema

```toml
# Template metadata
[template]
name = "django-hypermedia"           # Required: unique identifier
description = "Django + HTMX + TailwindCSS"  # Required: human-readable
version = "1.0.0"                     # Optional: semver
author = "Alex Cabrera"               # Optional
url = "https://github.com/..."        # Optional: source URL

# Scaffolding configuration
[scaffold]
# Files/directories to exclude when copying template to new project
exclude = [
    ".git",
    "dade.toml",
    ".dade",
    "__pycache__",
    "*.pyc",
    ".DS_Store"
]
# Script to run after copying files (relative to project root)
setup = "./setup.sh"
# Whether setup script is interactive (requires TTY)
setup_interactive = true

# Serving configuration
[serve]
# Type: "static" or "command"
type = "command"
# Command to run in dev mode (only for type=command)
dev = "./start.sh --dev"
# Command to run in prod mode (only for type=command)
prod = "./start.sh"
# Environment variable name for port
port_env = "PORT"
# Default port preference (optional, will be assigned if taken)
default_port = 8000

# Static serving fallback (used if type=static or as fallback)
[serve.static]
root = "."
# Additional file extensions to serve (beyond defaults)
extensions = []

# Optional: project marker file content
[project]
# Fields to include in .dade project file
marker_fields = ["template", "port", "created"]
```

## Minimal Valid Manifest

```toml
[template]
name = "my-template"
description = "My awesome template"

[serve]
type = "static"
```

## Validation Rules

1. [template] section required
2. template.name required, must match [a-z0-9-]+
3. template.description required
4. [serve] section required
5. serve.type must be "static" or "command"
6. If serve.type = "command", serve.dev or serve.prod required

## Acceptance Criteria

- [ ] Schema documented in dade README
- [ ] Validation function implemented
- [ ] Helpful error messages for invalid manifests
- [ ] Example manifests for static and command types

