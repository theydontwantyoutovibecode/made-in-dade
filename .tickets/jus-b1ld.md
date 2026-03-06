---
id: jus-b1ld
status: done
deps: []
links: [jus-fqr9, jus-83w2, jus-mmae, jus-iu0c]
created: 2026-03-01T22:30:00Z
type: feature
priority: 0
assignee: Alex Cabrera
tags: [cli, build, cross-platform, blocking]
---
# Implement dade build subcommand

Create a new `dade build` command that compiles projects producing executables (Go CLI/TUI, Swift iOS, Kotlin Android) or runs necessary build steps for other project types.

## Background

Unlike existing templates (Django, static HTML) which are interpreted/served directly, the new CLI, TUI, iOS, Android, and hybrid mobile templates produce compiled artifacts. The build command provides a consistent interface for building these projects.

## Command Specification

```bash
# Basic usage
dade build [name]           # Build a registered project
dade build                  # Build current directory project

# Flags
dade build --output DIR     # Output directory (default: ./bin)
dade build --os OS          # Target OS (darwin, windows, linux)
dade build --arch ARCH      # Target architecture (amd64, arm64)
dade build --all            # Build for all supported platforms
dade build --release        # Optimized release build
dade build --verbose        # Show build output
```

## Project Type Detection

| Indicator | Project Type | Build Command |
|-----------|--------------|---------------|
| `go.mod` + `main.go` | Go CLI/TUI | `go build -o ./bin/{name} .` |
| `*.xcodeproj` | iOS/macOS | `xcodebuild -scheme {scheme} -configuration Release` |
| `build.gradle.kts` + `AndroidManifest.xml` | Android | `./gradlew assembleRelease` |
| `pubspec.yaml` | Flutter | `flutter build apk --release` |

## dade.toml Schema Extension

Add new `[build]` section:

```toml
[build]
command = "go build -o ./bin/{{name}} ."
command_windows = "go build -o .\\bin\\{{name}}.exe ."
output = "./bin"
targets = [
    { os = "darwin", arch = "arm64" },
    { os = "darwin", arch = "amd64" },
    { os = "windows", arch = "amd64" },
]
release_flags = "-ldflags '-s -w'"
pre = ["go mod tidy"]
post = []
```

## Acceptance Criteria

- [ ] Command auto-detects project type from directory structure
- [ ] Reads build configuration from `dade.toml` if present
- [ ] Supports Go projects with cross-compilation (GOOS/GOARCH)
- [ ] Creates output directory (`./bin/`) automatically
- [ ] `--release` flag produces optimized builds
- [ ] `--all` flag builds for all configured targets
- [ ] Works for both registered projects and current directory
- [ ] Manifest schema updated to include `[build]` section
- [ ] Integrates with existing `[prod]` section for fallback

## Implementation Notes

- File: `dade/cmd/dade/cmd_build.go`
- Update `internal/manifest/manifest.go` for build schema
- See PLAN.md for detailed implementation code
