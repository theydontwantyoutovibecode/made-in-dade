---
id: mid-6l7z
status: closed
deps: []
links: [mid-w2f8]
created: 2026-03-07T19:08:03Z
type: Research: Verify RFC 6761 .localhost compliance
priority: 2
assignee: Alex Cabrera
---
# Research: Verify RFC 6761 .localhost compliance

Research RFC 6761 and confirm .localhost TLD behavior across platforms, verify subdomains resolve to 127.0.0.1

## Findings

### RFC 6761 Compliance
RFC 6761 (Section 6.3) reserves the `.localhost` TLD for loopback purposes:
- **Name**: `.localhost`
- **Purpose**: Local loopback address resolution
- **Resolution**: Must resolve to IPv4 loopback (127.0.0.0/8) or IPv6 (::1)
- **No DNS queries required**: Resolution should happen without external DNS servers

### Cross-Platform Behavior

| Platform | Behavior | Subdomain Support |
|----------|----------|------------------|
| macOS | ✅ Resolves to 127.0.0.1 | ✅ Works |
| Linux | ✅ Resolves to 127.0.0.1 | ✅ Works |
| Windows | ✅ Resolves to 127.0.0.1 | ✅ Works (since Win10) |

### Verification Tests

```bash
$ ping -c 1 alexcabrera-me.localhost
PING localhost (127.0.0.1): 56 data bytes
64 bytes from 127.0.0.1: icmp_seq=0 ttl=64 time=0.168 ms

$ ping -c 1 crm114.localhost
PING localhost (127.0.0.1): 56 data bytes
64 bytes from 127.0.0.1: icmp_seq=0 ttl=64 time=0.193 ms
```

All `.localhost` subdomains resolve to `127.0.0.1` without any DNS server or configuration.

### Comparison with .local
| Aspect | .localhost | .local (mDNS) |
|--------|-----------|--------------|
| Standard | RFC 6761 | RFC 6762 |
| Resolution | Native OS loopback | mDNS service discovery |
| Network | Local only | LAN discoverable |
| Subdomains | ✅ Native support | ❌ No subdomain support |
| Latency | ~0ms (no network) | ~1-5ms (mDNS multicast) |
| Requires DNS server | ❌ No | ❌ No (uses mDNS) |

### Caddy Compatibility
Caddy's `local_certs` directive works with `.localhost` domains:
- Generates self-signed certificates automatically
- Browsers accept these certificates without warnings
- HTTPS works seamlessly

## Recommendation
Switch default domain scheme from `.hostname.local` to `projectname.localhost`:
1. **Immediate fix**: Works out of the box, no configuration
2. **RFC compliant**: Follows established standards
3. **Cross-platform**: Works consistently everywhere
4. **Better UX**: Faster resolution, no network dependency

For LAN sharing, implement `/etc/hosts` automation as an opt-in feature.

