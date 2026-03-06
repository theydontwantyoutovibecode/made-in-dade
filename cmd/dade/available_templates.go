package main

import (
	"fmt"
	"strings"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/config"
	"github.com/charmbracelet/lipgloss"
)

func availableTemplatesText(templates config.Templates, configPath string, styled bool) string {
	lines := []string{""}
	if styled {
		header := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
		nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
		descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		urlStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("159"))

		lines = append(lines, header.Render("Available Templates"), "")
		for _, tpl := range templates.Ordered {
			display := tpl.DisplayName
			if display == "" {
				display = tpl.Name
			}
			lines = append(lines,
				nameStyle.Render(fmt.Sprintf("  %s", tpl.Name)),
				descStyle.Render(fmt.Sprintf("    %s", display)),
				urlStyle.Render(fmt.Sprintf("    URL: %s", tpl.URL)),
				"",
			)
		}
		lines = append(lines,
			"",
			descStyle.Render(fmt.Sprintf("To add custom templates, create %s:", configPath)),
			"",
			urlStyle.Render("  [templates]"),
			urlStyle.Render("  my-template = \"https://github.com/user/repo.git\""),
			"",
		)
		return strings.Join(lines, "\n")
	}

	lines = append(lines, "Available Templates", "")
	for _, tpl := range templates.Ordered {
		display := tpl.DisplayName
		if display == "" {
			display = tpl.Name
		}
		lines = append(lines,
			fmt.Sprintf("  %s", tpl.Name),
			fmt.Sprintf("    %s", display),
			fmt.Sprintf("    URL: %s", tpl.URL),
			"",
		)
	}
	lines = append(lines,
		"",
		fmt.Sprintf("To add custom templates, create %s:", configPath),
		"",
		"  [templates]",
		"  my-template = \"https://github.com/user/repo.git\"",
		"",
	)
	return strings.Join(lines, "\n")
}
