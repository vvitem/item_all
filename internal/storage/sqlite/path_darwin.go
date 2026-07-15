//go:build darwin

package sqlite

import (
	"os"
	"path/filepath"

	"github.com/vvitem/item_all/internal/storage"
)

func platformDataRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", storage.Wrap(storage.CodePathUnavailable, "resolve-macos-home", err, false)
	}
	return filepath.Join(home, "Library", "Application Support", "ItemAll", "data"), nil
}

func isNetworkPath(string) (bool, error) {
	return false, nil
}
