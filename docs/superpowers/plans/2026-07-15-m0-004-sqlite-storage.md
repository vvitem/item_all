# M0-004 SQLite Storage Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build ItemAll's pure-Go, fail-closed SQLite WAL storage foundation with embedded append-only migrations, stable error projection, serialized writes, and safe Wails startup integration.

**Architecture:** `internal/storage` defines SQLite-independent repository contracts and safe errors. `internal/storage/sqlite` owns paths, DSNs, connection pooling, PRAGMAs, migration history, transactions, and SQLite error classification. The Wails layer receives only a readiness summary; `*sql.DB`, `*sql.Tx`, DSNs, SQL text, local usernames, and raw driver errors stay inside storage.

**Tech Stack:** Go 1.26.5, `database/sql`, `modernc.org/sqlite` v1.53.0, `modernc.org/libc` v1.73.4, embedded SQL, Wails v2.13.0, React 19.1.0, TypeScript 5.6.3, Vitest, Ubuntu 24.04, Windows Server 2025.

## Global Constraints

- Base implementation work on design merge commit `e31f4b31f840bc8962be09d3f9f6b90f7a4da8b6`.
- Use `modernc.org/sqlite` exactly `v1.53.0` and `modernc.org/libc` exactly `v1.73.4`.
- Do not introduce CGO, ORM, goose, golang-migrate, generic CRUD repositories, or another migration framework.
- Production defaults: WAL, `foreign_keys=ON`, `busy_timeout=5000`, `synchronous=NORMAL`, `_txlock=immediate`, `_dqs=false`, `_time_format=sqlite`, `_timezone=UTC`.
- Pool defaults: `MaxOpenConns=4`, `MaxIdleConns=4`, `ConnMaxLifetime=0`.
- Use OS-standard local data paths; reject Windows UNC and mapped network-drive locations.
- Migrations are six-digit, contiguous, append-only SQL files with LF endings and SHA-256 checksums.
- The fixed runner bootstrap creates only `schema_migrations`; versioned migrations never create or mutate the ledger.
- Each pending migration runs in its own immediate transaction.
- Initialization is fail-closed; never delete, rename, replace, relocate, or automatically rebuild an existing database.
- No business tables or SSH, SFTP, remote database, Operation Bus, AI, credential, telemetry, backup, recovery, portable-mode, or network storage capability enters M0-004.
- `GetAppInfo` remains the only Wails binding.
- Follow TDD: focused failing test, verify RED, minimum implementation, verify GREEN, commit.
- Keep required check names unchanged: `quality`, `generated-and-security`, `windows-build`.

---

## File Map

### New contracts

- `internal/storage/errors.go` — stable codes, wrapped errors, safe projection.
- `internal/storage/querier.go` — `Row` and `Querier` contracts.
- `internal/storage/transaction.go` — `Transactor` contract.

### New SQLite implementation

- `internal/storage/sqlite/config.go` and `config_test.go` — defaults, validation, URI DSN.
- `internal/storage/sqlite/path.go` — default path orchestration and directory preparation.
- `internal/storage/sqlite/path_windows.go`, `path_windows_test.go` — LocalAppData and network-drive checks.
- `internal/storage/sqlite/path_darwin.go` — Application Support path.
- `internal/storage/sqlite/path_linux.go`, `path_linux_test.go` — XDG data path and fallback.
- `internal/storage/sqlite/errors.go` — SQLite result-code classification.
- `internal/storage/sqlite/migrate.go`, `migrate_test.go` — ledger, discovery, history, apply loop.
- `internal/storage/sqlite/migrations/migrations.go` — embedded migration filesystem.
- `internal/storage/sqlite/migrations/000001_storage_meta.sql` — storage-only metadata table.
- `internal/storage/sqlite/health.go` — multi-connection PRAGMA verification.
- `internal/storage/sqlite/store.go`, `store_test.go` — open, ready, query, execute, close lifecycle.
- `internal/storage/sqlite/transaction.go`, `transaction_test.go` — callback transactions.
- `internal/storage/sqlite/concurrency_test.go` — WAL read/write and lock-timeout tests.

### Application integration

- `bootstrap.go`, `bootstrap_test.go` — storage runtime startup and close.
- `main.go` — bootstrap storage and register shutdown.
- `app.go`, `app_test.go` — safe readiness DTO.
- `frontend/src/types.ts`, `frontend/src/App.tsx`, `frontend/src/App.test.tsx` — safe storage error state.
- `frontend/wailsjs/**` — regenerated models/bindings only.

### Policy and docs

- `go.mod`, `go.sum`, `.gitattributes`.
- `scripts/ci/dependencies.go`, `dependencies_test.go`.
- `scripts/ci/boundary.go`, `boundary_test.go`.
- `.github/workflows/ci.yml`.
- `docs/development/storage.md`, `README.md`.
- Seven project progress documents under `docs/project-management/04-progress` and `02-roadmap/backlog.md`.

