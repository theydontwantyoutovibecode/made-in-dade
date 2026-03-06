---
id: jus-nphc
status: closed
deps: [jus-u0zj]
links: []
created: 2026-03-02T01:57:29Z
type: epic
priority: 1
assignee: Alex Cabrera
tags: [mobile, setup, ws-4]
---
# WS-4: iOS & Android Auto-Setup

Make dade dev on iOS/Android projects fully automatic. Install all system dependencies via Homebrew, generate projects, launch emulators, build and run — all without user interaction with Xcode or Android Studio GUIs.

## Acceptance Criteria

1. dade dev on iOS project: installs Xcode CLI, XcodeGen, generates project, opens Simulator, builds and runs. 2. dade dev on Android project: installs JDK, Android SDK, creates AVD, starts emulator, builds and installs.

