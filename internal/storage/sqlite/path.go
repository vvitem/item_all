package sqlite

import (
	"os"
	"path/filepath"

	"github.com/vvitem/item_all/internal/storage"
)

// DefaultDatabasePath returns the per-user ItemAll SQLite path for the platform.
func DefaultDatabasePath() (string, error) {
	root, err := platformDataRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "itemall.db"), nil
}

func prepareDatabasePath(path string) (string, error) {
	if path == "" {
		var err error
		path, err = DefaultDatabasePath()
		if err != nil {
			return "", err
		}
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", storage.Wrap(storage.CodePathUnavailable, "resolve-absolute-database-path", err, false)
	}
	network, err := isNetworkPath(absolute)
	if err != nil {
		return "", err
	}
	if network {
		return "", storage.Wrap(storage.CodePathUnavailable, "reject-network-database-path", errNetworkPath, false)
	}
	parent := filepath.Dir(absolute)
	if err := os.MkdirAll(parent, 0o700); err != nil {
		return "", storage.Wrap(storage.CodePermissionDenied, "create-database-directory", err, true)
	}
	probe, err := os.CreateTemp(parent, ".itemall-write-probe-*")
	if err != nil {
		return "", storage.Wrap(storage.CodePermissionDenied, "probe-database-directory", err, true)
	}
	probePath := probe.Name()
	defer os.Remove(probePath)
	if err := probe.Chmod(0o600); err != nil {
		_ = probe.Close()
		return "", storage.Wrap(storage.CodePermissionDenied, "secure-database-probe", err, true)
	}
	if err := probe.Close(); err != nil {
		return "", storage.Wrap(storage.CodePermissionDenied, "close-database-probe", err, true)
	}
	if err := os.Remove(probePath); err != nil {
		return "", storage.Wrap(storage.CodePermissionDenied, "remove-database-probe", err, true)
	}
	return absolute, nil
}
