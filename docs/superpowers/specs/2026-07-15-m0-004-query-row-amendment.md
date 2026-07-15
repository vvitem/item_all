# M0-004 QueryRow Contract Amendment

> Status: Proposed with Implementation Plan
> Owner: vvitem
> Date: 2026-07-15
> Parent design: `2026-07-15-m0-004-sqlite-storage-design.md`
> Issue: #12

## Problem

The approved `storage.Querier` draft returned `*sql.Row` from `QueryRowContext`. In `database/sql`, `QueryRowContext` defers query errors until `Scan`. A closed Store therefore cannot return the required stable `storage_closed` error before or during `Scan` without exposing or fabricating a concrete `*sql.Row`.

This conflicts with two approved requirements:

1. Store methods after `Close` return `storage_closed`.
2. SQLite implementation details remain behind the storage contract.

## Amendment

Introduce a minimal row contract:

```go
type Row interface {
    Scan(dest ...any) error
}

type Querier interface {
    ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
    QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
    QueryRowContext(ctx context.Context, query string, args ...any) Row
}
```

`*sql.Row` already satisfies `Row`. A closed Store returns an internal `errorRow` whose `Scan` returns `storage_closed`. SQLite transactions are wrapped by a small adapter so `*sql.Tx` never escapes.

## Scope impact

- No new Wails method.
- No new dependency.
- No business schema or repository behavior.
- No change to transaction ownership, migration rules, or failure policy.
- This amendment narrows the interface and strengthens the original safety guarantee.
