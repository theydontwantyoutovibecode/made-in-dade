package manifest

import (
	"strings"
	"testing"
)

func TestParseAndValidateManifest(t *testing.T) {
	input := strings.Join([]string{
		"[template]",
		"name = \"hypertext\"",
		"description = \"HTML + HTMX\"",
		"version = \"1.0.0\"",
		"",
		"[scaffold]",
		"exclude = [\".git\", \".DS_Store\"]",
		"setup = \"./setup.sh\"",
		"setup_interactive = true",
		"",
		"[serve]",
		"type = \"command\"",
		"dev = \"./start.sh --dev\"",
		"prod = \"./start.sh\"",
		"port_env = \"PORT\"",
		"default_port = 8000",
		"",
		"[serve.static]",
		"root = \".\"",
		"extensions = [\".html\"]",
		"",
		"[project]",
		"marker_fields = [\"template\", \"port\"]",
	}, "\n")

	manifest, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if manifest.Template.Name != "hypertext" || manifest.Template.Description != "HTML + HTMX" {
		t.Fatalf("unexpected template values")
	}
	if manifest.Scaffold.Setup != "./setup.sh" || !manifest.Scaffold.SetupInteractive {
		t.Fatalf("unexpected scaffold values")
	}
	if manifest.Serve.Type != "command" || manifest.Serve.Dev == "" || manifest.Serve.Prod == "" {
		t.Fatalf("unexpected serve values")
	}
	if manifest.Serve.DefaultPort != 8000 {
		t.Fatalf("unexpected default port")
	}
	if len(manifest.Scaffold.Exclude) != 2 || manifest.Scaffold.Exclude[0] != ".git" {
		t.Fatalf("unexpected exclude list")
	}
	if len(manifest.Serve.Static.Extensions) != 1 || manifest.Serve.Static.Extensions[0] != ".html" {
		t.Fatalf("unexpected static extensions")
	}
	if len(manifest.Project.MarkerFields) != 2 {
		t.Fatalf("unexpected marker fields")
	}

	if err := Validate(manifest); err != nil {
		t.Fatalf("validate: %v", err)
	}
	if TemplateName(manifest) != "hypertext" {
		t.Fatalf("unexpected template name helper")
	}
	if ServeType(manifest) != "command" {
		t.Fatalf("unexpected serve type helper")
	}
	if ServeCommand(manifest, "dev") == "" || ServeCommand(manifest, "prod") == "" {
		t.Fatalf("expected serve command helpers")
	}
}

func TestValidateManifestErrors(t *testing.T) {
	manifest := Manifest{}
	if err := Validate(manifest); err == nil {
		t.Fatalf("expected validation error")
	}

	manifest.Template.Name = "Invalid"
	manifest.Template.Description = "desc"
	manifest.Serve.Type = "bad"
	if err := Validate(manifest); err == nil {
		t.Fatalf("expected validation error")
	}

	manifest.Template.Name = "valid-name"
	manifest.Serve.Type = "command"
	manifest.Serve.Dev = ""
	manifest.Serve.Prod = ""
	if err := Validate(manifest); err == nil {
		t.Fatalf("expected validation error for command serve")
	}
}

