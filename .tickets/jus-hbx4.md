---
id: jus-hbx4
status: closed
deps: [jus-nmo3]
links: []
created: 2026-03-02T01:57:38Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-nphc
tags: [ios, setup, ws-4]
---
# iOS template auto-setup with emulator launch

Update ios-made-in-dade setup.sh to: 1) Install Xcode CLI tools if missing (xcode-select --install). 2) Accept Xcode license (sudo xcodebuild -license accept). 3) Install XcodeGen via Homebrew. 4) Generate Xcode project (xcodegen generate). Update dade.toml [serve] dev command to: build for simulator and launch. Use xcrun simctl to boot simulator if not running, then xcodebuild to build and install. Caveats: Xcode itself must be installed from App Store — cannot be automated. setup.sh should detect if Xcode.app exists and give clear instructions if not. SimDevice selection should default to latest iPhone simulator.

## Acceptance Criteria

1. setup.sh installs all deps except Xcode.app. 2. dade dev builds and launches in Simulator. 3. Clear error if Xcode.app missing.

