package lifecycle

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// CleanupFunc is a function called during shutdown.
type CleanupFunc func()

// CleanupManager handles signal-based cleanup and shutdown.
type CleanupManager struct {
	cleanups []CleanupFunc
	mu       sync.Mutex
	done     chan struct{}
	started  bool
}

// NewCleanupManager creates a new cleanup manager.
func NewCleanupManager() *CleanupManager {
	return &CleanupManager{
		cleanups: make([]CleanupFunc, 0),
		done:     make(chan struct{}),
	}
}

// Register adds a cleanup function to be called on shutdown.
// Functions are called in reverse order of registration (LIFO).
func (m *CleanupManager) Register(fn CleanupFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cleanups = append(m.cleanups, fn)
}

// HandleSignals starts listening for SIGINT and SIGTERM.
// When received, runs all cleanup functions and returns.
// Call this in a goroutine or before your main blocking operation.
func (m *CleanupManager) HandleSignals(ctx context.Context) {
	m.mu.Lock()
	if m.started {
		m.mu.Unlock()
		return
	}
	m.started = true
	m.mu.Unlock()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		m.runCleanups()
	case <-ctx.Done():
		m.runCleanups()
	case <-m.done:
		// Manual shutdown
	}
}

// Shutdown manually triggers cleanup without waiting for signals.
func (m *CleanupManager) Shutdown() {
	m.mu.Lock()
	select {
	case <-m.done:
		// Already closed
	default:
		close(m.done)
	}
	m.mu.Unlock()
	m.runCleanups()
}

// Done returns a channel that's closed when shutdown is triggered.
func (m *CleanupManager) Done() <-chan struct{} {
	return m.done
}

func (m *CleanupManager) runCleanups() {
	m.mu.Lock()
	cleanups := make([]CleanupFunc, len(m.cleanups))
	copy(cleanups, m.cleanups)
	m.cleanups = nil // Clear to prevent double-run
	m.mu.Unlock()

	// Run in reverse order (LIFO)
	for i := len(cleanups) - 1; i >= 0; i-- {
		cleanups[i]()
	}
}

// WithCleanup creates a context that will trigger cleanup on cancellation.
func WithCleanup(parent context.Context, manager *CleanupManager) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)

	go func() {
		<-ctx.Done()
		manager.Shutdown()
	}()

	return ctx, cancel
}
