//go:build !windows

package sqlite

import (
	"os"

	"github.com/vvitem/item_all/internal/storage"
)

func secureDatabaseFile(path string) error {
	if err := os.Chmod(path, 0o600); err != nil {
		return storage.Wrap(storage.CodePermissionDenied, "secure-database-file", err, true)
	}
	return nil
}
