---
id: jus-d3wz
status: closed
deps: [jus-w9pp, jus-njz0]
links: []
created: 2026-02-17T15:13:55Z
type: task
priority: 3
assignee: Alex Cabrera
parent: jus-fxaf
tags: [cli, fang, cobra, documentation]
---
# Add manpage generation

Verify and document manpage generation provided by Fang.

## Background

dade is a CLI for scaffolding web projects. Fang adds a hidden 'man' command that generates manpages using mango for better roff output.

## Current State

Fang automatically adds:
- dade man (hidden command)

This generates a single combined manpage rather than one per subcommand.

## Implementation

1. Verify man command exists:
   - Run dade man
   - Check output is valid roff format

2. Test manpage rendering:
   - dade man | man -l -
   - Verify formatting looks correct
   - Check all commands and flags documented

3. Add to build/release process:
   - Generate manpage during build
   - Include in release artifacts
   - Document installation location

4. Optional: Add manpage to Homebrew formula

## Manpage Installation

```bash
# Generate manpage
dade man > dade.1

# Install (macOS)
sudo cp dade.1 /usr/local/share/man/man1/
sudo mandb  # or makewhatis on some systems

# View
man dade
```

## Acceptance Criteria

- [ ] dade man generates valid roff output
- [ ] Manpage renders correctly with man -l -
- [ ] All commands documented in manpage
- [ ] All flags documented in manpage
- [ ] Build process can generate manpage
- [ ] Installation instructions documented

