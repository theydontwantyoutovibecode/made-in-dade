---
id: jus-ovh7
status: closed
deps: [jus-jngy]
links: []
created: 2026-02-12T19:55:57Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-rg8c
tags: [commands, templates]
---
# Implement 'dade install' command

Implement the install command for adding template plugins from git repositories.

## Usage

```
dade install <git-url>              # Install template
dade install <git-url> --name foo   # Install with custom name
dade install --list-official        # Show official templates
```

## Implementation

```bash
cmd_install() {
    local url="$1"
    local name=""
    
    # Parse --name option
    # ...
    
    # Clone to temp directory first
    local tmp_dir
    tmp_dir=$(mktemp -d)
    git clone --depth 1 "$url" "$tmp_dir"
    
    # Validate manifest exists
    if [[ ! -f "$tmp_dir/dade.toml" ]]; then
        log_error "Template missing dade.toml manifest"
        rm -rf "$tmp_dir"
        exit 1
    fi
    
    # Parse name from manifest if not provided
    if [[ -z "$name" ]]; then
        name=$(parse_toml_value "$tmp_dir/dade.toml" "template.name")
    fi
    
    # Check if already installed
    local target="$DADE_TEMPLATES_DIR/$name"
    if [[ -d "$target" ]]; then
        if confirm "Template '$name' already installed. Update?"; then
            rm -rf "$target"
        else
            rm -rf "$tmp_dir"
            exit 0
        fi
    fi
    
    # Move to templates directory
    mv "$tmp_dir" "$target"
    
    # Store source URL for updates
    echo "$url" > "$target/.source"
    
    log_success "Installed template: $name"
}
```

## Official Templates

Maintain list of official templates:

```bash
OFFICIAL_TEMPLATES=(
    "https://github.com/theydontwantyoutovibecode/dade-with-django-and-hypermedia.git"
    "https://github.com/theydontwantyoutovibecode/dade-with-hypertext.git"
)
```

## Template Directory Structure

```
~/.config/dade/templates/
├── django-hypermedia/
│   ├── .source              # URL for updates
│   ├── dade.toml       # Manifest
│   └── ...                  # Template files
└── hypertext/
    ├── .source
    ├── dade.toml
    └── ...
```

## Acceptance Criteria

- [ ] Templates cloned from git URLs
- [ ] Manifest validated before install
- [ ] Name derived from manifest or --name flag
- [ ] Existing templates can be updated
- [ ] Source URL stored for later updates
- [ ] --list-official shows available templates

