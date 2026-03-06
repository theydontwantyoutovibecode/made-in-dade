package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/theydontwantyoutovibecode/dade/internal/config"
	"github.com/theydontwantyoutovibecode/dade/internal/logging"
	"github.com/theydontwantyoutovibecode/dade/internal/registry"
	"github.com/theydontwantyoutovibecode/dade/internal/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type listProject struct {
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Path     string `json:"path"`
	Template string `json:"template"`
	URL      string `json:"url"`
	Running  bool   `json:"running"`
}

func runListCmd(cmd *cobra.Command, _ []string) error {
	output := getOutputSettings(cmd)
	console := ui.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger := logging.New(cmd.OutOrStdout(), cmd.ErrOrStderr(), output.Styled)
	logger.SetSilent(output.Quiet)
	logger.SetVerbose(output.Verbose)

	runningOnly, _ := cmd.Flags().GetBool("running")

	projectsPath, err := config.ProjectsFile()
	if err != nil {
		logger.Error("Failed to resolve projects file")
		return errors.New("list command failed")
	}

	entries, err := registry.List(projectsPath)
	if err != nil {
		logger.Error("Failed to load projects")
		return errors.New("list command failed")
	}

	projects := make([]listProject, 0, len(entries))
	for _, entry := range entries {
		running := isPortInUse(entry.Project.Port)
		if runningOnly && !running {
			continue
		}
		projects = append(projects, listProject{
			Name:     entry.Name,
			Port:     entry.Project.Port,
			Path:     entry.Project.Path,
			Template: entry.Project.Template,
			URL:      fmt.Sprintf("https://%s", config.ProjectDomain(entry.Name)),
			Running:  running,
		})
	}

	if output.JSON {
		data, err := json.MarshalIndent(projects, "", "  ")
		if err != nil {
			logger.Error("Failed to encode JSON")
			return errors.New("list command failed")
		}
		console.PrintHelp(string(data))
		return nil
	}

	if len(projects) == 0 {
		if runningOnly {
			logger.Info("No running projects.")
		} else {
			logger.Info("No projects registered.")
			logger.Info("Create one: dade new <name>")
		}
		return nil
	}

	console.PrintHelp(listText(projects, output.Styled))
	return nil
}

func isPortInUse(port int) bool {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	conn, err := net.DialTimeout("tcp", addr, 100*1000*1000)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func listText(projects []listProject, styled bool) string {
	var lines []string
	lines = append(lines, "")

	runningCount := 0
	for _, p := range projects {
		if p.Running {
			runningCount++
		}
	}

	if styled {
		header := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
		nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true)
		urlStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("159"))
		metaStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		runningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
		stoppedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

		lines = append(lines, header.Render("Projects"), "")
		for _, p := range projects {
			statusIcon := "○"
			statusText := "stopped"
			statusStyled := stoppedStyle.Render(statusText)
			if p.Running {
				statusIcon = "●"
				statusText = "running"
				statusStyled = runningStyle.Render(statusText)
			}
			lines = append(lines,
				fmt.Sprintf("  %s %s", runningStyle.Render(statusIcon), nameStyle.Render(p.Name)),
				urlStyle.Render(fmt.Sprintf("    %s", p.URL)),
				metaStyle.Render(fmt.Sprintf("    Template: %s | Port: %d | %s", p.Template, p.Port, statusStyled)),
				metaStyle.Render(fmt.Sprintf("    Path: %s", p.Path)),
				"",
			)
		}
		lines = append(lines, metaStyle.Render(fmt.Sprintf("%d project(s) (%d running)", len(projects), runningCount)))
		return strings.Join(lines, "\n")
	}

	lines = append(lines, "Projects", "")
	for _, p := range projects {
		statusIcon := "○"
		statusText := "stopped"
		if p.Running {
			statusIcon = "●"
			statusText = "running"
		}
		lines = append(lines,
			fmt.Sprintf("  %s %s", statusIcon, p.Name),
			fmt.Sprintf("    %s", p.URL),
			fmt.Sprintf("    Template: %s | Port: %d | %s", p.Template, p.Port, statusText),
			fmt.Sprintf("    Path: %s", p.Path),
			"",
		)
	}
	lines = append(lines, fmt.Sprintf("%d project(s) (%d running)", len(projects), runningCount))
	return strings.Join(lines, "\n")
}
