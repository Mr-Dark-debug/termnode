# Architecture Decision Records

This document records the key architectural decisions made during TermNode's design, along with the reasoning and trade-offs for each.

---

## ADR-001: Bubble Tea as TUI Framework

**Decision**: Use charmbracelet/bubbletea as the TUI framework.

**Why**: Bubble Tea implements The Elm Architecture (Model/Update/View), which provides:
- Explicit state management (no hidden state)
- Composable components (each screen is a self-contained model)
- Natural async handling via tea.Cmd (goroutine-based)
- Beautiful rendering via Lip Gloss
- Active development and community (Charm ecosystem)

**Alternatives Considered**:
- `tview` (rivo/tview): More widget-heavy, but less elegant state management
- `termui` (gizak/termui): Dashboard-focused, but less flexible for interactive UIs
- Raw `tcell`: Too low-level, too much boilerplate

**Trade-off**: Bubble Tea requires more boilerplate for simple layouts but scales much better for complex multi-screen apps.

---

## ADR-002: modernc.org/sqlite for Database

**Decision**: Use modernc.org/sqlite as the SQLite driver.

**Why**: This is the only pure Go SQLite implementation. It:
- Requires **zero CGO** (critical for Termux where CGO toolchain may be missing)
- Supports linux/arm64 (Android aarch64)
- Implements standard database/sql interface
- Tracks upstream SQLite faithfully (version 3.51.3)

**Alternatives Considered**:
- `mattn/go-sqlite3`: Industry standard but requires CGO (C compiler). Non-starter for Termux.
- Flat files (JSON/CSV): No query capability, no concurrency, no indexing.
- BadgerDB/BoltDB: Key-value only, no SQL queries. Overkill for simple event storage.

**Trade-off**: modernc.org/sqlite adds ~10MB to binary size due to the pure Go SQLite reimplementation. This is acceptable for a device with GBs of storage.

---

## ADR-003: Channel Bridge Pattern for Goroutine Communication

**Decision**: Use a buffered Go channel to bridge events from HTTP handler goroutines to the Bubble Tea main goroutine.

**Why**:
- Bubble Tea's Update() runs on a single goroutine (no mutexes needed for UI state)
- tea.Cmd provides the natural bridge: a function that blocks on a channel read
- Buffered channel (capacity 64) ensures HTTP handlers never block
- Drop semantics prevent memory issues under load

**Alternatives Considered**:
- Shared state with mutex: Error-prone, violates Bubble Tea's single-goroutine principle
- Callback functions: Less type-safe, harder to reason about
- Polling (DB queries on timer): Wasteful, introduces latency

**Trade-off**: Events can be dropped if the channel fills up (64 events). In practice, this is unlikely for typical IoT use cases (events arrive every few seconds).

---

## ADR-004: Go Standard Library HTTP Server (No External Router)

**Decision**: Use Go 1.22+'s enhanced `http.ServeMux` with method+path routing. No external router library.

**Why**: Go 1.22 introduced native support for:
- Method-based routing: `mux.HandleFunc("POST /webhook/{topic}", handler)`
- Path parameters: `r.PathValue("topic")`
- No additional dependency

**Alternatives Considered**:
- `chi`: Excellent router, but overkill for 2 endpoints
- `gorilla/mux`: Archived/project archived, not recommended
- `gin`: Full framework, way too heavy for a webhook receiver

**Trade-off**: Standard library routing is less feature-rich than chi (no middleware chaining, no regex patterns). For 2 endpoints, this is perfectly fine.

---

## ADR-005: PID-File-Based Process Tracking

**Decision**: Track managed processes via PID files in `$HOME/.termnode/`.

**Why**:
- Simple and well-understood pattern on Unix/Linux
- Survives TermNode restarts (can detect previously started services)
- Works within Termux's process model (no systemd, no D-Bus)
- Minimal code complexity

**Alternatives Considered**:
- In-memory-only tracking: Lost on restart, orphaned processes
- Socket activation: Not available in Termux
- Supervisor pattern (respawn on death): Overkill for this scope, added later as optional feature

**Trade-off**: PID files can become stale if the process dies abnormally (crash without cleanup). Mitigated by checking `/proc/<pid>/` existence before using stored PIDs.

---

## ADR-006: Go Native File Server (Not Python)

**Decision**: Use Go's `net/http.FileServer` instead of spawning Python's `http.server`.

**Why**:
- Zero additional dependencies (no Python needed in Termux)
- Single binary handles everything
- Better performance (Go HTTP server vs Python)
- Can be extended with authentication, logging, etc.
- Consistent with the Go-native philosophy

**Alternatives Considered**:
- `python3 -m http.server`: Requires Python, separate process to manage
- `busybox httpd`: Additional dependency, limited features
- Custom static file handler: Reimplements what stdlib already provides

**Trade-off**: Go's FileServer serves the entire home directory by default. For future security, should be configurable to serve only specific directories.

---

## ADR-007: Theme Package to Break Circular Imports

**Decision**: Extract all Lip Gloss styles into a separate `internal/theme/` package with no internal imports.

**Why**:
- Both `app` and `screen` packages need access to styles
- `app` imports `screen` (for sub-models)
- `screen` cannot import `app` (would create a cycle)
- Solution: Both import `theme`, which has no internal dependencies

**Dependency graph:**
```
app → screen → theme    ✓ (no cycle)
app → theme             ✓
theme → (nothing)       ✓ (leaf package)
```

**Alternatives Considered**:
- Inline styles in each screen: DRY violation, inconsistent theming
- Pass styles via constructor injection: Boilerplate-heavy, over-engineering
- Use a global registry: Less explicit, harder to test

**Trade-off**: One extra package to maintain. Worth it for clean dependency graph and consistent theming.

---

## ADR-008: MQTT Behind Build Tag

**Decision**: Make MQTT support optional via Go build tags (`//go:build mqtt`).

**Why**:
- MQTT client adds ~2MB to binary size
- Most users only need HTTP webhooks
- Default build should be as lean as possible
- Users who need MQTT can opt in: `go build -tags mqtt`

**How**:
```go
// internal/iot/mqtt.go
//go:build mqtt

package iot

import "github.com/eclipse-paho/paho.mqtt.golang"
// ... MQTT subscriber implementation
```

**Trade-off**: Users need to know about the build tag to use MQTT. Mitigated by documentation in README and help screen.

---

## ADR-009: Direct SQL Migration (No Goose Library)

**Decision**: Use `//go:embed` to embed SQL files and execute them directly with `db.Exec()`. Not using Goose as a library.

**Why**:
- Only one migration file with one table
- Goose adds a dependency for minimal benefit at this scale
- Direct `db.Exec()` is simpler and more transparent

**Alternatives Considered**:
- `pressly/goose`: Robust migration tool, but overkill for one table
- `golang-migrate`: Similar, also overkill
- Hardcoded SQL strings in Go: Harder to maintain than separate files

**Trade-off**: If the schema grows to 5+ migrations, we should migrate to Goose. For now, simplicity wins.

---

## ADR-010: main.go at Root (No cmd/ Directory)

**Decision**: Place `main.go` at the project root, not in `cmd/termnode/`.

**Why**:
- Single-binary project — no other commands
- `cmd/` adds a directory level with no benefit
- Go team's official docs say `cmd/` is a "common convention" not a requirement
- Cleaner file paths on a phone

**Trade-off**: If a second binary is added later (e.g., a CLI tool), we'd need to restructure to `cmd/`. Acceptable risk given the project scope.
