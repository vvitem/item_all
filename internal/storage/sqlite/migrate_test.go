package sqlite

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"io/fs"
	"path/filepath"
	"testing"
	"testing/fstest"
	"time"

	"github.com/vvitem/item_all/internal/storage"
)

func TestLoadMigrationsAcceptsContiguousVersionsAndHashesRawBytes(t *testing.T) {
	files := fstest.MapFS{
		"000001_first.sql":  &fstest.MapFile{Data: []byte("CREATE TABLE first(id INTEGER);\n")},
		"000002_second.sql": &fstest.MapFile{Data: []byte("CREATE TABLE second(id INTEGER);\n")},
	}
	got, err := loadMigrations(files)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Version != 1 || got[0].Name != "first" || got[1].Version != 2 || got[1].Name != "second" {
		t.Fatalf("migrations = %#v", got)
	}
	sum := sha256.Sum256(files["000001_first.sql"].Data)
	if got[0].Checksum != hex.EncodeToString(sum[:]) {
		t.Fatalf("checksum = %q", got[0].Checksum)
	}
}

func TestLoadMigrationsRejectsInvalidSets(t *testing.T) {
	tests := map[string]fs.FS{
		"version gap": fstest.MapFS{
			"000001_first.sql": &fstest.MapFile{Data: []byte("SELECT 1;")},
			"000003_third.sql": &fstest.MapFile{Data: []byte("SELECT 3;")},
		},
		"malformed name": fstest.MapFS{
			"1_first.sql": &fstest.MapFile{Data: []byte("SELECT 1;")},
		},
		"unsupported extension": fstest.MapFS{
			"000001_first.txt": &fstest.MapFile{Data: []byte("SELECT 1;")},
		},
		"empty sql": fstest.MapFS{
			"000001_first.sql": &fstest.MapFile{Data: []byte(" \n\t")},
		},
		"duplicate name": fstest.MapFS{
			"000001_same.sql": &fstest.MapFile{Data: []byte("SELECT 1;")},
			"000002_same.sql": &fstest.MapFile{Data: []byte("SELECT 2;")},
		},
	}
	for name, files := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := loadMigrations(files)
			if !storage.IsCode(err, storage.CodeMigrationInvalid) {
				t.Fatalf("error = %v", err)
			}
		})
	}
}

func TestRunMigrationsFreshAndRepeated(t *testing.T) {
	db := openMigrationTestDB(t)
	files := migrationFS(
		"000001_storage_meta.sql", "CREATE TABLE storage_meta(key TEXT PRIMARY KEY, value TEXT NOT NULL) STRICT;",
	)
	now := func() time.Time { return time.Date(2026, 7, 15, 4, 0, 0, 0, time.FixedZone("offset", 8*60*60)) }
	if err := runMigrations(context.Background(), db, files, now); err != nil {
		t.Fatal(err)
	}
	if err := runMigrations(context.Background(), db, files, now); err != nil {
		t.Fatal(err)
	}
	assertRowCount(t, db, "schema_migrations", 1)
	var appliedAt string
	if err := db.QueryRow(`SELECT applied_at FROM schema_migrations WHERE version=1`).Scan(&appliedAt); err != nil {
		t.Fatal(err)
	}
	if appliedAt != "2026-07-14T20:00:00Z" {
		t.Fatalf("applied_at = %q", appliedAt)
	}
	assertTableExists(t, db, "storage_meta", true)
}

func TestRunMigrationsUpgradesFromPreviousVersion(t *testing.T) {
	db := openMigrationTestDB(t)
	first := migrationFS("000001_first.sql", "CREATE TABLE first(id INTEGER) STRICT;")
	if err := runMigrations(context.Background(), db, first, time.Now); err != nil {
		t.Fatal(err)
	}
	upgraded := migrationFS(
		"000001_first.sql", "CREATE TABLE first(id INTEGER) STRICT;",
		"000002_second.sql", "CREATE TABLE second(id INTEGER) STRICT;",
	)
	if err := runMigrations(context.Background(), db, upgraded, time.Now); err != nil {
		t.Fatal(err)
	}
	assertRowCount(t, db, "schema_migrations", 2)
	assertTableExists(t, db, "second", true)
}

func TestRunMigrationsRejectsChecksumMismatch(t *testing.T) {
	db := openMigrationTestDB(t)
	original := migrationFS("000001_first.sql", "CREATE TABLE first(id INTEGER) STRICT;")
	if err := runMigrations(context.Background(), db, original, time.Now); err != nil {
		t.Fatal(err)
	}
	changed := migrationFS("000001_first.sql", "CREATE TABLE first(id TEXT) STRICT;")
	err := runMigrations(context.Background(), db, changed, time.Now)
	if !storage.IsCode(err, storage.CodeMigrationChecksumMismatch) {
		t.Fatalf("error = %v", err)
	}
}

