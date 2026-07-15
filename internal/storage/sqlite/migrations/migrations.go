package migrations

import "embed"

// FS contains the append-only versioned SQLite migrations.
//
//go:embed *.sql
var FS embed.FS