---

### Task 1: Pin SQLite dependencies and narrow CI policy exceptions

**Files:**
- Modify: `go.mod`, `go.sum`, `.gitattributes`
- Modify: `scripts/ci/dependencies.go`, `scripts/ci/dependencies_test.go`
- Modify: `scripts/ci/boundary.go`, `scripts/ci/boundary_test.go`

**Interfaces:**
- Produces `expectedSQLiteVersion = "v1.53.0"` and `expectedSQLiteLibcVersion = "v1.73.4"`.
- Allows `database/sql` only under `internal/storage/`.
- Allows `modernc.org/sqlite` only under `internal/storage/sqlite/`.
- Allows only `os.LookupEnv("XDG_DATA_HOME")` in `internal/storage/sqlite/path_linux.go`.

- [ ] **Step 1: Add failing dependency-policy tests**

Append to `scripts/ci/dependencies_test.go`:

```go
func TestDependencyPolicyRequiresExactSQLiteVersions(t *testing.T) {
	root := dependencyFixture(t)
	if violations := checkDependencies(root); len(violations) != 0 {
		t.Fatalf("expected exact SQLite versions, got %s", violations.Error())
	}
}

func TestDependencyPolicyRejectsSQLiteVersionDrift(t *testing.T) {
	root := dependencyFixture(t)
	text, err := readText(filepath.Join(root, "go.mod"))
	if err != nil { t.Fatal(err) }
	text = strings.Replace(text, "modernc.org/sqlite v1.53.0", "modernc.org/sqlite v1.52.0", 1)
	writeFixture(t, filepath.Join(root, "go.mod"), text)
	assertViolation(t, checkDependencies(root), "sqlite-version")
}

func TestDependencyPolicyRejectsSQLiteLibcVersionDrift(t *testing.T) {
	root := dependencyFixture(t)
	text, err := readText(filepath.Join(root, "go.mod"))
	if err != nil { t.Fatal(err) }
	text = strings.Replace(text, "modernc.org/libc v1.73.4", "modernc.org/libc v1.73.3", 1)
	writeFixture(t, filepath.Join(root, "go.mod"), text)
	assertViolation(t, checkDependencies(root), "sqlite-libc-version")
}
```

Update the fixture `go.mod` to include:

```go
require (
	github.com/wailsapp/wails/v2 v2.13.0
	modernc.org/sqlite v1.53.0
	modernc.org/libc v1.73.4 // indirect
)
```

- [ ] **Step 2: Add failing boundary tests**

Append to `scripts/ci/boundary_test.go`:

```go
func TestBoundaryPolicyAllowsDatabaseSQLOnlyInsideStorage(t *testing.T) {
	root := boundaryFixture(t)
	writeFixture(t, filepath.Join(root, "internal", "storage", "querier.go"), "package storage\nimport \"database/sql\"\nvar _ sql.Result\n")
	if violations := checkBoundary(root); len(violations) != 0 {
		t.Fatalf("expected storage import to pass, got %s", violations.Error())
	}
}

func TestBoundaryPolicyRejectsDatabaseSQLOutsideStorage(t *testing.T) {
	root := boundaryFixture(t)
	writeFixture(t, filepath.Join(root, "app.go"), "package main\nimport \"database/sql\"\nvar _ sql.Result\n")
	assertViolation(t, checkBoundary(root), "scope-import")
}

func TestBoundaryPolicyAllowsModernSQLiteOnlyInsideSQLiteStorage(t *testing.T) {
	root := boundaryFixture(t)
	writeFixture(t, filepath.Join(root, "internal", "storage", "sqlite", "store.go"), "package sqlite\nimport _ \"modernc.org/sqlite\"\n")
	if violations := checkBoundary(root); len(violations) != 0 {
		t.Fatalf("expected driver import to pass, got %s", violations.Error())
	}
}

func TestBoundaryPolicyAllowsOnlyStandardXDGDataHomeLookup(t *testing.T) {
	root := boundaryFixture(t)
	writeFixture(t, filepath.Join(root, "internal", "storage", "sqlite", "path_linux.go"), "package sqlite\nimport \"os\"\nfunc root() string { v, _ := os.LookupEnv(\"XDG_DATA_HOME\"); return v }\n")
	if violations := checkBoundary(root); len(violations) != 0 {
		t.Fatalf("expected XDG lookup to pass, got %s", violations.Error())
	}
}

func TestBoundaryPolicyRejectsOtherEnvironmentReadInLinuxPath(t *testing.T) {
	root := boundaryFixture(t)
	writeFixture(t, filepath.Join(root, "internal", "storage", "sqlite", "path_linux.go"), "package sqlite\nimport \"os\"\nfunc root() string { return os.Getenv(\"HOME\") }\n")
	assertViolation(t, checkBoundary(root), "environment-read")
}
```

