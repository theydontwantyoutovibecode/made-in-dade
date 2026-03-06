---
id: jus-jngy
status: closed
deps: [jus-cys3]
links: []
created: 2026-02-12T19:55:44Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-rg8c
tags: [templates, parsing]
---
# Implement TOML parsing for manifests

Implement TOML parsing for dade.toml manifests in bash.

## Challenge

Bash doesn't have native TOML support. Options:

### Option A: Use a TOML parser tool
- Install 'tomlq' (part of yq) or 'dasel'
- Pros: Full TOML compliance
- Cons: Additional dependency

### Option B: Use Python one-liner
- Python has tomllib (3.11+) or tomli
- Pros: Usually available, accurate
- Cons: Requires Python

### Option C: Simple regex-based parsing
- Parse common patterns ourselves
- Pros: No dependencies
- Cons: Limited TOML support, fragile

## Recommended Approach

Use Python if available, fall back to simple parsing:

```bash
parse_toml_value() {
    local file="$1"
    local key="$2"  # e.g., "template.name" or "serve.type"
    
    if command -v python3 &>/dev/null; then
        python3 -c "
import tomllib
with open('$file', 'rb') as f:
    data = tomllib.load(f)
keys = '$key'.split('.')
val = data
for k in keys:
    val = val.get(k, '')
print(val if val else '')
"
    else
        # Fallback: simple grep-based parsing
        # Only works for simple key = "value" patterns
        ...
    fi
}
```

## Functions to Implement

```bash
parse_toml_value()       # Get single value by dotted key
parse_toml_array()       # Get array as newline-separated values
validate_manifest()      # Check required fields exist
get_template_name()      # Shorthand for template.name
get_serve_type()         # Shorthand for serve.type
get_serve_command()      # Get dev or prod command based on mode
```

## Caching

Consider caching parsed manifests to avoid repeated parsing:
- ~/.config/dade/cache/manifest-<template>.json
- Invalidate when template updated

## Acceptance Criteria

- [ ] TOML values can be read by dotted key path
- [ ] Arrays can be parsed
- [ ] Works with Python 3.11+
- [ ] Fallback works for simple manifests without Python
- [ ] Invalid TOML produces helpful error

