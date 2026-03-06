package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	osexec "os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
	execx "github.com/theydontwantyoutovibecode/made-in-dade/internal/exec"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/fsutil"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/manifest"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/proxy"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/registry"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/ui"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/version"
	"github.com/charmbracelet/huh"
	"golang.org/x/term"
)

var projectNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

type pluginTemplate struct {
	Name     string
	Path     string
	Manifest manifest.Manifest
}

type newCommand struct {
	runner       execx.Runner
	copyDir      func(src, dst string) error
	removeGitDir func(path string) error
	isExecutable func(path string) (bool, error)
	templatesDir func() (string, error)
	readDir      func(string) ([]os.DirEntry, error)
	readFile     func(string) ([]byte, error)
	writeMarker  func(projectDir, name, template string, port int) (registry.Marker, error)
	migrateSrv   func(projectDir string) (registry.Marker, bool, error)
	register     func(path, name string, port int, projectPath, template string) (registry.Project, error)
	nextPort     func(path string) (int, error)
	projectsFile func() (string, error)
	caddyfilePath func() (string, error)
	generateCaddy func(context.Context, execx.Runner, string, string) error
	reloadProxy   func(context.Context, execx.Runner, string) error
	spin         func(message string, work func() error) error
}

var newCommandFactory = defaultNewCommand

func defaultNewCommand() newCommand {
	return newCommand{
		runner:        execx.NewSystemRunner(),
		copyDir:       fsutil.CopyDir,
		removeGitDir:  fsutil.RemoveGitDir,
		isExecutable:  fsutil.IsExecutable,
		templatesDir:  config.TemplatesDir,
		readDir:       os.ReadDir,
		readFile:      os.ReadFile,
		writeMarker:   registry.WriteMarker,
		migrateSrv:    registry.MigrateSrvMarker,
		register:      registry.Register,
		nextPort:      registry.NextPort,
		projectsFile:  config.ProjectsFile,
		caddyfilePath: config.CaddyfilePath,
		generateCaddy: proxy.GenerateCaddyfile,
		reloadProxy:   proxy.ReloadProxy,
	}
}

func runNew(args []string, console *ui.UI, logger *logging.Logger) int {
	cmd := newCommandFactory()
	interactive := term.IsTerminal(int(os.Stdin.Fd()))
	return cmd.run(context.Background(), args, console, logger, interactive)
}