- [ ] **Step 3: Verify RED**

```bash
go test ./scripts/ci -run 'TestDependencyPolicyRequiresExactSQLite|TestDependencyPolicyRejectsSQLite|TestBoundaryPolicyAllowsDatabaseSQL|TestBoundaryPolicyRejectsDatabaseSQL|TestBoundaryPolicyAllowsModernSQLite|TestBoundaryPolicyAllowsOnlyStandardXDG|TestBoundaryPolicyRejectsOtherEnvironment' -v
```

Expected: FAIL because the rules are absent.

- [ ] **Step 4: Implement dependency rules**

Add to the constants in `scripts/ci/dependencies.go`:

```go
expectedSQLiteVersion     = "v1.53.0"
expectedSQLiteLibcVersion = "v1.73.4"
```

Add to `required` in `checkGoModule`:

```go
"sqlite-version":      "modernc.org/sqlite " + expectedSQLiteVersion,
"sqlite-libc-version": "modernc.org/libc " + expectedSQLiteLibcVersion,
```

- [ ] **Step 5: Implement scoped boundary helpers**

Replace the current environment/capability loops with calls to:

```go
func checkEnvironmentReads(rel, text string, violations *Violations) {
	sanitized := text
	if rel == "internal/storage/sqlite/path_linux.go" {
		sanitized = strings.ReplaceAll(sanitized, `os.LookupEnv("XDG_DATA_HOME")`, "")
	}
	for _, token := range []string{"os.Getenv(", "os.LookupEnv(", "process.env", "import.meta.env"} {
		if strings.Contains(sanitized, token) {
			violations.Add(rel, "environment-read", "only the standard XDG_DATA_HOME lookup is approved")
		}
	}
}

func checkCapabilityImports(rel, text string, violations *Violations) {
	rules := []struct{ token, allowedPrefix string }{
		{token: "golang.org/x/crypto/ssh"},
		{token: "database/sql", allowedPrefix: "internal/storage/"},
		{token: "modernc.org/sqlite", allowedPrefix: "internal/storage/sqlite/"},
		{token: "github.com/openai/"},
		{token: "go.opentelemetry.io/"},
	}
	for _, rule := range rules {
		if !strings.Contains(text, rule.token) { continue }
		if rule.allowedPrefix != "" && strings.HasPrefix(rel, rule.allowedPrefix) { continue }
		violations.Add(rel, "scope-import", "capability dependency is outside its approved package boundary: "+rule.token)
	}
}
```

Keep every existing forbidden capability directory unchanged.

- [ ] **Step 6: Add dependencies and LF rule**

```bash
go get modernc.org/sqlite@v1.53.0
go mod edit -require=modernc.org/libc@v1.73.4
go mod tidy
printf '\n*.sql text eol=lf\n' >> .gitattributes
go list -m modernc.org/sqlite modernc.org/libc
```

Expected module output:

```text
modernc.org/sqlite v1.53.0
modernc.org/libc v1.73.4
```

- [ ] **Step 7: Verify GREEN and commit**

```bash
go test ./scripts/ci -v
go run ./scripts/ci dependencies
go run ./scripts/ci boundary
go mod verify
git add go.mod go.sum .gitattributes scripts/ci
git commit -m "build: approve pure Go SQLite dependency"
```

---

### Task 2: Define storage contracts and safe errors

**Files:**
- Create: `internal/storage/errors.go`, `errors_test.go`, `querier.go`, `transaction.go`
- Consume amendment: `docs/superpowers/specs/2026-07-15-m0-004-query-row-amendment.md`

**Interfaces:**

```go
type Row interface { Scan(dest ...any) error }
type Querier interface {
    ExecContext(context.Context, string, ...any) (sql.Result, error)
    QueryContext(context.Context, string, ...any) (*sql.Rows, error)
    QueryRowContext(context.Context, string, ...any) Row
}
type Transactor interface {
    WithTx(context.Context, func(context.Context, Querier) error) error
}
```

- [ ] **Step 1: Write failing safe-projection tests**

Create `internal/storage/errors_test.go`:

