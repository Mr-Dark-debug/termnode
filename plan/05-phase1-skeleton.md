# Phase 1: Project Init + Skeleton TUI

## Goal

Create the project structure, initialize Go module, write CLAUDE.md, and build a running TUI with 4 switchable tabs (placeholder screens). The TUI should compile and show a tab bar that responds to F1-F4 keys.

---

## Steps

### Step 1: Create Project Directory
```bash
mkdir -p /storage/emulated/0/Projects/termnode/internal/{app,daemon,hardware,iot,db,screen,theme}
mkdir -p /storage/emulated/0/Projects/termnode/migrations
```

### Step 2: Initialize Go Module
```bash
cd /storage/emulated/0/Projects/termnode
go mod init github.com/Mr-Dark-debug/termnode
```

### Step 3: Write CLAUDE.md
Project documentation for AI-assisted development. Contains:
- Project overview
- Architecture summary (Bubble Tea, Elm Architecture)
- Key conventions (all code in internal/, screen models implement Init/Update/View)
- Operational rules (all files in this directory, no CGO)
- Build & run instructions
- Dependency list
- Testing approach

### Step 4: Write Theme Package
File: `internal/theme/styles.go`

Defines:
- `Theme` struct with color constants (Primary violet, Accent cyan, Success green, Danger red, Muted gray)
- `ActiveTabStyle` — bold white on violet background
- `InactiveTabStyle` — gray text on dark gray background
- `PanelStyle` — rounded border with padding
- `PanelTitleStyle` — bold accent color
- `LabelStyle` — muted color for labels
- `ValueStyle` — bright white for values
- `StatusOnStyle` — bold green for active states
- `StatusOffStyle` — red for inactive states

**Critical**: This package has NO imports from other internal packages. It is a leaf dependency that breaks potential circular imports.

### Step 5: Write App Package
Files:
- `internal/app/messages.go` — errMsg type
- `internal/app/keys.go` — KeyBindings struct with DefaultKeyBindings() returning F1-F4, q, j/k, Enter, r
- `internal/app/styles.go` — Re-exports from theme/ package
- `internal/app/app.go` — Root Model:
  - Screen enum: ScreenDashboard, ScreenServices, ScreenIoTLog, ScreenHelp
  - Model struct with currentScreen, dashboard, services, iotLog, help, keys, width, height, bridge, poller, version
  - New(repo, port, ver) constructor
  - Init() — starts poller and bridge listener
  - Update() — intercepts F1-F4 for tab switching, hardware.UpdateMsg for dashboard, iot.IoTEventMsg for IoT log
  - View() — renders tab bar + active screen

### Step 6: Write Screen Package (Placeholders)
Files:
- `internal/screen/dashboard.go` — DashboardModel with SetHardware(), placeholder View
- `internal/screen/services.go` — ServicesModel with service list, cursor navigation, toggle logic
- `internal/screen/iotlog.go` — IoTLogModel with viewport, AddEvent()
- `internal/screen/help.go` — HelpModel with static keybindings and features list

### Step 7: Write Main Entry Point
File: `main.go`

- Parse flags: -db, -port, -debug, -version
- Open SQLite database via db.Open()
- Create repository
- Create app model via app.New()
- Run Bubble Tea program with alt screen

### Step 8: Stub Packages
Files with minimal implementations:
- `internal/daemon/daemon.go` — EnableWakeLock(), DisableWakeLock()
- `internal/daemon/process.go` — Manager with Start/Stop/Status stubs
- `internal/hardware/` — All types and stub poll functions
- `internal/iot/bridge.go` — Bridge with ListenCmd() stub
- `internal/db/` — Open(), Repository with stub methods

---

## Verification

```bash
cd /storage/emulated/0/Projects/termnode
go mod tidy
go build -o termnode .
./termnode
```

Expected:
- TUI launches in alternate screen
- Tab bar visible at top with 4 tabs
- Pressing F1-F4 switches active tab
- Pressing q quits cleanly
- No crash on startup

---

## Files Created (This Phase)

| File | Purpose |
|------|---------|
| `CLAUDE.md` | Project documentation |
| `go.mod` | Module definition |
| `main.go` | Entry point |
| `Makefile` | Build shortcuts |
| `internal/theme/styles.go` | Visual theme |
| `internal/app/app.go` | Root model |
| `internal/app/keys.go` | Keybindings |
| `internal/app/messages.go` | Message types |
| `internal/app/styles.go` | Style re-exports |
| `internal/screen/dashboard.go` | Dashboard screen |
| `internal/screen/services.go` | Services screen |
| `internal/screen/iotlog.go` | IoT log screen |
| `internal/screen/help.go` | Help screen |

---

## Common Issues

### /tmp not available on Android
The Bash tool in Claude Code requires `/tmp` which doesn't exist in Termux. Workaround:
```bash
mkdir -p $PREFIX/tmp
ln -s $PREFIX/tmp /tmp  # requires root or specific setup
```
Or use the Write tool to create files directly.

### go.sum not generated
Run `go mod tidy` to generate the go.sum file with all dependency checksums.
