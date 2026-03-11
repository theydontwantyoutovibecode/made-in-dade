package hotreload

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	s := NewServer()
	if s == nil {
		t.Fatal("server is nil")
	}
	if s.clients == nil {
		t.Fatal("clients map is nil")
	}
	if s.messages == nil {
		t.Fatal("messages channel is nil")
	}
}

func TestServerStartAndBroadcast(t *testing.T) {
	s := NewServer()
	ctx := context.Background()
	s.Start(ctx)
	defer s.Close()

	// Give the server time to start
	time.Sleep(50 * time.Millisecond)

	// Broadcast a message
	s.Broadcast("test message")

	// Wait a bit for processing
	time.Sleep(50 * time.Millisecond)
}

func TestSSEConnection(t *testing.T) {
	s := NewServer()
	ctx := context.Background()
	s.Start(ctx)
	defer s.Close()

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(s.HandleSSE))
	defer ts.Close()

	// Connect as a client
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("failed to connect to SSE endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Check headers
	if resp.Header.Get("Content-Type") != "text/event-stream" {
		t.Errorf("expected Content-Type text/event-stream, got %s", resp.Header.Get("Content-Type"))
	}
	if resp.Header.Get("Cache-Control") != "no-cache" {
		t.Errorf("expected Cache-Control no-cache, got %s", resp.Header.Get("Cache-Control"))
	}

	// Read initial connection message
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if !strings.Contains(string(body), "data: connected") {
		t.Errorf("expected 'data: connected' in response, got: %s", string(body))
	}
}

func TestBroadcastToClients(t *testing.T) {
	s := NewServer()
	ctx := context.Background()
	s.Start(ctx)
	defer s.Close()

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(s.HandleSSE))
	defer ts.Close()

	// Connect as a client
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("failed to connect to SSE endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Give time for connection
	time.Sleep(50 * time.Millisecond)

	// Broadcast a message
	s.Broadcast("test message")

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if !strings.Contains(string(body), "data: test message") {
		t.Errorf("expected 'data: test message' in response, got: %s", string(body))
	}
}

func TestHTTPServer(t *testing.T) {
	hs := NewHTTPServer(0) // Use random port
	if hs == nil {
		t.Fatal("HTTPServer is nil")
	}

	hs.Start()
	defer hs.Stop()

	// Give time for server to start
	time.Sleep(100 * time.Millisecond)
}
