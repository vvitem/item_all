# M0-004 SQLite Storage Plan Self-Review

> Status: Normative companion to `2026-07-15-m0-004-sqlite-storage.md`
> Date: 2026-07-15
> Issue: #12
> Rule: Where this document is more specific than the implementation plan, this document controls.

## Review result

The plan covers every approved design section and contains no `TBD`, `TODO`, deferred implementation placeholder, or unrelated subsystem. The following clarifications remove ambiguities discovered during type, platform, and failure-path review.

## 1. Exact project-status files

Task 8 must update exactly these seven files and no inferred eighth status file:

```text
docs/project-management/02-roadmap/backlog.md
docs/project-management/04-progress/changelog.md
docs/project-management/04-progress/current-status.md
docs/project-management/04-progress/implementation-trace.md
docs/project-management/04-progress/milestone-status.md
docs/project-management/04-progress/task-tracking.md
docs/project-management/04-progress/weekly-log.md
```

## 2. CI boundary wording

Task 1 must replace stale M0-003-specific policy messages while preserving the same forbidden directories.

Use:

```go
violations.Add(rel, "scope-directory", "capability is outside the currently approved task scope")
```

The environment rule must say:

```go
violations.Add(rel, "environment-read", "only XDG_DATA_HOME in the Linux storage path resolver is approved")
```

No environment read is allowed in any other application source file.

## 3. Platform path implementation and test seam

### Windows

Use the existing pinned `golang.org/x/sys/windows` package:

```go
func platformDataRoot() (string, error) {
	root, err := windows.KnownFolderPath(windows.FOLDERID_LocalAppData, 0)
	if err != nil {
		return "", storage.Wrap(storage.CodePathUnavailable, "resolve-windows-local-app-data", err, false)
	}
	return filepath.Join(root, "ItemAll", "data"), nil
}

var getDriveType = windows.GetDriveType

func isNetworkPath(path string) (bool, error) {
	clean := filepath.Clean(path)
	if strings.HasPrefix(clean, `\\`) || strings.HasPrefix(clean, `//`) {
		return true, nil
	}
	volume := filepath.VolumeName(clean)
	if len(volume) != 2 || volume[1] != ':' {
		return false, nil
	}
	root := volume + `\`
	ptr, err := windows.UTF16PtrFromString(root)
	if err != nil {
		return false, storage.Wrap(storage.CodePathUnavailable, "encode-drive-root", err, false)
	}
	return getDriveType(ptr) == windows.DRIVE_REMOTE, nil
}
```

`path_windows_test.go` must restore `getDriveType` with `t.Cleanup` and cover both UNC and an injected `DRIVE_REMOTE` result. The test must not require a real mapped drive.

### Linux

`path_linux_test.go` must use `t.Setenv` and must cover:

```go
func TestPlatformDataRootUsesAbsoluteXDGDataHome(t *testing.T)
func TestPlatformDataRootFallsBackWhenXDGDataHomeIsRelative(t *testing.T)
```

Set `HOME` to `t.TempDir()` in the fallback test so no real user path is asserted or touched.

### macOS

`path_darwin.go` uses `os.UserHomeDir()` and appends `Library/Application Support/ItemAll/data`. No environment-variable exception is added to CI for macOS source.

## 4. Config validation and filesystem permissions

`Config.validate()` must reject:

```text
BusyTimeout <= 0
MaxOpenConns <= 0
MaxIdleConns < 0
MaxIdleConns > MaxOpenConns
```

An empty `Config.Path` is valid and means “resolve the platform default”. An explicit path is converted to an absolute path before DSN construction.

`prepareDatabasePath` creates the parent directory with `0700`, performs a create/close/remove write probe, and never truncates the target database. On Unix, after SQLite creates the database, call `os.Chmod(path, 0600)`. On Windows, rely on the per-user LocalAppData ACL inherited by the directory and do not claim encryption at rest.

## 5. Migration-ledger schema validation

