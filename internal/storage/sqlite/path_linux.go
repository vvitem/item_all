//go:build linux

package sqlite

import (
	"os"
	"path/filepath"

	"github.com/vvitem/item_all/internal/storage"
)

func platformDataRoot() (string, error) {
	if root, ok := os.LookupEnv("XDG_DATA_HOME"); ok && filepath.IsAbs(root) {
		return filepath.Join(root, "itemall", "data"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", storage.Wrap(storage.CodePathUnavailable, "resolve-linux-home", err, false)
	}
	return filepath.Join(home, ".local", "share", "itemall", "data"), nil
}

func isNetworkPath(string) (bool, error) {
	return false, nil
}
