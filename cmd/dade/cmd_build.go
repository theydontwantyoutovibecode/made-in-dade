package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/theydontwantyoutovibecode/dade/internal/config"
	"github.com/theydontwantyoutovibecode/dade/internal/logging"
	"github.com/theydontwantyoutovibecode/dade/internal/manifest"
	"github.com/theydontwantyoutovibecode/dade/internal/registry"
	"github.com/theydontwantyoutovibecode/dade/internal/ui"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [name]",
	Short: "Build a compiled project",
	Long: `Build a compiled project (Go, Swift, Kotlin, etc.) to produce
an executable binary or installable package.

If no name is provided, builds the project in the current directory.
The build command auto-detects the project type from directory contents
or reads the [build] section from the template manifest.`,
	Example: `dade build              # Build current project
dade build myapp        # Build a registered project
dade build --release    # Build with release optimizations
dade build --os linux   # Cross-compile for Linux`,
	GroupID: "dev",
	Args:    cobra.MaximumNArgs(1),
	RunE:    runBuildCmd,
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringP("output", "o", "", "Output directory (default: ./bin)")
	buildCmd.Flags().String("os", "", "Target OS (darwin, linux)")
	buildCmd.Flags().String("arch", "", "Target architecture (amd64, arm64)")
	buildCmd.Flags().Bool("all", false, "Build for all configured targets")
	buildCmd.Flags().Bool("release", false, "Optimized release build")
}

type buildCommand struct {
	templatesDir func() (string, error)
	projectsFile func() (string, error)
	readMarker   func(string) (registry.Marker, error)
	readFile     func(string) ([]byte, error)
	runCmd       func(ctx context.Context, dir, name string, args []string, env []string, stdout, stderr *strings.Builder) error
	runShell     func(ctx context.Context, dir, cmdStr string, extraEnv []string) error
	fileExists   func(string) bool
	globMatch    func(string) ([]string, error)
}

var buildCommandFactory = defaultBuildCommand

func defaultBuildCommand() buildCommand {
	return buildCommand{
		templatesDir: config.TemplatesDir,
		projectsFile: config.ProjectsFile,
		readMarker:   registry.ReadMarker,
		readFile:     os.ReadFile,
		runCmd:       defaultRunCmd,
		runShell:     defaultRunShellCmd,
		fileExists:   defaultFileExists,
		globMatch:    filepath.Glob,
	}
}

func defaultRunCmd(ctx context.Context, dir, name string, args []string, env []string, stdout, stderr *strings.Builder) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func defaultFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func runBuildCmd(cmd *cobra.Command, args []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	outputDir, _ := cmd.Flags().GetString("output")
	targetOS, _ := cmd.Flags().GetString("os")
	targetArch, _ := cmd.Flags().GetString("arch")
	buildAll, _ := cmd.Flags().GetBool("all")
	release, _ := cmd.Flags().GetBool("release")

	impl := buildCommandFactory()
	code := impl.run(context.Background(), args, console, logger, buildOptions{
		outputDir:  outputDir,
		targetOS:   targetOS,
		targetArch: targetArch,
		all:        buildAll,
		release:    release,
	})
	if code != 0 {
		return errors.New("build failed")
	}
	return nil
}

type buildOptions struct {
	outputDir  string
	targetOS   string
	targetArch string
	all        bool
	release    bool
}

func (c buildCommand) run(ctx context.Context, args []string, console *ui.UI, logger *logging.Logger, opts buildOptions) int {
	var projectDir string

	if len(args) > 0 {
		projectsPath, err := c.projectsFile()
		if err != nil {
			logger.Error("Failed to resolve projects file")
			return 1
		}
		project, ok, err := registry.Get(projectsPath, args[0])
		if err != nil {
			logger.Error("Failed to load project registry")
			return 1
		}
		if !ok {
			logger.Error(fmt.Sprintf("Project '%s' not found", args[0]))
			return 1
		}
		projectDir = project.Path
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			logger.Error("Failed to get current directory")
			return 1
		}
		projectDir = cwd
	}

	_ = console

	var mf manifest.Manifest
	if registry.MarkerExists(projectDir) {
		marker, err := c.readMarker(projectDir)
		if err == nil {
			templatesDir, err := c.templatesDir()
			if err == nil {
				manifestPath := filepath.Join(templatesDir, marker.Template, "dade.toml")
				if data, err := c.readFile(manifestPath); err == nil {
					if parsed, err := manifest.Parse(data); err == nil {
						mf = parsed
					}
				}
			}
		}
	}

	localManifest := filepath.Join(projectDir, "dade.toml")
	if data, err := c.readFile(localManifest); err == nil {
		if parsed, err := manifest.Parse(data); err == nil {
			if manifest.HasBuildSection(parsed) {
				mf = parsed
			}
		}
	}

	if manifest.HasBuildSection(mf) {
		return c.buildFromManifest(ctx, projectDir, mf, logger, opts)
	}

	return c.buildAutoDetect(ctx, projectDir, logger, opts)
}

