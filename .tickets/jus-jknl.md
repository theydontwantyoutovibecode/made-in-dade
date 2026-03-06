---
id: jus-jknl
status: open
deps: [jus-syr2]
links: []
created: 2026-02-23T16:55:01Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-3el5
tags: [hooks, extensibility]
---
# Add hook discovery and execution system

Allow templates to provide hook scripts in project directories.

Hook directory structure:
  .dade/
  ├── dev-setup.sh      # Custom setup (overrides [dev.setup])
  ├── dev-teardown.sh   # Cleanup on exit
  ├── share-setup.sh    # Additional share setup
  └── hooks.toml        # Hook configuration

Hook discovery:
1. Check for .dade/ directory in project
2. Look for known hook scripts
3. Parse hooks.toml for custom hooks

Hook execution:
- Run scripts with project directory as cwd
- Pass environment variables
- Capture output with prefixes

hooks.toml format:
  [dev]
  setup = "./custom-setup.sh"
  teardown = "./custom-teardown.sh"
  
  [share]
  setup = "./share-prep.sh"

Files to create:
- internal/lifecycle/hooks.go
- internal/lifecycle/hooks_test.go

Acceptance:
- Discovers hooks in .dade directory
- Executes hooks at appropriate lifecycle points
- Hooks override manifest commands when present