```go
package storage

import (
	"errors"
	"strings"
	"testing"
)

func TestProjectReturnsSafeStorageError(t *testing.T) {
	err := Wrap(CodePermissionDenied, "open", errors.New(`C:\Users\private\itemall.db`), true)
	got := Project(err)
	if got.Code != string(CodePermissionDenied) { t.Fatalf("Code = %q", got.Code) }
	if got.SafeMessage != "无法访问本地数据目录，请检查目录权限。" { t.Fatalf("message = %q", got.SafeMessage) }
	if !got.Retryable { t.Fatal("Retryable = false") }
	if strings.Contains(got.SafeMessage, "Users") || strings.Contains(got.SafeMessage, "itemall.db") { t.Fatalf("leak: %q", got.SafeMessage) }
}

func TestProjectUnknownErrorIsGeneric(t *testing.T) {
	got := Project(errors.New("raw driver detail"))
	if got.Code != string(CodeOpenFailed) || got.SafeMessage != "无法打开本地数据存储。" || got.Retryable {
		t.Fatalf("projection = %+v", got)
	}
}

func TestIsCodeTraversesWrappedError(t *testing.T) {
	if !IsCode(Wrap(CodeBusy, "write", errors.New("locked"), true), CodeBusy) { t.Fatal("missing code") }
}
```

- [ ] **Step 2: Verify RED**

```bash
go test ./internal/storage -v
```

Expected: package/types missing.

- [ ] **Step 3: Implement contracts**

Create `internal/storage/querier.go` and `transaction.go` exactly from the interfaces above.

Create `internal/storage/errors.go` with these codes:

```go
const (
	CodePathUnavailable Code = "storage_path_unavailable"
	CodePermissionDenied Code = "storage_permission_denied"
	CodeOpenFailed Code = "storage_open_failed"
	CodeBusy Code = "storage_busy"
	CodeCorrupt Code = "storage_corrupt"
	CodePragmaInvalid Code = "storage_pragma_invalid"
	CodeMigrationInvalid Code = "storage_migration_invalid"
	CodeMigrationFailed Code = "storage_migration_failed"
	CodeMigrationChecksumMismatch Code = "storage_migration_checksum_mismatch"
	CodeSchemaTooNew Code = "storage_schema_too_new"
	CodeTransactionFailed Code = "storage_transaction_failed"
	CodeCommitFailed Code = "storage_commit_failed"
	CodeClosed Code = "storage_closed"
)
```

Implement:

```go
type Error struct {
	Code Code
	Operation string
	Retryable bool
	Cause error
}
func (e *Error) Error() string { return fmt.Sprintf("%s: %s", e.Code, e.Operation) }
func (e *Error) Unwrap() error { return e.Cause }
func Wrap(code Code, operation string, cause error, retryable bool) error
func IsCode(err error, code Code) bool

type SafeError struct {
	Code string `json:"code"`
	SafeMessage string `json:"safeMessage"`
	Retryable bool `json:"retryable"`
}
func Project(err error) SafeError
```

Use fixed Chinese safe messages; never interpolate `Cause`, SQL, path, or DSN.

- [ ] **Step 4: Verify GREEN and commit**

```bash
gofmt -w internal/storage
go test ./internal/storage -v
git add internal/storage docs/superpowers/specs/2026-07-15-m0-004-query-row-amendment.md
git commit -m "feat: define storage contracts and safe errors"
```

---

### Task 3: Implement OS paths and exact DSN configuration

**Files:**
- Create: `internal/storage/sqlite/config.go`, `config_test.go`, `path.go`
- Create: `path_windows.go`, `path_windows_test.go`, `path_darwin.go`, `path_linux.go`, `path_linux_test.go`

**Interfaces:**

```go
type Config struct {
    Path string
    BusyTimeout time.Duration
    MaxOpenConns int
    MaxIdleConns int
}
func DefaultConfig(path string) Config
func DefaultDatabasePath() (string, error)
func buildDSN(Config) (string, error)
func prepareDatabasePath(string) (string, error)
```

- [ ] **Step 1: Write failing config tests**

Test exact defaults and parse the URI query to assert:

```text
mode=rwc
_pragma=journal_mode(WAL)
_pragma=foreign_keys(1)
_pragma=busy_timeout(5000)
_pragma=synchronous(NORMAL)
_txlock=immediate
_dqs=false
_time_format=sqlite
_timezone=UTC
```

Use a table test that rejects zero timeout, zero open connections, and idle connections greater than open connections.

- [ ] **Step 2: Write platform path tests**

Linux tests:

```go
func TestPlatformDataRootUsesAbsoluteXDGDataHome(t *testing.T)
func TestPlatformDataRootIgnoresRelativeXDGDataHome(t *testing.T)
```

Windows test:

```go
func TestIsNetworkPathRejectsUNC(t *testing.T)
```

Use `t.Setenv`, `t.TempDir`, and build tags; never write to real application data directories.

- [ ] **Step 3: Verify RED**

```bash
go test ./internal/storage/sqlite -run 'TestDefaultConfig|TestBuildDSN|TestConfig|TestPlatformDataRoot|TestIsNetworkPath' -v
```

- [ ] **Step 4: Implement config and DSN**

`DefaultConfig` returns 5 seconds, 4 open, 4 idle. `buildDSN` uses `url.URL`/`url.Values`, converts Windows drive paths to file URIs, and emits the exact query above.

