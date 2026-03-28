# Phase 4: Service Manager

## Goal

Implement the interactive service manager that allows users to toggle critical background services from the TUI — wake-lock, SSH server, and HTTP file server.

---

## Managed Services

| Service | Key | Start Command | Stop Method | Port |
|---------|-----|--------------|-------------|------|
| Wake Lock | `wakelock` | `termux-wake-lock` | `termux-wake-unlock` | N/A |
| SSH Server | `sshd` | `sshd` | SIGINT + PID file | 8022 |
| HTTP File Server | `httpfs` | `net/http.FileServer` | `server.Shutdown()` | 8081 |

---

## Files

### internal/daemon/daemon.go

Simple wrapper around Termux wake-lock commands:

```go
func EnableWakeLock() error {
    return exec.Command("termux-wake-lock").Run()
}

func DisableWakeLock() error {
    return exec.Command("termux-wake-unlock").Run()
}
```

**How wake-lock works:**
- `termux-wake-lock` acquires a partial wake lock from Android
- Prevents Android from killing Termux background processes
- The lock persists until `termux-wake-unlock` is called or the device reboots
- This is the #1 most-needed feature for headless Termux servers

### internal/daemon/process.go

The `Manager` struct handles all service lifecycle:

```go
type Service struct {
    Name    string
    Key     string
    Desc    string
    Running bool
}

type Manager struct {
    mu        sync.Mutex
    pids      map[string]int           // key → PID
    servers   map[string]*http.Server  // key → HTTP server
    listeners map[string]net.Listener  // key → listener
}
```

**Key Methods:**

| Method | Purpose |
|--------|---------|
| `Start(key)` | Routes to appropriate start function based on key |
| `Stop(key)` | Routes to appropriate stop function |
| `Status(key)` | Checks if service is currently running |

### SSH Server Management

```
Start("sshd"):
  1. Check if sshd already running (PID file or pgrep)
  2. exec.Command("sshd").Start()
  3. Record PID in m.pids["sshd"]
  4. Write PID to $HOME/.termnode/sshd.pid
  5. Background goroutine: c.Wait() → clean up PID

Stop("sshd"):
  1. Read PID from m.pids or PID file
  2. os.FindProcess(pid)
  3. proc.Signal(os.Interrupt) — graceful stop
  4. Fallback: proc.Kill() if interrupt fails
  5. Delete PID from map and file

Status("sshd"):
  1. Check m.pids map
  2. If no tracked PID, check pgrep for external sshd
  3. Verify PID is alive via proc.Signal(nil)
```

### HTTP File Server

Built with Go's standard `net/http.FileServer`:

```go
func (m *Manager) startFileServer(key string) error {
    home, _ := os.UserHomeDir()
    mux := http.NewServeMux()
    mux.Handle("/", http.FileServer(http.Dir(home)))

    srv := &http.Server{Addr: ":8081", Handler: mux}
    ln, _ := net.Listen("tcp", ":8081")

    m.servers[key] = srv
    m.listeners[key] = ln
    go srv.Serve(ln)
}
```

**Why Go's FileServer instead of Python?**
- Zero additional dependencies (no Python needed)
- Single binary handles everything
- Consistent with the project's Go-native philosophy
- Proper MIME type detection built-in
- Can be extended with auth, logging, etc.

### internal/screen/services.go

The TUI screen provides:
- **Service list** with cursor navigation (↑/↓ or j/k)
- **Status indicators**: `● ON` (green) / `● OFF` (red)
- **Toggle on Enter**: Starts or stops the selected service
- **Description text** for each service
- **Navigation hint** at the bottom

**View rendering:**
```
  Service Management

> Wake Lock            ● OFF   Prevent Android battery optimization
    Status: Stopped

  SSH Server           ● ON    OpenSSH remote access
    Status: Running

  HTTP File Server     ● OFF   Serve files over HTTP
    Status: Stopped

  Press Enter to toggle  |  ↑/↓ or j/k to navigate
```

---

## PID File Tracking

PID files are stored at `$HOME/.termnode/`:

```
~/.termnode/
├── termnode.db        # SQLite database
├── sshd.pid           # SSH server PID (when running)
└── httpfs.pid         # File server PID (when running)
```

**Why PID files?**
- Persist across TermNode restarts
- Can detect services started by previous instances
- Simple and reliable on Linux/Android
- Alternative: D-Bus or socket activation — overkill for Termux

---

## Process Lifecycle

```
                    Start
                      │
                      ▼
            ┌─────────────────┐
            │ exec.Command()  │
            │ .Start()        │
            └────────┬────────┘
                     │
                     ▼
            ┌─────────────────┐
            │ Record PID      │
            │ Write PID file  │
            └────────┬────────┘
                     │
                     ▼
            ┌─────────────────┐
            │ go proc.Wait()  │◄── Background reaper
            │ → clean up PID  │
            └─────────────────┘

                    Stop
                      │
                      ▼
            ┌─────────────────┐
            │ Read PID        │
            └────────┬────────┘
                     │
                     ▼
            ┌─────────────────┐
            │ SIGINT          │
            └────────┬────────┘
                     │
              ┌──────┴──────┐
              │             │
         Success        Failed
              │             │
              ▼             ▼
        Clean up PID   SIGKILL
                       Clean up PID
```

---

## Error Handling

| Error | User-Facing Message |
|-------|-------------------|
| `sshd` not installed | "Failed to start sshd: executable not found" |
| Port 8081 in use | "Failed to start file server: address already in use" |
| `termux-wake-lock` fails | "Failed to enable wake lock: termux-api not available" |
| Service already running | "sshd is already running (PID 12345)" |

---

## Prerequisites

```bash
# For SSH server
pkg install openssh

# For wake-lock
pkg install termux-api
# Also install Termux:API app from F-Droid

# Generate SSH key (first time only)
ssh-keygen -t ed25519
# Set password for SSH access
passwd
```

---

## Verification

```bash
# Build and run
go build -o termnode . && ./termnode

# Press F2 to go to Services tab

# Test wake-lock toggle
# - Navigate to "Wake Lock", press Enter
# - Status changes to "● ON Running"
# - Run in another terminal: pgrep -x termux-wake-lock → should show PID

# Test SSH server toggle
# - Navigate to "SSH Server", press Enter
# - Status changes to "● ON Running"
# - Run: ssh localhost -p 8022 → should prompt for password

# Test HTTP file server
# - Navigate to "HTTP File Server", press Enter
# - Status changes to "● ON Running"
# - Run: curl http://localhost:8081/ → should list home directory

# Toggle all off
# - Navigate to each, press Enter
# - All show "● OFF Stopped"
```
