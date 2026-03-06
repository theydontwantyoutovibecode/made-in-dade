---
id: jus-eb0x
status: closed
deps: []
links: []
created: 2026-03-02T02:59:15Z
type: task
priority: 1
assignee: Alex Cabrera
tags: [docs, ws-9, ios-app]
---
# Rewrite ios-app template README.md

Current README is a 3-line stub. Needs comprehensive documentation:
- What the template creates (SwiftUI app, MVVM, NavigationStack)
- Prerequisites (macOS, Xcode 16+, XcodeGen)
- What setup.sh does (installs XcodeGen, boots simulator)
- Project structure (Sources/, Resources/, Tests/, project.yml)
- How dade dev works (runs xcode-select, opens Xcode project)
- How dade build works (xcodebuild with simulator/release targets)
- That dade start, share, tunnel, open do NOT apply (no port, no server)
- XcodeGen project.yml structure
- .read-only manifest.txt contents
- AGENTS.md and .tickets workflow

