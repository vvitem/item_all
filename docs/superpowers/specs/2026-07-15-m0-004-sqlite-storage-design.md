# M0-004 SQLite Storage Design

> Status: Approved design  
> Owner: vvitem  
> Date: 2026-07-15  
> Issue: #12  
> Scope: SQLite WAL, Migration runner, Store/transaction boundary, repository foundation

## 1. Purpose

M0-004 establishes ItemAll's local-first persistence foundation. It must support future identity, asset, CredentialRef, Operation Bus, AuditEvent, Task Engine, Evidence and Runbook modules without creating those business schemas yet.

The storage layer must be deterministic, Windows-friendly, fail-closed, testable and inaccessible from React or Wails bindings.

## 2. Goals

- Use SQLite WAL through a pure-Go driver.
- Resolve a stable per-user application data path on each supported OS.
- Apply explicit, append-only embedded SQL migrations at startup.
- Verify migration history with SHA-256 checksums.
- Expose a narrow Store and transaction callback boundary.
- Keep `*sql.DB`, `*sql.Tx`, DSNs and driver errors inside `internal/storage/sqlite`.
- Provide stable storage error codes and safe UI-facing error projection.
- Validate fresh install, repeated startup, upgrades, rollback, lock contention and corruption paths.

## 3. Non-goals

M0-004 does not implement:

- workspace, identity, asset, credential, task, operation, audit, evidence or runbook business tables;
- OS Keychain or CredentialRef behavior;
- SSH, SFTP, MySQL, PostgreSQL or AI features;
- database encryption;
- automatic backup, restore or destructive repair;
- portable mode or network-share databases;
- ORM, generic CRUD repository or Active Record;
- a database administration UI.

## 4. Selected approach

### 4.1 SQLite driver

Use `modernc.org/sqlite` and avoid CGO.

Windows is the primary supported desktop platform. Avoiding GCC/MinGW and CGO reduces local setup, CI and release complexity. The implementation plan must pin the selected `modernc.org/sqlite` version and verify the resolved `modernc.org/libc` dependency in `go.mod` and `go.sum`.

### 4.2 Migration mechanism

Use embedded SQL plus a project-local minimal runner. Do not introduce goose, golang-migrate or another general migration framework.

The runner is limited to:

- fixed migration-ledger bootstrap;
- strict discovery and ordering;
- metadata and checksum validation;
- applying one migration per transaction;
- recording successful versions;
- refusing unsafe or ambiguous histories.

### 4.3 Data directory

Use OS-standard application data paths:

- Windows: `%LOCALAPPDATA%\ItemAll\data\itemall.db`
- macOS: `~/Library/Application Support/ItemAll/data/itemall.db`
- Linux: `${XDG_DATA_HOME:-~/.local/share}/itemall/itemall.db`

The default path must never depend on the current working directory or executable directory. Tests and explicit development wiring may inject a database path. Ordinary environment variables must not silently override the production location.

### 4.4 Transaction model

Use `Store.WithTx` with callback-managed transactions. Callers do not manually commit or roll back and do not receive `*sql.Tx`. Application-level write transactions are serialized inside Store to reduce SQLite writer contention.

## 5. Package structure

```text
internal/storage/
├── errors.go
├── querier.go
├── transaction.go
└── sqlite/
    ├── config.go
    ├── path.go
    ├── path_windows.go
    ├── path_darwin.go
    ├── path_linux.go
    ├── store.go
    ├── transaction.go
    ├── migrate.go
    ├── health.go
    ├── errors.go
    ├── store_test.go
    ├── transaction_test.go
    ├── migrate_test.go
    ├── concurrency_test.go
    └── migrations/
        ├── migrations.go
        └── 000001_storage_meta.sql
```

`internal/storage` defines portable contracts. `internal/storage/sqlite` owns every SQLite-specific detail. Future domain repositories depend on `storage.Querier`, not SQLite types.

## 6. Internal contracts

### 6.1 Querier

