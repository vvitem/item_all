//go:build windows

package sqlite

import (
	"testing"

	"golang.org/x/sys/windows"
)

func TestIsNetworkPathRejectsUNC(t *testing.T) {
	got, err := isNetworkPath(`\\server\share\itemall.db`)
	if err != nil {
		t.Fatal(err)
	}
	if !got {
		t.Fatal("UNC path accepted")
	}
}

func TestIsNetworkPathRejectsMappedRemoteDrive(t *testing.T) {
	original := getDriveType
	getDriveType = func(*uint16) uint32 { return windows.DRIVE_REMOTE }
	t.Cleanup(func() { getDriveType = original })
	got, err := isNetworkPath(`Z:\ItemAll\data\itemall.db`)
	if err != nil {
		t.Fatal(err)
	}
	if !got {
		t.Fatal("remote drive accepted")
	}
}
