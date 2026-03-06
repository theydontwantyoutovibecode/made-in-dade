---
id: jus-3ibt
status: closed
deps: [jus-w9pp, jus-njz0]
links: []
created: 2026-02-17T15:13:46Z
type: task
priority: 3
assignee: Alex Cabrera
parent: jus-fxaf
tags: [cli, fang, cobra, completions]
---
# Add shell completion generation

Verify and document shell completion generation provided by Fang/Cobra.

## Background

dade is a CLI for scaffolding web projects. Cobra provides built-in shell completion generation, and Fang adds a 'completion' command automatically.

## Current State

Fang automatically adds:
- dade completion bash
- dade completion zsh
- dade completion fish
- dade completion powershell

These should work out of the box after Fang integration.

## Implementation

1. Verify completion command exists:
   - Run dade completion --help
   - Run dade completion bash
   - Run dade completion zsh
   - Run dade completion fish

2. Test completions actually work:
   - Source bash completion and test
   - Source zsh completion and test
   - Verify flag and subcommand completion

3. Add documentation to help text or README:
   - Installation instructions for each shell
   - Where to put completion scripts

4. Optional: Customize completion descriptions if needed

## Completion Installation Instructions (for docs)

### Bash
```bash
# Add to ~/.bashrc
source <(dade completion bash)

# Or save to file
dade completion bash > /usr/local/etc/bash_completion.d/dade
```

### Zsh
```zsh
# Add to ~/.zshrc (before compinit)
source <(dade completion zsh)

# Or save to fpath
dade completion zsh > "${fpath[1]}/_dade"
```

### Fish
```fish
dade completion fish | source

# Or save permanently
dade completion fish > ~/.config/fish/completions/dade.fish
```

## Acceptance Criteria

- [ ] dade completion bash generates valid script
- [ ] dade completion zsh generates valid script
- [ ] dade completion fish generates valid script
- [ ] Completions work for commands (new, templates, install, proxy, setup)
- [ ] Completions work for flags (--template, --json, etc.)
- [ ] Documentation added to help or README