func (c buildCommand) buildFromManifest(ctx context.Context, projectDir string, mf manifest.Manifest, logger *logging.Logger, opts buildOptions) int {
	outputDir := manifest.BuildOutput(mf)
	if opts.outputDir != "" {
		outputDir = opts.outputDir
	}
	absOutput := outputDir
	if !filepath.IsAbs(absOutput) {
		absOutput = filepath.Join(projectDir, outputDir)
	}
	if err := os.MkdirAll(absOutput, 0755); err != nil {
		logger.Error(fmt.Sprintf("Failed to create output directory: %v", err))
		return 1
	}

	preCmds := manifest.BuildPreCommands(mf)
	if len(preCmds) > 0 {
		logger.Info("Running pre-build commands...")
		for _, cmd := range preCmds {
			if err := c.runShell(ctx, projectDir, cmd, nil); err != nil {
				logger.Error(fmt.Sprintf("Pre-build command failed: %s: %v", cmd, err))
				return 1
			}
		}
	}

	if opts.all && len(manifest.BuildTargets(mf)) > 0 {
		targets := manifest.BuildTargets(mf)
		for _, target := range targets {
			if code := c.runBuildCmd(ctx, projectDir, mf, logger, opts, target.OS, target.Arch); code != 0 {
				return code
			}
		}
	} else {
		targetOS := opts.targetOS
		targetArch := opts.targetArch
		if targetOS == "" {
			targetOS = runtime.GOOS
		}
		if targetArch == "" {
			targetArch = runtime.GOARCH
		}
		if code := c.runBuildCmd(ctx, projectDir, mf, logger, opts, targetOS, targetArch); code != 0 {
			return code
		}
	}

	postCmds := manifest.BuildPostCommands(mf)
	if len(postCmds) > 0 {
		logger.Info("Running post-build commands...")
		for _, cmd := range postCmds {
			if err := c.runShell(ctx, projectDir, cmd, nil); err != nil {
				logger.Error(fmt.Sprintf("Post-build command failed: %s: %v", cmd, err))
				return 1
			}
		}
	}

	return 0
}

func (c buildCommand) runBuildCmd(ctx context.Context, projectDir string, mf manifest.Manifest, logger *logging.Logger, opts buildOptions, targetOS, targetArch string) int {
	buildCmdStr := manifest.BuildCommand(mf)
	if opts.release && manifest.BuildReleaseFlags(mf) != "" {
		buildCmdStr = buildCmdStr + " " + manifest.BuildReleaseFlags(mf)
	}

	name := filepath.Base(projectDir)

	outputDir := manifest.BuildOutput(mf)
	if opts.outputDir != "" {
		outputDir = opts.outputDir
	}

	buildCmdStr = strings.ReplaceAll(buildCmdStr, "{{name}}", name)
	buildCmdStr = strings.ReplaceAll(buildCmdStr, "{{output}}", outputDir)
	buildCmdStr = strings.ReplaceAll(buildCmdStr, "{{os}}", targetOS)
	buildCmdStr = strings.ReplaceAll(buildCmdStr, "{{arch}}", targetArch)

	env := []string{
		"GOOS=" + targetOS,
		"GOARCH=" + targetArch,
	}

	logger.Info(fmt.Sprintf("Building for %s/%s...", targetOS, targetArch))
	if err := c.runShell(ctx, projectDir, buildCmdStr, env); err != nil {
		logger.Error(fmt.Sprintf("Build failed: %v", err))
		return 1
	}

	logger.Success(fmt.Sprintf("Built for %s/%s", targetOS, targetArch))
	return 0
}

