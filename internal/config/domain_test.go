package config

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestNormalizeDomainTLD(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{".localhost", ".localhost"},
		{"localhost", ".localhost"},
		{".local", ".local"},
		{"local", ".local"},
		{"", ".localhost"},
		{".test", ".test"},
		{"test", ".test"},
		{"  .example  ", ".example"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeDomainTLD(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeDomainTLD(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDomainTLD_Default(t *testing.T) {
	// Clear environment variable
	os.Unsetenv("DADE_DOMAIN_TLD")

	// Reset cached value
	domainTLD = ""
	domainTLDOnce = sync.Once{}

	// Test default behavior
	// Note: If an existing dade installation is detected (projects.json exists),
	// it will use .local for backward compatibility. Otherwise, it uses .localhost.
	result := DomainTLD()
	if result != defaultDomainTLD && result != legacyDomainTLD {
		t.Errorf("DomainTLD() = %q, want either %q (new) or %q (legacy)", result, defaultDomainTLD, legacyDomainTLD)
	}
}

func TestDomainTLD_FromEnv(t *testing.T) {
	// Set environment variable
	os.Setenv("DADE_DOMAIN_TLD", ".test")
	defer os.Unsetenv("DADE_DOMAIN_TLD")

	// Reset cached value
	domainTLD = ""
	domainTLDOnce = sync.Once{}

	result := DomainTLD()
	if result != ".test" {
		t.Errorf("DomainTLD() with env = %q, want %q", result, ".test")
	}
}

func TestDomainTLD_FromConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Create config file with custom TLD
	configPath := filepath.Join(tmpDir, configFileName)
	configContent := `# dade configuration
domain_tld = ".test"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create projects.json to make it look like an existing install
	projectsPath := filepath.Join(tmpDir, ProjectsFileName)
	if err := os.WriteFile(projectsPath, []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}

	// Cannot override ConfigDir function, so skip this test
	// In real code, we'd use dependency injection
	t.Skip("Cannot override ConfigDir for testing")
}

func TestLocalDomain(t *testing.T) {
	// Set environment variable
	os.Setenv("DADE_DOMAIN_TLD", ".test")
	defer os.Unsetenv("DADE_DOMAIN_TLD")

	// Reset cached value
	domainTLD = ""
	domainTLDOnce = sync.Once{}

	result := LocalDomain()
	if result == "" {
		t.Error("LocalDomain() returned empty string")
	}

	if !strings.Contains(result, ".test") {
		t.Errorf("LocalDomain() = %q, should contain TLD", result)
	}
}

func TestProjectDomain(t *testing.T) {
	// Set environment variable
	os.Setenv("DADE_DOMAIN_TLD", ".test")
	defer os.Unsetenv("DADE_DOMAIN_TLD")

	// Reset cached value
	domainTLD = ""
	domainTLDOnce = sync.Once{}

	projectName := "myapp"
	result := ProjectDomain(projectName)

	if result == "" {
		t.Error("ProjectDomain() returned empty string")
	}

	if !strings.Contains(result, projectName) {
		t.Errorf("ProjectDomain(%q) = %q, should contain project name", projectName, result)
	}

	if !strings.Contains(result, ".test") {
		t.Errorf("ProjectDomain(%q) = %q, should contain TLD", projectName, result)
	}
}

func TestSetDomainTLD(t *testing.T) {
	// Cannot override ConfigDir function, so skip this test
	// In real code, we'd use dependency injection
	t.Skip("Cannot override ConfigDir for testing")
}
