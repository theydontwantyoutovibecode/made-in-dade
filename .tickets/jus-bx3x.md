---
id: jus-bx3x
status: closed
deps: []
links: []
created: 2026-02-12T20:34:31Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-swo4
tags: [templates, django]
---
# Create dade.toml for django-hypermedia

Create the dade.toml manifest file for the Django template.

## File: dade.toml

```toml
[template]
name = "django-hypermedia"
description = "Django + HTMX + TailwindCSS full-stack web application"
version = "1.0.0"
author = "Alex Cabrera"
url = "https://github.com/theydontwantyoutovibecode/dade-with-django-and-hypermedia"

[scaffold]
exclude = [
    ".git",
    "dade.toml",
    ".dade",
    "__pycache__",
    "*.pyc",
    ".venv",
    "db.sqlite3",
    "staticfiles",
    "node_modules",
    ".read-only"
]
setup = "./setup.sh"
setup_interactive = true

[serve]
type = "command"
dev = "./start.sh --dev"
prod = "./start.sh"
port_env = "PORT"
default_port = 8000
```

## Location

/Users/acabrera/Code/dade/dade-with-django-and-hypermedia/dade.toml

## Acceptance Criteria

- [ ] dade.toml created
- [ ] All fields populated correctly
- [ ] Validated with validate_manifest()

