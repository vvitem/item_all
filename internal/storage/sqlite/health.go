package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/vvitem/item_all/internal/storage"
)

const verifiedConnectionCount = 4

func verifyPool(ctx context.Context, db *sql.DB, config Config) error {
	count := verifiedConnectionCount
	if config.MaxOpenConns < count {
		count = config.MaxOpenConns
	}
	connections := make([]*sql.Conn, 0, count)
	defer func() {
		for _, conn := range connections {
			_ = conn.Close()
		}
	}()
	for range count {
		conn, err := db.Conn(ctx)
		if err != nil {
			return classifySQLiteError("acquire-verification-connection", err, storage.CodePragmaInvalid, false)
		}
		connections = append(connections, conn)
	}
	for index, conn := range connections {
		if err := verifyConnection(ctx, conn, config, index); err != nil {
			return err
		}
	}
	return nil
}

func verifyConnection(ctx context.Context, conn *sql.Conn, config Config, index int) error {
	var journal string
	var foreignKeys, timeout, synchronous int
	checks := []struct {
		query string
		dest  any
	}{
		{query: `PRAGMA journal_mode`, dest: &journal},
		{query: `PRAGMA foreign_keys`, dest: &foreignKeys},
		{query: `PRAGMA busy_timeout`, dest: &timeout},
		{query: `PRAGMA synchronous`, dest: &synchronous},
	}
	for _, check := range checks {
		if err := conn.QueryRowContext(ctx, check.query).Scan(check.dest); err != nil {
			return classifySQLiteError("read-sqlite-pragma", err, storage.CodePragmaInvalid, false)
		}
	}
	if !strings.EqualFold(journal, "wal") || foreignKeys != 1 || timeout != int(config.BusyTimeout.Milliseconds()) || synchronous != 1 {
		cause := fmt.Errorf("connection %d has incompatible SQLite settings", index)
		return storage.Wrap(storage.CodePragmaInvalid, "verify-sqlite-pragmas", cause, false)
	}
	return nil
}
