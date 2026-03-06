---
id: jus-ikny
status: closed
deps: []
links: []
created: 2026-02-12T20:29:27Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-jngy
tags: [templates, parsing]
---
# Implement parse_toml_value() function

Implement the function that parses values from TOML files.

## Function Implementation

```bash
parse_toml_value() {
    local file="$1"
    local key="$2"  # e.g., "template.name" or "serve.type"
    
    if [[ ! -f "$file" ]]; then
        return 1
    fi
    
    # Use Python if available (most reliable)
    if command -v python3 &>/dev/null; then
        python3 << EOF
import sys
try:
    import tomllib
except ImportError:
    import tomli as tomllib

try:
    with open('$file', 'rb') as f:
        data = tomllib.load(f)
    keys = '$key'.split('.')
    val = data
    for k in keys:
        if isinstance(val, dict):
            val = val.get(k, '')
        else:
            val = ''
            break
    if val and not isinstance(val, (list, dict)):
        print(val)
except Exception:
    pass
EOF
        return
    fi
    
    # Fallback: simple grep-based parsing for key = "value"
    local section=""
    local target_section=""
    local target_key=""
    
    if [[ "$key" == *.* ]]; then
        target_section="${key%%.*}"
        target_key="${key#*.}"
    else
        target_key="$key"
    fi
    
    local in_section=false
    while IFS= read -r line; do
        # Check for section header
        if [[ "$line" =~ ^\[([a-zA-Z0-9._-]+)\] ]]; then
            section="${BASH_REMATCH[1]}"
            if [[ "$section" == "$target_section" ]]; then
                in_section=true
            else
                in_section=false
            fi
            continue
        fi
        
        # Check for key = value
        if [[ -z "$target_section" ]] || $in_section; then
            if [[ "$line" =~ ^\ *$target_key\ *=\ *\"(.*)\" ]]; then
                echo "${BASH_REMATCH[1]}"
                return
            fi
        fi
    done < "$file"
}
```

## Acceptance Criteria

- [ ] Parses simple key = "value" patterns
- [ ] Handles nested sections (template.name)
- [ ] Uses Python when available
- [ ] Falls back to grep for simple cases
- [ ] Returns empty string for missing keys

