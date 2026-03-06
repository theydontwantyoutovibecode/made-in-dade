---
id: jus-b7gh
status: closed
deps: [jus-nmo3]
links: []
created: 2026-03-02T01:58:53Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-g3fm
tags: [read-only, binary, ws-7]
---
# Add .read-only sync to dade binary

Create a syncReadOnlyDeps function in the dade binary (internal/fsutil or new internal/readonly package). It reads .read-only/manifest.txt from the project directory, shallow-clones each listed repo into .read-only/<name>/. Call this from cmd_dev.go before starting the dev server. Show spinner per clone. Skip already-cloned repos. Gracefully handle network errors (warn, don't fail). Remove the sync_read_only_deps() function from the django template's start.sh after this is in the binary.

## Acceptance Criteria

1. Function in binary clones repos from manifest.txt. 2. Called during dade dev. 3. Django start.sh sync removed. 4. Tests: manifest exists, no manifest, partial clones, git failure.

