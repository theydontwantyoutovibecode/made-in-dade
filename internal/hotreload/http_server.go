package hotreload

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"
)

// HTTPServer wraps the SSE server with HTTP endpoints
type HTTPServer struct {
	sseServer   *Server
	httpServer  *http.Server
	reloadChan  chan string
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewHTTPServer creates a new HTTP server for hot-reload
// If fileHandler is provided, it will be wrapped with script injection middleware
func NewHTTPServer(port int, fileHandler http.Handler) *HTTPServer {
	sseServer := NewServer()
	mux := http.NewServeMux()
	hs := &HTTPServer{
		sseServer:  sseServer,
		reloadChan: make(chan string, 100),
	}

	hs.ctx, hs.cancel = context.WithCancel(context.Background())

	// Register SSE endpoint
	mux.HandleFunc("/_dade/events", sseServer.Handler())

	// Register reload trigger endpoint
	mux.HandleFunc("/_dade/reload", hs.handleReload)

	// If file handler provided, wrap with script injection and register as default
	if fileHandler != nil {
		handler := InjectScriptMiddleware(fileHandler)
		mux.Handle("/", handler)
	}

	// Create HTTP server
	hs.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return hs
}

// handleReload handles requests to trigger a page reload
func (hs *HTTPServer) handleReload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Broadcast reload event to all clients
	hs.sseServer.Broadcast("reload")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

// Start starts the HTTP server in the background
func (hs *HTTPServer) Start() error {
	go func() {
		hs.sseServer.Start(hs.ctx)
	}()

	go func() {
		if err := hs.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	return nil
}

// Stop stops the HTTP server
func (hs *HTTPServer) Stop() error {
	hs.cancel()
	hs.sseServer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return hs.httpServer.Shutdown(ctx)
}

// TriggerReload sends a reload event to all connected clients
func (hs *HTTPServer) TriggerReload() {
	hs.sseServer.Broadcast("reload")
}

// Port returns the server's port as a string
func (hs *HTTPServer) Port() string {
	return hs.httpServer.Addr[1:] // Remove leading ":"
}

// NewFileServer creates a simple file server for serving static files
func NewFileServer(root string) http.Handler {
	fileServer := http.FileServer(http.Dir(root))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set cache-control headers for CSS, HTML, and JS files to prevent caching
		if filepath.Ext(r.URL.Path) == ".css" || filepath.Ext(r.URL.Path) == ".html" || filepath.Ext(r.URL.Path) == ".js" {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		}

		// Try to serve HTML file if directory
		if r.URL.Path == "/" {
			// Serve index.html directly without modifying the request
			http.ServeFile(w, r, filepath.Join(root, "index.html"))
			return
		}
		fileServer.ServeHTTP(w, r)
	})
}
