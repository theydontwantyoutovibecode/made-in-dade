---
id: jus-hrzu
status: closed
deps: [jus-lz4t, jus-a97f]
links: []
created: 2026-02-12T20:35:26Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-fgp8
tags: [testing, django]
---
# Test django-hypermedia with dade

Test the Django template end-to-end with dade.

## Test Steps

```bash
# 1. Install template
dade install https://github.com/theydontwantyoutovibecode/dade-with-django-and-hypermedia.git
dade templates  # Verify listed

# 2. Create project
cd /tmp
dade new testdjango --template django-hypermedia
# Verify: setup.sh runs, prompts work

# 3. Start server
cd testdjango
dade start
# Verify: https://testdjango.localhost loads
# Verify: Homepage shows
# Verify: HTMX demo works

# 4. Test commands
dade stop
dade list  # Shows stopped
dade start --bg
dade list  # Shows running
dade open  # Browser opens

# 5. Test tunnel
dade tunnel
# Verify: Public URL generated

# 6. Cleanup
dade stop
dade remove testdjango --files
```

## Acceptance Criteria

- [ ] Template installs correctly
- [ ] Project creates with setup
- [ ] Server starts and site loads
- [ ] HTTPS works
- [ ] Stop/start cycle works
- [ ] Tunnel works
- [ ] Cleanup works

