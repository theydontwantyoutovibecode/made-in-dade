---
id: jus-nsda
status: closed
deps: [jus-q6kc]
links: []
created: 2026-02-17T15:14:13Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-fxaf
tags: [cli, fang, cobra, ci, testing]
---
# Add CI test for headless operation

Add CI tests that verify all commands work in headless/non-TTY mode without prompts.

## Background

dade is a CLI for scaffolding web projects. A key goal of the Cobra/Fang migration is enabling headless operation in CI/CD pipelines. This ticket adds CI-specific tests.

## Current State

Tests run via go test, which doesn't have a TTY. But tests don't explicitly verify headless behavior.

## Implementation

1. Create cmd/dade/headless_test.go:

```go
package main

import (
    "os"
    "testing"
)

// TestHeadlessNew verifies 'new' command works without TTY
func TestHeadlessNew(t *testing.T) {
    // Ensure no TTY
    if term.IsTerminal(int(os.Stdin.Fd())) {
        t.Skip("test requires non-TTY stdin")
    }
    
    tmpDir := t.TempDir()
    t.Setenv("XDG_CONFIG_HOME", tmpDir)
    
    // Should succeed with all required args
    stdout, stderr, err := executeCommand(rootCmd, 
        "new", "testproj", "--template", "django-hypermedia")
    if err != nil {
        t.Errorf("new failed: %v\nstderr: %s", err, stderr)
    }
    
    // Should fail without required args (no prompt fallback)
    _, _, err = executeCommand(rootCmd, "new")
    if err == nil {
        t.Error("new without args should fail in non-TTY")
    }
}

func TestHeadlessSetup(t *testing.T) {
    // Test setup with --yes flag
    stdout, stderr, err := executeCommand(rootCmd,
        "setup", "--check")
    // etc.
}

func TestHeadlessTemplates(t *testing.T) {
    // Test templates with --json flag
}
```

2. Add CI workflow job that explicitly tests headless:

```yaml
headless-test:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v4
      with:
        go-version: '1.22'
    - name: Run headless tests
      run: |
        # Ensure no TTY
        go test -v -run 'Headless' ./cmd/dade/
```

3. Document headless usage in README or help

## Test Matrix

| Command | Headless Args | Expected |
|---------|---------------|----------|
| new | name + --template | Success |
| new | (none) | Error: name required |
| templates | --json | JSON output |
| templates | (none) | Text output |
| install | url | Success |
| install | --list-official | Success |
| install | (none) | Error: url required |
| setup | --check | Success |
| setup | -y | Success (if deps available) |
| setup | (none) | Error or skip prompts |
| proxy status | --json | JSON output |

## Acceptance Criteria

- [ ] Headless tests exist for all commands
- [ ] Tests verify correct error messages when args missing
- [ ] Tests pass in CI (no TTY)
- [ ] CI workflow includes headless test job
- [ ] README documents headless usage for CI/CD

