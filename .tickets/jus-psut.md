---
id: jus-psut
status: closed
deps: []
links: []
created: 2026-02-12T20:35:00Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-zf02
tags: [templates, static]
---
# Create dade.toml for hypertext

Create the dade.toml manifest file for the hypertext template.

## File: dade.toml

```toml
[template]
name = "hypertext"
description = "Vanilla HTML/CSS/JS with HTMX - simple static site"
version = "1.0.0"
author = "Alex Cabrera"
url = "https://github.com/theydontwantyoutovibecode/dade-with-hypertext"

[scaffold]
exclude = [
    ".git",
    "dade.toml",
    ".dade",
    ".DS_Store"
]

[serve]
type = "static"

[serve.static]
root = "."
```

## Location

/Users/acabrera/Code/dade/dade-with-hypertext/dade.toml

## Acceptance Criteria

- [ ] dade.toml created
- [ ] serve.type is static
- [ ] Validated correctly