- [ ] **Step 5: Implement common path preparation**

`prepareDatabasePath` must:

1. resolve the default only when `Path` is empty;
2. convert to an absolute path;
3. call platform `isNetworkPath`;
4. create the parent directory with `0700` where meaningful;
5. create/close/remove a temporary write probe;
6. map failures to `storage_path_unavailable` or `storage_permission_denied`.

- [ ] **Step 6: Implement platform files**

Windows:

```go
root, err := windows.KnownFolderPath(windows.FOLDERID_LocalAppData, 0)
return filepath.Join(root, "ItemAll", "data"), err
```

Reject paths beginning `\\` or `//`; for drive letters, call `windows.GetDriveType` and reject `windows.DRIVE_REMOTE`.

macOS:

```go
home, err := os.UserHomeDir()
return filepath.Join(home, "Library", "Application Support", "ItemAll", "data"), err
```

Linux:

```go
if root, ok := os.LookupEnv("XDG_DATA_HOME"); ok && filepath.IsAbs(root) {
    return filepath.Join(root, "itemall", "data"), nil
}
home, err := os.UserHomeDir()
return filepath.Join(home, ".local", "share", "itemall", "data"), err
```

- [ ] **Step 7: Verify GREEN and commit**

```bash
gofmt -w internal/storage/sqlite
go test ./internal/storage/sqlite -run 'TestDefaultConfig|TestBuildDSN|TestConfig|TestPlatformDataRoot' -v
go run ./scripts/ci boundary
git add internal/storage/sqlite
git commit -m "feat: resolve local SQLite paths and configuration"
```

---

### Task 4: Build the append-only migration runner

**Files:**
- Create: `internal/storage/sqlite/migrations/migrations.go`
- Create: `internal/storage/sqlite/migrations/000001_storage_meta.sql`
- Create: `internal/storage/sqlite/migrate.go`, `migrate_test.go`

**Interfaces:**

```go
func loadMigrations(fs.FS) ([]migration, error)
func runMigrations(context.Context, *sql.DB, fs.FS, func() time.Time) error
```

- [ ] **Step 1: Add embedded files**

`migrations.go`:

```go
package migrations
import "embed"
//go:embed *.sql
var FS embed.FS
```

`000001_storage_meta.sql`:

```sql
CREATE TABLE storage_meta (
    key        TEXT PRIMARY KEY,
    value      TEXT NOT NULL,
    updated_at TEXT NOT NULL
) STRICT;
```

- [ ] **Step 2: Write failing discovery tests**

Use `testing/fstest.MapFS` to cover:

- valid versions 1 and 2;
- version gap 1 then 3;
- malformed filename;
- unsupported extension;
- empty SQL;
- duplicate name;
- deterministic SHA-256 checksum.

- [ ] **Step 3: Implement discovery**

Use regex:

```go
var migrationFilePattern = regexp.MustCompile(`^([0-9]{6})_([a-z][a-z0-9_]*)\.sql$`)
```

Sort by integer version; require version `index+1`; hash raw embedded bytes with SHA-256.

- [ ] **Step 4: Write failing history/apply tests**

Cover:

```go
func TestRunMigrationsFreshAndRepeated(t *testing.T)
func TestRunMigrationsUpgradesFromPreviousVersion(t *testing.T)
func TestRunMigrationsRejectsChecksumMismatch(t *testing.T)
func TestRunMigrationsRejectsChangedName(t *testing.T)
func TestRunMigrationsRejectsSchemaTooNew(t *testing.T)
func TestRunMigrationsRejectsInvalidLedgerSchema(t *testing.T)
func TestRunMigrationsRollsBackFailingVersion(t *testing.T)
```

- [ ] **Step 5: Implement fixed ledger bootstrap**

Use exactly:

```sql
CREATE TABLE IF NOT EXISTS schema_migrations (
    version    INTEGER PRIMARY KEY,
    name       TEXT NOT NULL,
    checksum   TEXT NOT NULL,
    applied_at TEXT NOT NULL
) STRICT;
```

Immediately verify `PRAGMA table_info(schema_migrations)` returns exactly `version,name,checksum,applied_at` in that order.

- [ ] **Step 6: Implement history verification**

Read records ordered by version. Reject:

- applied gaps;
- database version greater than embedded version count;
- changed name;
- changed checksum.

Map to `storage_migration_invalid`, `storage_schema_too_new`, or `storage_migration_checksum_mismatch` as appropriate.

- [ ] **Step 7: Implement one transaction per migration**

For each pending migration:

```go
tx, err := db.BeginTx(ctx, nil)
_, err = tx.ExecContext(ctx, string(item.SQL))
_, err = tx.ExecContext(ctx,
    `INSERT INTO schema_migrations(version,name,checksum,applied_at) VALUES(?,?,?,?)`,
    item.Version, item.Name, item.Checksum, now().UTC().Format(time.RFC3339Nano),
)
err = tx.Commit()
```

