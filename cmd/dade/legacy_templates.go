package main

import (
	"github.com/theydontwantyoutovibecode/dade/internal/config"
	"github.com/theydontwantyoutovibecode/dade/internal/logging"
	"github.com/theydontwantyoutovibecode/dade/internal/ui"
)

func runTemplates(args []string, console *ui.UI, logger *logging.Logger, styled bool) int {
	jsonOutput := false
	for _, arg := range args {
		if arg == "--json" {
			jsonOutput = true
			break
		}
	}

	templatesDir, err := config.TemplatesDir()
	if err != nil {
		logger.Error("Failed to resolve templates directory")
		return 1
	}
	installed, err := loadInstalledTemplates(templatesDir)
	if err != nil {
		logger.Error("Failed to load installed templates")
		return 1
	}

	if len(installed) == 0 {
		logger.Info("No templates installed.")
		logger.Info("Install with: dade install <git-url>")
		logger.Info("Or see official: dade install --list-official")
		return 0
	}

	if jsonOutput {
		payload, err := templatesJSON(installed)
		if err != nil {
			logger.Error("Failed to render JSON")
			return 1
		}
		console.PrintHelp(payload)
		return 0
	}

	console.PrintHelp(templatesText(installed, styled))
	return 0
}
