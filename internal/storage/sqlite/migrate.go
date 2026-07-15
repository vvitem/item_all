package sqlite

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/vvitem/item_all/internal/storage"
)

const migrationLedgerDDL = `CREATE TABLE IF NOT EXISTS schema_migrations (
    version    INTEGER PRIMARY KEY,
    name       TEXT NOT NULL,
    checksum   TEXT NOT NULL,
    applied_at TEXT NOT NULL
) STRICT;`

const migrationLedgerStoredDDL = `CREATE TABLE schema_migrations (
    version    INTEGER PRIMARY KEY,
    name       TEXT NOT NULL,
    checksum   TEXT NOT NULL,
    applied_at TEXT NOT NULL
) STRICT`

var migrationFilePattern = regexp.MustCompile(`^([0-9]{6})_([a-z][a-z0-9_]*)\.sql$`)

type migration struct {
	Version  int
	Name     string
	Checksum string
	SQL      []byte
}

func loadMigrations(source fs.FS) ([]migration, error) {
	entries, err := fs.ReadDir(source, ".")
	if err != nil {
		return nil, storage.Wrap(storage.CodeMigrationInvalid, "read-embedded-migrations", err, false)
	}
	items := make([]migration, 0, len(entries))
	names := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			return nil, storage.Wrap(storage.CodeMigrationInvalid, "validate-migration-entry", fmt.Errorf("migration entry %q is a directory", entry.Name()), false)
		}
		match := migrationFilePattern.FindStringSubmatch(entry.Name())
		if match == nil {
			return nil, storage.Wrap(storage.CodeMigrationInvalid, "validate-migration-filename", fmt.Errorf("invalid migration filename %q", entry.Name()), false)
		}
		version, err := strconv.Atoi(match[1])
		if err != nil || version <= 0 {
			return nil, storage.Wrap(storage.CodeMigrationInvalid, "parse-migration-version", err, false)
		}
		name := match[2]
		if _, exists := names[name]; exists {
			return nil, storage.Wrap(storage.CodeMigrationInvalid, "validate-migration-name", fmt.Errorf("duplicate migration name %q", name), false)
		}
		names[name] = struct{}{}
		sqlBytes, err := fs.ReadFile(source, entry.Name())
		if err != nil {
			return nil, storage.Wrap(storage.CodeMigrationInvalid, "read-migration", err, false)
		}
		if strings.TrimSpace(string(sqlBytes)) == "" {
			return nil, storage.Wrap(storage.CodeMigrationInvalid, "validate-migration-sql", fmt.Errorf("migration %q is empty", entry.Name()), false)
		}
		sum := sha256.Sum256(sqlBytes)
		items = append(items, migration{
			Version:  version,
			Name:     name,
			Checksum: hex.EncodeToString(sum[:]),
			SQL:      append([]byte(nil), sqlBytes...),
		})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Version < items[j].Version })
	for index, item := range items {
		if item.Version != index+1 {
			return nil, storage.Wrap(storage.CodeMigrationInvalid, "validate-migration-sequence", fmt.Errorf("version %d is not contiguous", item.Version), false)
		}
	}
	return items, nil
}

func runMigrations(ctx context.Context, db *sql.DB, source fs.FS, now func() time.Time) error {
	if db == nil || source == nil || now == nil {
		return storage.Wrap(storage.CodeMigrationInvalid, "validate-migration-runner", errors.New("database, source, and clock are required"), false)
	}
	if _, err := db.ExecContext(ctx, migrationLedgerDDL); err != nil {
		return storage.Wrap(storage.CodeMigrationFailed, "bootstrap-migration-ledger", err, false)
	}
	if err := verifyMigrationLedger(ctx, db); err != nil {
		return err
	}
	items, err := loadMigrations(source)
	if err != nil {
		return err
	}
	applied, err := readAppliedMigrations(ctx, db)
	if err != nil {
		return err
	}
	if err := validateAppliedMigrations(applied, items); err != nil {
		return err
	}
	for _, item := range items[len(applied):] {
		if err := applyMigration(ctx, db, item, now); err != nil {
			return err
		}
	}
	return nil
}

type ledgerColumn struct {
	Name       string
	Type       string
	NotNull    int
	PrimaryKey int
	HasDefault bool
}

