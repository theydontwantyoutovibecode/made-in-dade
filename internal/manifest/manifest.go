package manifest

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Template struct {
	Name        string
	Description string
	Version     string
	Author      string
	URL         string
	Aliases     []string
}

type Scaffold struct {
	Exclude          []string
	Setup            string
	SetupInteractive bool
}

type Serve struct {
	Type        string
	Dev         string
	Prod        string
	PortEnv     string
	DefaultPort int
	Proxy       *bool
	Static      ServeStatic
}

type ServeStatic struct {
	Root       string
	Extensions []string
}

type Project struct {
	MarkerFields []string
}

type DevMessages struct {
	Ready   string
	Running string
}

type Dev struct {
	Setup       []string
	Background  []string
	Env         []string
	SetupScript string
	Messages    DevMessages
}

type Share struct {
	Env          []string
	TunnelName   string
	TunnelDomain string
}

type Prod struct {
	Setup       []string
	Env         []string
	SetupScript string
}

type BuildTarget struct {
	OS   string
	Arch string
}

type Build struct {
	Command      string
	Output       string
	Targets      []BuildTarget
	ReleaseFlags string
	Pre          []string
	Post         []string
}

type Manifest struct {
	Template Template
	Scaffold Scaffold
	Serve    Serve
	Project  Project
	Dev      Dev
	Share    Share
	Prod     Prod
	Build    Build
}

var namePattern = regexp.MustCompile(`^[a-z0-9-]+$`)

func Parse(data []byte) (Manifest, error) {
	manifest := Manifest{}
	section := ""
	lines := strings.Split(string(data), "\n")
	for i, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") {
			if !strings.HasSuffix(line, "]") {
				return Manifest{}, fmt.Errorf("invalid section header at line %d", i+1)
			}
			section = strings.TrimSuffix(strings.TrimPrefix(line, "["), "]")
			if section == "" {
				return Manifest{}, fmt.Errorf("invalid section header at line %d", i+1)
			}
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return Manifest{}, fmt.Errorf("invalid assignment at line %d", i+1)
		}
		key := strings.TrimSpace(parts[0])
		if key == "" {
			return Manifest{}, fmt.Errorf("invalid assignment at line %d", i+1)
		}
		value := strings.TrimSpace(parts[1])
		if value == "" {
			return Manifest{}, fmt.Errorf("invalid assignment at line %d", i+1)
		}

		switch section {
		case "template":
			assignTemplate(&manifest.Template, key, value)
		case "scaffold":
			assignScaffold(&manifest.Scaffold, key, value)
		case "serve":
			assignServe(&manifest.Serve, key, value)
		case "serve.static":
			assignServeStatic(&manifest.Serve.Static, key, value)
		case "project":
			assignProject(&manifest.Project, key, value)
		case "dev":
			assignDev(&manifest.Dev, key, value)
		case "dev.messages":
			assignDevMessages(&manifest.Dev.Messages, key, value)
		case "share":
			assignShare(&manifest.Share, key, value)
		case "prod":
			assignProd(&manifest.Prod, key, value)
		case "build":
			assignBuild(&manifest.Build, key, value)
		case "build.targets":
			assignBuildTarget(&manifest.Build, key, value)
		}
	}
	return manifest, nil
}

