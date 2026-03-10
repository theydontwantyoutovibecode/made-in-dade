# dade

dade is a macOS CLI that creates projects from templates, runs development servers, and manages a local HTTPS proxy. It supports six built-in templates covering web apps, static sites, iOS apps, Android apps, CLI tools, and TUI programs.

## Installation

```bash
brew tap theydontwantyoutovibecode/tap
brew install dade
```

Homebrew installs dade and its required dependencies (Caddy, jq, git). The first time you run any command, dade automatically configures itself: it creates `~/.config/dade/`, generates a Caddyfile, starts the HTTPS proxy, and installs the default templates.

To build from source instead:

```bash
go install github.com/theydontwantyoutovibecode/dade/cmd/dade@latest
dade setup
```

Building from source requires Go 1.22+, macOS, and git. You must run `dade setup` manually after `go install` to install dependencies and configure the proxy.

## Setup

Setup runs automatically on first use. You can also run it explicitly:

```bash
dade setup
```

It does the following:

1. Creates `~/.config/dade/` and its subdirectories
2. Checks for required dependencies (Caddy, jq) and offers to install them
3. Generates a Caddyfile and starts Caddy as a launchd service
4. Trusts the local CA certificate so browsers accept local HTTPS
5. Installs the six default templates

Setup is safe to re-run. It skips work that is already complete.

Run `dade setup --check` to verify dependencies without changing anything.

## Creating a Project

```bash
dade new myproject
```

This prompts you to pick a template. You can also specify one directly:

```bash
dade new myproject --template web-app
```

What happens during `dade new`:

1. Clones the template into a new directory
2. Removes the template's `.git` directory
3. Runs the template's `setup.sh` if present (installs language runtimes, SDKs, etc.)
4. Initializes a new git repository
5. Creates a `.tickets/` directory for ticket tracking
6. Registers the project in `~/.config/dade/projects.json`
7. Assigns a port and updates the Caddyfile so the project has an HTTPS URL

## Templates

dade ships with six templates. They are automatically installed the first time you run any command. Each template is a git repository containing a `dade.toml` manifest, a `setup.sh` script, source code, and an `AGENTS.md` file for AI-assisted development.

You can refer to any template by its name or one of its aliases when running `dade new --template <name>`. Use `dade new --inspect` to see details about a template before creating a project.

Not every dade command applies to every template. The table below shows which commands are relevant to each one:

| Command | web-app | web-site | ios-app | android-app | cli | tui |
|---------|---------|----------|---------|-------------|-----|-----|
| `dev` | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| `build` | — | — | ✓ | ✓ | ✓ | ✓ |
| `project start` | ✓ | — | — | — | — | — |
| `share` | ✓ | ✓ | — | — | — | — |
| `share --attach` | ✓ | ✓ | — | — | — | — |
| `dev --open` | ✓ | ✓ | — | — | — | — |

Commands marked `—` are not applicable because the template either has no server (mobile, CLI, TUI) or no compiled output (web). Running them will not produce useful results.

### web-app

A full-stack web application using Django and HTMX. Django handles routing, database access, and server-side rendering. HTMX handles dynamic page updates by fetching HTML fragments from Django views and swapping them into the page without JavaScript frameworks.

Styling uses vanilla CSS with no build step. There is no Tailwind, Sass, or CSS framework.

