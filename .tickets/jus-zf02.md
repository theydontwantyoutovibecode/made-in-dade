---
id: jus-zf02
status: closed
deps: [jus-cys3]
links: []
created: 2026-02-12T19:59:57Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-2e3m
tags: [templates, static]
---
# Add dade.toml to hypertext template

Add dade.toml manifest to the dade-with-hypertext template.

## Repository

theydontwantyoutovibecode/dade-with-hypertext

## Manifest Content

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
# No setup script needed for static sites

[serve]
type = "static"
# No commands needed - dade serves directly via Caddy

[serve.static]
root = "."
```

## Changes to Make

1. **Remove serve.sh**
   - No longer needed, dade handles serving
   - Static sites served directly by Caddy

2. **Update AGENTS.md**
   - Replace serve.sh references with dade commands
   - Document: dade start, dade open, dade tunnel

3. **Update README.md**
   - Installation via dade
   - Development workflow with dade
   - Keep fallback instructions (python -m http.server)

4. **Add .dade to .gitignore**

## Standalone Fallback

For users without dade, document:
```bash
# Without dade
python3 -m http.server 8000
# Then visit http://localhost:8000
```

## Acceptance Criteria

- [ ] dade.toml created with static type
- [ ] serve.sh removed
- [ ] AGENTS.md updated with dade commands
- [ ] README.md documents dade and standalone usage
- [ ] .dade added to .gitignore