func Validate(manifest Manifest) error {
	var errs []string
	if manifest.Template.Name == "" {
		errs = append(errs, "template.name is required")
	} else if !namePattern.MatchString(manifest.Template.Name) {
		errs = append(errs, "template.name must match [a-z0-9-]+")
	}
	if manifest.Template.Description == "" {
		errs = append(errs, "template.description is required")
	}
	if manifest.Serve.Type == "" {
		errs = append(errs, "serve.type is required")
	} else if manifest.Serve.Type != "static" && manifest.Serve.Type != "command" {
		errs = append(errs, "serve.type must be static or command")
	}
	if manifest.Serve.Type == "command" {
		if manifest.Serve.Dev == "" && manifest.Serve.Prod == "" {
			errs = append(errs, "serve.dev or serve.prod is required for command templates")
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

func TemplateName(manifest Manifest) string {
	return manifest.Template.Name
}

func TemplateAliases(manifest Manifest) []string {
	return manifest.Template.Aliases
}

func ServeType(manifest Manifest) string {
	return manifest.Serve.Type
}

func HasServeSection(manifest Manifest) bool {
	return manifest.Serve.Type != ""
}

func ServeCommand(manifest Manifest, mode string) string {
	if mode == "prod" {
		if manifest.Serve.Prod != "" {
			return manifest.Serve.Prod
		}
		return manifest.Serve.Dev
	}
	if manifest.Serve.Dev != "" {
		return manifest.Serve.Dev
	}
	return manifest.Serve.Prod
}

func NeedsProxy(manifest Manifest) bool {
	if manifest.Serve.Proxy != nil {
		return *manifest.Serve.Proxy
	}
	return manifest.Serve.DefaultPort > 0
}

func DevReadyMessage(manifest Manifest) string {
	return manifest.Dev.Messages.Ready
}

func DevRunningMessage(manifest Manifest) string {
	return manifest.Dev.Messages.Running
}

func HasDevSection(manifest Manifest) bool {
	return len(manifest.Dev.Setup) > 0 ||
		len(manifest.Dev.Background) > 0 ||
		len(manifest.Dev.Env) > 0 ||
		manifest.Dev.SetupScript != ""
}

func HasShareSection(manifest Manifest) bool {
	return len(manifest.Share.Env) > 0 ||
		manifest.Share.TunnelName != "" ||
		manifest.Share.TunnelDomain != ""
}

func DevSetupCommands(manifest Manifest) []string {
	return manifest.Dev.Setup
}

func DevBackgroundCommands(manifest Manifest) []string {
	return manifest.Dev.Background
}

func DevEnv(manifest Manifest) []string {
	return manifest.Dev.Env
}

func DevSetupScript(manifest Manifest) string {
	return manifest.Dev.SetupScript
}

func ShareEnv(manifest Manifest) []string {
	return manifest.Share.Env
}

func ShareTunnelName(manifest Manifest) string {
	return manifest.Share.TunnelName
}

func ShareTunnelDomain(manifest Manifest) string {
	return manifest.Share.TunnelDomain
}

func HasProdSection(manifest Manifest) bool {
	return len(manifest.Prod.Setup) > 0 ||
		len(manifest.Prod.Env) > 0 ||
		manifest.Prod.SetupScript != ""
}

func ProdSetupCommands(manifest Manifest) []string {
	return manifest.Prod.Setup
}

func ProdEnv(manifest Manifest) []string {
	return manifest.Prod.Env
}

func ProdSetupScript(manifest Manifest) string {
	return manifest.Prod.SetupScript
}

func HasBuildSection(manifest Manifest) bool {
	return manifest.Build.Command != ""
}

func BuildCommand(manifest Manifest) string {
	return manifest.Build.Command
}

func BuildOutput(manifest Manifest) string {
	if manifest.Build.Output != "" {
		return manifest.Build.Output
	}
	return "./bin"
}

func BuildReleaseFlags(manifest Manifest) string {
	return manifest.Build.ReleaseFlags
}

func BuildPreCommands(manifest Manifest) []string {
	return manifest.Build.Pre
}

func BuildPostCommands(manifest Manifest) []string {
	return manifest.Build.Post
}

func BuildTargets(manifest Manifest) []BuildTarget {
	return manifest.Build.Targets
}

func assignTemplate(template *Template, key, value string) {
	switch key {
	case "name":
		template.Name = trimString(value)
	case "description":
		template.Description = trimString(value)
	case "version":
		template.Version = trimString(value)
	case "author":
		template.Author = trimString(value)
	case "url":
		template.URL = trimString(value)
	case "aliases":
		template.Aliases = trimArray(value)
	}
}

func assignScaffold(scaffold *Scaffold, key, value string) {
	switch key {
	case "exclude":
		scaffold.Exclude = trimArray(value)
	case "setup":
		scaffold.Setup = trimString(value)
	case "setup_interactive":
		scaffold.SetupInteractive = trimBool(value)
	}
}

func assignServe(serve *Serve, key, value string) {
	switch key {
	case "type":
		serve.Type = trimString(value)
	case "dev":
		serve.Dev = trimString(value)
	case "prod":
		serve.Prod = trimString(value)
	case "port_env":
		serve.PortEnv = trimString(value)
	case "default_port":
		serve.DefaultPort = trimInt(value)
	case "proxy":
		v := trimBool(value)
		serve.Proxy = &v
	}
}

func assignServeStatic(static *ServeStatic, key, value string) {
	switch key {
	case "root":
		static.Root = trimString(value)
	case "extensions":
		static.Extensions = trimArray(value)
	}
}

func assignProject(project *Project, key, value string) {
	switch key {
	case "marker_fields":
		project.MarkerFields = trimArray(value)
	}
}

func assignDev(dev *Dev, key, value string) {
	switch key {
	case "setup":
		dev.Setup = trimArray(value)
	case "background":
		dev.Background = trimArray(value)
	case "env":
		dev.Env = trimArray(value)
	case "setup_script":
		dev.SetupScript = trimString(value)
	}
}

func assignDevMessages(msg *DevMessages, key, value string) {
	switch key {
	case "ready":
		msg.Ready = trimString(value)
	case "running":
		msg.Running = trimString(value)
	}
}

func assignShare(share *Share, key, value string) {
	switch key {
	case "env":
		share.Env = trimArray(value)
	case "tunnel_name":
		share.TunnelName = trimString(value)
	case "tunnel_domain":
		share.TunnelDomain = trimString(value)
	}
}

func assignProd(prod *Prod, key, value string) {
	switch key {
	case "setup":
		prod.Setup = trimArray(value)
	case "env":
		prod.Env = trimArray(value)
	case "setup_script":
		prod.SetupScript = trimString(value)
	}
}

func assignBuild(build *Build, key, value string) {
	switch key {
	case "command":
		build.Command = trimString(value)
	case "output":
		build.Output = trimString(value)
	case "release_flags":
		build.ReleaseFlags = trimString(value)
	case "pre":
		build.Pre = trimArray(value)
	case "post":
		build.Post = trimArray(value)
	}
}

func assignBuildTarget(build *Build, key, value string) {
	switch key {
	case "os":
		build.Targets = append(build.Targets, BuildTarget{OS: trimString(value)})
	case "arch":
		if len(build.Targets) > 0 {
			build.Targets[len(build.Targets)-1].Arch = trimString(value)
		}
	}
}

func trimString(value string) string {
	value = strings.TrimSpace(value)
	return strings.Trim(value, "\"")
}

func trimBool(value string) bool {
	value = trimString(value)
	return value == "true"
}

func trimInt(value string) int {
	value = trimString(value)
	if value == "" {
		return 0
	}
	var parsed int
	_, err := fmt.Sscanf(value, "%d", &parsed)
	if err != nil {
		return 0
	}
	return parsed
}

func trimArray(value string) []string {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "[")
	value = strings.TrimSuffix(value, "]")
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		item := trimString(part)
		if item != "" {
			items = append(items, item)
		}
	}
	return items
}
