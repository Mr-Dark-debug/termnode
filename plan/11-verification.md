# Testing & Verification Plan

## Overview

This document describes how to verify that TermNode works correctly end-to-end, from build to runtime behavior.

---

## 1. Build Verification

### Prerequisites Check

```bash
# Go version (must be 1.22+)
go version

# Termux tools
which termux-battery-status
which termux-wake-lock
which sshd

# Available disk space
df -h /storage/emulated/0
```

### Build Steps

```bash
cd /storage/emulated/0/Projects/termnode

# 1. Download dependencies
go mod tidy
# Expected: downloads bubbletea, lipgloss, bubbles, modernc.org/sqlite
# Creates go.sum file

# 2. Compile
go build -o termnode .
# Expected: no errors, creates ~15-20MB binary

# 3. Verify binary
file termnode
# Expected: ELF 64-bit LSB executable, ARM aarch64

ls -lh termnode
# Expected: 15-25MB

# 4. Check version
./termnode -version
# Expected: termnode dev
```

### Build Failure Troubleshooting

| Error | Cause | Fix |
|-------|-------|-----|
| `undefined: someType` | Import path mismatch | Check module name in go.mod matches imports |
| `cgo: exec gcc` | CGO dependency found | Ensure all deps are pure Go |
| `cannot find package` | Missing dependency | Run `go mod tidy` |
| `go: module requires ...` | Go version too old | Update Go: `pkg upgrade golang` |

---

## 2. Unit Tests

### Running Tests

```bash
# All tests
go test ./...

# Verbose output
go test -v ./...

# Specific package
go test ./internal/db/...

# With coverage
go test -cover ./...
```

### Database Tests

Database tests use SQLite's `:memory:` mode for isolation:

```go
func TestRepository(t *testing.T) {
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()

    // Run migration manually
    db.Exec(schema)

    repo := NewRepository(db)

    // Test Insert
    id, err := repo.Insert(IoTEvent{Topic: "test", Payload: "hello", Source: "http"})
    assert.NoError(err)
    assert.Equal(int64(1), id)

    // Test Recent
    events, err := repo.Recent(10)
    assert.NoError(err)
    assert.Len(events, 1)

    // Test ByTopic
    events, err = repo.ByTopic("test", 10)
    assert.NoError(err)
    assert.Len(events, 1)

    // Test Count
    count, err := repo.Count()
    assert.NoError(err)
    assert.Equal(int64(1), count)

    // Test Purge
    purged, err := repo.Purge(0) // purge everything older than 0
    assert.NoError(err)
    assert.Equal(int64(1), purged)
}
```

---

## 3. TUI Smoke Test

### Launch Test

```bash
./termnode
```

**Expected behavior:**
1. Terminal switches to alternate screen
2. Tab bar appears at top: `1:Dashboard  2:Services  3:IoT Log  4:Help`
3. Dashboard tab is active (highlighted)
4. If termux-api is installed: real data appears
5. If termux-api is NOT installed: "termux-api unavailable" shown

### Tab Switching Test

| Action | Expected |
|--------|----------|
| Press `F1` | Dashboard tab active, hardware panels visible |
| Press `F2` | Services tab active, 3 services listed |
| Press `F3` | IoT Log tab active, event log or "No events" |
| Press `F4` | Help tab active, keybindings table visible |
| Press `1` | Same as F1 |
| Press `2` | Same as F2 |
| Press `3` | Same as F3 |
| Press `?` | Same as F4 |
| Press `q` | Clean exit, terminal restored |

### Window Resize Test

```bash
# Launch in a terminal, then resize the window
# Expected: TUI reflows to fit new dimensions
# Panels resize appropriately
# No text overflow or clipping
```

---

## 4. Hardware Dashboard Test

### With termux-api installed

```bash
# Verify termux-api works
termux-battery-status
termux-wifi-connectioninfo

# Launch TermNode
./termnode

# Expected on Dashboard (F1):
# Left panel:
#   Battery
#   Level: XX%
#   Status: CHARGING/DISCHARGING
#   Health: GOOD
#   Temp: XX.X°C
#   Current: XXX mA
#   [████████░░░░░]  (color-coded bar)
#
# Right panel:
#   System
#   WiFi SSID: XXXXX
#   Local IP:  192.168.X.X
#   BSSID:     XX:XX:XX:XX:XX:XX
#
#   CPU
#   Usage:     XX.X%
#   Cores:     8
#   Temp:      XX.X°C
#   Arch:      aarch64
```

