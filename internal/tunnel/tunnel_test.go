package tunnel

import (
	"context"
	"testing"
	"time"
)

func TestIsAvailable(t *testing.T) {
	// This test just checks the function doesn't panic
	// The result depends on whether cloudflared is installed
	_ = IsAvailable()
}

func TestQuickTunnelURLPattern(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "Your quick Tunnel has been created! Visit it at https://abc-def-123.trycloudflare.com",
			expected: "https://abc-def-123.trycloudflare.com",
		},
		{
			input:    "https://my-test-tunnel.trycloudflare.com is ready",
			expected: "https://my-test-tunnel.trycloudflare.com",
		},
		{
			input:    "No tunnel URL here",
			expected: "",
		},
		{
			input:    "2026/02/19 INF https://some-tunnel-name.trycloudflare.com",
			expected: "https://some-tunnel-name.trycloudflare.com",
		},
	}

	for _, tc := range tests {
		match := quickTunnelURLPattern.FindString(tc.input)
		if match != tc.expected {
			t.Errorf("input %q: expected %q, got %q", tc.input, tc.expected, match)
		}
	}
}

func TestTunnelStopNilProcess(t *testing.T) {
	tunnel := &Tunnel{}
	// Should not panic
	err := tunnel.Stop()
	if err != nil {
		t.Errorf("expected no error stopping nil process, got: %v", err)
	}
}

func TestStartQuickWithoutCloudflared(t *testing.T) {
	if IsAvailable() {
		t.Skip("skipping: cloudflared is installed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := StartQuick(ctx, 8000)
	if err == nil {
		t.Error("expected error when cloudflared not available")
	}
}