```go
type Querier interface {
    ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
    QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
    QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
```

This is an internal repository dependency, not a Wails binding contract.

### 6.2 Store transaction callback

```go
func (s *Store) WithTx(
    ctx context.Context,
    fn func(context.Context, Querier) error,
) error
```

Behavior:

- nil callback is rejected;
- callback returning nil commits;
- callback returning an error rolls back and preserves that error;
- panic triggers rollback and is then re-panicked;
- context cancellation stops work and rolls back;
- commit failure maps to a stable storage error;
- rollback failure is attached to internal diagnostics without replacing the primary error;
- nested transactions are unsupported;
- transactions must not escape the callback;
- Store methods after Close return `storage_closed`.

## 7. SQLite runtime configuration

Required effective configuration:

```text
journal_mode = WAL
foreign_keys = ON
busy_timeout = 5000 ms
synchronous = NORMAL
MaxOpenConns = 4
MaxIdleConns = 4
ConnMaxLifetime = 0
```

Use driver-supported DSN options so each created connection receives required connection-scoped settings. The implementation must explicitly verify effective values after opening.

Additional decisions:

- use immediate write transaction locking where supported;
- disable permissive double-quoted string compatibility when supported by the pinned driver;
- use UTC for stored timestamps;
- do not add aggressive cache, mmap or checkpoint tuning in M0-004;
- do not automatically run VACUUM or REINDEX;
- reject Windows UNC paths and known network-share locations for the default database.

Reads may execute concurrently. Application writes are serialized through Store. Cross-process contention waits up to the configured busy timeout and then maps to `storage_busy`.

## 8. Startup lifecycle

```text
resolve platform path
  → create and validate data directory
  → build SQLite configuration
  → open database pool
  → ping database
  → verify effective pragmas
  → bootstrap migration ledger
  → load and validate embedded migrations
  → validate recorded migration history
  → apply pending migrations
  → verify schema and pragmas again
  → mark Store ready
  → initialize storage-dependent services
```

Until the final step succeeds, Store is not ready and storage-dependent services must not initialize.

The app may display a storage startup error state, but it must not expose DSNs, SQL, local usernames, stack traces or raw driver errors to React.

## 9. Migration design

### 9.1 Ledger bootstrap

The runner owns one fixed, idempotent bootstrap DDL statement:

```sql
CREATE TABLE IF NOT EXISTS schema_migrations (
    version    INTEGER PRIMARY KEY,
    name       TEXT NOT NULL,
    checksum   TEXT NOT NULL,
    applied_at TEXT NOT NULL
) STRICT;
```

This ledger bootstrap runs before versioned migration discovery and is the only non-versioned DDL allowed in the runner. It must remain semantically stable. Any future ledger evolution requires an explicit design update and compatibility test; it must not be changed casually.

The bootstrap does not create business or storage-domain tables.

### 9.2 Versioned migrations

Migration file names use six-digit strictly increasing versions:

```text
000001_storage_meta.sql
000002_example.sql
```

The runner rejects:

- duplicate versions or names;
- missing versions in the sequence;
- malformed or non-six-digit versions;
- unsupported file extensions;
- empty SQL files;
- changed name or checksum for an applied version;
- a database schema version newer than the application knows;
- invalid migration metadata.

Checksum rules:

- SHA-256 over embedded raw SQL bytes;
- SQL files use LF line endings;
- applied migration files are append-only and must never be edited.

Each migration executes in its own immediate transaction:

1. begin transaction;
2. execute the complete SQL file without custom semicolon splitting;
3. insert the migration record;
4. commit;
5. on failure, roll back that migration and stop startup.

Earlier successfully committed versions remain recorded. The application still fails closed until the failure is corrected. Applied history is never rewritten.

`000001_storage_meta.sql` creates only the minimal versioned storage metadata needed to prove migration behavior. It must not create future domain tables or recreate `schema_migrations`.

## 10. Error model

Stable internal codes:

