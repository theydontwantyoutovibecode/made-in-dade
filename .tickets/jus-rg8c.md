---
id: jus-rg8c
status: closed
deps: [jus-cdr3]
links: []
created: 2026-02-12T19:53:19Z
type: epic
priority: 1
assignee: Alex Cabrera
parent: jus-nq0k
tags: [templates, plugins, architecture]
---
# Template Plugin System

Implement a plugin architecture where templates are installed from git repositories and stored locally. Each template includes a manifest (dade.toml) that defines metadata and how the project should be served.

## Plugin Storage

Templates installed to: ~/.config/dade/templates/<name>/

## Template Manifest (dade.toml)

```toml
[template]
name = "django-hypermedia"
description = "Django + HTMX + TailwindCSS"
version = "1.0.0"

[scaffold]
# Files/dirs to exclude when copying template
exclude = [".git", "dade.toml", ".dade"]
# Script to run after scaffolding (optional)
setup = "./setup.sh"

[serve]
# How to serve this project type
type = "command"  # or "static"
dev = "./start.sh --dev"
prod = "./start.sh"
port_env = "PORT"  # env var to pass port

[serve.static]
# Used if type = "static" or as fallback
root = "."
```

## Plugin Commands

- dade install <git-url> [--name <name>]
- dade uninstall <name>
- dade templates (list installed)
- dade update <name> (git pull)
- dade update --all

## Template Discovery

When running 'dade new':
1. List installed templates
2. If multiple, show picker (gum choose)
3. Copy template to target directory
4. Write .dade project marker
5. Run scaffold.setup if defined

## Acceptance Criteria

- [ ] Templates installed from git URLs to ~/.config/dade/templates/
- [ ] dade.toml manifest parsed correctly
- [ ] Template picker works with multiple templates
- [ ] Templates can be updated via git pull
- [ ] Templates can be uninstalled cleanly

