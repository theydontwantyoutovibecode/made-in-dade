package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
	execx "github.com/theydontwantyoutovibecode/made-in-dade/internal/exec"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/logging"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/proxy"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/registry"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/ui"
	"github.com/theydontwantyoutovibecode/made-in-dade/internal/version"
	"github.com/charmbracelet/huh"
	"golang.org/x/term"
)

type setupCommand struct {
	runner          execx.Runner
	initConfig      func() (bool, error)
	detectSrv       func() (bool, error)
	migrateSrv      func(context.Context, execx.Runner, string, *logging.Logger) error
	confirm         func(string) (bool, error)
	spin            func(message string, work func() error) error
	generateCaddy   func(context.Context, execx.Runner, string, string) error
	createPlist     func(context.Context, execx.Runner, string, string, string, string) error
	installProxy    func(context.Context, execx.Runner, string) error
	trustCA         func(context.Context, execx.Runner) error
	installTemplate func(context.Context, string, *ui.UI, *logging.Logger, bool) int
	templatesDir    func() (string, error)
	projectsFile    func() (string, error)
	caddyfilePath   func() (string, error)
	plistPath       func() (string, error)
	logPath         func() (string, error)
	errPath         func() (string, error)
}

func defaultSetupCommand() setupCommand {
	return setupCommand{
		runner:        execx.NewSystemRunner(),
		initConfig:    config.InitConfig,
		detectSrv:     config.DetectSrvConfig,
		migrateSrv:    migrateFromSrv,
		generateCaddy: proxy.GenerateCaddyfile,
		createPlist:   proxy.CreatePlist,
		installProxy:  proxy.InstallProxyService,
		trustCA:       trustCaddyCA,
		templatesDir:  config.TemplatesDir,
		projectsFile:  config.ProjectsFile,
		caddyfilePath: config.CaddyfilePath,
		plistPath:     config.ProxyPlistPath,
		logPath:       config.ProxyLogPath,
		errPath:       config.ProxyErrPath,
		installTemplate: func(ctx context.Context, url string, console *ui.UI, logger *logging.Logger, styled bool) int {
			cmd := defaultInstallCommand()
			return cmd.run(ctx, []string{url}, console, logger, styled)
		},
	}
}

var setupCommandFactory = defaultSetupCommand

func runSetupCommand(args []string, console *ui.UI, logger *logging.Logger, styled bool) int {
	cmd := setupCommandFactory()
	return cmd.run(context.Background(), args, console, logger, styled)
}

