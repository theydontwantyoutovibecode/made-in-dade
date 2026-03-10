# Migration Guide: .local → .localhost

## What Changed?

dade now uses `.localhost` domains by default instead of `.local`. This is because `.localhost` resolves to `127.0.0.1` automatically (per RFC 6761) without requiring any `/etc/hosts` configuration.

### Why This Change?

**The old approach (`.local`):**
- Required `/etc/hosts` entries for each project
- **Did NOT work across LAN** (mDNS only resolves base hostname, not subdomains)
- Required manual configuration or root privileges

**The new approach (`.localhost`):**
- Works immediately without `/etc/hosts`
- Works automatically on any system (RFC 6761 compliant)
- Faster resolution (loopback, no network traffic)
- Cross-platform compatible

## Do I Need to Do Anything?

### Most Users: No!

If you're using dade for local development only (which is the default use case), you don't need to do anything. Your existing projects will continue to work with their new `.localhost` URLs.

### Legacy Installations: Preserved

If you have an existing dade installation (with `projects.json` but no `config.toml`), dade will continue using `.local` to preserve backward compatibility. To switch to `.localhost`, follow the steps below.

## Switching to `.localhost`

To migrate from `.local` to `.localhost`, create a config file:

```bash
echo 'domain_tld = ".localhost"' > ~/.config/dade/config.toml
```

Then reload the proxy:

```bash
dade proxy reload
```

That's it! Your projects will now use `.localhost` domains.

## Switching to `.local` (for LAN Access)

If you need LAN access to your development projects and understand the limitations:

```bash
echo 'domain_tld = ".local"' > ~/.config/dade/config.toml
```

Then add each project's domain to `/etc/hosts`:

```bash
# For each project
sudo bash -c "echo '127.0.0.1\t<project-name>.<hostname>.local' >> /etc/hosts"
```

**Important:** `.local` subdomains will NOT work across LAN unless you manually configure each device's `/etc/hosts` file or set up a local DNS server. mDNS only resolves the base hostname (e.g., `crm114.local`), not subdomains (e.g., `myapp.crm114.local`).

## Environment Variable Override

You can also set the domain TLD via environment variable for testing:

```bash
DADE_DOMAIN_TLD=".localhost" dade dev
DADE_DOMAIN_TLD=".test" dade dev
```

## What About My Bookmarks?

If you have browser bookmarks to `.local` URLs, you'll need to update them to use `.localhost` (or manually add `/etc/hosts` entries if you want to keep using `.local`).

## Troubleshooting

### "Project already running on port X"

If you see this warning after switching domains:

```bash
dade stop
dade dev
```

The old process might still be running from the previous session. Stopping and restarting should resolve this.

### Browser Can't Connect to .localhost

If your browser doesn't load `.localhost` URLs:

1. Check that the proxy is running:
   ```bash
   dade proxy status
   ```

2. Try accessing via `http://` instead of `https://` (should redirect to HTTPS)
   ```bash
   http://myproject.hostname.localhost
   ```

3. Clear browser cache and try again

### Want to Use a Custom Domain TLD?

You can use any TLD you want:

```bash
echo 'domain_tld = ".test"' > ~/.config/dade/config.toml
```

**Note:** Custom TLDs will require `/etc/hosts` entries unless you configure DNS.

## Summary

| TLD | Works Immediately | LAN Access | Requires /etc/hosts |
|------|------------------|--------------|-------------------|
| `.localhost` (default) | ✓ | ✗ | ✗ |
| `.local` | ✗ (requires hosts) | ✗ (mDNS limitation) | ✓ |
| `.test` (custom) | ✗ (requires hosts) | ✗ | ✓ |
| Any custom TLD | ✗ (requires hosts) | ✗ | ✓ |

**Recommendation:** Use `.localhost` for local development. Use `.local` only if you understand its limitations and manually configure `/etc/hosts`.
