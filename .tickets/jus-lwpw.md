---
id: jus-lwpw
status: closed
deps: [jus-t35w]
links: []
created: 2026-02-12T20:30:07Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-ovh7
tags: [commands, templates]
---
# Implement show_official_templates() function

Implement the function that displays official templates available for installation.

## Function Implementation

```bash
show_official_templates() {
    echo ""
    if has_gum; then
        gum style --bold "Official Templates"
    else
        echo "Official Templates"
        echo "=================="
    fi
    echo ""
    
    echo "  django-hypermedia"
    echo "    Django + HTMX + TailwindCSS full-stack web application"
    echo "    dade install https://github.com/theydontwantyoutovibecode/dade-with-django-and-hypermedia.git"
    echo ""
    echo "  hypertext"
    echo "    Vanilla HTML/CSS/JS with HTMX - simple static site"
    echo "    dade install https://github.com/theydontwantyoutovibecode/dade-with-hypertext.git"
    echo ""
}

offer_official_templates() {
    local templates=(
        "django-hypermedia|Django + HTMX + TailwindCSS|https://github.com/theydontwantyoutovibecode/dade-with-django-and-hypermedia.git"
        "hypertext|Vanilla HTML/CSS/JS + HTMX|https://github.com/theydontwantyoutovibecode/dade-with-hypertext.git"
    )
    
    for entry in "${templates[@]}"; do
        IFS='|' read -r name desc url <<< "$entry"
        
        if [[ -d "$DADE_TEMPLATES_DIR/$name" ]]; then
            log_info "$name already installed"
        else
            if confirm "Install $name? ($desc)"; then
                cmd_install "$url"
            fi
        fi
    done
}
```

## Acceptance Criteria

- [ ] Lists all official templates
- [ ] Shows name, description, install command
- [ ] offer_official_templates prompts for each
- [ ] Skips already installed templates

