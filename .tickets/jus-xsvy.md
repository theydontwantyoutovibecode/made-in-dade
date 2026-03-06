---
id: jus-xsvy
status: closed
deps: []
links: []
created: 2026-02-17T15:10:33Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-fxaf
tags: [cli, fang, cobra, dependencies]
---
# Add Cobra and Fang dependencies to go.mod

Add github.com/spf13/cobra and github.com/charmbracelet/fang as dependencies to the dade Go module.

## Background

dade is a CLI tool for scaffolding web application projects. It currently uses manual argument parsing in cmd/dade/main.go. We are migrating to Cobra (command framework) wrapped by Fang (Charm's CLI starter kit with styled help, manpages, and completions).

## Current State

- go.mod located at /Users/acabrera/Code/dade/go.mod
- Current dependencies include charmbracelet/huh, charmbracelet/lipgloss, charmbracelet/log, charmbracelet/bubbletea
- CLI entry point: cmd/dade/main.go

## Implementation

1. Run: go get github.com/spf13/cobra@latest
2. Run: go get github.com/charmbracelet/fang@latest
3. Run: go mod tidy
4. Verify go.mod and go.sum are updated correctly
5. Run: go build ./cmd/dade to verify no build errors

## Acceptance Criteria

- [ ] github.com/spf13/cobra appears in go.mod
- [ ] github.com/charmbracelet/fang appears in go.mod
- [ ] go build ./cmd/dade succeeds
- [ ] go test ./... passes (existing tests unaffected)