func (c buildCommand) buildAutoDetect(ctx context.Context, projectDir string, logger *logging.Logger, opts buildOptions) int {
	if c.fileExists(filepath.Join(projectDir, "go.mod")) {
		return c.buildGo(ctx, projectDir, logger, opts)
	}

	if matches, _ := c.globMatch(filepath.Join(projectDir, "*.xcodeproj")); len(matches) > 0 {
		return c.buildXcode(ctx, projectDir, logger, opts, matches[0])
	}

	if c.fileExists(filepath.Join(projectDir, "build.gradle.kts")) || c.fileExists(filepath.Join(projectDir, "build.gradle")) {
		return c.buildGradle(ctx, projectDir, logger, opts)
	}

	logger.Error("Could not detect project type")
	logger.Info("Add a [build] section to your dade.toml or ensure the project has a recognized structure (go.mod, *.xcodeproj, build.gradle.kts)")
	return 1
}

func (c buildCommand) buildGo(ctx context.Context, projectDir string, logger *logging.Logger, opts buildOptions) int {
	targetOS := opts.targetOS
	targetArch := opts.targetArch
	if targetOS == "" {
		targetOS = runtime.GOOS
	}
	if targetArch == "" {
		targetArch = runtime.GOARCH
	}

	outputDir := opts.outputDir
	if outputDir == "" {
		outputDir = "./bin"
	}
	absOutput := outputDir
	if !filepath.IsAbs(absOutput) {
		absOutput = filepath.Join(projectDir, outputDir)
	}
	if err := os.MkdirAll(absOutput, 0755); err != nil {
		logger.Error(fmt.Sprintf("Failed to create output directory: %v", err))
		return 1
	}

	name := filepath.Base(projectDir)
	outPath := filepath.Join(outputDir, name)
	args := []string{"build", "-o", outPath}
	if opts.release {
		args = append(args, "-ldflags", "-s -w")
	}
	args = append(args, ".")

	env := []string{
		"GOOS=" + targetOS,
		"GOARCH=" + targetArch,
	}

	logger.Info(fmt.Sprintf("Building %s for %s/%s...", name, targetOS, targetArch))

	var stdout, stderr strings.Builder
	if err := c.runCmd(ctx, projectDir, "go", args, env, &stdout, &stderr); err != nil {
		logger.Error(fmt.Sprintf("Build failed: %v", err))
		if stderr.Len() > 0 {
			logger.Error(stderr.String())
		}
		return 1
	}

	logger.Success(fmt.Sprintf("Built: %s", outPath))
	return 0
}

func (c buildCommand) buildXcode(ctx context.Context, projectDir string, logger *logging.Logger, opts buildOptions, xcodeproj string) int {
	scheme := strings.TrimSuffix(filepath.Base(xcodeproj), ".xcodeproj")
	configuration := "Debug"
	if opts.release {
		configuration = "Release"
	}

	args := []string{
		"-project", xcodeproj,
		"-scheme", scheme,
		"-configuration", configuration,
		"build",
	}

	logger.Info(fmt.Sprintf("Building %s (%s)...", scheme, configuration))

	var stdout, stderr strings.Builder
	if err := c.runCmd(ctx, projectDir, "xcodebuild", args, nil, &stdout, &stderr); err != nil {
		logger.Error(fmt.Sprintf("xcodebuild failed: %v", err))
		return 1
	}

	logger.Success(fmt.Sprintf("Built %s", scheme))
	return 0
}

func (c buildCommand) buildGradle(ctx context.Context, projectDir string, logger *logging.Logger, opts buildOptions) int {
	task := "assembleDebug"
	if opts.release {
		task = "assembleRelease"
	}

	gradlew := filepath.Join(projectDir, "gradlew")
	cmd := gradlew
	if !c.fileExists(cmd) {
		cmd = "gradle"
	}

	logger.Info(fmt.Sprintf("Building with Gradle (%s)...", task))

	var stdout, stderr strings.Builder
	if err := c.runCmd(ctx, projectDir, cmd, []string{task}, nil, &stdout, &stderr); err != nil {
		logger.Error(fmt.Sprintf("Gradle build failed: %v", err))
		return 1
	}

	logger.Success("Gradle build complete")
	return 0
}

func defaultRunShellCmd(ctx context.Context, dir, cmdStr string, extraEnv []string) error {
	cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), extraEnv...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
