---
id: jus-3bvb
status: closed
deps: [jus-hbgh]
links: []
created: 2026-02-12T20:31:07Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-keyx
tags: [commands, templates]
---
# Implement copy_template() function

Implement the function that copies template files to a new project.

## Function Implementation

```bash
copy_template() {
    local src="$1"
    local dest="$2"
    local manifest="$src/dade.toml"
    
    mkdir -p "$dest"
    
    # Get exclude patterns
    local excludes=()
    excludes+=(".git")
    excludes+=("dade.toml")
    excludes+=(".source")
    
    # Add patterns from manifest
    while IFS= read -r pattern; do
        [[ -n "$pattern" ]] && excludes+=("$pattern")
    done < <(parse_toml_array "$manifest" "scaffold.exclude")
    
    # Build rsync exclude args
    local rsync_excludes=()
    for pattern in "${excludes[@]}"; do
        rsync_excludes+=("--exclude=$pattern")
    done
    
    # Copy files
    rsync -a "${rsync_excludes[@]}" "$src/" "$dest/"
}
```

## Alternative Without rsync

```bash
copy_template_simple() {
    local src="$1"
    local dest="$2"
    
    # Copy everything except known excludes
    mkdir -p "$dest"
    
    for item in "$src"/*; do
        local name=$(basename "$item")
        case "$name" in
            dade.toml|.git|.source)
                continue
                ;;
            *)
                cp -R "$item" "$dest/"
                ;;
        esac
    done
    
    # Copy hidden files except .git
    for item in "$src"/.*; do
        local name=$(basename "$item")
        case "$name" in
            .|..|.git|.source)
                continue
                ;;
            *)
                cp -R "$item" "$dest/"
                ;;
        esac
    done
}
```

## Acceptance Criteria

- [ ] Copies template files to destination
- [ ] Excludes .git directory
- [ ] Excludes dade.toml
- [ ] Respects scaffold.exclude patterns
- [ ] Copies hidden files (except excluded)