func verifyMigrationLedger(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, `PRAGMA table_info(schema_migrations)`)
	if err != nil {
		return storage.Wrap(storage.CodeMigrationInvalid, "inspect-migration-ledger", err, false)
	}
	defer rows.Close()
	columns := make([]ledgerColumn, 0, 4)
	for rows.Next() {
		var cid int
		var name, columnType string
		var notNull, primaryKey int
		var defaultValue sql.NullString
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			return storage.Wrap(storage.CodeMigrationInvalid, "scan-migration-ledger", err, false)
		}
		if cid != len(columns) {
			return storage.Wrap(storage.CodeMigrationInvalid, "validate-migration-ledger-order", fmt.Errorf("unexpected column id %d", cid), false)
		}
		columns = append(columns, ledgerColumn{Name: name, Type: columnType, NotNull: notNull, PrimaryKey: primaryKey, HasDefault: defaultValue.Valid})
	}
	if err := rows.Err(); err != nil {
		return storage.Wrap(storage.CodeMigrationInvalid, "read-migration-ledger", err, false)
	}
	want := []ledgerColumn{
		{Name: "version", Type: "INTEGER", NotNull: 0, PrimaryKey: 1},
		{Name: "name", Type: "TEXT", NotNull: 1, PrimaryKey: 0},
		{Name: "checksum", Type: "TEXT", NotNull: 1, PrimaryKey: 0},
		{Name: "applied_at", Type: "TEXT", NotNull: 1, PrimaryKey: 0},
	}
	if len(columns) != len(want) {
		return storage.Wrap(storage.CodeMigrationInvalid, "validate-migration-ledger-shape", fmt.Errorf("ledger has %d columns", len(columns)), false)
	}
	for index := range want {
		got := columns[index]
		expected := want[index]
		if got.Name != expected.Name || !strings.EqualFold(got.Type, expected.Type) || got.NotNull != expected.NotNull || got.PrimaryKey != expected.PrimaryKey || got.HasDefault {
			return storage.Wrap(storage.CodeMigrationInvalid, "validate-migration-ledger-column", fmt.Errorf("column %d is incompatible", index), false)
		}
	}
	var schemaSQL sql.NullString
	if err := db.QueryRowContext(ctx, `SELECT sql FROM sqlite_schema WHERE type='table' AND name='schema_migrations'`).Scan(&schemaSQL); err != nil {
		return storage.Wrap(storage.CodeMigrationInvalid, "read-migration-ledger-schema", err, false)
	}
	if !schemaSQL.Valid || canonicalSchemaSQL(schemaSQL.String) != canonicalSchemaSQL(migrationLedgerStoredDDL) {
		return storage.Wrap(storage.CodeMigrationInvalid, "validate-migration-ledger-ddl", errors.New("migration ledger DDL is incompatible"), false)
	}
	return nil
}

func canonicalSchemaSQL(value string) string {
	value = strings.TrimSuffix(strings.TrimSpace(value), ";")
	return strings.ToUpper(strings.Join(strings.Fields(value), " "))
}

type appliedMigration struct {
	Version  int
	Name     string
	Checksum string
}

func readAppliedMigrations(ctx context.Context, db *sql.DB) ([]appliedMigration, error) {
	rows, err := db.QueryContext(ctx, `SELECT version, name, checksum FROM schema_migrations ORDER BY version`)
	if err != nil {
		return nil, storage.Wrap(storage.CodeMigrationInvalid, "read-migration-history", err, false)
	}
	defer rows.Close()
	var applied []appliedMigration
	for rows.Next() {
		var item appliedMigration
		if err := rows.Scan(&item.Version, &item.Name, &item.Checksum); err != nil {
			return nil, storage.Wrap(storage.CodeMigrationInvalid, "scan-migration-history", err, false)
		}
		applied = append(applied, item)
	}
	if err := rows.Err(); err != nil {
		return nil, storage.Wrap(storage.CodeMigrationInvalid, "iterate-migration-history", err, false)
	}
	return applied, nil
}

func validateAppliedMigrations(applied []appliedMigration, available []migration) error {
	for index, item := range applied {
		if item.Version > len(available) {
			return storage.Wrap(storage.CodeSchemaTooNew, "validate-migration-version", fmt.Errorf("database version %d exceeds application version %d", item.Version, len(available)), false)
		}
		if item.Version != index+1 {
			return storage.Wrap(storage.CodeMigrationInvalid, "validate-migration-history-sequence", fmt.Errorf("applied version %d is not contiguous", item.Version), false)
		}
		expected := available[item.Version-1]
		if item.Name != expected.Name {
			return storage.Wrap(storage.CodeMigrationInvalid, "validate-migration-history-name", fmt.Errorf("migration %d name changed", item.Version), false)
		}
		if item.Checksum != expected.Checksum {
			return storage.Wrap(storage.CodeMigrationChecksumMismatch, "validate-migration-history-checksum", fmt.Errorf("migration %d checksum changed", item.Version), false)
		}
	}
	return nil
}

func applyMigration(ctx context.Context, db *sql.DB, item migration, now func() time.Time) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return storage.Wrap(storage.CodeMigrationFailed, "begin-migration", err, false)
	}
	rollback := func(cause error) error {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			return storage.Wrap(storage.CodeMigrationFailed, "rollback-migration", errors.Join(cause, rollbackErr), false)
		}
		return storage.Wrap(storage.CodeMigrationFailed, "apply-migration", cause, false)
	}
	if _, err := tx.ExecContext(ctx, string(item.SQL)); err != nil {
		return rollback(err)
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO schema_migrations(version,name,checksum,applied_at) VALUES(?,?,?,?)`,
		item.Version,
		item.Name,
		item.Checksum,
		now().UTC().Format(time.RFC3339Nano),
	); err != nil {
		return rollback(err)
	}
	if err := tx.Commit(); err != nil {
		return storage.Wrap(storage.CodeMigrationFailed, "commit-migration", err, false)
	}
	return nil
}
