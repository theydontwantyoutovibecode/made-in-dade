---
id: jus-k85j
status: open
deps: []
links: []
created: 2026-03-02T15:41:04Z
type: epic
priority: 1
assignee: Alex Cabrera
tags: [template, macos, swiftui]
---
# macOS App Template

Create a new default template for native macOS applications using SwiftUI, targeting latest macOS on supported hardware. Follows the same patterns as the iOS template: XcodeGen for project generation, dev.sh for terminal-driven build+launch, setup.sh for prereqs, no Xcode GUI required.

## Design

Template structure mirrors ios-app:
- project.yml: XcodeGen spec targeting macOS 14+ with SwiftUI App lifecycle
- Sources/App/: MyApp.swift (entry), ContentView.swift (nav), Feature.swift (enum)
- Sources/Features/: Home, Counter, Settings views + view models
- Resources/: Assets.xcassets with AppIcon and AccentColor
- Tests/: Unit tests for view models
- dev.sh: builds with xcodebuild, launches app directly (no simulator needed for macOS)
- setup.sh: checks Xcode, installs xcodegen, generates project
- dade.toml: proxy=false, dev.messages for macOS context
- AGENTS.md: documents stack and conventions
- .read-only/manifest.txt: reference libraries (swift-collections, swift-algorithms)

Key difference from iOS: macOS builds run natively, no simulator. dev.sh builds and launches the .app directly via `open`.

## Acceptance Criteria

1. dade new mac-app creates a working project. 2. dade dev builds and launches the app natively. 3. App displays welcome screen with navigation. 4. All Swift code compiles on Xcode 26 / macOS 26. 5. Template registered in DefaultTemplates. 6. All tests pass.


## Notes

**2026-03-02T15:42:57Z**

## Critical Review — Gaps and Issues

### 1. iOS-isms that won't work on macOS
- `.navigationBarTitleDisplayMode(.inline)` — iOS-only API, doesn't exist on macOS. Must use `.navigationTitle()` alone or macOS-specific navigation patterns.
- `NavigationStack` works on macOS but the UX pattern is wrong — macOS apps typically use `NavigationSplitView` with a sidebar, not a mobile-style push stack.
- `Form` works but looks different on macOS — this is fine, actually idiomatic.
- `.textContentType(.username)` — iOS-only. Remove for macOS.
- iPhone-specific tips in HomeView ('Build with: xcodebuild or dade build') should reference macOS.

### 2. macOS-specific considerations missing from plan
- **Window sizing**: macOS apps need `.defaultSize()` or `frame()` on WindowGroup. Without it, the window starts tiny.
- **Menu bar**: macOS apps should have proper menu commands. At minimum, Settings should be in the app menu via `Settings { SettingsView() }` scene.
- **Toolbar**: macOS uses toolbar items differently. Consider sidebar toggle.
- **App icon**: macOS requires a 1024x1024 AppIcon in the asset catalog, different from iOS.
- **Sandbox & entitlements**: macOS apps may need entitlements file. XcodeGen handles this differently.

### 3. dev.sh gaps
- 'Kill previous instance before re-launching' — need to specify HOW. Use `pkill -f MyApp` or `killall MyApp` before launching. This is important because unlike iOS simulator, you can't reinstall over a running macOS app.
- Build destination should be `'platform=macOS'` not simulator.
- App path in DerivedData will be `Debug/` not `Debug-iphonesimulator/`.
- No platform download needed (macOS SDK ships with Xcode), so the ensure_platform step from iOS is unnecessary.

### 4. project.yml differences
- Platform must be `macOS` not `iOS`.
- Deployment target should be `macOS: '14.0'` (or '15.0' to match latest patterns).
- No UIKit-specific plist keys (UILaunchScreen, UISupportedInterfaceOrientations, etc).
- Need macOS-specific plist keys: `LSMinimumSystemVersion`, potentially `NSMainStoryboardFile` set to empty.

### 5. Alias concern
- `appkit` is misleading — this is a SwiftUI app, not AppKit. Better aliases: `mac`, `macos`, `mac-app`, `desktop`.

### 6. Missing from plan: what about the .read-only manifest?
- swift-collections and swift-algorithms as read-only refs — fine, mirrors iOS template.

### 7. Template test gap
- default_templates_test.go — need to verify this test exists and what it checks. If it asserts on template count (currently 6), it'll break when we add a 7th.

**2026-03-02T15:43:40Z**

## Review Verdict

Plan is sound but had 7 issues found during review:

1. **iOS API bleed** — NavigationStack/navigationBarTitleDisplayMode/textContentType are iOS patterns. macOS needs NavigationSplitView + sidebar. FIXED in ticket notes.
2. **dev.sh over-engineered** — No simulator, no platform download, no simctl. Just build + open. Much simpler than iOS. FIXED in ticket notes.
3. **Wrong alias** — 'appkit' is misleading for SwiftUI. Changed to 'desktop'. FIXED in ticket notes.
4. **Test breakage** — templates_test.go hardcodes count=6 in TWO places. Must update to 7. FIXED in ticket notes.
5. **Window sizing** — macOS apps need explicit default size or they start tiny. Added to ticket notes.
6. **Settings scene** — macOS convention is Settings in app menu, not in-app navigation. Added to ticket notes.
7. **macOS plist keys** — Different from iOS, no UIKit keys. Added to ticket notes.

All issues documented in respective tickets. Plan is ready for execution.