Rollback the current version on any failure and stop startup. Do not custom-split SQL by semicolons.

- [ ] **Step 8: Verify GREEN and commit**

```bash
gofmt -w internal/storage/sqlite
go test ./internal/storage/sqlite -run 'TestLoadMigrations|TestRunMigrations' -v
git add internal/storage/sqlite/migrate.go internal/storage/sqlite/migrate_test.go internal/storage/sqlite/migrations
git commit -m "feat: add append-only SQLite migrations"
```

---

### Task 5: Open, verify, classify, and close Store

**Files:**
- Create: `internal/storage/sqlite/errors.go`, `health.go`, `store.go`, `store_test.go`

**Interfaces:**

```go
func Open(context.Context, Config) (*Store, error)
func (s *Store) Path() string
func (s *Store) ExecContext(...) (sql.Result, error)
func (s *Store) QueryContext(...) (*sql.Rows, error)
func (s *Store) QueryRowContext(...) storage.Row
func (s *Store) Close() error
```

- [ ] **Step 1: Write failing integration tests**

Cover:

```go
func TestOpenCreatesDatabaseAndAppliesMigration(t *testing.T)
func TestOpenRepeatedIsIdempotent(t *testing.T)
func TestOpenEnforcesForeignKeys(t *testing.T)
func TestOpenVerifiesApprovedPragmasAcrossFourConnections(t *testing.T)
func TestOpenClassifiesCorruptDatabase(t *testing.T)
func TestStoreRejectsUseAfterClose(t *testing.T)
```

`QueryRowContext(...).Scan(...)` after close must return `storage_closed`.

- [ ] **Step 2: Implement SQLite classification**

Use `errors.As` with `*modernsqlite.Error` and primary result codes:

```text
3 PERM -> storage_permission_denied
5 BUSY -> storage_busy
6 LOCKED -> storage_busy
8 READONLY -> storage_permission_denied
11 CORRUPT -> storage_corrupt
14 CANTOPEN -> storage_open_failed
26 NOTADB -> storage_corrupt
```

Preserve raw errors only as internal causes.

- [ ] **Step 3: Implement `verifyPool`**

Acquire four `*sql.Conn` simultaneously so the pool cannot reuse one physical connection. On each connection verify:

```text
PRAGMA journal_mode == wal
PRAGMA foreign_keys == 1
PRAGMA busy_timeout == configured milliseconds
PRAGMA synchronous == 1
```

Map any mismatch to `storage_pragma_invalid`.

- [ ] **Step 4: Implement `Open`**

Order:

```text
validate Config
prepare path
build DSN
sql.Open("sqlite", dsn)
set pool
PingContext
verifyPool
runMigrations using migrations.FS
verifyPool again
return ready Store
```

Close the DB on every failure path; never return a partial Store.

- [ ] **Step 5: Implement Store methods and `errorRow`**

`Store` contains private `db`, path, config, write mutex, and atomic closed flag. `errorRow.Scan` returns a prebuilt `storage_closed` error. `ExecContext` takes the write mutex; reads do not.

- [ ] **Step 6: Verify GREEN and commit**

```bash
gofmt -w internal/storage/sqlite
go test ./internal/storage/sqlite -run 'TestOpen|TestStoreRejectsUseAfterClose' -v
git add internal/storage/sqlite/errors.go internal/storage/sqlite/health.go internal/storage/sqlite/store.go internal/storage/sqlite/store_test.go
git commit -m "feat: open and verify SQLite Store"
```

---

### Task 6: Add callback transactions and concurrency guarantees

**Files:**
- Create: `internal/storage/sqlite/transaction.go`, `transaction_test.go`, `concurrency_test.go`

**Interfaces:**

```go
func (s *Store) WithTx(context.Context, func(context.Context, storage.Querier) error) error
```

- [ ] **Step 1: Write failing transaction tests**

Cover:

```go
func TestWithTxCommitsOnSuccess(t *testing.T)
func TestWithTxRollsBackCallbackError(t *testing.T)
func TestWithTxRollsBackAndRepanics(t *testing.T)
func TestWithTxHonorsContextCancellation(t *testing.T)
func TestWithTxRejectsNilCallback(t *testing.T)
func TestWithTxRejectsNestedTransaction(t *testing.T)
func TestRunTransactionMapsCommitFailure(t *testing.T)
func TestRunTransactionPreservesCallbackAndRollbackErrors(t *testing.T)
```

Use a fake raw transaction for commit/rollback failure tests.

- [ ] **Step 2: Implement transaction adapter**

Wrap raw `*sql.Tx` so `QueryRowContext` returns `storage.Row`. Pass the adapter to callbacks; never expose `*sql.Tx`.

