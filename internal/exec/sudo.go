package exec

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// RunWithSudo executes a command with sudo privileges
func RunWithSudo(name string, args ...string) error {
	cmd := exec.Command("sudo", append([]string{name}, args...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunWithSudoOutput executes a command with sudo and returns output
func RunWithSudoOutput(name string, args ...string) (string, error) {
	cmd := exec.Command("sudo", append([]string{name}, args...)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("sudo command failed: %w", err)
	}
	return string(output), nil
}

// IsAdmin checks if the current process has admin/root privileges
func IsAdmin() bool {
	if runtime.GOOS == "windows" {
		// Windows admin check
		cmd := exec.Command("net", "session")
		err := cmd.Run()
		return err == nil
	}

	// Unix-like systems: check if we're root (UID 0)
	return os.Getuid() == 0
}

// RequestSudo prompts for sudo password and caches credentials
func RequestSudo() error {
	// Run a simple sudo command to cache credentials
	cmd := exec.Command("sudo", "-v")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SudoCommand builds a sudo command string for display purposes
func SudoCommand(name string, args ...string) string {
	return "sudo " + strings.Join(append([]string{name}, args...), " ")
}

// CanUseSudo checks if sudo is available on the system
func CanUseSudo() bool {
	if runtime.GOOS == "windows" {
		return true // Windows uses UAC, not sudo
	}

	_, err := exec.LookPath("sudo")
	return err == nil
}
