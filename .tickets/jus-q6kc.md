---
id: jus-q6kc
status: closed
deps: [jus-w9pp]
links: []
created: 2026-02-17T15:13:17Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-fxaf
tags: [cli, fang, cobra, testing]
---
# Update test suite for Cobra-based CLI

Update the test suite to work with the Cobra-based CLI architecture, ensuring all commands are testable in headless mode.

## Background

dade is a CLI for scaffolding web projects. After migrating to Cobra/Fang, tests need to be updated to use Cobra's testing patterns while maintaining good coverage.

## Current State

Test files:
- cmd/dade/main_test.go: Tests via runWithIO()
- cmd/dade/new_test.go: Tests via newCommand struct
- cmd/dade/install_test.go: Tests via installCommand struct
- cmd/dade/proxy_test.go: Tests via proxyCommand struct
- cmd/dade/setup_test.go: Tests via setupCommand struct
- cmd/dade/templates_test.go: Tests helpers and runWithIO
- cmd/dade/choose_template_test.go: Tests template selection

Current test patterns:
- Command structs with injectable dependencies (runner, spin, confirm, etc.)
- runWithIO for integration-style tests
- t.TempDir() and t.Setenv() for isolation

## Implementation

1. Create cmd/dade/testing_test.go with test helpers:

```go
package main

import (
    "bytes"
    "github.com/spf13/cobra"
)

func executeCommand(root *cobra.Command, args ...string) (string, string, error) {
    stdout := new(bytes.Buffer)
    stderr := new(bytes.Buffer)
    
    root.SetOut(stdout)
    root.SetErr(stderr)
    root.SetArgs(args)
    
    err := root.Execute()
    return stdout.String(), stderr.String(), err
}

func resetRootCmd() {
    // Reset any persistent state between tests
    rootCmd.SetArgs(nil)
    // Reset flag values to defaults
}
```

2. Update main_test.go:
   - Replace runWithIO calls with executeCommand(rootCmd, ...)
   - Test help, version, unknown commands via Cobra

3. Update command-specific tests:
   - Keep using command structs with dependency injection
   - Add tests for new flags (--yes, --quiet, etc.)
   - Test flag validation and mutual exclusion

4. Add flag-specific tests for each command:
   - Test that each flag is recognized
   - Test flag defaults
   - Test mutually exclusive flags error correctly

## Test Cases to Add

new command:
- [ ] --template flag selects template without prompt
- [ ] --name flag provides project name
- [ ] Missing name in non-TTY errors appropriately

setup command:
- [ ] --yes flag skips all prompts
- [ ] --check flag only checks dependencies
- [ ] --install-deps installs without prompting
- [ ] --migrate migrates without prompting
- [ ] Conflicting flags produce errors

templates command:
- [ ] --json flag outputs valid JSON

proxy command:
- [ ] Each subcommand is accessible
- [ ] proxy status --json outputs JSON

Global flags:
- [ ] --quiet suppresses output
- [ ] --verbose increases output
- [ ] --no-color disables ANSI

## Acceptance Criteria

- [ ] All existing tests pass
- [ ] New tests cover all flag behaviors
- [ ] Tests can run in CI (no TTY required)
- [ ] go test ./cmd/dade -v shows all tests
- [ ] Coverage maintained or improved

