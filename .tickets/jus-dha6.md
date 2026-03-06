---
id: jus-dha6
status: closed
deps: []
links: []
created: 2026-03-02T02:59:21Z
type: task
priority: 1
assignee: Alex Cabrera
tags: [docs, ws-9, android-app]
---
# Rewrite android-app template README.md

Current README is a 3-line stub. Needs comprehensive documentation:
- What the template creates (Kotlin + Jetpack Compose, Material 3, MVVM)
- Prerequisites (JDK 17+, Android SDK)
- What setup.sh does (installs SDK, creates AVD, launches emulator)
- Project structure (app/src/main/kotlin/..., gradle files, version catalog)
- How dade dev works (runs assembleDebug, then installDebug)
- How dade build works (gradlew assembleDebug/Release)
- That dade start, share, tunnel, open do NOT apply (no port, no server)
- Gradle version catalog (libs.versions.toml)
- .read-only manifest.txt contents
- AGENTS.md and .tickets workflow

