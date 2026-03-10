---
id: mid-w2f8
status: closed
deps: []
links: [mid-6l7z]
created: 2026-03-07T19:07:52Z
type: Research: Document mDNS subdomain limitation
priority: 2
assignee: Alex Cabrera
---
# Research: Document mDNS subdomain limitation

Investigate and document why *.hostname.local subdomains don't resolve via mDNS, test alternative domain schemes (.localhost, .local, .test)

## Findings

### mDNS/Bonjour Limitation
- mDNS (Multicast DNS) only resolves the **base hostname** (e.g., `crm114.local`)
- **Subdomains do not resolve** (e.g., `project.crm114.local` fails with NXDOMAIN)
- This is a fundamental limitation of mDNS, not a configuration issue
- mDNS works on the local network for other devices to find `crm114.local`

### Why Subdomains Don't Work
- mDNS registers only one DNS record per hostname: `A` record for `crm114.local` -> `192.168.x.x`
- It does **not** register wildcard or subdomain records
- DNS resolution for `project.crm114.local` would require either:
  - A wildcard record (`*.crm114.local`)
  - Individual subdomain records (`project.crm114.local`)
  - A proper DNS server with zone configuration

### Alternative Domain Schemes Tested

| Scheme | Resolution | Notes |
|--------|------------|-------|
| `*.crm114.local` | ❌ Fails (NXDOMAIN) | mDNS limitation |
| `*.crm114.localhost` | ✅ Works (127.0.0.1) | RFC 6761 compliant |
| `*.crm114.test` | ❓ Untested | Reserved but requires configuration |
| `project.localhost` | ✅ Works (127.0.0.1) | RFC 6761 compliant |

### Current Caddyfile Configuration
The Caddyfile at `~/.config/dade/Caddyfile` is correctly configured:
```
{
  local_certs
}

https://alexcabrera-me.crm114.local {
  reverse_proxy localhost:3001
}
```
The configuration is valid - Caddy is running and would accept connections if the domain resolved.

### Current Status
- Proxy is running: ✅
- Caddyfile is valid: ✅
- Domain resolves: ❌ (blocked by mDNS limitation)

## Conclusion
The README's claim that *.hostname.local works across the LAN is incorrect. Need to either:
1. Use .localhost TLD (works immediately, no LAN sharing)
2. Implement /etc/hosts automation for .local subdomains (enables LAN sharing)

Both approaches should be implemented for flexibility.

