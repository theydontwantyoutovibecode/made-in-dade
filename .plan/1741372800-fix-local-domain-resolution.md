# Goal: Fix Local Domain Resolution for dade Projects

## Background
When running `dade dev` in `~/Code/alexcabrera-me`, the project URL `https://alexcabrera-me.crm114.local` is not accessible. The Caddy proxy is running correctly and configured to proxy to `localhost:3001`, but DNS resolution fails for `alexcabrera-me.crm114.local`.

## Context
The dade project assumes that wildcard subdomains of `<hostname>.local` (e.g., `*.crm114.local`) work via mDNS/Bonjour. However, mDNS only resolves the base hostname (`crm114.local`), not subdomains. This is a fundamental limitation of mDNS.

The README incorrectly states:
> "For example, if your hostname is `macbook` and you create a project called `myapp`, it is available at `https://myapp.macbook.local`. These URLs work from any device on the same local network."

This is false - subdomains do not resolve via mDNS.

## Strategy
Implement a multi-tiered solution:

1. **Primary**: Switch to `.localhost` TLD for local development (RFC 2606 compliant, subdomains resolve to 127.0.0.1 by default)
2. **Fallback**: Implement `/etc/hosts` file automation for `.local` subdomains when users need LAN sharing
3. **Documentation**: Update README to accurately describe domain resolution behavior and limitations

This approach provides immediate fixes while preserving the original intent of network sharing when needed.

## Passes

### Pass 1: Analysis & Validation

