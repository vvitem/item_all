package sqlite

import (
	"errors"

	"github.com/vvitem/item_all/internal/storage"
	modernsqlite "modernc.org/sqlite"
)

func classifySQLiteError(operation string, err error, fallback storage.Code, fallbackRetryable bool) error {
	if err == nil {
		return nil
	}
	code := fallback
	retryable := fallbackRetryable
	var sqliteErr *modernsqlite.Error
	if errors.As(err, &sqliteErr) {
		switch sqliteErr.Code() & 0xff {
		case 3, 8: // SQLITE_PERM, SQLITE_READONLY
			code = storage.CodePermissionDenied
			retryable = true
		case 5, 6: // SQLITE_BUSY, SQLITE_LOCKED
			code = storage.CodeBusy
			retryable = true
		case 11, 26: // SQLITE_CORRUPT, SQLITE_NOTADB
			code = storage.CodeCorrupt
			retryable = false
		case 14: // SQLITE_CANTOPEN
			code = storage.CodeOpenFailed
			retryable = false
		}
	}
	return storage.Wrap(code, operation, err, retryable)
}