func (c newCommand) run(ctx context.Context, args []string, console *ui.UI, logger *logging.Logger, interactive bool) int {
	if _, err := config.InitConfig(); err != nil {
		logger.Error("Failed to initialize config")
		return 1
	}

	if err := ensureDefaultTemplates(ctx, logger); err != nil {
		logger.Info(fmt.Sprintf("Default template install: %v", err))
	}

	projectName := ""
	localPath := ""
	templateName := ""

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--local":
			if i+1 >= len(args) {
				logger.Error("Missing value for --local")
				return 1
			}
			localPath = args[i+1]
			i++
		case "--template":
			if i+1 >= len(args) {
				logger.Error("Missing value for --template")
				return 1
			}
			templateName = args[i+1]
			i++
		default:
			if strings.HasPrefix(arg, "-") {
				logger.Error(fmt.Sprintf("Unknown option: %s", arg))
				return 1
			}
			if projectName == "" {
				projectName = arg
			}
		}
	}

	console.PrintHeader("dade", version.Version)

	if projectName == "" {
		if !interactive {
			logger.Error("Project name is required")
			return 1
		}
		name, err := promptProjectName()
		if err != nil {
			logger.Error("Failed to read project name")
			return 1
		}
		projectName = name
	}

	if projectName == "" {
		logger.Error("Project name is required")
		return 1
	}

	if !projectNamePattern.MatchString(projectName) {
		logger.Error("Invalid project name. Use letters, numbers, hyphens, and underscores. Must start with a letter.")
		return 1
	}

	if info, err := os.Stat(projectName); err == nil && info != nil {
		logger.Error(fmt.Sprintf("Directory '%s' already exists", projectName))
		return 1
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Error("Failed to check target directory")
		return 1
	}

	scaffold := manifest.Scaffold{}

	if localPath == "" {
		templatesDir, err := c.templatesDir()
		if err != nil {
			logger.Error("Failed to resolve templates directory")
			return 1
		}
		installed, err := loadPluginTemplates(templatesDir, c.readDir, c.readFile)
		if err != nil {
			logger.Error("Failed to load installed templates")
			return 1
		}
		if len(installed) == 0 {
			logger.Error("No templates installed")
			logger.Info("Install one: dade install --list-official")
			return 1
		}
		if templateName == "" {
			if len(installed) == 1 {
				templateName = installed[0].Name
			} else {
				picked, err := choosePluginTemplate(installed, interactive)
				if err != nil {
					logger.Error(err.Error())
					return 1
				}
				templateName = picked
			}
		}
		match, ok := findPluginTemplate(installed, templateName)
		if !ok {
			logger.Error(fmt.Sprintf("Template '%s' not found", templateName))
			return 1
		}
		scaffold = match.Manifest.Scaffold
		localPath = match.Path
	} else if templateName == "" {
		templateName = filepath.Base(localPath)
	}

	logger.Info(fmt.Sprintf("Creating project: %s", projectName))
	if templateName != "" {
		logger.Info(fmt.Sprintf("Template: %s", templateName))
	}

	spin := c.spin
	if spin == nil {
		spinnerEnabled := interactive && term.IsTerminal(int(os.Stdout.Fd()))
		spinner := ui.NewSpinner(os.Stdout, spinnerEnabled)
		spin = spinner.Run
	}

	if localPath != "" {
		if info, err := os.Stat(localPath); err != nil || !info.IsDir() {
			logger.Error(fmt.Sprintf("Local template path does not exist: %s", localPath))
			return 1
		}
		excludes := normalizeExcludes(scaffold.Exclude)
		if err := spin("Copying template", func() error { return copyTemplate(localPath, projectName, excludes) }); err != nil {
			logger.Error("Failed to copy template")
			return 1
		}
	}

	if err := c.removeGitDir(projectName); err != nil {
		logger.Error("Failed to remove .git directory")
		return 1
	}

	logger.Success("Template copied")

	ticketsDir := filepath.Join(projectName, ".tickets")
	if err := os.MkdirAll(ticketsDir, 0755); err != nil {
		logger.Warn(fmt.Sprintf("Failed to create .tickets directory: %v", err))
	}

	if !execx.CommandAvailable(c.runner, "git") {
		logger.Error("git is required to initialize repository")
		return 1
	}
	if err := spin("Initializing git repository", func() error {
		return c.runner.Run(ctx, "git", "-C", projectName, "init")
	}); err != nil {
		logger.Error("Failed to initialize git repository")
		return 1
	}
	logger.Success("Git repository initialized")

	projectsPath, err := c.projectsFile()
	if err != nil {
		logger.Error("Failed to resolve projects registry")
		return 1
	}
	port, err := c.nextPort(projectsPath)
	if err != nil {
		logger.Error("Failed to assign project port")
		return 1
	}
	fullPath, err := filepath.Abs(projectName)
	if err != nil {
		logger.Error("Failed to resolve project path")
		return 1
	}
	if c.register != nil {
		if _, err := c.register(projectsPath, projectName, port, fullPath, templateName); err != nil {
			logger.Error("Failed to register project")
			return 1
		}
	}
	if c.writeMarker != nil {
		if _, err := c.writeMarker(projectName, projectName, templateName, port); err != nil {
			logger.Error("Failed to write .dade marker")
			return 1
		}
	}
	if c.migrateSrv != nil {
		if _, migrated, err := c.migrateSrv(projectName); err != nil {
			logger.Error("Failed to migrate .srv marker")
			return 1
		} else if migrated {
			logger.Info("Migrated .srv marker to .dade")
		}
	}

	caddyfilePath, err := c.caddyfilePath()
	if err != nil {
		logger.Error("Failed to resolve Caddyfile path")
		return 1
	}
	if c.generateCaddy != nil {
		if err := c.generateCaddy(ctx, c.runner, projectsPath, caddyfilePath); err != nil {
			logger.Error("Failed to generate Caddyfile")
			return 1
		}
	}
	if c.reloadProxy != nil {
		if err := c.reloadProxy(ctx, c.runner, caddyfilePath); err != nil {
			logger.Error("Failed to reload proxy")
			return 1
		}
	}

	if scaffold.Setup != "" {
		if !interactive && scaffold.SetupInteractive {
			logger.Info("Skipping setup (requires a TTY)")
		} else if err := runScaffoldSetup(ctx, projectName, scaffold.Setup, interactive); err != nil {
			logger.Error("Failed to run setup")
			return 1
		}
	}

	logger.Success(fmt.Sprintf("Project '%s' created successfully!", projectName))
	console.PrintHelp(nextStepsText(projectName))
	return 0
}

