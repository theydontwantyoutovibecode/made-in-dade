---
id: jus-96eu
status: closed
deps: [jus-nmo3]
links: []
created: 2026-03-02T01:57:47Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-nphc
tags: [android, setup, ws-4]
---
# Android template auto-setup with emulator launch

Update android-made-in-dade setup.sh to: 1) Install OpenJDK 21 via Homebrew. 2) Install Android SDK command-line tools via Homebrew (brew install --cask android-commandlinetools). 3) Use sdkmanager to install platform-tools, build-tools, platform API 35, system-images for emulator. 4) Create AVD with avdmanager. 5) Generate Gradle wrapper. Update dade.toml [serve] dev command to: start emulator if not running (emulator -avd), then ./gradlew installDebug. Caveats: Android SDK download is large (~2GB). First build takes a long time. sdkmanager license acceptance must be automated (yes | sdkmanager --licenses). Emulator requires hardware acceleration (HAXM or HVF on Apple Silicon).

## Acceptance Criteria

1. setup.sh installs JDK + Android SDK + creates AVD. 2. dade dev starts emulator and installs app. 3. Works on Apple Silicon Macs.

