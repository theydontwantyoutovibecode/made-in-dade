---
id: jus-vwdw
status: closed
deps: [jus-55ji]
links: []
created: 2026-02-17T15:12:42Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-fxaf
tags: [cli, fang, cobra, global-flags]
---
# Add global flags for output control

Add global/persistent flags to the root command for controlling output format and verbosity across all commands.

## Background

dade is a CLI for scaffolding web projects. To support headless/scripted usage, we need global flags that control output behavior consistently across all commands.

## Current State

File: cmd/dade/root.go (created in earlier ticket)

Current global behavior:
- Styled output auto-detected via term.IsTerminal
- No quiet mode
- No verbose mode
- JSON output only on specific commands (templates --json)

## Implementation

1. Update cmd/dade/root.go to add persistent flags:

```go
var (
    flagQuiet   bool
    flagVerbose bool
    flagNoColor bool
    flagJSON    bool
)

func init() {
    rootCmd.PersistentFlags().BoolVarP(&flagQuiet, "quiet", "q", false, "Suppress non-essential output")
    rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "Enable verbose output")
    rootCmd.PersistentFlags().Bool("no-color", false, "Disable colored output")
    rootCmd.PersistentFlags().Bool("json", false, "Output in JSON format where supported")
    
    // Mark mutually exclusive
    rootCmd.MarkFlagsMutuallyExclusive("quiet", "verbose")
}
```

2. Create helper function to get output settings:

```go
type OutputSettings struct {
    Quiet   bool
    Verbose bool
    Styled  bool
    JSON    bool
}

func getOutputSettings(cmd *cobra.Command) OutputSettings {
    noColor, _ := cmd.Flags().GetBool("no-color")
    jsonOut, _ := cmd.Flags().GetBool("json")
    
    styled := term.IsTerminal(int(os.Stdout.Fd())) && !noColor
    
    return OutputSettings{
        Quiet:   flagQuiet,
        Verbose: flagVerbose,
        Styled:  styled,
        JSON:    jsonOut,
    }
}
```

3. Update each command's RunE to use getOutputSettings

4. Modify logging.Logger to support quiet/verbose modes

## Flag Documentation (shown in --help for any command)

```
Global Flags:
      --json        Output in JSON format where supported
      --no-color    Disable colored output
  -q, --quiet       Suppress non-essential output
  -v, --verbose     Enable verbose output
```

## Behavior Matrix

| Flag | Effect |
|------|--------|
| (none) | Normal output, auto-detect color |
| --quiet | Only errors and essential results |
| --verbose | Extra detail (git commands, paths, etc.) |
| --no-color | Plain text, no ANSI codes |
| --json | Machine-readable JSON (where supported) |

## Acceptance Criteria

- [ ] dade --help shows global flags
- [ ] dade new --help shows global flags inherited
- [ ] dade -q new myproject suppresses progress messages
- [ ] dade --verbose new myproject shows extra detail
- [ ] dade --no-color new myproject has no ANSI codes
- [ ] dade --json templates outputs JSON
- [ ] -q and -v together produce error
- [ ] Global flags work on all commands

