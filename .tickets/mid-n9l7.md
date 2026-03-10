---
id: mid-n9l7
status: closed
deps: [mid-w2f8, mid-6l7z]
links: []
created: 2026-03-07T19:07:56Z
type: Config: Add domain TLD configuration option
priority: 1
assignee: Alex Cabrera
---
# Config: Add domain TLD configuration option

## Implementation Complete

### What was implemented:
1. Added `internal/config/domain.go` with:
   - `DomainTLD()` function that reads from env, config, or defaults to `.localhost`
   - `SetDomainTLD()` function to update config file
   - `Legacy detection`: Existing installations use `.local` for backward compatibility
   - Environment variable support: `DADE_DOMAIN_TLD` can override

2. Updated `LocalDomain()` and `ProjectDomain()` in `paths.go` to use configured TLD
3. Added comprehensive tests in `domain_test.go`

### How it works:
- New installations: Use `.localhost` (RFC 6761 compliant, no /etc/hosts needed)
- Legacy installations: Use `.local` for backward compatibility
- Override: Set `DADE_DOMAIN_TLD` env var or add `domain_tld = ".tld" to `~/.config/dade/config.toml`

### Why .localhost works without /etc/hosts:
RFC 6761 Section 6.3 mandates `.localhost` TLD resolve to loopback (127.0.0.1). All modern OSes handle this natively - no DNS or hosts file required.

### What's NOT needed for .localhost:
- /etc/hosts automation (not required, .localhost is OS-native)
- DNS server (not required, OS handles resolution)
- LAN sharing (localhost is by definition local-only)

### When /etc/hosts IS needed:
Only if users want `.local` domains for LAN sharing. That's a future enhancement (tickets mid-bjyk, mid-zuxa, mid-a4p7, mid-z748, mid-qjgk).

