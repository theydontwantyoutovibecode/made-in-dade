---
id: jus-jjl8
status: closed
deps: [jus-fgp8]
links: []
created: 2026-02-12T20:00:40Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-k78p
tags: [documentation, github]
---
# Update dade README for unified architecture

Rewrite the dade README to document the unified CLI architecture.

## New README Structure

```markdown
# dade

A CLI for scaffolding and serving web projects with automatic HTTPS.

## Features

- 🎨 **Template Plugins** - Install project templates from git
- 🚀 **One-Command Setup** - Scaffold and configure projects instantly
- 🔒 **Automatic HTTPS** - Every project gets https://name.localhost
- 🌐 **Cloudflare Tunnels** - Share projects publicly with one command
- 🎯 **Framework-Aware** - Serves Django, Node, static sites, and more

## Installation

\`\`\`bash
brew tap theydontwantyoutovibecode/tap
brew install dade

# First-time setup
dade setup
\`\`\`

## Quick Start

\`\`\`bash
# Install a template
dade install https://github.com/theydontwantyoutovibecode/dade-with-hypertext.git

# Create a new project
dade new mysite

# Start serving
cd mysite
dade start

# Open in browser
dade open

# Share publicly
dade tunnel
\`\`\`

## Commands

| Command | Description |
|---------|-------------|
| \`setup\` | First-time setup (install deps, configure proxy) |
| \`install <url>\` | Install template from git repository |
| \`templates\` | List installed templates |
| \`new [name]\` | Create new project from template |
| \`start\` | Start serving current project |
| \`stop\` | Stop serving current project |
| \`list\` | Show all projects and status |
| \`open\` | Open project in browser |
| \`tunnel\` | Create public tunnel |
| \`proxy <cmd>\` | Manage central HTTPS proxy |

## Official Templates

| Template | Description |
|----------|-------------|
| [django-hypermedia](https://github.com/theydontwantyoutovibecode/dade-with-django-and-hypermedia) | Django + HTMX + TailwindCSS |
| [hypertext](https://github.com/theydontwantyoutovibecode/dade-with-hypertext) | Vanilla HTML/CSS/JS + HTMX |

## Creating Templates

Templates are git repositories with a \`dade.toml\` manifest.
See [Template Development Guide](docs/templates.md).

## How It Works

dade runs a central Caddy proxy that routes \`*.localhost\` 
domains to your projects via HTTPS. Each project runs on its own 
port, and the proxy handles SSL termination with locally-trusted 
certificates.

## Requirements

- macOS (Linux support planned)
- bash 4.0+
- caddy, jq (installed via \`dade setup\`)
- gum (optional, for pretty UI)

## License

MIT
\`\`\`

## Acceptance Criteria

- [ ] README rewritten for unified architecture
- [ ] All commands documented
- [ ] Official templates listed
- [ ] Installation instructions clear
- [ ] Quick start guide works

