---
id: jus-lvnb
status: closed
deps: [jus-ikny]
links: []
created: 2026-02-12T20:29:45Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-jngy
tags: [templates, validation]
---
# Implement validate_manifest() function

Implement the function that validates a template manifest.

## Function Implementation

```bash
validate_manifest() {
    local manifest="$1"
    local errors=()
    
    if [[ ! -f "$manifest" ]]; then
        log_error "Manifest not found: $manifest"
        return 1
    fi
    
    # Required: template.name
    local name=$(parse_toml_value "$manifest" "template.name")
    if [[ -z "$name" ]]; then
        errors+=("Missing required field: template.name")
    elif [[ ! "$name" =~ ^[a-z0-9-]+$ ]]; then
        errors+=("Invalid template.name: must be lowercase alphanumeric with hyphens")
    fi
    
    # Required: template.description
    local desc=$(parse_toml_value "$manifest" "template.description")
    if [[ -z "$desc" ]]; then
        errors+=("Missing required field: template.description")
    fi
    
    # Required: serve.type
    local serve_type=$(parse_toml_value "$manifest" "serve.type")
    if [[ -z "$serve_type" ]]; then
        errors+=("Missing required field: serve.type")
    elif [[ "$serve_type" != "static" && "$serve_type" != "command" ]]; then
        errors+=("Invalid serve.type: must be 'static' or 'command'")
    fi
    
    # If command type, require serve.dev or serve.prod
    if [[ "$serve_type" == "command" ]]; then
        local dev=$(parse_toml_value "$manifest" "serve.dev")
        local prod=$(parse_toml_value "$manifest" "serve.prod")
        if [[ -z "$dev" && -z "$prod" ]]; then
            errors+=("Command type requires serve.dev or serve.prod")
        fi
    fi
    
    # Report errors
    if [[ ${#errors[@]} -gt 0 ]]; then
        log_error "Invalid manifest: $manifest"
        for err in "${errors[@]}"; do
            echo "  - $err"
        done
        return 1
    fi
    
    return 0
}
```

## Acceptance Criteria

- [ ] Validates required fields
- [ ] Checks template.name format
- [ ] Validates serve.type values
- [ ] Checks command type has dev/prod
- [ ] Reports all errors clearly

