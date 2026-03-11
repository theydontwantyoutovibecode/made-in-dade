---
id: mid-0jq8
status: closed
deps: [mid-yz6r]
links: []
created: 2026-03-11T18:13:16Z
type: task
priority: 1
assignee: Alex Cabrera
tags: [hot-reload, sse, server]
---
# Implement SSE server for browser notifications

Create an internal/hotreload package with an SSE server that runs alongside the main dev server. The server should: (1) Listen on a dynamic port or a dedicated port, (2) Provide /_dade/events endpoint for SSE connections, (3) Broadcast reload events when files change, (4) Send periodic keep-alive pings. This replaces the Node.js SSE server.

