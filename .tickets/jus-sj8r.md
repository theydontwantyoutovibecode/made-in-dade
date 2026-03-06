---
id: jus-sj8r
status: closed
deps: []
links: []
created: 2026-03-02T15:41:19Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-k85j
tags: [template, swift]
---
# Create macOS template project structure and SwiftUI sources

Create the template directory at /Users/acabrera/Code/dade/mac-made-in-dade/ with:
- project.yml (XcodeGen, macOS 14+, Swift 6, SwiftUI App lifecycle)
- Sources/App/MyApp.swift, ContentView.swift, Feature.swift
- Sources/Features/Home/HomeView.swift (welcome screen)
- Sources/Features/Counter/CounterView.swift, CounterViewModel.swift
- Sources/Features/Settings/SettingsView.swift
- Resources/Assets.xcassets (AppIcon, AccentColor)
- Tests/CounterTests.swift
- .gitignore
All Swift code must compile on Xcode 26 / Swift 6 (no .accent, use .tint).


## Notes

**2026-03-02T15:43:30Z**

## Reviewed — required changes vs iOS template

Must NOT copy iOS code blindly. macOS-specific adaptations:
1. Use `NavigationSplitView` with sidebar, not `NavigationStack` with push
2. Remove all `.navigationBarTitleDisplayMode()` calls (iOS-only)
3. Remove `.textContentType(.username)` (iOS-only)
4. Add `.defaultSize(width: 900, height: 600)` to WindowGroup
5. Add `Settings { SettingsView() }` scene for proper macOS Settings menu
6. Use `.tint` not `.accent` (Xcode 26 compat)
7. project.yml: platform macOS, deployment target macOS 14.0, NO UIKit plist keys
8. Ensure macOS-specific plist: NSMainStoryboardFile empty, no UILaunchScreen