```text
storage_path_unavailable
storage_permission_denied
storage_open_failed
storage_busy
storage_corrupt
storage_pragma_invalid
storage_migration_invalid
storage_migration_failed
storage_migration_checksum_mismatch
storage_schema_too_new
storage_transaction_failed
storage_commit_failed
storage_closed
```

Internal errors may contain Code, Operation, Cause, an optional local Path and diagnostic metadata safe for local logs.

UI projection contains only:

- Code;
- SafeMessage;
- Retryable.

Raw SQL, DSNs, filesystem usernames and driver text must not cross the Wails boundary.

## 11. Failure and recovery policy

Use strict fail-closed behavior.

For migration failure, permission failure, lock exhaustion, corruption or incompatible schema:

- do not enter Ready state;
- do not delete, replace, rename or automatically rebuild the database;
- do not silently open a new database elsewhere;
- do not downgrade into a partial read-only product mode;
- present retry and safe diagnostic options only;
- require an explicit future recovery workflow for backup, restore or destructive repair.

An older application opening a newer schema must fail with `storage_schema_too_new`.

## 12. Testing strategy

Required tests:

- platform path resolution and explicit path injection;
- first database creation and ledger bootstrap;
- repeated open and idempotent migration;
- upgrade from the previous migration version;
- duplicate version/name and version-gap rejection;
- malformed and empty migration rejection;
- checksum mismatch and schema-too-new rejection;
- migration failure rollback;
- real foreign-key enforcement;
- WAL, busy timeout and synchronous mode verification;
- settings verified across multiple physical connections;
- concurrent reads while a write path is active;
- serialized application write transactions;
- independent connection lock contention producing `storage_busy`;
- callback commit, rollback, panic and context cancellation;
- commit/rollback error handling through test doubles where needed;
- corrupt database classification;
- use after Store.Close;
- no leakage of SQLite implementation through Wails-facing APIs.

CI requirements:

- `go test ./internal/storage/...`;
- race coverage for pure-Go storage code where compatible with the pinned driver and current CI policy;
- existing `quality`, `generated-and-security` and `windows-build` jobs remain required;
- dependency policy, drift checks and full-history secret scanning remain green;
- Windows production Wails build remains green.

Tests use `t.TempDir()` and explicit database paths. They must not write into real user application data directories.

## 13. Security and privacy boundaries

- No secret values are stored in SQLite in M0-004.
- Future credential tables store references only.
- No DSN or database content is sent to AI, telemetry or remote services.
- No network database access is added.
- Migration and storage diagnostics are local-only and redacted before UI projection.
- File permissions use the most restrictive practical mode on each platform without claiming encryption at rest.

## 14. Future compatibility

The foundation must support future domain repositories such as:

```text
identity.Repository
asset.Repository
credential.Repository
task.Repository
audit.Repository
```

These repositories share `storage.Querier` and Store transaction boundaries but own their schemas and queries.

Future migrations add workspace, actor, device and version fields as defined by their owning tasks. M0-004 does not speculate on those schemas beyond preserving clean migration and repository extension points.

## 15. Acceptance criteria

M0-004 is complete only when:

- a fresh install creates the database, migration ledger and all versioned migrations;
- repeated startup is idempotent;
- WAL, foreign keys and busy timeout are verified in tests;
- transaction commit, rollback, panic and cancellation paths are verified;
- lock contention returns a stable error after the configured timeout;
- invalid, changed or too-new migration histories block startup;
- migration failure is atomic for the current version and fails closed;
- `*sql.DB` and `*sql.Tx` do not escape the SQLite package;
- no business-domain tables or out-of-scope capabilities are introduced;
- all three required CI checks pass on the implementation PR;
- architecture, development and project status documentation is synchronized with evidence.

## 16. Implementation constraints

- Follow test-driven development.
- Keep changes in small reviewable commits.
- Do not weaken existing CI or repository boundary checks.
- Do not add unrelated refactors.
- Do not begin M0-005 or later domain work inside this task.
- Any deviation from this design requires a documented design update before implementation.
