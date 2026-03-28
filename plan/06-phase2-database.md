# Phase 2: Database Layer

## Goal

Implement the SQLite database layer with pure Go (modernc.org/sqlite), embedded migrations, and a repository pattern for IoT event CRUD operations.

---

## Why SQLite?

- Single-file database, no server process
- Perfect for embedded/edge devices (Android phones)
- Pure Go implementation available (no CGO)
- WAL mode for concurrent read/write performance
- Auto-creates database file if it doesn't exist

---

## Schema

```sql
CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    topic TEXT NOT NULL,
    payload TEXT NOT NULL,
    source TEXT NOT NULL DEFAULT 'http',
    timestamp DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_events_topic ON events(topic);
CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp);
```

### Column Descriptions

| Column | Type | Purpose |
|--------|------|---------|
| `id` | INTEGER PK AUTO | Auto-incrementing event ID |
| `topic` | TEXT NOT NULL | Webhook topic (e.g., "temperature", "humidity") |
| `payload` | TEXT NOT NULL | Raw request body (JSON string) |
| `source` | TEXT | Origin: "http" or "mqtt" |
| `timestamp` | DATETIME | When the event was received |

### Indexes

- `idx_events_topic` — Fast lookups by topic (filter sensor types)
- `idx_events_timestamp` — Fast ordering and time-range queries

---

## Files

### internal/db/db.go

**Responsibilities:**
- Open SQLite connection with WAL mode and foreign key pragmas
- Ensure parent directory exists (`$HOME/.termnode/`)
- Embed and execute SQL migrations via `//go:embed migrations/*.sql`
- Return `*sql.DB` ready for use

**Key Implementation:**
```go
//go:embed migrations/*.sql
var migrationsFS embed.FS

func Open(path string) (*sql.DB, error) {
    // 1. os.MkdirAll for parent directory
    // 2. sql.Open("sqlite", dsn) with WAL + FK pragmas
    // 3. migrationsFS.ReadFile("migrations/001_init.sql")
    // 4. db.Exec(string(migration))
    // 5. return db
}
```

**DSN Format:**
```
file:/path/to/termnode.db?_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)
```

### internal/db/models.go

```go
type IoTEvent struct {
    ID        int64
    Topic     string
    Payload   string
    Source    string    // "http" or "mqtt"
    Timestamp time.Time
}
```

### internal/db/repository.go

**Repository Pattern:**
```go
type Repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) *Repository
func (r *Repository) Insert(event IoTEvent) (int64, error)
func (r *Repository) Recent(limit int) ([]IoTEvent, error)
func (r *Repository) ByTopic(topic string, limit int) ([]IoTEvent, error)
func (r *Repository) Count() (int64, error)
func (r *Repository) Purge(olderThan time.Duration) (int64, error)
```

### internal/db/migrations/001_init.sql

The actual SQL migration file. Embedded into the binary at compile time via `//go:embed`.

---

## Design Decisions

### Why embed migrations instead of external files?

Using `//go:embed migrations/*.sql` means:
- Migration SQL is compiled into the binary
- No external files needed at runtime
- Database self-migrates on every startup
- Single-binary distribution remains intact

### Why not use Goose as a library?

While Goose provides a robust migration framework, for a single-table database the overhead isn't justified. Direct `db.Exec()` of embedded SQL files is simpler and has fewer dependencies. We can migrate to Goose later if the schema grows complex.

### Why not use an ORM?

Go's `database/sql` interface with hand-written queries is:
- More performant
- Easier to debug
- More explicit about what SQL runs
- No reflection magic or hidden queries

### Why WAL mode?

Write-Ahead Logging allows:
- Concurrent reads while writing (no reader-writer locks)
- Better performance for the IoT bridge (many small writes)
- Atomic transactions
- Slight increase in file count (WAL + SHM files) but worth it

---

## Testing

Tests use SQLite's `:memory:` DSN for fast, isolated tests:

```go
func TestRepository(t *testing.T) {
    db, _ := sql.Open("sqlite", ":memory:")
    defer db.Close()
    // Run migration
    // Test Insert, Recent, ByTopic, Count, Purge
}
```

---

## Verification

```bash
# Unit tests
go test ./internal/db/...

# Manual test
go run . &
# In another terminal:
curl -X POST http://localhost:8080/webhook/test -d '{"hello":"world"}'
# Kill the process
# Check database:
sqlite3 ~/.termnode/termnode.db "SELECT * FROM events;"
```
