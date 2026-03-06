---
id: jus-nq0k
status: open
deps: []
links: []
created: 2026-02-12T19:52:57Z
type: epic
priority: 0
assignee: Alex Cabrera
tags: [architecture, planning]
---
# Unified dade CLI Architecture

Merge srv and dade into a single unified CLI that handles both project scaffolding and serving. Templates become installable plugins that define how projects are created and served.

## Vision

dade becomes a complete local development environment manager:
- Install templates as plugins from git repositories
- Scaffold new projects from installed templates
- Serve projects with automatic HTTPS via Caddy (*.localhost)
- Framework-aware serving (static, Django, Node, etc.)
- Cloudflare tunnel support for sharing

## Current State

- **dade** (theydontwantyoutovibecode/dade): CLI for scaffolding from templates
- **srv** (theydontwantyoutovibecode/srv): Static site dev server with Caddy/HTTPS
- **dade-with-django-and-hypermedia**: Django template with complex start.sh
- **dade-with-hypertext**: Simple HTML/HTMX template

## Target State

Single dade CLI that:
1. Replaces both dade and srv
2. Has plugin architecture for templates
3. Serves any project type via template-defined runners
4. Maintains HTTPS routing via central Caddy proxy

## Key Commands (Target)

```
dade install <git-url>     # Install template plugin
dade uninstall <name>      # Remove template plugin
dade templates             # List installed templates
dade new [name]            # Scaffold project (picker if needed)
dade start                 # Serve current project
dade stop                  # Stop server
dade open                  # Open in browser
dade list                  # Show all projects
dade tunnel                # Cloudflare tunnel
dade proxy start|stop      # Manage central proxy
```

## Affected Repositories

- theydontwantyoutovibecode/dade (primary changes)
- theydontwantyoutovibecode/srv (deprecated, merged)
- theydontwantyoutovibecode/dade-with-django-and-hypermedia (template manifest)
- theydontwantyoutovibecode/dade-with-hypertext (template manifest)
- theydontwantyoutovibecode/homebrew-tap (formula updates)

## Acceptance Criteria

- [ ] Single dade CLI replaces both dade and srv
- [ ] Templates installable as plugins from git URLs
- [ ] Projects served with HTTPS via *.localhost
- [ ] Framework-aware serving via template manifests
- [ ] All existing functionality preserved
- [ ] Homebrew formula updated
- [ ] srv repository archived with redirect

