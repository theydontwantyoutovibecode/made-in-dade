---
id: jus-syr2
status: closed
deps: []
links: []
created: 2026-02-23T16:54:28Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-3el5
tags: [core, lifecycle]
---
# Create internal/lifecycle package for setup/background process management

New package to manage dev server lifecycle.

Components:
1. setup.go - Sequential command runner for setup phase
   - Run commands in order, fail on first error
   - Support both command strings and script files
   - Pass environment variables
   - Log output with prefixes

2. background.go - Background process manager  
   - Start multiple background processes
   - Track PIDs
   - Forward signals appropriately
   - Cleanup on exit

3. cleanup.go - Signal handling and cleanup
   - SIGINT/SIGTERM handlers
   - Kill background processes
   - Run teardown hooks
   - Clean exit

Files to create:
- internal/lifecycle/setup.go
- internal/lifecycle/background.go  
- internal/lifecycle/cleanup.go
- internal/lifecycle/lifecycle_test.go

Acceptance:
- Can run setup commands sequentially
- Can manage background processes
- Clean shutdown on Ctrl+C
- Unit tests for core logic

