---
id: jus-3cof
status: closed
deps: []
links: []
created: 2026-02-12T20:27:25Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-f3br
tags: [infrastructure, setup]
---
# Define constants and global variables

Define all global constants and variables at the top of the dade script.

## Variables to Define

```bash
# Version
DADE_VERSION="1.0.0"

# Configuration paths
DADE_CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/dade"
DADE_PROJECTS_FILE="${DADE_CONFIG_DIR}/projects.json"
DADE_TEMPLATES_DIR="${DADE_CONFIG_DIR}/templates"
DADE_CADDYFILE="${DADE_CONFIG_DIR}/Caddyfile"
DADE_PLIST="${HOME}/Library/LaunchAgents/land.charm.dade.proxy.plist"
DADE_LOG="${DADE_CONFIG_DIR}/proxy.log"
DADE_ERR="${DADE_CONFIG_DIR}/proxy.err"

# Defaults
DADE_BASE_PORT=3000
DADE_PROXY_LABEL="land.charm.dade.proxy"

# Official templates
DADE_OFFICIAL_TEMPLATES=(
    "https://github.com/theydontwantyoutovibecode/dade-with-django-and-hypermedia.git"
    "https://github.com/theydontwantyoutovibecode/dade-with-hypertext.git"
)
```

## Considerations

- Respect XDG_CONFIG_HOME if set
- Use consistent naming convention
- Group related variables together

## Acceptance Criteria

- [ ] All paths defined as variables
- [ ] XDG_CONFIG_HOME respected
- [ ] Official templates array defined
- [ ] Version variable matches release plan

