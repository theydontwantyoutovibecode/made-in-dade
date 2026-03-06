---
id: jus-lkra
status: closed
deps: [jus-55ji]
links: []
created: 2026-02-17T15:11:57Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-fxaf
tags: [cli, fang, cobra, proxy-command]
---
# Migrate 'proxy' command to Cobra with subcommands

Migrate the 'dade proxy' command from manual arg parsing to Cobra with proper subcommands.

## Background

dade is a CLI for scaffolding web projects. The 'proxy' command manages a local Caddy HTTPS proxy service via launchd.

## Current State

File: cmd/dade/proxy.go

Current structure (positional subcommands):
- proxy start: Start proxy service
- proxy stop: Stop proxy service  
- proxy restart: Restart proxy service
- proxy status: Show proxy status (default if no action)
- proxy logs: Tail proxy logs

No flags currently. No interactive prompts - already headless-capable.

## Implementation

1. Create cmd/dade/cmd_proxy.go with Cobra command and subcommands:

```go
package main

import (
    "github.com/spf13/cobra"
)

var proxyCmd = &cobra.Command{
    Use:   "proxy",
    Short: "Manage the local HTTPS proxy service",
    Long:  "Control the Caddy-based HTTPS proxy that provides local .localhost domains for your projects.",
}

var proxyStartCmd = &cobra.Command{
    Use:   "start",
    Short: "Start the proxy service",
    RunE:  runProxyStartCmd,
}

var proxyStopCmd = &cobra.Command{
    Use:   "stop",
    Short: "Stop the proxy service",
    RunE:  runProxyStopCmd,
}

var proxyRestartCmd = &cobra.Command{
    Use:   "restart",
    Short: "Restart the proxy service",
    RunE:  runProxyRestartCmd,
}

var proxyStatusCmd = &cobra.Command{
    Use:   "status",
    Short: "Show proxy status",
    Long:  "Display whether the proxy is running, number of registered projects, and port range.",
    RunE:  runProxyStatusCmd,
}

var proxyLogsCmd = &cobra.Command{
    Use:   "logs",
    Short: "Tail proxy logs",
    Long:  "Stream the proxy service logs in real-time.",
    RunE:  runProxyLogsCmd,
}

func init() {
    rootCmd.AddCommand(proxyCmd)
    
    proxyCmd.AddCommand(proxyStartCmd)
    proxyCmd.AddCommand(proxyStopCmd)
    proxyCmd.AddCommand(proxyRestartCmd)
    proxyCmd.AddCommand(proxyStatusCmd)
    proxyCmd.AddCommand(proxyLogsCmd)
    
    // Add flags to status subcommand
    proxyStatusCmd.Flags().Bool("json", false, "Output status as JSON")
    
    // Add flags to logs subcommand
    proxyLogsCmd.Flags().IntP("lines", "n", 0, "Number of lines to show (default: follow)")
    proxyLogsCmd.Flags().BoolP("follow", "f", true, "Follow log output")
}
```

2. Implement each RunE function calling existing proxyCommand methods

3. Add --json flag to status for scripted status checks

4. Remove old arg parsing from proxy.go, keep business logic

## Flag Documentation (shown in --help)

Main command:
```
Usage:
  dade proxy [command]

Available Commands:
  logs        Tail proxy logs
  restart     Restart the proxy service
  start       Start the proxy service
  status      Show proxy status
  stop        Stop the proxy service
```

Status subcommand:
```
Flags:
      --json   Output status as JSON
  -h, --help   help for status
```

Logs subcommand:
```
Flags:
  -f, --follow        Follow log output (default true)
  -n, --lines int     Number of lines to show (default: follow)
  -h, --help          help for logs
```

## Acceptance Criteria

- [ ] dade proxy --help shows available subcommands
- [ ] dade proxy start/stop/restart/status/logs work
- [ ] dade proxy status --json outputs valid JSON
- [ ] dade proxy logs -n 50 shows last 50 lines
- [ ] dade proxy logs -f=false shows logs without following
- [ ] Existing tests pass or are updated appropriately