| | |
|-|-|
| **Aliases** | `webapp`, `django` |
| **Prerequisites** | Python 3.13+, [uv](https://docs.astral.sh/uv/) |
| **Stack** | Django 5.x, HTMX, Gunicorn, uv |
| **Serve type** | Command (`runserver` in dev, `gunicorn` in prod) |
| **Default port** | 8000 |

##### What `dade dev` does

1. Runs `uv sync --dev` to install Python dependencies from `pyproject.toml`
2. Runs `python manage.py migrate --no-input` to apply pending database migrations
3. Starts `python manage.py runserver 127.0.0.1:$PORT`
4. Sets `DJANGO_SETTINGS_MODULE=config.settings.development`

##### What `dade project start` does

1. Runs `uv sync --extra prod` for production dependencies
2. Runs migrations and `collectstatic`
3. Starts Gunicorn with 4 workers

##### What `dade share` does

Starts the dev server and creates a Cloudflare tunnel. Automatically sets `DJANGO_EXTRA_ALLOWED_HOSTS` and `DJANGO_CSRF_TRUSTED_ORIGINS` so Django accepts requests from the tunnel domain.

##### Settings structure

Settings are split into `config/settings/base.py`, `development.py`, and `production.py`. The base file reads `SECRET_KEY`, `DATABASE_URL`, and other values from environment variables. Development uses SQLite. Production enforces HTTPS and requires `SECRET_KEY` to be set explicitly.

---

### web-site

A static website using plain HTML, CSS, and JavaScript with HTMX for dynamic interactions. There is no build step, no bundler, and no package manager. The HTMX library is vendored in `lib/htmx.min.js`.

| | |
|-|-|
| **Aliases** | `website`, `static`, `html` |
| **Prerequisites** | None |
| **Stack** | HTML, CSS, JavaScript, HTMX |
| **Serve type** | Static (dade's built-in file server) |

##### What `dade dev` does

Starts a Caddy static file server on the project's assigned port, serving files from the project root directory. The server runs in the foreground and supports:

- Automatic HTTPS via the dade proxy (available at `https://<name>.localhost`)
- Live reloading of static files
- Proper MIME type handling
- Directory listing for missing index files

There are no setup commands or dependencies to install beyond what dade provides.

##### HTMX partials

The `partials/` directory contains HTML fragments. Pages load these fragments dynamically using HTMX attributes like `hx-get` and `hx-swap`. This lets you build interactive pages without writing JavaScript.

##### Deployment

Copy the entire project to any static hosting service (GitHub Pages, Netlify, Vercel, Cloudflare Pages, or any HTTP server). No build step is required.

---

### ios-app

A native iOS application using SwiftUI with MVVM architecture. The Xcode project is generated from a `project.yml` file using XcodeGen, which avoids `.xcodeproj` merge conflicts.

| | |
|-|-|
| **Aliases** | `ios`, `swift`, `swiftui` |
| **Prerequisites** | macOS, Xcode 16+ |
| **Stack** | Swift, SwiftUI, XcodeGen |
| **Minimum target** | iOS 17 |

##### What `setup.sh` does

1. Installs XcodeGen via Homebrew
2. Runs `xcodegen generate` to create `MyApp.xcodeproj` from `project.yml`
3. Boots an iPhone simulator

##### What `dade dev` does

Opens the Xcode project. Development and debugging happen in Xcode using the iOS Simulator.

##### What `dade build` does

Runs `xcodebuild` targeting the iOS Simulator. Use `--release` for a Release configuration build.

##### Architecture

- Views are SwiftUI structs organized by feature in `Sources/Features/`
- ViewModels use `@Observable` (Swift 5.9+) for state management
- Navigation uses `NavigationStack` with a `Feature` enum for type-safe routing

##### XcodeGen

The `.xcodeproj` file is generated, not checked into git. After editing `project.yml`, run `xcodegen generate` to regenerate it.

---

### android-app

A native Android application using Kotlin and Jetpack Compose with Material 3 theming. The setup script installs the Android SDK, creates an emulator, and launches it automatically.

| | |
|-|-|
| **Aliases** | `android`, `kotlin`, `compose` |
| **Prerequisites** | JDK 17+ |
| **Stack** | Kotlin, Jetpack Compose, Material 3, Gradle |
| **Minimum SDK** | 26 (Android 8.0) |
| **Target SDK** | 35 |

##### What `setup.sh` does

1. Installs JDK 17 via Homebrew
2. Installs Android command-line tools via `brew install --cask android-commandlinetools`
3. Runs `sdkmanager` to install platform-tools, emulator, system images, and build tools
4. Creates an AVD named `dade_pixel`
5. Launches the emulator in the background

##### What `dade dev` does

1. Runs `./gradlew assembleDebug` to build the app
2. Runs `./gradlew installDebug` to install it on the running emulator

##### What `dade build` does

Runs `./gradlew assembleDebug`. The APK is output to `app/build/outputs/apk/`. Use `--release` for `assembleRelease`.

##### Architecture

- Screens are composable functions in `app/src/main/kotlin/.../ui/screens/`
- ViewModels extend `ViewModel` and expose state via `StateFlow`
- Navigation uses Navigation Compose with a `NavHost`
- Dependencies are managed with a Gradle version catalog in `gradle/libs.versions.toml`

---

### cli

A command-line application using Go and Charm libraries. The template includes a root command, an example subcommand, configuration management, and styled terminal output.

| | |
|-|-|
| **Aliases** | `cli-app` |
| **Prerequisites** | Go 1.22+ |
| **Stack** | Cobra, Fang, Lipgloss, Huh, Log |

##### What `dade dev` does

1. Runs `go mod download` to fetch dependencies
2. Runs `go run .` to start the program

##### What `dade build` does

Runs `go build` and outputs the binary to `./bin/`. Runs `go mod tidy` first.

- `--release` strips debug symbols with `-ldflags '-s -w'`
- `--os linux --arch amd64` cross-compiles for a specific platform
- `--all` builds for all supported OS/architecture combinations

##### Libraries included

| Library | Purpose |
|---------|---------|
| [Cobra](https://github.com/spf13/cobra) | Command structure and argument parsing |
| [Fang](https://github.com/charmbracelet/fang) | Configuration from flags, env vars, and config files |
| [Lipgloss](https://github.com/charmbracelet/lipgloss) | Terminal colors, borders, padding |
| [Huh](https://github.com/charmbracelet/huh) | Interactive forms (text inputs, selects, confirms) |
| [Log](https://github.com/charmbracelet/log) | Structured logging with styled output |

##### Adding a command

Create a new file in `cmd/` with a `cobra.Command` variable and register it in an `init()` function using `rootCmd.AddCommand()`. See `cmd/hello.go` for an example.

---

### tui

A terminal UI application using Go and Charm Bubbletea v2. The template follows the Elm Architecture: all state lives in a single `Model` struct, `Update()` handles messages, and `View()` renders the terminal output.

| | |
|-|-|
| **Aliases** | `tui-app`, `terminal` |
| **Prerequisites** | Go 1.22+ |
| **Stack** | Bubbletea v2, Bubbles v2, Lipgloss |

##### What `dade dev` does

1. Runs `go mod download` to fetch dependencies
2. Runs `go run .` with `DEBUG=1` in the environment
3. When `DEBUG=1` is set, the program logs to `debug.log` instead of stdout (stdout is used by the TUI)

##### What `dade build` does

Same as the `cli` template: `go build` to `./bin/`, with `--release`, `--os/--arch`, and `--all` flags available.

##### Elm Architecture

The program loop has three phases that repeat continuously:

1. **Init** (`model.go`) — Creates the initial `Model` and returns any startup commands
2. **Update** (`update.go`) — Receives a `tea.Msg` (keyboard input, timer tick, custom event), updates the `Model`, and optionally returns a `tea.Cmd` for async work
3. **View** (`view.go`) — Renders the `Model` as a string using Lipgloss styles from `styles.go`

Key bindings are defined in `keys.go` using `key.NewBinding`.

##### Bubbletea v2 vs v1

This template uses Bubbletea v2, which has breaking API changes from v1. Most online documentation and tutorials still reference v1. The `.read-only/` reference libraries contain the actual v2 source code. Use those as the authoritative API reference, not external docs.

Key differences from v1:
- `tea.Model` uses generics
- `Init()` returns `(Model, tea.Cmd)` instead of just `tea.Cmd`
- `Update()` and `View()` method signatures changed
- Bubbles components updated to match v2

## Commands

Commands are organized into three groups: Development, Management, and System.

### Development

#### `dade new [name]`

Creates a new project from a template. See [Creating a Project](#creating-a-project) for the full flow.

| Flag | Description |
|------|-------------|
| `--template, -t` | Template name or alias (default: `web-app`) |
| `--local` | Use a local directory as the template source instead of an installed template |
| `--name, -n` | Project name (alternative to the positional argument) |
| `--inspect` | Show template metadata (name, aliases, description, serve type, source URL) without creating a project |

#### `dade dev [name]`

Starts a development server with full orchestration.

1. Syncs `.read-only` reference libraries if a `manifest.txt` exists in the project
2. Runs setup commands defined in the `[dev]` section of the template's `dade.toml`
3. Starts any background processes (asset watchers, etc.)
4. Starts the main dev server
5. Ensures the HTTPS proxy is running and the project's domain is configured

For web templates, the project is available at `https://<name>.<hostname>.localhost` (or `.local` if configured). For non-web templates (ios, android, cli, tui), the command runs the project directly.

| Flag | Description |
|------|-------------|
| `--skip-setup` | Skip the setup commands (steps 1-2) |
| `--open` | Open the project URL in the default browser after starting |
| `--port, -p` | Override the port assigned to this project |

#### `dade build [name]`

Compiles the project. Applicable to templates that produce a binary or artifact: `cli`, `tui`, `ios-app`, and `android-app`. Not applicable to `web-app` or `web-site`.

dade auto-detects the project type by looking for `go.mod` (Go), `*.xcodeproj` (Xcode), or `build.gradle.kts` (Gradle). The `[build]` section in `dade.toml` can override the default build command.

| Flag | Description |
|------|-------------|
| `--output, -o` | Output directory for the built artifact |
| `--os` | Target OS for cross-compilation (Go only) |
| `--arch` | Target architecture for cross-compilation (Go only) |
| `--all` | Build for all supported OS/architecture combinations (Go only) |
| `--release` | Build with release optimizations (strips symbols for Go, uses Release configuration for Xcode, runs `assembleRelease` for Gradle) |

#### `dade share [name]`

Starts the dev server and creates a Cloudflare tunnel so the project is accessible from the public internet. Only applicable to web templates.

Requires `cloudflared`. Install it with `brew install cloudflared`.

The `[share]` section in `dade.toml` can define additional environment variables needed for the tunnel (for example, Django's `ALLOWED_HOSTS` and `CSRF_TRUSTED_ORIGINS`).

| Flag | Description |
|------|-------------|
| `--skip-setup` | Skip setup commands before starting the server |
| `--quick` | Force a quick tunnel even if a named tunnel is configured |
| `--attach` | Attach a tunnel to an already-running server (skip server startup) |
| `--port, -p` | Override the port assigned to this project |

### Management

#### `dade project <subcommand>`

Manage the project registry.

| Subcommand | Description |
|------------|-------------|
| `list` | List all registered projects with status, port, template, path, and HTTPS URL. Use `--running` to filter to running projects. |
| `register [name]` | Register an existing directory as a dade project. Use `-t` to specify the template type. |
| `remove <name>` | Remove a project from the registry. Use `--files` to also delete project files. Alias: `rm`. |
| `port` | Show or update the project port. Use `--set` to change it. |
| `sync [path]` | Rebuild the registry by scanning for `.dade` marker files. Use `--clean` to remove stale entries. |
| `start [name]` | Start a production server using the `[prod]` manifest section. |
| `stop [name]` | Stop a running server by sending SIGTERM. |

Running `dade project` with no subcommand defaults to `list`.

#### `dade template <subcommand>`

Manage installed templates.

| Subcommand | Description |
|------------|-------------|
| `list` | List installed templates with name, description, serve type, and source URL. |
| `add <git-url>` | Install a template from a git repository. Use `--name` to override the name, `--list-official` to browse curated templates. |
| `remove <name>` | Remove an installed template. Use `--all` to remove all. |
| `update <name>` | Re-clone a template from its source. Use `--all` to update all. |

Running `dade template` with no subcommand defaults to `list`.

#### `dade proxy <subcommand>`

Manages the local Caddy HTTPS proxy service that runs as a launchd daemon.

| Subcommand | Description |
|------------|-------------|
| `start` | Start the proxy service |
| `stop` | Stop the proxy service |
| `restart` | Stop and restart the proxy service |
| `status` | Show whether the proxy is running, how many projects are configured, and the Caddyfile path |
| `logs` | Tail the proxy log file. Use `--lines/-n` to control how many lines to show and `--follow/-f` to stream new lines |
| `reload` | Regenerate the Caddyfile from the project registry and reload Caddy. Use `--list` to print all project URLs after reloading. |

### System

#### `dade setup`

First-time setup. See [Setup](#setup).

| Flag | Description |
|------|-------------|
| `--check` | Only check dependencies, don't run setup |
| `--yes, -y` | Answer yes to all prompts |
| `--install-deps` | Install missing dependencies via Homebrew without prompting |
| `--skip-deps` | Skip dependency installation prompts |
| `--no-templates` | Skip installing default templates |

### Backward Compatibility

The following old command forms still work as hidden aliases but no longer appear in help output:

| Old command | New equivalent |
|-------------|----------------|
| `dade templates` | `dade template list` |
| `dade install <url>` | `dade template add <url>` |
| `dade uninstall <name>` | `dade template remove <name>` |
| `dade update <name>` | `dade template update <name>` |
| `dade list` | `dade project list` |
| `dade register [name]` | `dade project register [name]` |
| `dade remove <name>` | `dade project remove <name>` |
| `dade port` | `dade project port` |
| `dade sync [path]` | `dade project sync [path]` |
| `dade start [name]` | `dade project start [name]` |
| `dade stop [name]` | `dade project stop [name]` |
| `dade refresh` | `dade proxy reload` |
| `dade tunnel [name]` | `dade share --attach [name]` |
| `dade open [name]` | `dade dev --open [name]` |

### Global Flags

These flags work with every command:

| Flag | Description |
|------|-------------|
| `--quiet, -q` | Suppress non-essential output |
| `--verbose, -v` | Enable detailed output |
| `--json` | Output in JSON format (supported by `project list`, `template list`, `proxy status`) |
| `--no-color` | Disable colored terminal output |

`--quiet` and `--verbose` cannot be used together.

## HTTPS Proxy

dade runs [Caddy](https://caddyserver.com/) as a launchd service. Caddy generates local TLS certificates and routes HTTPS requests to project ports.

### Domains

By default, dade uses `.localhost` domains which resolve to `127.0.0.1` automatically (no `/etc/hosts` needed):

```
https://<project-name>.<hostname>.localhost
```

For example, if your hostname is `macbook` and you create a project called `myapp`, it is available at `https://myapp.macbook.localhost`.

#### `.localhost` vs `.local`

- **`.localhost`** (default): Resolves to `127.0.0.1` per RFC 6761. Works immediately on this machine. No configuration required.
- **`.local`**: Requires `/etc/hosts` entries for subdomains. Only works on this machine unless you manually configure DNS. Useful if you need LAN access to your development projects.

To switch to `.local`:

```bash
echo 'domain_tld = ".local"' > ~/.config/dade/config.toml
dade proxy reload
```

Then add each project's domain to `/etc/hosts`:

```bash
sudo bash -c "echo '127.0.0.1\t<project-name>.<hostname>.local' >> /etc/hosts"
```

The Caddyfile is stored at `~/.config/dade/Caddyfile` and is regenerated automatically whenever projects are added, removed, or have their ports changed. You should not edit the Caddyfile manually.

## Reference Libraries (.read-only)

Each template includes a `.read-only/manifest.txt` file that lists git repositories relevant to that template's stack. When you run `dade dev`, dade shallow-clones each listed repository into `.read-only/` in your project directory. Repos that are already cloned are skipped.

The purpose is to give AI coding agents access to the actual source code of the libraries your project uses. The cloned repos are read-only reference material, not project dependencies.

| Template | Reference repos |
|----------|----------------|
| web-app | Django, HTMX, django-htmx |
| web-site | HTMX |
| ios-app | swift-collections, swift-algorithms |
| android-app | compose-samples, architecture-samples |
| cli | Lipgloss, Huh, Fang |
| tui | Bubbletea, Bubbles, Lipgloss |

## Ticket Tracking (.tickets)

Every project created with `dade new` includes a `.tickets/` directory for use with the [tk](https://github.com/wedow/tk) CLI. `tk` stores tickets as markdown files and supports creating, listing, closing, and linking them.

Install `tk` with:

```bash
brew install wedow/tap/ticket
```

Each template's `AGENTS.md` file instructs AI agents to use `tk` for planning work before writing code. The intended workflow is:

1. Break a request into granular tickets with `tk create`
2. Work on one ticket at a time
3. One ticket = one commit
4. Close tickets with `tk close` when done

## Template Manifest (dade.toml)

Templates are git repositories with a `dade.toml` file at the root. This manifest tells dade how to scaffold, serve, and build projects created from the template.

### [template]

Metadata about the template itself.

```toml
[template]
name = "my-template"
description = "What this template creates"
version = "1.0.0"
author = "Your Name"
url = "https://github.com/you/my-template"
aliases = ["alt-name", "another-name"]
```

The `aliases` field lets users refer to the template by alternative names (e.g., `dade new --template alt-name`).

### [scaffold]

Controls what happens when `dade new` copies the template into a new project.

```toml
[scaffold]
exclude = [".git", "dade.toml", ".dade", "node_modules"]
setup = "./setup.sh"
setup_interactive = true
```

- `exclude` — Files and directories that should not be copied into the new project
- `setup` — A script to run after copying (installs dependencies, generates project files, etc.)
- `setup_interactive` — Set to `true` if the setup script requires terminal input

### [serve]

Defines how the project is served during development and production.

For templates with a server process:

```toml
[serve]
type = "command"
dev = "python manage.py runserver 127.0.0.1:$PORT"
prod = "gunicorn app:app --bind 0.0.0.0:$PORT"
port_env = "PORT"
default_port = 8000
```

For static file templates:

```toml
[serve]
type = "static"

[serve.static]
root = "."
```

For templates without a server (CLI, TUI, mobile apps), use `type = "command"` with `port_env = ""` and `default_port = 0`. The dev command in `[serve]` runs the program directly (e.g., `go run .` or `open MyApp.xcodeproj`).

### [dev]

Commands and environment variables for `dade dev`.

```toml
[dev]
setup = ["uv sync --dev", "python manage.py migrate"]
env = ["DJANGO_SETTINGS_MODULE=config.settings.development"]
```

- `setup` — Commands run sequentially before the dev server starts
- `env` — Environment variables set in the server process

### [share]

Additional environment variables applied when running `dade share`. Typically used for web frameworks that validate hostnames.

```toml
[share]
env = [
    "DJANGO_EXTRA_ALLOWED_HOSTS=.trycloudflare.com",
    "DJANGO_CSRF_TRUSTED_ORIGINS=https://*.trycloudflare.com"
]
```

### [build]

Controls how `dade build` compiles the project.

```toml
[build]
command = "go build -o {{output}}/{{name}} ."
output = "./bin"
release_flags = "-ldflags '-s -w'"
pre = ["go mod tidy"]
post = []
```

- `command` — The build command. Supports `{{name}}`, `{{output}}`, `{{os}}`, and `{{arch}}` placeholders
- `output` — Directory for build artifacts
- `release_flags` — Appended to the command when `--release` is used
- `pre` / `post` — Commands run before and after the build

### [prod]

Commands and environment variables for `dade project start` (production mode).

```toml
[prod]
setup = ["uv sync --extra prod", "python manage.py collectstatic --no-input"]
env = ["DJANGO_SETTINGS_MODULE=config.settings.production"]
```

## Custom Templates

You can install any git repository as a template:

```bash
dade template add https://github.com/someone/their-template.git
```

The repository must contain a `dade.toml` manifest. After installing, the template appears in `dade new`.

You can also override default template URLs in `~/.config/dade/templates.toml`:

```toml
[templates]
my-template = "https://github.com/you/your-template.git"
```

## Configuration Paths

| Path | Purpose |
|------|---------|
| `~/.config/dade/` | Base configuration directory |
| `~/.config/dade/config.toml` | Optional configuration file (domain TLD, etc.) |
| `~/.config/dade/templates/` | Installed template directories |
| `~/.config/dade/projects.json` | Registry of all dade projects (name, port, path, template) |
| `~/.config/dade/Caddyfile` | Generated Caddy proxy configuration |
| `~/.config/dade/templates.toml` | Optional file to override or add template git URLs |
| `~/.config/dade/proxy.log` | Caddy proxy stdout log |
| `~/.config/dade/proxy.err` | Caddy proxy stderr log |
| `~/Library/LaunchAgents/land.charm.dade.proxy.plist` | launchd service definition for the Caddy proxy |

## License

MIT
