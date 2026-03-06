package serve

import (
	"errors"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

type PortLookup interface {
	CurrentPort(projectDir string) (int, error)
	ProjectPort(name string) (int, error)
}

type portLookupFunc struct {
	current func(string) (int, error)
	project func(string) (int, error)
}

func (p portLookupFunc) CurrentPort(projectDir string) (int, error) {
	if p.current == nil {
		return 0, errors.New("current port lookup not configured")
	}
	return p.current(projectDir)
}

func (p portLookupFunc) ProjectPort(name string) (int, error) {
	if p.project == nil {
		return 0, errors.New("project port lookup not configured")
	}
	return p.project(name)
}

func IsProjectRunning(projectDir, name string, lookup PortLookup) (bool, error) {
	pidFile := filepath.Join(projectDir, DefaultPIDFile)
	if pid, err := readPID(pidFile); err == nil {
		if running, _ := pidRunning(pid); running {
			return true, nil
		}
	}

	if lookup == nil {
		return false, nil
	}

	var port int
	var err error
	if name == "" {
		port, err = lookup.CurrentPort(projectDir)
	} else {
		port, err = lookup.ProjectPort(name)
	}
	if err != nil || port == 0 {
		return false, err
	}
	return portInUse(port), nil
}

func readPID(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, err
	}
	return pid, nil
}

func pidRunning(pid int) (bool, error) {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false, err
	}
	if err := process.Signal(syscall.Signal(0)); err != nil {
		return false, err
	}
	return true, nil
}

func portInUse(port int) bool {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return true
	}
	_ = listener.Close()
	return false
}