func loadPluginTemplates(templatesDir string, readDir func(string) ([]os.DirEntry, error), readFile func(string) ([]byte, error)) ([]pluginTemplate, error) {
	entries, err := readDir(templatesDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	templates := make([]pluginTemplate, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		path := filepath.Join(templatesDir, entry.Name())
		manifestPath := filepath.Join(path, "dade.toml")
		data, err := readFile(manifestPath)
		if err != nil {
			continue
		}
		parsed, err := manifest.Parse(data)
		if err != nil {
			continue
		}
		if err := manifest.Validate(parsed); err != nil {
			continue
		}
		templates = append(templates, pluginTemplate{Name: entry.Name(), Path: path, Manifest: parsed})
	}
	return templates, nil
}

func findPluginTemplate(templates []pluginTemplate, name string) (pluginTemplate, bool) {
	for _, tpl := range templates {
		if tpl.Name == name {
			return tpl, true
		}
	}
	for _, tpl := range templates {
		for _, alias := range tpl.Manifest.Template.Aliases {
			if alias == name {
				return tpl, true
			}
		}
	}
	return pluginTemplate{}, false
}

func choosePluginTemplate(templates []pluginTemplate, interactive bool) (string, error) {
	if !interactive {
		return "", errors.New("template selection requires a TTY")
	}
	options := pluginTemplateOptions(templates)
	choice := ""
	selectField := huh.NewSelect[string]().Title("Select a template:").Options(options...).Value(&choice)
	if err := selectField.Run(); err == nil && choice != "" {
		return choice, nil
	}
	return promptTemplateNumberPlugin(templates, os.Stdin, os.Stdout)
}

func pluginTemplateOptions(templates []pluginTemplate) []huh.Option[string] {
	var defaults, userInstalled []pluginTemplate
	for _, tpl := range templates {
		if isDefaultTemplate(tpl.Path) {
			defaults = append(defaults, tpl)
		} else {
			userInstalled = append(userInstalled, tpl)
		}
	}

	options := make([]huh.Option[string], 0, len(templates)+2)
	if len(defaults) > 0 {
		for _, tpl := range defaults {
			display := tpl.Name
			if tpl.Manifest.Template.Description != "" {
				display = fmt.Sprintf("%s — %s", tpl.Name, tpl.Manifest.Template.Description)
			}
			options = append(options, huh.NewOption(display, tpl.Name))
		}
	}
	if len(userInstalled) > 0 {
		for _, tpl := range userInstalled {
			display := tpl.Name
			if tpl.Manifest.Template.Description != "" {
				display = fmt.Sprintf("%s — %s", tpl.Name, tpl.Manifest.Template.Description)
			}
			options = append(options, huh.NewOption(display, tpl.Name))
		}
	}
	return options
}

func promptTemplateNumberPlugin(templates []pluginTemplate, input io.Reader, output io.Writer) (string, error) {
	if len(templates) == 0 {
		return "", errors.New("no templates available")
	}

	var defaults, userInstalled []pluginTemplate
	for _, tpl := range templates {
		if isDefaultTemplate(tpl.Path) {
			defaults = append(defaults, tpl)
		} else {
			userInstalled = append(userInstalled, tpl)
		}
	}

	fmt.Fprintln(output, "Select a template:")
	idx := 1
	ordered := make([]pluginTemplate, 0, len(templates))
	if len(defaults) > 0 {
		fmt.Fprintln(output, "\n  Default Templates:")
		for _, tpl := range defaults {
			display := tpl.Name
			if tpl.Manifest.Template.Description != "" {
				display = fmt.Sprintf("%s — %s", tpl.Name, tpl.Manifest.Template.Description)
			}
			fmt.Fprintf(output, "  %d) %s\n", idx, display)
			ordered = append(ordered, tpl)
			idx++
		}
	}
	if len(userInstalled) > 0 {
		fmt.Fprintln(output, "\n  User-Installed Templates:")
		for _, tpl := range userInstalled {
			display := tpl.Name
			if tpl.Manifest.Template.Description != "" {
				display = fmt.Sprintf("%s — %s", tpl.Name, tpl.Manifest.Template.Description)
			}
			fmt.Fprintf(output, "  %d) %s\n", idx, display)
			ordered = append(ordered, tpl)
			idx++
		}
	}
	fmt.Fprint(output, "\nChoice: ")
	reader := bufio.NewReader(input)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)
	choice, err := strconv.Atoi(line)
	if err != nil || choice < 1 || choice > len(ordered) {
		return "", errors.New("invalid selection")
	}
	return ordered[choice-1].Name, nil
}

func normalizeExcludes(excludes []string) []string {
	base := []string{".git", "dade.toml", ".source", ".dade"}
	combined := append([]string{}, base...)
	for _, item := range excludes {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		combined = append(combined, item)
	}
	return combined
}

func shouldExclude(rel string, excludes []string) bool {
	rel = filepath.Clean(rel)
	for _, pattern := range excludes {
		if rel == pattern {
			return true
		}
		if strings.HasPrefix(rel, pattern+string(os.PathSeparator)) {
			return true
		}
		if match, _ := filepath.Match(pattern, filepath.Base(rel)); match {
			return true
		}
	}
	return false
}

func copyTemplate(src, dst string, excludes []string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("source is not a directory")
	}
	if err := os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if shouldExclude(rel, excludes) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		target := filepath.Join(dst, rel)
		info, err := d.Info()
		if err != nil {
			return err
		}
		if d.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(path, target, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func runScaffoldSetup(ctx context.Context, projectDir, setupCmd string, interactive bool) error {
	cmd := osexec.CommandContext(ctx, "bash", "-c", setupCmd)
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if interactive {
		cmd.Stdin = os.Stdin
	}
	return cmd.Run()
}

func promptProjectName() (string, error) {
	name := ""
	input := huh.NewInput().Title("Project name:").Placeholder("my-project").Value(&name)
	if err := input.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(name), nil
}

func nextStepsText(projectName string) string {
	return strings.Join([]string{
		"",
		"Next steps:",
		fmt.Sprintf("  cd %s", projectName),
		"  dade dev",
	}, "\n")
}