Checking column names alone is insufficient. After the idempotent bootstrap, read `PRAGMA table_info(schema_migrations)` and reject extra, missing, reordered, or incompatible columns.

Required shape:

```text
version     INTEGER  primary-key=1
name        TEXT     not-null=1 primary-key=0
checksum    TEXT     not-null=1 primary-key=0
applied_at  TEXT     not-null=1 primary-key=0
```

SQLite may report `notnull=0` for the `INTEGER PRIMARY KEY` column because primary-key semantics already imply non-null; do not reject that value for `version`.

Also read `sqlite_schema.sql` for `schema_migrations` and require the table to be `STRICT`. Any mismatch maps to `storage_migration_invalid`; the runner must not alter or repair the table.

## 6. SQLite error classification

For `*modernc.org/sqlite.Error`, classify the primary result code, not the full extended code:

```go
primary := sqliteErr.Code() & 0xff
```

Map:

```text
3  SQLITE_PERM      -> storage_permission_denied
5  SQLITE_BUSY      -> storage_busy
6  SQLITE_LOCKED    -> storage_busy
8  SQLITE_READONLY  -> storage_permission_denied
11 SQLITE_CORRUPT   -> storage_corrupt
14 SQLITE_CANTOPEN  -> storage_open_failed
26 SQLITE_NOTADB    -> storage_corrupt
```

Unknown SQLite codes retain their internal cause and map to the operation's fallback storage code. Raw driver text never enters `SafeError`.

## 7. Transaction nesting rule

Nested transaction detection is context-based. The callback must use the `context.Context` supplied to it for every repository call and for any attempted nested `WithTx` call.

The normative test is:

```go
func TestWithTxRejectsNestedTransaction(t *testing.T) {
	store := openTestStore(t)
	err := store.WithTx(context.Background(), func(txCtx context.Context, _ storage.Querier) error {
		return store.WithTx(txCtx, func(context.Context, storage.Querier) error { return nil })
	})
	if !storage.IsCode(err, storage.CodeTransactionFailed) {
		t.Fatalf("error = %v", err)
	}
}
```

Calling a nested transaction with an unrelated context is outside the contract and must be explicitly prohibited in `docs/development/storage.md`; domain repositories must always propagate the callback context.

## 8. Store close semantics

`QueryRowContext` returns the approved `storage.Row` interface. After `Close`, it returns `errorRow{err: storage_closed}` and `Scan` returns that stable error. `ExecContext`, `QueryContext`, and `WithTx` perform the same closed-state check before touching `database/sql`.

`Close` is idempotent. Concurrent use during shutdown is not promised to finish successfully, but it must never panic or expose a raw `sql.ErrConnDone` through the Wails boundary.

## 9. Connection verification

`verifyPool` must hold four `*sql.Conn` values simultaneously, verify each, then close all four even when one verification fails. It must use a bounded context inherited from `Open`.

Required effective values:

```text
journal_mode = wal
foreign_keys = 1
busy_timeout = Config.BusyTimeout in milliseconds
synchronous = 1
```

The implementation test must prove four distinct physical connections were held by setting `MaxOpenConns=4`, acquiring all four before querying, and failing if any acquisition or PRAGMA check fails.

## 10. TDD and failure evidence

Every implementation task follows this order:

```text
focused test added
focused command run and observed failing for the expected reason
minimum implementation added
same focused command observed passing
broader package tests observed passing
commit
```

A test that passes before its implementation does not count as RED evidence. In Task 1, the two drift tests provide the expected RED; the baseline-acceptance test is a post-change regression test.

For CI failures, inspect the exact job and step. Rerunning without diagnosis is prohibited. A product-code defect requires a focused regression test before its fix.

## 11. Scope and type consistency result

The reviewed signatures are consistent across all tasks:

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

`*sql.DB` and `*sql.Tx` remain private to `internal/storage/sqlite`. `*sql.Rows` and `sql.Result` are repository-level standard-library abstractions and are not Wails DTOs. `GetAppInfo` remains the sole generated binding.