func (c setupCommand) run(ctx context.Context, args []string, console *ui.UI, logger *logging.Logger, styled bool) int {
	console.PrintHeader("dade", version.Version)
	checkOnly := false
	for _, arg := range args {
		switch arg {
		case "--check":
			checkOnly = true
		default:
			logger.Error(fmt.Sprintf("Unknown option: %s", arg))
			return 1
		}
	}

	if c.runner == nil {
		c.runner = execx.NewSystemRunner()
	}
	if c.initConfig == nil {
		c.initConfig = config.InitConfig
	}
	if c.detectSrv == nil {
		c.detectSrv = config.DetectSrvConfig
	}
	if c.migrateSrv == nil {
		c.migrateSrv = migrateFromSrv
	}
	if c.generateCaddy == nil {
		c.generateCaddy = proxy.GenerateCaddyfile
	}
	if c.createPlist == nil {
		c.createPlist = proxy.CreatePlist
	}
	if c.installProxy == nil {
		c.installProxy = proxy.InstallProxyService
	}
	if c.trustCA == nil {
		c.trustCA = trustCaddyCA
	}
	if c.templatesDir == nil {
		c.templatesDir = config.TemplatesDir
	}
	if c.projectsFile == nil {
		c.projectsFile = config.ProjectsFile
	}
	if c.caddyfilePath == nil {
		c.caddyfilePath = config.CaddyfilePath
	}
	if c.plistPath == nil {
		c.plistPath = config.ProxyPlistPath
	}
	if c.logPath == nil {
		c.logPath = config.ProxyLogPath
	}
	if c.errPath == nil {
		c.errPath = config.ProxyErrPath
	}
	if c.installTemplate == nil {
		c.installTemplate = func(ctx context.Context, url string, console *ui.UI, logger *logging.Logger, styled bool) int {
			cmd := defaultInstallCommand()
			return cmd.run(ctx, []string{url}, console, logger, styled)
		}
	}
	if c.spin == nil {
		spinnerEnabled := term.IsTerminal(int(os.Stdout.Fd()))
		spinner := ui.NewSpinner(os.Stdout, spinnerEnabled)
		c.spin = spinner.Run
	}
	if c.confirm == nil {
		interactive := term.IsTerminal(int(os.Stdin.Fd()))
		c.confirm = func(question string) (bool, error) {
			if !interactive {
				return false, nil
			}
			confirmed := false
			prompt := huh.NewConfirm().Title(question).Value(&confirmed)
			if err := prompt.Run(); err != nil {
				return false, err
			}
			return confirmed, nil
		}
	}

	logger.Info("Setting up dade...")
	if !c.checkAllDependencies(ctx, logger) {
		logger.Error("Missing required dependencies")
		return 1
	}
	if checkOnly {
		logger.Success("All dependencies OK")
		return 0
	}

	if _, err := c.initConfig(); err != nil {
		logger.Error("Failed to initialize config")
		return 1
	}
	logger.Success("Configuration initialized")

	srvDetected, err := c.detectSrv()
	if err != nil {
		logger.Error("Failed to detect srv installation")
		return 1
	}
	if srvDetected {
		logger.Info("Existing srv installation detected")
		migrate, err := c.confirm("Migrate from srv?")
		if err != nil {
			logger.Error("Failed to read migration confirmation")
			return 1
		}
		if migrate {
			projectsPath, err := c.projectsFile()
			if err != nil {
				logger.Error("Failed to resolve projects registry")
				return 1
			}
			if err := c.migrateSrv(ctx, c.runner, projectsPath, logger); err != nil {
				logger.Error("Failed to migrate from srv")
				return 1
			}
		}
	}

	projectsPath, err := c.projectsFile()
	if err != nil {
		logger.Error("Failed to resolve projects registry")
		return 1
	}
	caddyfilePath, err := c.caddyfilePath()
	if err != nil {
		logger.Error("Failed to resolve Caddyfile path")
		return 1
	}
	if err := c.generateCaddy(ctx, c.runner, projectsPath, caddyfilePath); err != nil {
		logger.Error("Failed to generate Caddyfile")
		return 1
	}
	logger.Success("Caddyfile generated")

	plistPath, err := c.plistPath()
	if err != nil {
		logger.Error("Failed to resolve proxy plist path")
		return 1
	}
	logPath, err := c.logPath()
	if err != nil {
		logger.Error("Failed to resolve proxy log path")
		return 1
	}
	errPath, err := c.errPath()
	if err != nil {
		logger.Error("Failed to resolve proxy error path")
		return 1
	}
	if err := c.createPlist(ctx, c.runner, plistPath, caddyfilePath, logPath, errPath); err != nil {
		logger.Error("Failed to create proxy plist")
		return 1
	}
	if err := c.installProxy(ctx, c.runner, plistPath); err != nil {
		logger.Error("Failed to start proxy service")
		return 1
	}
	logger.Success("Proxy service started")

	trust, err := c.confirm("Trust Caddy CA? (requires sudo)")
	if err != nil {
		logger.Error("Failed to read CA trust confirmation")
		return 1
	}
	if trust {
		if err := c.trustCA(ctx, c.runner); err != nil {
			logger.Error("Failed to trust Caddy CA")
			return 1
		}
		logger.Success("Caddy CA trusted")
	}

	installTemplates, err := c.confirm("Install official templates?")
	if err != nil {
		logger.Error("Failed to read templates confirmation")
		return 1
	}
	if installTemplates {
		if err := ensureDefaultTemplates(ctx, logger); err != nil {
			logger.Error(fmt.Sprintf("Failed to install default templates: %v", err))
			return 1
		}
	}

	logger.Success("Setup complete!")
	logger.Info("Create your first project: dade new myproject")
	return 0
}

func (c setupCommand) checkAllDependencies(ctx context.Context, logger *logging.Logger) bool {
	allOK := true
	if !c.checkDependency(ctx, logger, "jq", "jq", true) {
		allOK = false
	}
	if !c.checkDependency(ctx, logger, "caddy", "caddy", true) {
		allOK = false
	}
	if !c.checkDependency(ctx, logger, "tailwindcss", "tailwindcss", true) {
		allOK = false
	}
	c.checkDependency(ctx, logger, "cloudflared", "cloudflared", false)
	c.checkDependency(ctx, logger, "tk", "wedow/tap/ticket", false)
	return allOK
}

