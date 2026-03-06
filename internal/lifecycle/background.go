package lifecycle

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

// Process represents a running background process.
type Process struct {
	Cmd     *exec.Cmd
	Command string
	Done    chan error
}

// BackgroundManager manages multiple background processes.
type BackgroundManager struct {
	Dir       string
	Env       []string
	Stdout    io.Writer
	Stderr    io.Writer
	processes []*Process
	mu        sync.Mutex
}

// NewBackgroundManager creates a manager for background processes.
func NewBackgroundManager(dir string) *BackgroundManager {
	return &BackgroundManager{
		Dir:       dir,
		Env:       os.Environ(),
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
		processes: make([]*Process, 0),
	}
}

// WithEnv adds environment variables to all background processes.
func (m *BackgroundManager) WithEnv(env []string) *BackgroundManager {
	m.Env = append(m.Env, env...)
	return m
}

// WithOutput sets custom stdout/stderr writers for all processes.
func (m *BackgroundManager) WithOutput(stdout, stderr io.Writer) *BackgroundManager {
	m.Stdout = stdout
	m.Stderr = stderr
	return m
}

// Start launches a background command and tracks it.
func (m *BackgroundManager) Start(ctx context.Context, cmdStr string) (*Process, error) {
	cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)
	cmd.Dir = m.Dir
	cmd.Env = m.Env
	cmd.Stdout = m.Stdout
	cmd.Stderr = m.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start %q: %w", cmdStr, err)
	}

	proc := &Process{
		Cmd:     cmd,
		Command: cmdStr,
		Done:    make(chan error, 1),
	}

	// Monitor process completion in goroutine
	go func() {
		proc.Done <- cmd.Wait()
		close(proc.Done)
	}()

	m.mu.Lock()
	m.processes = append(m.processes, proc)
	m.mu.Unlock()

	return proc, nil
}

// StartAll launches multiple background commands.
func (m *BackgroundManager) StartAll(ctx context.Context, commands []string) error {
	for _, cmdStr := range commands {
		if _, err := m.Start(ctx, cmdStr); err != nil {
			// If one fails to start, stop the others
			m.StopAll()
			return err
		}
	}
	return nil
}

// StopAll terminates all running background processes.
func (m *BackgroundManager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, proc := range m.processes {
		if proc.Cmd.Process != nil {
			_ = proc.Cmd.Process.Kill()
		}
	}
	m.processes = nil
}

// Wait waits for all background processes to complete.
func (m *BackgroundManager) Wait() error {
	m.mu.Lock()
	procs := make([]*Process, len(m.processes))
	copy(procs, m.processes)
	m.mu.Unlock()

	var firstErr error
	for _, proc := range procs {
		if err := <-proc.Done; err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// Processes returns a snapshot of currently tracked processes.
func (m *BackgroundManager) Processes() []*Process {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]*Process, len(m.processes))
	copy(result, m.processes)
	return result
}
