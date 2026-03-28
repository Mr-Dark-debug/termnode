package daemon

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// Service represents a manageable service with its current state.
type Service struct {
	Name    string
	Key     string
	Desc    string
	Running bool
}

// Manager handles starting, stopping, and querying services.
type Manager struct {
	mu        sync.Mutex
	pids      map[string]int              // service key -> PID
	servers   map[string]*http.Server     // service key -> HTTP server
	listeners map[string]net.Listener     // service key -> listener
}

// NewManager creates a new service manager.
func NewManager() *Manager {
	return &Manager{
		pids:      make(map[string]int),
		servers:   make(map[string]*http.Server),
		listeners: make(map[string]net.Listener),
	}
}

// Start starts a service by key.
func (m *Manager) Start(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch key {
	case "wakelock":
		return EnableWakeLock()
	case "sshd":
		return m.startProcess(key, "sshd", nil)
	case "httpfs":
		return m.startFileServer(key)
	default:
		return fmt.Errorf("unknown service: %s", key)
	}
}

// Stop stops a service by key.
func (m *Manager) Stop(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch key {
	case "wakelock":
		return DisableWakeLock()
	case "sshd":
		return m.stopProcess(key)
	case "httpfs":
		return m.stopFileServer(key)
	default:
		return fmt.Errorf("unknown service: %s", key)
	}
}

// Status returns whether a service is currently running.
func (m *Manager) Status(key string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch key {
	case "wakelock":
		return isProcessRunning("termux-wake-lock"), nil
	case "sshd":
		pid, exists := m.pids[key]
		if !exists {
			return isProcessRunning("sshd"), nil
		}
		return isPIDAlive(pid), nil
	case "httpfs":
		_, exists := m.servers[key]
		return exists, nil
	default:
		return false, fmt.Errorf("unknown service: %s", key)
	}
}

func (m *Manager) startProcess(key, cmd string, args []string) error {
	if pid, exists := m.pids[key]; exists && isPIDAlive(pid) {
		return fmt.Errorf("%s is already running (PID %d)", key, pid)
	}

	c := exec.Command(cmd, args...)
	if err := c.Start(); err != nil {
		return fmt.Errorf("failed to start %s: %w", key, err)
	}

	m.pids[key] = c.Process.Pid

	pidDir := pidDir()
	os.MkdirAll(pidDir, 0755)
	pidFile := filepath.Join(pidDir, key+".pid")
	os.WriteFile(pidFile, []byte(strconv.Itoa(c.Process.Pid)), 0644)

	go func() {
		c.Wait()
		m.mu.Lock()
		delete(m.pids, key)
		m.mu.Unlock()
	}()

	return nil
}

func (m *Manager) stopProcess(key string) error {
	pid, exists := m.pids[key]
	if !exists {
		pidFile := filepath.Join(pidDir(), key+".pid")
		data, err := os.ReadFile(pidFile)
		if err != nil {
			return fmt.Errorf("%s is not running", key)
		}
		pid, _ = strconv.Atoi(string(data))
	}

	if pid == 0 || !isPIDAlive(pid) {
		delete(m.pids, key)
		return nil
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	if err := proc.Signal(os.Interrupt); err != nil {
		proc.Kill()
	}

	delete(m.pids, key)
	os.Remove(filepath.Join(pidDir(), key+".pid"))
	return nil
}

func (m *Manager) startFileServer(key string) error {
	if _, exists := m.servers[key]; exists {
		return fmt.Errorf("file server already running")
	}

	home, _ := os.UserHomeDir()
	addr := ":8081"

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(home)))

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start file server: %w", err)
	}

	m.servers[key] = srv
	m.listeners[key] = ln

	go srv.Serve(ln)
	return nil
}

func (m *Manager) stopFileServer(key string) error {
	srv, exists := m.servers[key]
	if !exists {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	delete(m.servers, key)
	delete(m.listeners, key)
	return nil
}

func pidDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".termnode")
}

func isPIDAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = proc.Signal(os.Signal(nil))
	return err == nil
}

func isProcessRunning(name string) bool {
	_, err := exec.Command("pgrep", "-x", name).Output()
	return err == nil
}