- [ ] **Step 3: Implement callback lifecycle**

Rules:

- closed Store -> `storage_closed`;
- nil callback -> `storage_transaction_failed`;
- context marker detects nested calls;
- lock `writeMu` before `BeginTx`;
- callback nil -> commit;
- callback error -> rollback and preserve callback error;
- rollback error -> `errors.Join(callbackErr, rollbackErr)`;
- panic -> rollback, then re-panic;
- commit failure -> `storage_commit_failed`.

- [ ] **Step 4: Write concurrency tests**

Cover:

```go
func TestWALAllowsReadWhileWriteTransactionIsOpen(t *testing.T)
func TestStoreSerializesApplicationWrites(t *testing.T)
func TestIndependentStoreMapsLockTimeoutToBusy(t *testing.T)
```

Use a 200 ms test timeout in the second Store; production remains 5 seconds.

- [ ] **Step 5: Verify GREEN and commit**

```bash
gofmt -w internal/storage/sqlite
go test ./internal/storage/sqlite -run 'TestWithTx|TestRunTransaction|TestWAL|TestStoreSerializes|TestIndependentStore' -v
go test -race ./internal/storage/...
git add internal/storage/sqlite/transaction.go internal/storage/sqlite/transaction_test.go internal/storage/sqlite/concurrency_test.go
git commit -m "feat: enforce SQLite transaction boundaries"
```

---

### Task 7: Integrate fail-closed startup into Wails and React

**Files:**
- Create: `bootstrap.go`, `bootstrap_test.go`
- Modify: `main.go`, `app.go`, `app_test.go`
- Modify: `frontend/src/types.ts`, `frontend/src/App.tsx`, `frontend/src/App.test.tsx`
- Regenerate: `frontend/wailsjs/go/main/App.js`, `App.d.ts`, `frontend/wailsjs/go/models.ts`

**Interfaces:**

```go
type storageRuntime struct { store *sqlite.Store; err error }
func bootstrapStorage(context.Context, storageOpener, sqlite.Config) storageRuntime
```

`AppInfo` adds optional:

```go
Error *storage.SafeError `json:"error,omitempty"`
```

- [ ] **Step 1: Write failing Go tests**

Add:

```go
func TestGetAppInfoReturnsSafeStorageFailure(t *testing.T)
func TestBootstrapStoragePreservesStartupFailure(t *testing.T)
func TestStorageRuntimeCloseIsSafeWithoutStore(t *testing.T)
```

Assert no Windows user path or driver detail enters `SafeMessage`.

- [ ] **Step 2: Implement bootstrap runtime**

Use a small injectable opener interface for tests. `main` calls:

```go
runtime := bootstrapStorage(context.Background(), defaultStorageOpener{}, sqlite.DefaultConfig(""))
app := NewApp(runtime.err)
```

Register `OnShutdown` to close the Store. Storage failure does not terminate Wails; it produces an error-state app.

- [ ] **Step 3: Update App DTO**

Statuses:

```go
const appStatusReady = "ready"
const appStatusStorageError = "storage_error"
```

`GetAppInfo` returns normal build metadata plus `storage.Project(startupErr)` when initialization failed.

- [ ] **Step 4: Write failing frontend test**

Create an `AppInfo` fixture:

```ts
{
  ...readyInfo,
  status: 'storage_error',
  error: {
    code: 'storage_permission_denied',
    safeMessage: '无法访问本地数据目录，请检查目录权限。',
    retryable: true,
  },
}
```

Assert the alert contains the safe message and code but not `Users`, `itemall.db`, SQL, or `SQLITE_`.

- [ ] **Step 5: Implement frontend error state**

Render:

```tsx
<section role="alert" className="status-card status-card--error">
  <h1>本地数据存储未能启动</h1>
  <p>{info.error?.safeMessage ?? '无法打开本地数据存储。'}</p>
  {info.error?.code ? <code>{info.error.code}</code> : null}
  <p>{info.error?.retryable
    ? '请检查目录权限或关闭其他 ItemAll 实例，然后重新启动应用。'
    : '为保护现有数据，应用不会自动重建数据库。'}</p>
</section>
```

Keep rejected/synchronous binding error handling unchanged.

- [ ] **Step 6: Regenerate and verify**

```bash
pnpm --dir frontend install --frozen-lockfile
pnpm --dir frontend build
go install github.com/wailsapp/wails/v2/cmd/wails@v2.13.0
GOFLAGS='-tags=webkit2_41' wails generate module
find frontend/wailsjs -type f -exec chmod 0644 {} +
go test . -v
pnpm --dir frontend test:run
pnpm --dir frontend typecheck
pnpm --dir frontend lint
go run ./scripts/ci boundary
```

Verify generated `App.js` still exports only `GetAppInfo`.

- [ ] **Step 7: Commit**