### Auto-refresh Test

1. Note the CPU usage percentage
2. Wait 5 seconds
3. CPU usage should update (it's a live value)

### Without termux-api

```bash
# Uninstall termux-api temporarily
pkg uninstall termux-api

./termnode
# Expected: Dashboard shows "termux-api unavailable" gracefully (no crash)
```

---

## 5. Service Manager Test

### Wake Lock Test

```bash
./termnode

# Press F2 → Services tab
# Navigate to "Wake Lock" → Press Enter

# In another terminal:
pgrep -x termux-wake-lock
# Expected: shows a PID

# Back in TermNode:
# "Wake Lock" shows ● ON Running

# Press Enter again to toggle off
# "Wake Lock" shows ● OFF Stopped
```

### SSH Server Test

```bash
# Prerequisite: pkg install openssh && passwd

./termnode

# Press F2 → Navigate to "SSH Server" → Press Enter
# Status: ● ON Running

# In another terminal:
ssh localhost -p 8022
# Expected: password prompt (SSH server is running)

# Back in TermNode:
# Press Enter to stop
# Status: ● OFF Stopped
```

### HTTP File Server Test

```bash
./termnode

# Press F2 → Navigate to "HTTP File Server" → Press Enter
# Status: ● ON Running

# In another terminal:
curl http://localhost:8081/
# Expected: HTML listing of home directory files

# Or open in browser: http://<phone-ip>:8081/

# Back in TermNode:
# Press Enter to stop
# Status: ● OFF Stopped
```

---

## 6. IoT Bridge Test

### Webhook Reception Test

```bash
# Terminal 1: Run TermNode
./termnode
# Press F3 → IoT Log tab

# Terminal 2: Send test events
curl -X POST http://localhost:8080/webhook/temperature -d '{"sensor":"DHT22","value":23.5}'
# Expected response: {"status":"ok"}

curl -X POST http://localhost:8080/webhook/humidity -d '{"sensor":"DHT22","value":65.2}'
curl -X POST http://localhost:8080/webhook/motion -d '{"zone":"kitchen","detected":true}'

# Back in TermNode IoT Log:
# Expected: 3 events displayed with timestamp, topic, and payload
# Auto-scroll to latest event
# Status bar: "Events: 3"
```

### Health Check Test

```bash
curl http://localhost:8080/health
# Expected: {"status":"ok"}
```

### Persistence Test

```bash
# Send some events
curl -X POST http://localhost:8080/webhook/test -d '{"persist":"me"}'

# Quit TermNode (q)

# Restart TermNode
./termnode
# Press F3 → IoT Log tab

# Expected: Previous event is still visible (loaded from SQLite)

# Verify database directly:
sqlite3 ~/.termnode/termnode.db "SELECT * FROM events;"
```

### Large Payload Test

```bash
# Send 1MB payload (should be accepted)
python3 -c "print('x' * 1000000)" | curl -X POST http://localhost:8080/webhook/big -d @-

# Send >1MB payload (should be truncated or rejected)
python3 -c "print('x' * 2000000)" | curl -X POST http://localhost:8080/webhook/huge -d @-
```

---

## 7. Graceful Shutdown Test

```bash
# Start TermNode with services running
./termnode

# F2 → Start SSH server and HTTP file server
# F3 → Send some webhook events

# Press q to quit

# Verify:
# 1. Terminal restored to normal (not stuck in alt screen)
# 2. SSH server is still running (it's a detached process with PID file)
# 3. HTTP file server is stopped (it's a Go goroutine, dies with process)
# 4. No zombie processes: ps aux | grep termnode
```

---

## 8. Edge Cases

| Test | How | Expected |
|------|-----|----------|
| No database directory | `rm -rf ~/.termnode` then run | Auto-creates directory and database |
| Corrupted database | `echo "garbage" > ~/.termnode/termnode.db` then run | Fails with clear error message |
| Port already in use | Start another server on 8080, then run | Webhook server error shown in log |
| No termux-api | Uninstall termux-api, run | Dashboard shows "unavailable", no crash |
| No sshd | Don't install openssh, try to start SSH | Error shown in TUI |
| Very small terminal | Resize to 40x10 | UI degrades gracefully, no crash |
| No internet | Disconnect WiFi, run | Network info empty, everything else works |
