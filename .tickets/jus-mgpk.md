---
id: jus-mgpk
status: closed
deps: [jus-52py]
links: []
created: 2026-02-17T02:05:36Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-gyx3
---
# Implement 'new' command behavior in Go

Port the 'new' command: parse args [project-name] [--local PATH] [--template NAME]. Print the header immediately after parsing to match Bash (before prompts/validation). If project-name is missing, prompt via UI layer (huh input with placeholder “my-project”) when TTY; otherwise error. Validate project name with regex ^[a-zA-Z][a-zA-Z0-9_-]*$. Unknown flags should error. If directory exists, fail. When --local PATH provided, copy directory to project_name, then remove project_name/.git. Otherwise clone remote template URL from config/defaults with git clone --depth 1 into project_name; if template unknown, error. After copy/clone, rm -rf project_name/.git. Initialize new git repo (git init). If project_name/setup.sh exists and is executable, run it; if only file (non-exec), run via bash and error clearly if bash is missing. Provide progress/status logs (INFO/SUCCESS/WARN/ERROR analogs). Return exit codes on failure; on success print next steps (cd project_name; ./start.sh --dev). Include non-interactive behavior matching Bash; interactive template selection handled in separate ticket. Add tests using temp dirs, mocking exec for git/cp? design for abstraction so we can inject command runner. Document behavior parity with Bash script.

## Acceptance Criteria

- Header printed before prompts/validation\n- new command validates names and existing dir guard\n- Prompts for project name when missing (TTY only); errors when non-interactive\n- Unknown options error with exit 1\n- Supports --local copy with .git removal\n- Supports remote clone via configured template map with unknown template error\n- Runs setup.sh if present (exec vs bash based on file mode) and errors if bash missing\n- Initializes git repo\n- Prints next steps\n- Tests cover validation, local copy, remote clone failure path, setup.sh exec/batch

