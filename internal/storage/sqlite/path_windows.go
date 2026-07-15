//go:build windows

package sqlite

import (
	"path/filepath"
	"strings"

	"github.com/vvitem/item_all/internal/storage"
	"golang.org/x/sys/windows"
)

func platformDataRoot() (string, error) {
	root, err := windows.KnownFolderPath(windows.FOLDERID_LocalAppData, 0)
	if err != nil {
		return "", storage.Wrap(storage.CodePathUnavailable, "resolve-windows-local-app-data", err, false)
	}
	return filepath.Join(root, "ItemAll", "data"), nil
}

var getDriveType = windows.GetDriveType

func isNetworkPath(path string) (bool, error) {
	clean := filepath.Clean(path)
	if strings.HasPrefix(clean, `\\`) || strings.HasPrefix(clean, `//`) {
		return true, nil
	}
	volume := filepath.VolumeName(clean)
	if len(volume) != 2 || volume[1] != ':' {
		return false, nil
	}
	root := volume + `\`
	ptr, err := windows.UTF16PtrFromString(root)
	if err != nil {
		return false, storage.Wrap(storage.CodePathUnavailable, "encode-drive-root", err, false)
	}
	return getDriveType(ptr) == windows.DRIVE_REMOTE, nil
}
