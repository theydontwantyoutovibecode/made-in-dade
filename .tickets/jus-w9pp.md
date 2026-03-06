---
id: jus-w9pp
status: closed
deps: [jus-sibp, jus-qvpa, jus-woz2, jus-lkra, jus-jw2r, jus-vwdw]
links: []
created: 2026-02-17T15:12:55Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-fxaf
tags: [cli, fang, cobra, cleanup]
---
# Remove legacy arg parsing from main.go

Remove the legacy manual argument parsing from main.go after all commands have been migrated to Cobra.

## Background

dade is a CLI for scaffolding web projects. After migrating all commands to Cobra/Fang, the old switch-based dispatch code in main.go is no longer needed.

## Current State

File: cmd/dade/main.go

Contains:
- main() function
- run(args []string) int - old entry point
- runWithIO(args []string, stdout, stderr io.Writer, styled bool) int - switch dispatch
- helpText(styled bool) string - old help text generator

After Cobra migration, only main() calling Execute() is needed.

## Implementation

1. Simplify main.go to:

```go
package main

func main() {
    Execute()
}
```

2. Remove:
   - run() function
   - runWithIO() function (or move to test helper file if still needed)
   - helpText() function
   - All imports no longer needed

3. Update main_test.go:
   - Tests should use Cobra's testing patterns
   - Or keep runWithIO in a test helper file for backwards compatibility
   - Update to test via rootCmd.Execute() or similar

4. Verify no other files depend on removed functions

## Files to Modify

- cmd/dade/main.go: Simplify to just main() -> Execute()
- cmd/dade/main_test.go: Update test approach
- Possibly create cmd/dade/testing.go for test helpers

## Acceptance Criteria

- [ ] main.go contains only main() -> Execute()
- [ ] No unused code remains in main.go
- [ ] All tests pass
- [ ] go build ./cmd/dade succeeds
- [ ] dade --help works
- [ ] All subcommands work

