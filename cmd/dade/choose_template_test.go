package main

import (
	"strings"
	"testing"

	"github.com/theydontwantyoutovibecode/dade/internal/manifest"
)

func TestChooseTemplateNonInteractive(t *testing.T) {
	_, err := choosePluginTemplate(samplePluginTemplates(), false)
	if err == nil {
		t.Fatalf("expected error for non-interactive selection")
	}
}

func TestPromptTemplateNumberMapsChoice(t *testing.T) {
	templates := samplePluginTemplates()
	input := strings.NewReader("2\n")
	output := &strings.Builder{}
	name, err := promptTemplateNumberPlugin(templates, input, output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != templates[1].Name {
		t.Fatalf("expected %s, got %s", templates[1].Name, name)
	}
}

func samplePluginTemplates() []pluginTemplate {
	return []pluginTemplate{
		{
			Name: "alpha",
			Manifest: manifest.Manifest{
				Template: manifest.Template{Description: "Alpha template"},
			},
		},
		{
			Name: "bravo",
			Manifest: manifest.Manifest{
				Template: manifest.Template{Description: "Bravo template"},
			},
		},
	}
}
