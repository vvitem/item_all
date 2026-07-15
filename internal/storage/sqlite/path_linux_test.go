//go:build linux

package sqlite

import (
	"path/filepath"
	"testing"
)

func TestPlatformDataRootUsesAbsoluteXDGDataHome(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_DATA_HOME", root)
	got, err := platformDataRoot()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(root, "itemall", "data")
	if got != want {
		t.Fatalf("root = %q, want %q", got, want)
	}
}

func TestPlatformDataRootFallsBackWhenXDGDataHomeIsRelative(t *testing.T) {
	home := t.TempDir()
	t.Setenv("XDG_DATA_HOME", "relative")
	t.Setenv("HOME", home)
	got, err := platformDataRoot()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, ".local", "share", "itemall", "data")
	if got != want {
		t.Fatalf("root = %q, want %q", got, want)
	}
}
