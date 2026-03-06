---
id: jus-fgp8
status: closed
deps: [jus-swo4, jus-zf02, jus-ocjp]
links: []
created: 2026-02-12T20:00:11Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-2e3m
tags: [templates, testing]
---
# Test templates with new dade CLI

Comprehensive testing of both templates with the new unified dade CLI.

## Test Matrix

### Django Template Tests

1. **Install template**
   ```bash
   dade install https://github.com/theydontwantyoutovibecode/dade-with-django-and-hypermedia.git
   dade templates  # Verify listed
   ```

2. **Create project**
   ```bash
   dade new testdjango --template django-hypermedia
   # Verify: setup.sh runs, prompts work
   ```

3. **Start project**
   ```bash
   cd testdjango
   dade start
   # Verify: https://testdjango.localhost works
   # Verify: Django dev server running
   # Verify: TailwindCSS watcher running
   ```

4. **Stop project**
   ```bash
   dade stop
   # Verify: server stopped
   # Verify: port released
   ```

5. **List and open**
   ```bash
   dade list  # Shows testdjango
   dade open  # Opens browser
   ```

### Hypertext Template Tests

1. **Install template**
   ```bash
   dade install https://github.com/theydontwantyoutovibecode/dade-with-hypertext.git
   ```

2. **Create project**
   ```bash
   dade new testhtml --template hypertext
   # Verify: files copied
   # Verify: no setup script (static)
   ```

3. **Start project**
   ```bash
   cd testhtml
   dade start
   # Verify: https://testhtml.localhost works
   # Verify: index.html served
   # Verify: HTMX interactions work
   ```

4. **Test partials**
   ```bash
   # Click HTMX demo buttons
   # Verify: partials loaded correctly
   ```

### Cross-Template Tests

1. **Multiple projects simultaneously**
   ```bash
   dade start  # in testdjango
   dade start  # in testhtml (different terminal)
   dade list   # Both show running
   ```

2. **Tunnel test**
   ```bash
   dade tunnel  # Quick tunnel
   # Verify: public URL works
   ```

3. **Remove and cleanup**
   ```bash
   dade remove testdjango
   dade remove testhtml --files
   dade list  # Empty
   ```

## Acceptance Criteria

- [ ] Django template installs and creates projects
- [ ] Django projects start with all services
- [ ] Hypertext template installs and creates projects
- [ ] Static serving works for hypertext
- [ ] Multiple projects run simultaneously
- [ ] Tunnel works for both project types
- [ ] Cleanup removes projects correctly

