# TermNode — Master Plan

## What Is This?

This folder contains the complete planning documentation for the TermNode project — a unified TUI for headless Android management in Termux. Every phase, decision, and architectural detail is documented here.

## Files In This Folder

| File | Description |
|------|-------------|
| `01-overview.md` | This file — project overview, goals, and plan index |
| `02-research.md` | Technology research findings (Bubble Tea, Termux, Go patterns) |
| `03-architecture.md` | Complete system architecture and data flow |
| `04-directory-structure.md` | File/folder layout with descriptions of every file |
| `05-phase1-skeleton.md` | Phase 1: Project init + skeleton TUI |
| `06-phase2-database.md` | Phase 2: SQLite database layer |
| `07-phase3-dashboard.md` | Phase 3: Hardware dashboard |
| `08-phase4-services.md` | Phase 4: Service manager |
| `09-phase5-iot-bridge.md` | Phase 5: IoT bridge + log viewer |
| `10-phase6-polish.md` | Phase 6: Polish, Makefile, MQTT |
| `11-verification.md` | Testing and verification plan |
| `12-decisions.md` | Key architectural decisions and trade-offs |

---

## Project Goals

### The Problem

Users repurposing old Android phones as headless servers via Termux face fragmented management:

- Raw SSH commands to check battery status
- Bash scripts to keep services alive
- Android's battery optimization randomly killing processes
- No unified way to monitor hardware or manage services
- Fragmented IoT data collection

### The Solution

A **single, compiled Go binary** that drops a beautiful TUI dashboard into the terminal, providing:

1. **Hardware Dashboard** — Real-time battery, CPU, network monitoring
2. **Service Manager** — One-toggle controls for wake-lock, sshd, file server
3. **IoT Bridge** — HTTP webhook receiver with SQLite storage
4. **Help System** — Always-available keybinding reference

### Success Criteria

- Single binary, zero runtime dependencies (beyond `termux-api`)
- Compiles with pure Go (no CGO)
- Launches in under 1 second
- Works on any aarch64 Android device with Termux
- Database auto-creates and self-migrates on first run

---

## Module Path

```
github.com/Mr-Dark-debug/termnode
```

## Project Location

```
/storage/emulated/0/Projects/termnode/
```

All files live exclusively inside this directory.

---

## Operational Rules

1. All project files remain inside `/storage/emulated/0/Projects/termnode/`
2. No CGO — all dependencies must be pure Go
3. Database auto-migrates on startup via embedded Goose migrations
4. PID files for managed services stored at `$HOME/.termnode/*.pid`
5. Errors displayed in TUI, never via `fmt.Println`
6. Screen models implement `Init()`, `Update()`, `View()` (Bubble Tea pattern)
7. All async work via `tea.Cmd` — never touch UI state from goroutines

---

## Implementation Timeline

| Phase | Description | Status |
|-------|-------------|--------|
| Phase 1 | Project init + skeleton TUI | Done |
| Phase 2 | Database layer (SQLite + migrations) | Done |
| Phase 3 | Hardware dashboard (termux-api polling) | Done |
| Phase 4 | Service manager (wake-lock, sshd, file server) | Done |
| Phase 5 | IoT bridge + log viewer | Done |
| Phase 6 | Polish, Makefile, graceful shutdown | Done |

---

## Next Steps

After the initial build compiles successfully:

1. Run `go mod tidy` to resolve all dependencies
2. Run `go build -o termnode .` to compile
3. Test each tab manually
4. Verify `termux-api` integration with real device data
5. Test IoT webhook with `curl`
6. Push to GitHub and set up CI/CD
