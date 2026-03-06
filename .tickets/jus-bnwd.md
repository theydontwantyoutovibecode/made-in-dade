---
id: jus-bnwd
status: closed
deps: [jus-1p4n]
links: []
created: 2026-02-12T20:01:59Z
type: feature
priority: 3
assignee: Alex Cabrera
parent: jus-nq0k
tags: [ux, completion]
---
# Add shell completion for dade

Add shell completion scripts for bash and zsh.

## Completion Features

- Command completion (new, start, stop, list, etc.)
- Project name completion for commands that take project names
- Template name completion for --template flag
- Flag completion

## Bash Completion

```bash
_dade_completions() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev="${COMP_WORDS[COMP_CWORD-1]}"
    
    local commands="setup install uninstall templates update new start stop list open tunnel proxy register remove sync help"
    
    case "$prev" in
        dade)
            COMPREPLY=($(compgen -W "$commands" -- "$cur"))
            ;;
        start|stop|open|tunnel|remove)
            # Complete with project names
            local projects=$(jq -r 'keys[]' ~/.config/dade/projects.json 2>/dev/null)
            COMPREPLY=($(compgen -W "$projects" -- "$cur"))
            ;;
        --template|-t)
            # Complete with installed templates
            local templates=$(ls ~/.config/dade/templates/ 2>/dev/null)
            COMPREPLY=($(compgen -W "$templates" -- "$cur"))
            ;;
        proxy)
            COMPREPLY=($(compgen -W "start stop restart status logs" -- "$cur"))
            ;;
        *)
            ;;
    esac
}
complete -F _dade_completions dade
```

## Installation

Add to formula caveats:

```
For bash completion, add to ~/.bashrc:
  source /opt/homebrew/etc/bash_completion.d/dade

For zsh completion, add to ~/.zshrc:
  fpath=(/opt/homebrew/share/zsh/site-functions $fpath)
```

## Acceptance Criteria

- [ ] Bash completion script created
- [ ] Zsh completion script created
- [ ] Commands complete correctly
- [ ] Project names complete
- [ ] Template names complete
- [ ] Installed via Homebrew formula

