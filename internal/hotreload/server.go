package hotreload

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Server manages SSE connections and broadcasts events to connected clients
type Server struct {
	clients     map[*client]bool
	clientsMux  sync.RWMutex
	messages    chan []byte
	newClients  chan chan []byte
	closed      chan bool
	broadcast   chan []byte
}

type client struct {
	messageChan chan []byte
}

// NewServer creates a new SSE server
func NewServer() *Server {
	return &Server{
		clients:    make(map[*client]bool),
		messages:   make(chan []byte),
		newClients: make(chan chan []byte),
		closed:     make(chan bool),
		broadcast:  make(chan []byte),
	}
}

// Start begins the SSE server in a background goroutine
func (s *Server) Start(ctx context.Context) {
	go s.run(ctx)
}

// run handles the main event loop for the SSE server
func (s *Server) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case clientChan := <-s.newClients:
			// Register new client
			s.clientsMux.Lock()
			c := &client{messageChan: clientChan}
			s.clients[c] = true
			s.clientsMux.Unlock()

		case msg := <-s.broadcast:
			// Broadcast message to all clients
			s.clientsMux.RLock()
			for c := range s.clients {
				select {
				case c.messageChan <- msg:
				default:
					// Client channel is full, skip this client
				}
			}
			s.clientsMux.RUnlock()

		case <-s.closed:
			// Server is shutting down
			s.clientsMux.Lock()
			for c := range s.clients {
				close(c.messageChan)
				delete(s.clients, c)
			}
			s.clientsMux.Unlock()
			return
		}
	}
}

// HandleSSE handles Server-Sent Events connections
func (s *Server) HandleSSE(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create message channel for this client
	messageChan := make(chan []byte, 10)

	// Register client
	s.newClients <- messageChan

	// Ensure client is cleaned up on disconnect
	defer func() {
		s.clientsMux.Lock()
		for c := range s.clients {
			if c.messageChan == messageChan {
				delete(s.clients, c)
				break
			}
		}
		s.clientsMux.Unlock()
		close(messageChan)
	}()

	// Send initial connection message
	fmt.Fprintf(w, "data: connected\n\n")
	flusher.Flush()

	// Start keep-alive ticker
	keepAlive := time.NewTicker(30 * time.Second)
	defer keepAlive.Stop()

	// Listen for messages
	for {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()

		case <-keepAlive.C:
			// Send keep-alive ping
			fmt.Fprintf(w, ": keep-alive\n\n")
			flusher.Flush()

		case <-r.Context().Done():
			return
		}
	}
}

// Broadcast sends a message to all connected clients
func (s *Server) Broadcast(message string) {
	s.broadcast <- []byte(message)
}

// Close shuts down the SSE server
func (s *Server) Close() {
	s.closed <- true
}

// Handler returns an HTTP handler for the SSE endpoint
func (s *Server) Handler() http.HandlerFunc {
	return s.HandleSSE
}
