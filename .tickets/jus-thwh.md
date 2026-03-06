---
id: jus-thwh
status: closed
deps: [jus-ohdb, jus-c6us]
links: []
created: 2026-02-12T20:35:32Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-fgp8
tags: [testing, static]
---
# Test hypertext with dade

Test the hypertext template end-to-end with dade.

## Test Steps

```bash
# 1. Install template
dade install https://github.com/theydontwantyoutovibecode/dade-with-hypertext.git
dade templates  # Verify listed

# 2. Create project
cd /tmp
dade new testhtml --template hypertext
# Verify: No setup script runs (static)
# Verify: Files copied

# 3. Start server
cd testhtml
dade start
# Verify: https://testhtml.localhost loads
# Verify: index.html shows
# Verify: HTMX demo works (partials load)

# 4. Test commands
dade stop
dade start --bg
dade open

# 5. Cleanup
dade remove testhtml --files
```

## Acceptance Criteria

- [ ] Template installs correctly
- [ ] Project creates without setup
- [ ] Static serving works
- [ ] HTMX partials work
- [ ] HTTPS works
- [ ] Cleanup works

