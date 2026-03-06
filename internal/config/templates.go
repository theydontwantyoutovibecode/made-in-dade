package config

import (
	"errors"
	"os"
	"strings"
)

type Template struct {
	Name        string
	URL         string
	DisplayName string
}

type Templates struct {
	Ordered []Template
	ByName  map[string]Template
}

func DefaultTemplates() Templates {
	ordered := []Template{
		{
			Name:        "web-app",
			URL:         "https://github.com/theydontwantyoutovibecode/web-app-made-in-dade.git",
			DisplayName: "Django + Hypermedia (HTMX, TailwindCSS)",
		},
		{
			Name:        "web-site",
			URL:         "https://github.com/theydontwantyoutovibecode/web-site-made-in-dade.git",
			DisplayName: "HTML + Hypertext (Vanilla HTML/CSS/JS with HTMX)",
		},
		{
			Name:        "ios-app",
			URL:         "https://github.com/theydontwantyoutovibecode/ios-made-in-dade.git",
			DisplayName: "iOS App (Swift + SwiftUI)",
		},
		{
			Name:        "mac-app",
			URL:         "https://github.com/theydontwantyoutovibecode/mac-made-in-dade.git",
			DisplayName: "macOS App (Swift + SwiftUI)",
		},
		{
			Name:        "android-app",
			URL:         "https://github.com/theydontwantyoutovibecode/android-made-in-dade.git",
			DisplayName: "Android App (Kotlin + Jetpack Compose)",
		},
		{
			Name:        "cli",
			URL:         "https://github.com/theydontwantyoutovibecode/cli-made-in-dade.git",
			DisplayName: "CLI App (Go + Charm: Fang, Cobra, Huh, Lipgloss)",
		},
		{
			Name:        "tui",
			URL:         "https://github.com/theydontwantyoutovibecode/tui-made-in-dade.git",
			DisplayName: "TUI App (Go + Charm: Bubbletea v2, Bubbles, Lipgloss)",
		},
	}

	byName := make(map[string]Template, len(ordered))
	for _, tpl := range ordered {
		byName[tpl.Name] = tpl
	}

	return Templates{Ordered: ordered, ByName: byName}
}

func LoadTemplates(path string) (Templates, error) {
	result := DefaultTemplates()

	if path == "" {
		return result, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return result, nil
		}
		return Templates{}, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || line == "[templates]" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"")
		if key == "" || value == "" {
			continue
		}

		result = upsertTemplate(result, key, value)
	}

	return result, nil
}

func upsertTemplate(templates Templates, name, url string) Templates {
	tpl := Template{Name: name, URL: url}
	if existing, ok := templates.ByName[name]; ok {
		tpl.DisplayName = existing.DisplayName
		for i := range templates.Ordered {
			if templates.Ordered[i].Name == name {
				templates.Ordered[i] = tpl
				break
			}
		}
		templates.ByName[name] = tpl
		return templates
	}

	templates.Ordered = append(templates.Ordered, tpl)
	templates.ByName[name] = tpl
	return templates
}
