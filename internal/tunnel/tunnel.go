package tunnel

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Tunnel represents a running cloudflared tunnel.
type Tunnel struct {
	Cmd       *exec.Cmd
	URL       string
	IsNamed   bool
	Name      string
	LocalPort int

	urlChan chan string
	done    chan error
	mu      sync.Mutex
}

// quickTunnelURLPattern matches cloudflare quick tunnel URLs.
var quickTunnelURLPattern = regexp.MustCompile(`https://[a-z0-9-]+\.trycloudflare\.com`)

// IsAvailable checks if cloudflared is installed.
func IsAvailable() bool {
	_, err := exec.LookPath("cloudflared")
	return err == nil
}

// StartQuick creates a quick (anonymous) tunnel to the given local port.
func StartQuick(ctx context.Context, port int) (*Tunnel, error) {
	localURL := fmt.Sprintf("http://localhost:%d", port)
	cmd := exec.CommandContext(ctx, "cloudflared", "tunnel", "--url", localURL)

	// Create pipes for capturing output while forwarding
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	t := &Tunnel{
		Cmd:       cmd,
		LocalPort: port,
		IsNamed:   false,
		urlChan:   make(chan string, 1),
		done:      make(chan error, 1),
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start cloudflared: %w", err)
	}

	// Monitor both stdout and stderr for the tunnel URL
	go t.scanForURL(stdoutPipe, os.Stdout)
	go t.scanForURL(stderrPipe, os.Stderr)

	// Monitor process completion
	go func() {
		t.done <- cmd.Wait()
		close(t.done)
	}()

	return t, nil
}

// StartNamed starts a named tunnel with the given configuration.
func StartNamed(ctx context.Context, name string, port int) (*Tunnel, error) {
	localURL := fmt.Sprintf("http://localhost:%d", port)
	cmd := exec.CommandContext(ctx, "cloudflared", "tunnel", "--url", localURL, "run", name)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	t := &Tunnel{
		Cmd:       cmd,
		LocalPort: port,
		IsNamed:   true,
		Name:      name,
		urlChan:   make(chan string, 1),
		done:      make(chan error, 1),
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start cloudflared: %w", err)
	}

	// Monitor process completion
	go func() {
		t.done <- cmd.Wait()
		close(t.done)
	}()

	return t, nil
}

// WaitForURL waits for the tunnel URL to be captured.
// Returns empty string if timeout or tunnel exits.
func (t *Tunnel) WaitForURL(timeout time.Duration) string {
	select {
	case url := <-t.urlChan:
		return url
	case <-time.After(timeout):
		return ""
	case <-t.done:
		return ""
	}
}

// Stop terminates the tunnel process.
func (t *Tunnel) Stop() error {
	if t.Cmd == nil || t.Cmd.Process == nil {
		return nil
	}
	return t.Cmd.Process.Kill()
}

// Wait waits for the tunnel process to exit.
func (t *Tunnel) Wait() error {
	return <-t.done
}

// Done returns a channel that's closed when the tunnel exits.
func (t *Tunnel) Done() <-chan error {
	return t.done
}

func (t *Tunnel) scanForURL(r io.Reader, forward io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		// Forward to output
		if forward != nil {
			fmt.Fprintln(forward, line)
		}

		// Check for URL
		if match := quickTunnelURLPattern.FindString(line); match != "" {
			t.mu.Lock()
			if t.URL == "" {
				t.URL = match
				select {
				case t.urlChan <- match:
				default:
				}
			}
			t.mu.Unlock()
		}
	}
}

// ListNamedTunnels returns a list of configured named tunnels.
func ListNamedTunnels(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, "cloudflared", "tunnel", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list tunnels: %w", err)
	}

	var tunnels []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines[1:] { // Skip header
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			tunnels = append(tunnels, fields[1]) // Name is second column
		}
	}
	return tunnels, nil
}

// HasNamedTunnel checks if a tunnel with the given name exists.
func HasNamedTunnel(ctx context.Context, name string) bool {
	tunnels, err := ListNamedTunnels(ctx)
	if err != nil {
		return false
	}
	for _, t := range tunnels {
		if strings.EqualFold(t, name) {
			return true
		}
	}
	return false
}
