---
id: jus-8m18
status: closed
deps: []
links: []
created: 2026-02-17T02:05:58Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-gyx3
---
# Implement process execution and filesystem helpers

Abstract shell interactions for Go port: run commands (git clone --depth 1, git init, cp -R), detect executable files, remove .git directories, create dirs, and copy trees. Provide a small exec runner with injectable interfaces for testing (allow fakes). Ensure safe path handling and clear error messages. Include helpers for checking command availability (git, bash when needed). Tests should cover success/failure cases, path handling, .git removal, exec detection logic.

## Acceptance Criteria

- Exec helper runs git clone/init and returns errors with context\n- File helpers copy directories and remove .git safely\n- Detect executable setup.sh to decide direct exec vs bash\n- Checks for git and bash when required\n- Tests cover success/failure and edge cases

