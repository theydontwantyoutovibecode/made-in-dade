---
id: jus-9h5b
status: closed
deps: []
links: []
created: 2026-03-02T06:14:44Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-yxme
tags: [manifest]
---
# Add serve.proxy bool and dev.messages section to manifest

Add `proxy` bool field to Serve struct (defaults true for web, false when default_port=0). Add `Messages` struct to Dev with `ready` and `running` string fields. Parse from `[serve]` and `[dev.messages]` sections. Add accessor functions: NeedsProxy(), DevReadyMessage(), DevRunningMessage().

