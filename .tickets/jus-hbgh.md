---
id: jus-hbgh
status: closed
deps: [jus-ikny]
links: []
created: 2026-02-12T20:29:35Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-jngy
tags: [templates, parsing]
---
# Implement parse_toml_array() function

Implement the function that parses array values from TOML files.

## Function Implementation

```bash
parse_toml_array() {
    local file="$1"
    local key="$2"  # e.g., "scaffold.exclude"
    
    if [[ ! -f "$file" ]]; then
        return 1
    fi
    
    # Use Python (arrays are complex to parse in bash)
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
            val = val.get(k, [])
        else:
            val = []
            break
    if isinstance(val, list):
        for item in val:
            print(item)
except Exception:
    pass
EOF
    else
        # Fallback: can't reliably parse arrays in bash
        log_warn "Python required for array parsing"
    fi
}
```

## Usage

```bash
# Get exclude patterns
while IFS= read -r pattern; do
    echo "Exclude: $pattern"
done < <(parse_toml_array "dade.toml" "scaffold.exclude")
```

## Acceptance Criteria

- [ ] Parses TOML arrays
- [ ] Returns one item per line
- [ ] Handles empty arrays
- [ ] Requires Python (documented)

