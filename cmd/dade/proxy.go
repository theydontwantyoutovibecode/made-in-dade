package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/theydontwantyoutovibecode/dade/internal/config"
	execx "github.com/theydontwantyoutovibecode/dade/internal/exec"
	"github.com/theydontwantyoutovibecode/dade/internal/logging"
	"github.com/theydontwantyoutovibecode/dade/internal/proxy"
	"github.com/theydontwantyoutovibecode/dade/internal/registry"
	"github.com/theydontwantyoutovibecode/dade/internal/ui"
)

type proxyCommand struct {
	runner       execx.Runner
	initConfig   func() (bool, error)
	createPlist  func(context.Context, execx.Runner, string, string, string, string) error
	install      func(context.Context, execx.Runner, string) error
	uninstall    func(context.Context, execx.Runner, string) error
	restart      func(context.Context, execx.Runner, string) error
	isRunning    func(context.Context, execx.Runner) bool
	loadRegistry func(string) (map[string]registry.Project, error)
	caddyfilePath func() (string, error)
	plistPath    func() (string, error)
	logPath      func() (string, error)
	errPath      func() (string, error)
	projectsFile func() (string, error)
	tail         func(string) error
}

func defaultProxyCommand() proxyCommand {
	return proxyCommand{
		runner:        execx.NewSystemRunner(),
		initConfig:    config.InitConfig,
		createPlist:   proxy.CreatePlist,
		install:       proxy.InstallProxyService,
		uninstall:     proxy.UninstallProxyService,
		restart:       proxy.RestartProxyService,
		isRunning:     proxy.IsProxyRunning,
		loadRegistry:  registry.Load,
		caddyfilePath: config.CaddyfilePath,
		plistPath:     config.ProxyPlistPath,
		logPath:       config.ProxyLogPath,
		errPath:       config.ProxyErrPath,
		projectsFile:  config.ProjectsFile,
		tail:          tailFile,
	}
}

var proxyCommandFactory = defaultProxyCommand

func runProxy(args []string, console *ui.UI, logger *logging.Logger) int {
	cmd := proxyCommandFactory()
	return cmd.run(context.Background(), args, console, logger)
}

func (c proxyCommand) run(ctx context.Context, args []string, console *ui.UI, logger *logging.Logger) int {
	if c.initConfig == nil {
		c.initConfig = config.InitConfig
	}
	if _, err := c.initConfig(); err != nil {
		logger.Error("Failed to initialize config")
		return 1
	}

	action := "status"
	if len(args) > 0 {
		action = args[0]
	}
	if len(args) > 1 {
		logger.Error("Too many arguments")
		return 1
	}

	switch action {
	case "start":
		return c.start(ctx, logger)
	case "stop":
		return c.stop(ctx, logger)
	case "restart":
		return c.restartService(ctx, logger)
	case "status":
		return c.status(ctx, logger)
	case "logs":
		return c.logs(logger)
	default:
		logger.Error(fmt.Sprintf("Unknown action: %s", action))
		return 1
	}
}

func (c proxyCommand) start(ctx context.Context, logger *logging.Logger) int {
	if c.isRunning(ctx, c.runner) {
		logger.Warn("Proxy already running")
		return 0
	}

	plistPath, err := c.plistPath()
	if err != nil {
		logger.Error("Failed to resolve proxy plist path")
		return 1
	}
	caddyfilePath, err := c.caddyfilePath()
	if err != nil {
		logger.Error("Failed to resolve Caddyfile path")
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
	if err := c.install(ctx, c.runner, plistPath); err != nil {
		logger.Error("Failed to start proxy")
		return 1
	}
	logger.Success("Proxy started")
	return 0
}

func (c proxyCommand) stop(ctx context.Context, logger *logging.Logger) int {
	if !c.isRunning(ctx, c.runner) {
		logger.Warn("Proxy not running")
		return 0
	}
	plistPath, err := c.plistPath()
	if err != nil {
		logger.Error("Failed to resolve proxy plist path")
		return 1
	}
	if err := c.uninstall(ctx, c.runner, plistPath); err != nil {
		logger.Error("Failed to stop proxy")
		return 1
	}
	logger.Success("Proxy stopped")
	return 0
}

func (c proxyCommand) restartService(ctx context.Context, logger *logging.Logger) int {
	plistPath, err := c.plistPath()
	if err != nil {
		logger.Error("Failed to resolve proxy plist path")
		return 1
	}
	if err := c.restart(ctx, c.runner, plistPath); err != nil {
		logger.Error("Failed to restart proxy")
		return 1
	}
	logger.Success("Proxy restarted")
	return 0
}

func (c proxyCommand) status(ctx context.Context, logger *logging.Logger) int {
	running := c.isRunning(ctx, c.runner)
	if running {
		logger.Success("Proxy running")
	} else {
		logger.Warn("Proxy not running")
	}

	projectsPath, err := c.projectsFile()
	if err != nil {
		logger.Error("Failed to resolve projects registry")
		return 1
	}
	projects, err := c.loadRegistry(projectsPath)
	if err != nil {
		logger.Error("Failed to load project registry")
		return 1
	}
	count, minPort, maxPort, ok := projectStats(projects)
	logger.Info(fmt.Sprintf("Projects: %d", count))
	if ok {
		logger.Info(fmt.Sprintf("Ports: %d-%d", minPort, maxPort))
	} else {
		logger.Info("Ports: none")
	}

	caddyfilePath, err := c.caddyfilePath()
	if err != nil {
		logger.Error("Failed to resolve Caddyfile path")
		return 1
	}
	logger.Info(fmt.Sprintf("Caddyfile: %s", caddyfilePath))
	return 0
}

func (c proxyCommand) logs(logger *logging.Logger) int {
	logPath, err := c.logPath()
	if err != nil {
		logger.Error("Failed to resolve proxy log path")
		return 1
	}
	if c.tail == nil {
		logger.Error("Log streaming unavailable")
		return 1
	}
	if err := c.tail(logPath); err != nil {
		logger.Error("Failed to stream proxy logs")
		return 1
	}
	return 0
}

func projectStats(projects map[string]registry.Project) (int, int, int, bool) {
	count := len(projects)
	minPort := 0
	maxPort := 0
	for _, project := range projects {
		if project.Port <= 0 {
			continue
		}
		if minPort == 0 || project.Port < minPort {
			minPort = project.Port
		}
		if project.Port > maxPort {
			maxPort = project.Port
		}
	}
	if minPort == 0 || maxPort == 0 {
		return count, 0, 0, false
	}
	return count, minPort, maxPort, true
}

func tailFile(path string) error {
	if path == "" {
		return errors.New("path required")
	}
	cmd := exec.Command("tail", "-f", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
