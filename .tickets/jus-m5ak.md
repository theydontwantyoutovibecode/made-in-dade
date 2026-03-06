---
id: jus-m5ak
status: closed
deps: [jus-nsda, jus-3ibt, jus-d3wz, jus-njz0, jus-q6kc]
links: []
created: 2026-02-17T15:14:30Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-fxaf
tags: [cli, fang, cobra, verification]
---
# Final integration verification and cleanup

Perform final verification that the Cobra/Fang migration is complete and working correctly.

## Background

dade is a CLI for scaffolding web projects. After all migration tickets are complete, this ticket performs final verification.

## Verification Checklist

### Build & Test
- [ ] go build ./cmd/dade succeeds
- [ ] go test ./... passes
- [ ] No compiler warnings
- [ ] No linter issues (if configured)

### Command Verification
Run each command and verify output:

```bash
# Help
dade --help
dade new --help
dade templates --help
dade install --help
dade proxy --help
dade setup --help

# Version
dade --version

# Templates (no setup required)
dade templates
dade templates --json

# Install (requires git)
dade install --list-official

# New (full test)
dade new testproj --template django-hypermedia
rm -rf testproj

# Proxy
dade proxy status
dade proxy status --json

# Setup (careful - modifies system)
dade setup --check
```

### Headless Verification
Run in non-TTY context:

```bash
echo 'dade new testproj -t django-hypermedia' | bash
dade templates --json | jq .
dade proxy status --json | jq .
```

### Completion Verification
```bash
dade completion bash > /dev/null
dade completion zsh > /dev/null
dade completion fish > /dev/null
```

### Manpage Verification
```bash
dade man | man -l -
```

## Cleanup Tasks

- [ ] Remove any TODO comments from migration
- [ ] Remove deprecated code paths
- [ ] Update README with new CLI documentation
- [ ] Update CHANGELOG if exists

## Acceptance Criteria

- [ ] All verification checklist items pass
- [ ] No migration-related TODOs remain
- [ ] README reflects new CLI structure
- [ ] Ready for release

