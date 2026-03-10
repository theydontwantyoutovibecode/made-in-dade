package exec

import (
	"context"
	"os/exec"
	"runtime"
	"testing"
)

func TestCanUseSudo(t *testing.T) {
	if runtime.GOOS == "windows" {
		// Windows doesn't use sudo
		return
	}

	result := CanUseSudo()
	// On most Unix systems, sudo should be available
	// We can't assert this is always true, but we can check it doesn't panic
	if !result {
		// If sudo is not available, that's unusual but not necessarily an error
		t.Logf("sudo not available on this system")
	}
}

func TestIsAdmin(t *testing.T) {
	result := IsAdmin()
	// Most tests won't run as root, so this should be false
	// We just check it doesn't panic
	if result {
		t.Logf("Running as root/admin")
	} else {
		t.Logf("Not running as root/admin")
	}
}

func TestSudoCommand(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		args     []string
		expected string
	}{
		{
			name:     "simple command",
			cmd:      "ls",
			args:     nil,
			expected: "sudo ls",
		},
		{
			name:     "command with args",
			cmd:      "echo",
			args:     []string{"hello", "world"},
			expected: "sudo echo hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SudoCommand(tt.cmd, tt.args...)
			if result != tt.expected {
				t.Errorf("SudoCommand() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRunWithSudo(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	// This test would require sudo password, so we skip it in automated tests
	t.Skip("Requires interactive sudo password")
}

func TestRunWithSudoOutput(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	// This test would require sudo password, so we skip it in automated tests
	t.Skip("Requires interactive sudo password")
}

func TestRequestSudo(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	// This test would require sudo password, so we skip it in automated tests
	t.Skip("Requires interactive sudo password")
}

// Test that non-sudo commands still work via the system runner
func TestSystemRunnerIntegration(t *testing.T) {
	// Verify we can still use the normal system runner
	runner := NewSystemRunner()
	ctx := context.Background()

	// Test a simple command that doesn't require sudo
	err := runner.Run(ctx, "echo", "test")
	if err != nil {
		t.Errorf("SystemRunner.Run() failed: %v", err)
	}

	// Test Output
	output, err := runner.Output(ctx, "echo", "test")
	if err != nil {
		t.Errorf("SystemRunner.Output() failed: %v", err)
	}
	// Output might or might not include trailing newline depending on implementation
	if output != "test" && output != "test\n" {
		t.Errorf("SystemRunner.Output() = %q, want %q or %q", output, "test", "test\n")
	}
}

func TestLookPath(t *testing.T) {
	// Test that common commands are available
	_, err := exec.LookPath("ls")
	if err != nil {
		t.Errorf("ls not found in PATH")
	}

	if runtime.GOOS != "windows" {
		_, err = exec.LookPath("sudo")
		if err != nil {
			t.Logf("sudo not found in PATH (expected on some systems)")
		}
	}
}
