---
id: jus-s5is
status: closed
deps: [jus-sj8r]
links: []
created: 2026-03-02T15:41:19Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-k85j
tags: [template, macos]
---
# Create macOS template dev.sh

dev.sh for macOS: generate xcodeproj if missing, build with xcodebuild for macOS destination, find .app in DerivedData, launch via open command. No simulator needed — runs natively. Kill previous instance before re-launching.


## Notes

**2026-03-02T15:43:30Z**

## Reviewed — dev.sh differences from iOS

1. NO simulator — builds and runs natively
2. NO ensure_platform step (macOS SDK ships with Xcode)
3. NO simctl usage at all
4. Build destination: `'platform=macOS'`
5. App path: `DerivedData/Build/Products/Debug/MyApp.app` (not Debug-iphonesimulator)
6. Launch via: `open DerivedData/Build/Products/Debug/MyApp.app`
7. Kill previous instance: `pkill -x MyApp 2>/dev/null || true` before launch
8. Much simpler script than iOS