```bash
git add main.go app.go app_test.go bootstrap.go bootstrap_test.go frontend/src frontend/wailsjs
git commit -m "feat: fail closed when local storage cannot start"
```

---

### Task 8: Harden CI, document operation, and synchronize evidence

**Files:**
- Modify: `.github/workflows/ci.yml`
- Create: `docs/development/storage.md`
- Modify: `README.md`
- Modify: seven project-management status files

- [ ] **Step 1: Extend race coverage without renaming jobs**

Change the quality job race command to:

```yaml
go test -race ./internal/buildinfo ./internal/storage/...
```

Do not change Action SHAs, runner versions, permissions, triggers, or job names.

- [ ] **Step 2: Write `docs/development/storage.md`**

Document exact production locations, PRAGMAs, pool settings, migration naming/checksum rules, local commands, and strict non-destructive failure policy.

- [ ] **Step 3: Update README**

Add a Local Storage subsection linking the development document. State clearly that M0-004 adds infrastructure only, not business tables or remote database management.

- [ ] **Step 4: Open the implementation PR as Draft**

Title:

```text
feat: establish SQLite WAL storage foundation
```

Body must reference Issue #12, design PR #17, the amendment, exact dependency versions, and the three required checks. State all non-goals.

- [ ] **Step 5: Synchronize status truthfully**

Update all seven project files consistently:

```text
M0-004 = IN_REVIEW
DONE=3
IN_REVIEW=1
NOT_STARTED=46
M0-foundation completed=3, in review=1
overall completion=3/50=6%
M0-005 remains NOT_STARTED
```

Record the actual implementation PR URL and only evidence that has already succeeded. Do not mark M0-004 `DONE` before merge and the later status-sync PR.

- [ ] **Step 6: Run full local verification**

```bash
pnpm --dir frontend install --frozen-lockfile
pnpm --dir frontend build
go test ./scripts/ci -v
go run ./scripts/ci all
go mod tidy
git diff --exit-code -- go.mod go.sum
go mod verify
go vet -tags=webkit2_41 ./...
go test -tags=webkit2_41 ./...
go test -race ./internal/buildinfo ./internal/storage/...
pnpm --dir frontend test:run
pnpm --dir frontend typecheck
pnpm --dir frontend lint
GOFLAGS='-tags=webkit2_41' wails generate module
find frontend/wailsjs -type f -exec chmod 0644 {} +
git diff --exit-code -- frontend/wailsjs
```

Expected: every command exits 0.

- [ ] **Step 7: Commit docs and CI**

```bash
git add .github/workflows/ci.yml README.md docs/development/storage.md docs/project-management
git commit -m "docs: record M0-004 storage verification"
```

- [ ] **Step 8: Push and verify required checks**

Required outcomes:

```text
quality: success
generated-and-security: success
windows-build: success
```

Investigate exact failing steps; do not blindly rerun. For a code defect, add a focused regression test before the fix.

- [ ] **Step 9: Final scope review**

```bash
git diff --check main...HEAD
git diff --stat main...HEAD
git diff --name-only main...HEAD
```

Confirm absent:

```text
internal/ssh
internal/sftp
internal/database
internal/operation
internal/ai
internal/credential
internal/telemetry
workspace/asset/task/audit business tables
backup/recovery/portable mode
additional Wails methods
```

- [ ] **Step 10: Mark the implementation PR Ready**

Only after all three checks are green and the final diff is clean. Do not merge before code review.

---

## Final Acceptance Checklist

- [ ] Exact SQLite and libc versions are policy-enforced.
- [ ] `database/sql` and the modernc driver are restricted to approved packages.
- [ ] Windows, macOS, and Linux data paths are implemented.
- [ ] Windows UNC and mapped network drives are rejected.
- [ ] Four physical connections verify WAL, foreign keys, timeout, and NORMAL synchronous mode.
- [ ] Fixed bootstrap creates only `schema_migrations`.
- [ ] `000001_storage_meta.sql` is embedded, LF-normalized, and checksum recorded.
- [ ] Fresh, repeated, upgrade, invalid, checksum, too-new, and rollback paths pass.
- [ ] Transactions commit, roll back, re-panic, reject nesting, and map failures.
- [ ] Application writes serialize and WAL reads continue during an open write transaction.
- [ ] External lock timeout maps to `storage_busy`.
- [ ] Corruption maps to `storage_corrupt`.
- [ ] Use after close maps to `storage_closed`, including row `Scan`.
- [ ] `GetAppInfo` remains the only Wails binding.
- [ ] UI receives only code, safe message, and retryability.
- [ ] No domain tables or out-of-scope capabilities are introduced.
- [ ] `quality`, `generated-and-security`, and `windows-build` are green.
- [ ] M0-004 remains `IN_REVIEW` until implementation merge and completion sync.