func (c setupCommand) checkDependency(ctx context.Context, logger *logging.Logger, name, brewName string, required bool) bool {
	if execx.CommandAvailable(c.runner, name) {
		logger.Success(fmt.Sprintf("%s installed", name))
		return true
	}
	if !required {
		logger.Info(fmt.Sprintf("%s not found (optional)", name))
		return true
	}

	logger.Warn(fmt.Sprintf("%s not found", name))
	install, err := c.confirm(fmt.Sprintf("Install %s via Homebrew?", name))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to confirm %s install", name))
		return false
	}
	if !install {
		logger.Error(fmt.Sprintf("Could not install %s", name))
		return false
	}
	if !execx.CommandAvailable(c.runner, "brew") {
		logger.Error("Homebrew not found")
		return false
	}
	if brewName == "" {
		brewName = name
	}
	work := func() error {
		return c.runner.Run(ctx, "brew", "install", brewName)
	}
	if err := c.spin(fmt.Sprintf("Installing %s", name), work); err != nil {
		logger.Error(fmt.Sprintf("Could not install %s", name))
		return false
	}
	if execx.CommandAvailable(c.runner, name) {
		logger.Success(fmt.Sprintf("%s installed", name))
		return true
	}
	logger.Error(fmt.Sprintf("Could not install %s", name))
	return false
}

func (c setupCommand) offerOfficialTemplates(ctx context.Context, console *ui.UI, logger *logging.Logger, styled bool) error {
	templatesDir, err := c.templatesDir()
	if err != nil {
		return err
	}
	for _, tpl := range config.DefaultTemplates().Ordered {
		target := filepath.Join(templatesDir, tpl.Name)
		if info, err := os.Stat(target); err == nil && info.IsDir() {
			logger.Info(fmt.Sprintf("%s already installed", tpl.Name))
			continue
		}
		description := tpl.DisplayName
		if description == "" {
			description = tpl.Name
		}
		install, err := c.confirm(fmt.Sprintf("Install %s? (%s)", tpl.Name, description))
		if err != nil {
			return err
		}
		if !install {
			continue
		}
		if code := c.installTemplate(ctx, tpl.URL, console, logger, styled); code != 0 {
			return errors.New("template install failed")
		}
	}
	return nil
}

func trustCaddyCA(ctx context.Context, runner execx.Runner) error {
	if runner == nil {
		runner = execx.NewSystemRunner()
	}
	return runner.Run(ctx, "sudo", "caddy", "trust")
}

type srvProject struct {
	Port int    `json:"port"`
	Path string `json:"path"`
}

func migrateFromSrv(ctx context.Context, runner execx.Runner, projectsPath string, logger *logging.Logger) error {
	if runner == nil {
		runner = execx.NewSystemRunner()
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	srvConfigDir := filepath.Join(home, ".config", "srv")
	srvProjects := filepath.Join(srvConfigDir, "projects.json")
	srvPlist := filepath.Join(home, "Library", "LaunchAgents", "land.charm.srv.proxy.plist")

	if output, err := runner.Output(ctx, "launchctl", "list"); err == nil && strings.Contains(output, "land.charm.srv.proxy") {
		uid := fmt.Sprintf("%d", os.Getuid())
		_ = runner.Run(ctx, "launchctl", "bootout", "gui/"+uid, srvPlist)
		logger.Success("Stopped srv proxy")
	}

	data, err := os.ReadFile(srvProjects)
	if err == nil {
		var projects map[string]srvProject
		if err := json.Unmarshal(data, &projects); err != nil {
			return err
		}
		count := 0
		for name, project := range projects {
			if name == "" || project.Path == "" || project.Port <= 0 {
				continue
			}
			info, err := os.Stat(project.Path)
			if err != nil || !info.IsDir() {
				continue
			}
			if _, migrated, err := registry.MigrateSrvMarker(project.Path); err != nil {
				return err
			} else if !migrated {
				if _, err := registry.WriteMarker(project.Path, name, "unknown", project.Port); err != nil {
					return err
				}
			}
			if _, err := registry.Register(projectsPath, name, project.Port, project.Path, "unknown"); err != nil {
				return err
			}
			logger.Success(fmt.Sprintf("Migrated: %s", name))
			count++
		}
		if count > 0 {
			logger.Success(fmt.Sprintf("Migrated %d project(s)", count))
		}
	}

	if info, err := os.Stat(srvConfigDir); err == nil && info.IsDir() {
		backup := fmt.Sprintf("%s.backup.%s", srvConfigDir, time.Now().UTC().Format("20060102150405"))
		if err := os.Rename(srvConfigDir, backup); err != nil {
			return err
		}
		logger.Info(fmt.Sprintf("Old config backed up to %s", backup))
	}
	_ = os.Remove(srvPlist)
	logger.Success("Migration complete")
	return nil
}
