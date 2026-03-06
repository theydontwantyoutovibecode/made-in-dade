package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/theydontwantyoutovibecode/made-in-dade/internal/manifest"
	"github.com/charmbracelet/lipgloss"
)

type installedTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ServeType   string `json:"type"`
	Source      string `json:"source"`
}

func loadInstalledTemplates(templatesDir string) ([]installedTemplate, error) {
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	templates := make([]installedTemplate, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		templatePath := filepath.Join(templatesDir, name)
		manifestPath := filepath.Join(templatePath, "dade.toml")
		sourcePath := filepath.Join(templatePath, ".source")

		description := "unknown"
		serveType := "unknown"
		if data, err := os.ReadFile(manifestPath); err == nil {
			parsed, err := manifest.Parse(data)
			if err == nil {
				if strings.TrimSpace(parsed.Template.Description) != "" {
					description = parsed.Template.Description
				}
				if strings.TrimSpace(parsed.Serve.Type) != "" {
					serveType = parsed.Serve.Type
				}
			}
		}

		source := "unknown"
		if data, err := os.ReadFile(sourcePath); err == nil {
			trimmed := strings.TrimSpace(string(data))
			if trimmed != "" {
				source = trimmed
			}
		}

		templates = append(templates, installedTemplate{
			Name:        name,
			Description: description,
			ServeType:   serveType,
			Source:      source,
		})
	}

	sort.Slice(templates, func(i, j int) bool {
		return templates[i].Name < templates[j].Name
	})

	return templates, nil
}

func templatesText(templates []installedTemplate, styled bool) string {
	lines := []string{""}
	if styled {
		header := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
		nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true)
		descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		metaStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("159"))

		lines = append(lines, header.Render("Installed Templates"), "")
		for _, tpl := range templates {
			lines = append(lines,
				nameStyle.Render(fmt.Sprintf("  %s", tpl.Name)),
				descStyle.Render(fmt.Sprintf("    %s", tpl.Description)),
				metaStyle.Render(fmt.Sprintf("    Type: %s", tpl.ServeType)),
				metaStyle.Render(fmt.Sprintf("    Source: %s", tpl.Source)),
				"",
			)
		}
		lines = append(lines,
			"",
			descStyle.Render("To install more templates:"),
			metaStyle.Render("  dade install <git-url>"),
			metaStyle.Render("  dade install --list-official"),
			"",
		)
		return strings.Join(lines, "\n")
	}

	lines = append(lines, "Installed Templates", "")
	for _, tpl := range templates {
		lines = append(lines,
			fmt.Sprintf("  %s", tpl.Name),
			fmt.Sprintf("    %s", tpl.Description),
			fmt.Sprintf("    Type: %s", tpl.ServeType),
			fmt.Sprintf("    Source: %s", tpl.Source),
			"",
		)
	}
	lines = append(lines,
		"",
		"To install more templates:",
		"  dade install <git-url>",
		"  dade install --list-official",
		"",
	)
	return strings.Join(lines, "\n")
}

func templatesJSON(templates []installedTemplate) (string, error) {
	payload := make([]installedTemplate, 0, len(templates))
	for _, tpl := range templates {
		payload = append(payload, tpl)
	}
	encoded, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