#### Research: Domain Resolution Behavior
**Problem**: Understand current domain resolution limitations
**Approach**: Document current behavior and test alternatives
**Definition of Done**:
- [x] Document mDNS limitation (base hostname resolves, subdomains don't)
- [x] Test `.localhost` subdomain resolution
- [x] Test current Caddy configuration
- [x] Document findings

#### Research: RFC Compliance
**Problem**: Verify `.localhost` TLD behavior across platforms
**Approach**: Research RFC 2606 and RFC 6761
**Definition of Done**:
- [x] Verify `.localhost` is reserved for loopback
- [x] Confirm subdomains resolve to 127.0.0.1
- [x] Document platform differences (macOS vs Linux)

### Pass 2: Core Implementation

#### Config: Domain Scheme Migration
**Problem**: Current domain scheme uses `*.hostname.local` which doesn't resolve
**Approach**: Update config package to support configurable domain TLD (`.localhost` default)
**Definition of Done**:
- [ ] Add config option for domain TLD (default: `.localhost`)
- [ ] Update `LocalDomain()` and `ProjectDomain()` to use new scheme
- [ ] Add backward compatibility mode for `.local` domains
- [ ] Update Caddyfile generation to use new domains
- [ ] Add tests for domain resolution functions

#### Proxy: Update Caddyfile Generation
**Problem**: Caddyfile must use new domain scheme
**Approach**: Modify `buildCaddyfile()` to use configured TLD
**Definition of Done**:
- [ ] Update `buildCaddyfile()` to use new domain scheme
- [ ] Ensure Caddyfile validates correctly
- [ ] Test Caddyfile generation with new domains
- [ ] Add test coverage

#### Setup: Initialize Default Domain TLD
**Problem**: New installations should use `.localhost` by default
**Approach**: Add domain TLD to config initialization
**Definition of Done**:
- [ ] Add `domain_tld` field to config initialization
- [ ] Set default to `.localhost`
- [ ] Preserve `.local` for existing installations (migration path)
- [ ] Update setup command to handle domain configuration

### Pass 3: Hosts File Automation (LAN Sharing)

#### Hosts: Read/Parse /etc/hosts
**Problem**: Need to read and parse hosts file entries
**Approach**: Create internal/hosts package for hosts file operations
**Definition of Done**:
- [ ] Create `internal/hosts` package
- [ ] Implement `ReadHostsFile()` function
- [ ] Implement `ParseHostsEntries()` function
- [ ] Add tests for parsing various host entry formats
- [ ] Handle commented entries and whitespace

#### Hosts: Add Project Domains
**Problem**: Add project subdomains to hosts file for LAN sharing
**Approach**: Implement function to add entries with proper formatting
**Definition of Done**:
- [ ] Implement `AddDomainEntry(domain, ip)` function
- [ ] Add comment markers for dade-managed entries
- [ ] Handle duplicate entries gracefully
- [ ] Add tests for adding entries
- [ ] Verify sudo integration works

#### Hosts: Remove Project Domains
**Problem**: Clean up hosts file when projects are removed
**Approach**: Implement function to remove dade-managed entries
**Definition of Done**:
- [ ] Implement `RemoveDomainEntry(domain)` function
- [ ] Identify dade-managed entries by comment markers
- [ ] Clean up empty comment sections
- [ ] Add tests for removing entries

#### Hosts: Sudo Integration
**Problem**: Hosts file modifications require root privileges
**Approach**: Implement secure sudo wrapper for hosts file operations
**Definition of Done**:
- [ ] Create `internal/exec/sudo.go` for privileged operations
- [ ] Implement `RunWithSudo()` helper
- [ ] Add confirmation prompts for privileged operations
- [ ] Handle sudo password prompts gracefully
- [ ] Add tests with mock sudo

#### Proxy: Integrate Hosts File Management
**Problem**: Update proxy to manage hosts file entries
**Approach**: Call hosts file functions in proxy lifecycle
**Definition of Done**:
- [ ] Update `GenerateCaddyfile()` to add hosts entries for `.local` domains
- [ ] Update proxy removal to clean up hosts entries
- [ ] Add warning when hosts file update fails
- [ ] Add `--no-hosts` flag to skip hosts file updates
- [ ] Add tests

### Pass 4: Command Updates

#### Dev: Display Correct Domain
**Problem**: `dade dev` should display the working domain URL
**Approach**: Update success message to use new domain
**Definition of Done**:
- [ ] Update cmd_dev.go success message to use new domain scheme
- [ ] Add note about LAN sharing requiring hosts file
- [ ] Test dev command output
- [ ] Update help text

#### Proxy: Status Command
**Problem**: Users need to know which domains are configured
**Approach**: Update proxy status to show domain resolution status
**Definition of Done**:
- [ ] Show domain TLD in proxy status
- [ ] Show hosts file status (managed vs unmanaged)
- [ ] List all configured domains
- [ ] Add warnings for unresolvable domains
- [ ] Update tests

#### Proxy: Reload Command
**Problem**: Reload should update hosts file entries
**Approach**: Ensure hosts file is synchronized on reload
**Definition of Done**:
- [ ] Call hosts file update in reload path
- [ ] Show updated domain list after reload
- [ ] Handle errors gracefully
- [ ] Update tests

### Pass 5: Documentation & Migration

#### Docs: Update README
**Problem**: README incorrectly states `.local` subdomains work on LAN
**Approach**: Fix documentation and explain new behavior
**Definition of Done**:
- [ ] Update domain resolution section to explain `.localhost` vs `.local`
- [ ] Add FAQ entry about domain resolution
- [ ] Document hosts file automation for LAN sharing
- [ ] Update examples to use `.localhost`
- [ ] Add troubleshooting section

#### Docs: Update AGENTS.md
**Problem**: Agent guidelines need domain scheme documentation
**Approach**: Document new domain configuration
**Definition of Done**:
- [ ] Add domain TLD configuration to AGENTS.md
- [ ] Document hosts file package usage
- [ ] Update architecture section with new packages
- [ ] Add testing guidelines for hosts operations

#### Docs: Migration Guide
**Problem**: Existing users need migration path
**Approach**: Create migration documentation
**Definition of Done**:
- [ ] Document how to migrate existing projects to new domain scheme
- [ ] Explain backward compatibility mode
- [ ] Provide step-by-step migration instructions
- [ ] Document rollback procedure

### Pass 6: Testing & Validation

#### Test: E2E Domain Resolution
**Problem**: Verify end-to-end domain resolution works
**Approach**: Integration tests for full workflow
**Definition of Done**:
- [ ] Test `dade new` with `.localhost` domain
- [ ] Test `dade dev` with new domain
- [ ] Test browser access to new domain
- [ ] Test LAN sharing with hosts file
- [ ] Add to test suite

#### Test: Edge Cases
**Problem**: Handle various edge cases in domain resolution
**Approach**: Test malformed domains, special characters, etc.
**Definition of Done**:
- [ ] Test project names with special characters
- [ ] Test very long project names
- [ ] Test duplicate project names
- [ ] Test hosts file corruption scenarios
- [ ] Test permission errors

### Pass 7: Release Preparation

#### Changelog: Document Changes
**Problem**: Users need to know what changed
**Approach**: Prepare changelog entry
**Definition of Done**:
- [ ] Write changelog entry with breaking changes
- [ ] Document migration requirements
- [ ] Add upgrade instructions
- [ ] Tag release version

## Won't Do
- [ ] Implement dnsmasq or other local DNS resolvers (too complex, not cross-platform)
- [ ] Support custom DNS servers (out of scope for initial fix)
- [ ] IPv6 support for localhost (future enhancement)
- [ ] Automatic detection of local network devices (out of scope)

## Observations
- `.localhost` subdomains resolve to `127.0.0.1` automatically on macOS (RFC 6761)
- Current Caddyfile is correctly formatted but domains don't resolve
- Hosts file has existing entries for `crm114.app.local` and `crm114.computer.local` suggesting previous attempts at解决这个问题
- mDNS/Bonjour is running and working for base hostname only
- The proxy is running correctly (confirmed via `dade proxy status`)

## Debrief
- Outcomes: TBD
- Tests: TBD
- Research: Documented mDNS limitation, verified `.localhost` behavior
- Git: N/A (planning phase)