func TestRunMigrationsRejectsChangedName(t *testing.T) {
	db := openMigrationTestDB(t)
	original := migrationFS("000001_first.sql", "CREATE TABLE first(id INTEGER) STRICT;")
	if err := runMigrations(context.Background(), db, original, time.Now); err != nil {
		t.Fatal(err)
	}
	changed := migrationFS("000001_renamed.sql", "CREATE TABLE first(id INTEGER) STRICT;")
	err := runMigrations(context.Background(), db, changed, time.Now)
	if !storage.IsCode(err, storage.CodeMigrationInvalid) {
		t.Fatalf("error = %v", err)
	}
}

func TestRunMigrationsRejectsSchemaTooNew(t *testing.T) {
	db := openMigrationTestDB(t)
	newer := migrationFS(
		"000001_first.sql", "CREATE TABLE first(id INTEGER) STRICT;",
		"000002_second.sql", "CREATE TABLE second(id INTEGER) STRICT;",
	)
	if err := runMigrations(context.Background(), db, newer, time.Now); err != nil {
		t.Fatal(err)
	}
	older := migrationFS("000001_first.sql", "CREATE TABLE first(id INTEGER) STRICT;")
	err := runMigrations(context.Background(), db, older, time.Now)
	if !storage.IsCode(err, storage.CodeSchemaTooNew) {
		t.Fatalf("error = %v", err)
	}
}

func TestRunMigrationsRejectsLedgerWithAdditionalConstraint(t *testing.T) {
	db := openMigrationTestDB(t)
	_, err := db.Exec(`CREATE TABLE schema_migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL CHECK(name <> ''),
		checksum TEXT NOT NULL,
		applied_at TEXT NOT NULL
	) STRICT`)
	if err != nil {
		t.Fatal(err)
	}
	err = runMigrations(context.Background(), db, migrationFS("000001_first.sql", "SELECT 1;"), time.Now)
	if !storage.IsCode(err, storage.CodeMigrationInvalid) {
		t.Fatalf("error = %v", err)
	}
}

func TestRunMigrationsRejectsInvalidLedgerSchema(t *testing.T) {
	db := openMigrationTestDB(t)
	_, err := db.Exec(`CREATE TABLE schema_migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		checksum TEXT NOT NULL,
		applied_at TEXT NOT NULL
	)`)
	if err != nil {
		t.Fatal(err)
	}
	err = runMigrations(context.Background(), db, migrationFS("000001_first.sql", "SELECT 1;"), time.Now)
	if !storage.IsCode(err, storage.CodeMigrationInvalid) {
		t.Fatalf("error = %v", err)
	}
}

func TestRunMigrationsRollsBackFailingVersion(t *testing.T) {
	db := openMigrationTestDB(t)
	files := migrationFS(
		"000001_first.sql", "CREATE TABLE first(id INTEGER) STRICT;",
		"000002_broken.sql", "CREATE TABLE should_rollback(id INTEGER) STRICT; INSERT INTO missing_table VALUES (1);",
	)
	err := runMigrations(context.Background(), db, files, time.Now)
	if !storage.IsCode(err, storage.CodeMigrationFailed) {
		t.Fatalf("error = %v", err)
	}
	assertRowCount(t, db, "schema_migrations", 1)
	assertTableExists(t, db, "first", true)
	assertTableExists(t, db, "should_rollback", false)
}

func openMigrationTestDB(t *testing.T) *sql.DB {
	t.Helper()
	path := filepath.Join(t.TempDir(), "migration.db")
	dsn, err := buildDSN(DefaultConfig(path))
	if err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
	return db
}

func migrationFS(entries ...string) fstest.MapFS {
	files := fstest.MapFS{}
	for index := 0; index < len(entries); index += 2 {
		files[entries[index]] = &fstest.MapFile{Data: []byte(entries[index+1])}
	}
	return files
}

func assertRowCount(t *testing.T, db *sql.DB, table string, want int) {
	t.Helper()
	var got int
	if err := db.QueryRow(`SELECT COUNT(*) FROM ` + table).Scan(&got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("%s count = %d, want %d", table, got, want)
	}
}

func assertTableExists(t *testing.T, db *sql.DB, table string, want bool) {
	t.Helper()
	var got int
	if err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_schema WHERE type='table' AND name=?`, table).Scan(&got); err != nil {
		t.Fatal(err)
	}
	if (got == 1) != want {
		t.Fatalf("table %s exists = %v, want %v", table, got == 1, want)
	}
}