func TestParseInvalidManifestReturnsError(t *testing.T) {
	input := "[template\nname = \"oops\"\n"
	if _, err := Parse([]byte(input)); err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestParseDevAndShareSections(t *testing.T) {
	input := strings.Join([]string{
		"[template]",
		"name = \"test-template\"",
		"description = \"Test template with dev and share\"",
		"",
		"[serve]",
		"type = \"command\"",
		"dev = \"npm run dev\"",
		"",
		"[dev]",
		`setup = ["uv sync --dev", "uv run python manage.py migrate"]`,
		`background = ["bin/tailwindcss --watch"]`,
		`env = ["DJANGO_SETTINGS_MODULE=config.settings.development"]`,
		`setup_script = ".dade/dev-setup.sh"`,
		"",
		"[share]",
		`env = ["DJANGO_EXTRA_ALLOWED_HOSTS=.trycloudflare.com", "DJANGO_CSRF_TRUSTED_ORIGINS=https://*.trycloudflare.com"]`,
		`tunnel_name = "my-tunnel"`,
		`tunnel_domain = "myapp.example.com"`,
	}, "\n")

	manifest, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	// Validate [dev] section
	if !HasDevSection(manifest) {
		t.Fatalf("expected HasDevSection to return true")
	}
	if len(manifest.Dev.Setup) != 2 {
		t.Fatalf("expected 2 setup commands, got %d", len(manifest.Dev.Setup))
	}
	if manifest.Dev.Setup[0] != "uv sync --dev" {
		t.Fatalf("unexpected setup command: %s", manifest.Dev.Setup[0])
	}
	if len(manifest.Dev.Background) != 1 || manifest.Dev.Background[0] != "bin/tailwindcss --watch" {
		t.Fatalf("unexpected background commands")
	}
	if len(manifest.Dev.Env) != 1 || manifest.Dev.Env[0] != "DJANGO_SETTINGS_MODULE=config.settings.development" {
		t.Fatalf("unexpected env vars")
	}
	if manifest.Dev.SetupScript != ".dade/dev-setup.sh" {
		t.Fatalf("unexpected setup script: %s", manifest.Dev.SetupScript)
	}

	// Validate helper functions
	if len(DevSetupCommands(manifest)) != 2 {
		t.Fatalf("DevSetupCommands returned wrong count")
	}
	if len(DevBackgroundCommands(manifest)) != 1 {
		t.Fatalf("DevBackgroundCommands returned wrong count")
	}
	if len(DevEnv(manifest)) != 1 {
		t.Fatalf("DevEnv returned wrong count")
	}
	if DevSetupScript(manifest) != ".dade/dev-setup.sh" {
		t.Fatalf("DevSetupScript returned wrong value")
	}

	// Validate [share] section
	if !HasShareSection(manifest) {
		t.Fatalf("expected HasShareSection to return true")
	}
	if len(manifest.Share.Env) != 2 {
		t.Fatalf("expected 2 share env vars, got %d", len(manifest.Share.Env))
	}
	if manifest.Share.TunnelName != "my-tunnel" {
		t.Fatalf("unexpected tunnel name: %s", manifest.Share.TunnelName)
	}
	if manifest.Share.TunnelDomain != "myapp.example.com" {
		t.Fatalf("unexpected tunnel domain: %s", manifest.Share.TunnelDomain)
	}

	// Validate helper functions
	if len(ShareEnv(manifest)) != 2 {
		t.Fatalf("ShareEnv returned wrong count")
	}
	if ShareTunnelName(manifest) != "my-tunnel" {
		t.Fatalf("ShareTunnelName returned wrong value")
	}
	if ShareTunnelDomain(manifest) != "myapp.example.com" {
		t.Fatalf("ShareTunnelDomain returned wrong value")
	}

	// Validate manifest overall
	if err := Validate(manifest); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

func TestEmptyDevShareSections(t *testing.T) {
	input := strings.Join([]string{
		"[template]",
		"name = \"minimal\"",
		"description = \"Minimal template\"",
		"",
		"[serve]",
		"type = \"command\"",
		"dev = \"npm run dev\"",
	}, "\n")

	manifest, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if HasDevSection(manifest) {
		t.Fatalf("expected HasDevSection to return false for empty dev section")
	}
	if HasShareSection(manifest) {
		t.Fatalf("expected HasShareSection to return false for empty share section")
	}
}

func TestParseBuildSection(t *testing.T) {
	input := strings.Join([]string{
		"[template]",
		"name = \"cli\"",
		"description = \"CLI application\"",
		"",
		"[serve]",
		"type = \"command\"",
		"dev = \"go run .\"",
		"",
		"[build]",
		`command = "go build -o {{output}}/{{name}} ."`,
		`output = "./bin"`,
		`release_flags = "-ldflags '-s -w'"`,
		`pre = ["go mod tidy", "go generate ./..."]`,
		`post = ["echo done"]`,
	}, "\n")

	manifest, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if !HasBuildSection(manifest) {
		t.Fatalf("expected HasBuildSection to return true")
	}
	if BuildCommand(manifest) != "go build -o {{output}}/{{name}} ." {
		t.Fatalf("unexpected build command: %s", BuildCommand(manifest))
	}
	if BuildOutput(manifest) != "./bin" {
		t.Fatalf("unexpected build output: %s", BuildOutput(manifest))
	}
	if BuildReleaseFlags(manifest) != "-ldflags '-s -w'" {
		t.Fatalf("unexpected release flags: %s", BuildReleaseFlags(manifest))
	}
	if len(BuildPreCommands(manifest)) != 2 {
		t.Fatalf("expected 2 pre commands, got %d", len(BuildPreCommands(manifest)))
	}
	if BuildPreCommands(manifest)[0] != "go mod tidy" {
		t.Fatalf("unexpected pre command: %s", BuildPreCommands(manifest)[0])
	}
	if len(BuildPostCommands(manifest)) != 1 {
		t.Fatalf("expected 1 post command, got %d", len(BuildPostCommands(manifest)))
	}
}

func TestBuildSectionDefaults(t *testing.T) {
	input := strings.Join([]string{
		"[template]",
		"name = \"minimal\"",
		"description = \"Minimal\"",
		"",
		"[serve]",
		"type = \"command\"",
		"dev = \"go run .\"",
	}, "\n")

	manifest, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if HasBuildSection(manifest) {
		t.Fatalf("expected HasBuildSection to return false")
	}
	if BuildOutput(manifest) != "./bin" {
		t.Fatalf("expected default output ./bin, got: %s", BuildOutput(manifest))
	}
	if len(BuildPreCommands(manifest)) != 0 {
		t.Fatalf("expected empty pre commands")
	}
	if len(BuildPostCommands(manifest)) != 0 {
		t.Fatalf("expected empty post commands")
	}
}

func TestNeedsProxy(t *testing.T) {
	// default_port > 0 implies proxy needed
	webApp := Manifest{Serve: Serve{DefaultPort: 8000}}
	if !NeedsProxy(webApp) {
		t.Fatalf("expected NeedsProxy=true for default_port=8000")
	}

	// default_port = 0 implies no proxy
	cli := Manifest{Serve: Serve{DefaultPort: 0}}
	if NeedsProxy(cli) {
		t.Fatalf("expected NeedsProxy=false for default_port=0")
	}

	// explicit proxy=true overrides
	proxyTrue := true
	explicit := Manifest{Serve: Serve{DefaultPort: 0, Proxy: &proxyTrue}}
	if !NeedsProxy(explicit) {
		t.Fatalf("expected NeedsProxy=true when explicitly set")
	}

	// explicit proxy=false overrides
	proxyFalse := false
	explicitOff := Manifest{Serve: Serve{DefaultPort: 8000, Proxy: &proxyFalse}}
	if NeedsProxy(explicitOff) {
		t.Fatalf("expected NeedsProxy=false when explicitly set to false")
	}
}

func TestParseDevMessages(t *testing.T) {
	input := strings.Join([]string{
		"[template]",
		"name = \"ios-app\"",
		"description = \"iOS app\"",
		"",
		"[serve]",
		"type = \"command\"",
		"dev = \"./dev.sh\"",
		"proxy = false",
		"default_port = 0",
		"",
		"[dev.messages]",
		"ready = \"App running on simulator\"",
		"running = \"Press Ctrl+C to stop\"",
	}, "\n")

	manifest, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if NeedsProxy(manifest) {
		t.Fatalf("expected NeedsProxy=false for proxy=false")
	}
	if DevReadyMessage(manifest) != "App running on simulator" {
		t.Fatalf("unexpected ready message: %s", DevReadyMessage(manifest))
	}
	if DevRunningMessage(manifest) != "Press Ctrl+C to stop" {
		t.Fatalf("unexpected running message: %s", DevRunningMessage(manifest))
	}
}

func TestParseProxyField(t *testing.T) {
	input := strings.Join([]string{
		"[template]",
		"name = \"web-app\"",
		"description = \"Web app\"",
		"",
		"[serve]",
		"type = \"command\"",
		"dev = \"npm run dev\"",
		"proxy = true",
		"default_port = 3000",
	}, "\n")

	manifest, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if !NeedsProxy(manifest) {
		t.Fatalf("expected NeedsProxy=true for proxy=true")
	}
	if manifest.Serve.Proxy == nil {
		t.Fatalf("expected Proxy to be non-nil")
	}
}
