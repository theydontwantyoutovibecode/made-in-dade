---
id: jus-x52t
status: closed
deps: [jus-4cer]
links: []
created: 2026-02-17T02:05:24Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-gyx3
---
# Scaffold Go module and CLI entrypoint for dade

Initialize Go module to replace the Bash script. Create cmd/dade/main.go with CLI entrypoint using standard library flag or cobra-free minimal parsing. Expose commands: new, templates, --help/-h, --version/-v. Wire version constant to match current DADE_VERSION=0.1.0. Provide structured project layout (cmd/, internal/ or pkg/) for reusable logic. Include graceful error handling with exit codes matching Bash (unknown command -> exit 1). Ensure root command prints header/help similar to Bash using lipgloss for styling, with plain fallback when non-tty. Add go.mod with module path github.com/theydontwantyoutovibecode/dade or existing repo path, set Go version (>=1.21). Add minimal README stub or comment noting Go build instructions placeholder to be completed later. No business logic yet—just skeleton wiring ready for subcommands to plug in.

## Acceptance Criteria

- go.mod created with correct module path and Go version\n- cmd/dade/main.go implements CLI dispatch for new/templates/help/version with proper exit codes\n- Version flag returns 0 and prints v0.1.0\n- Unknown command exits 1 with error message and help\n- Header/help styled via lipgloss when TTY; plain fallback otherwise\n- Version/help output includes full template list and config snippet\n- Project builds with go build ./...\n- Layout ready for internal packages for config/templates/commands

