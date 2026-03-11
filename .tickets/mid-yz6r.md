---
id: mid-yz6r
status: closed
deps: []
links: []
created: 2026-03-11T18:13:14Z
type: task
priority: 1
assignee: Alex Cabrera
tags: [hot-reload, go, watcher]
---
# Add Go-based file watcher for hot-reload

Create a Go package internal/watcher that monitors file changes (HTML, CSS, JS) and triggers callbacks. This replaces the bash-based fswatch approach with a native Go solution using fsnotify or similar library. The watcher should support multiple file patterns and be able to start/stop cleanly.

