package buildinfo

import (
	"runtime"
	"testing"
)

func TestCurrentUsesStableDefaults(t *testing.T) {
	originalVersion := Version
	originalCommit := Commit
	originalBuildTime := BuildTime
	t.Cleanup(func() {
		Version = originalVersion
		Commit = originalCommit
		BuildTime = originalBuildTime
	})

	Version = ""
	Commit = ""
	BuildTime = ""

	got := Current()

	if got.Name != "ItemAll" {
		t.Fatalf("Name = %q, want ItemAll", got.Name)
	}
	if got.Version != "dev" {
		t.Fatalf("Version = %q, want dev", got.Version)
	}
	if got.Commit != "unknown" {
		t.Fatalf("Commit = %q, want unknown", got.Commit)
	}
	if got.BuildTime != "unknown" {
		t.Fatalf("BuildTime = %q, want unknown", got.BuildTime)
	}
	if got.GoVersion != runtime.Version() {
		t.Fatalf("GoVersion = %q, want %q", got.GoVersion, runtime.Version())
	}
	if got.OS != runtime.GOOS {
		t.Fatalf("OS = %q, want %q", got.OS, runtime.GOOS)
	}
	if got.Arch != runtime.GOARCH {
		t.Fatalf("Arch = %q, want %q", got.Arch, runtime.GOARCH)
	}
}

func TestCurrentUsesInjectedValues(t *testing.T) {
	originalVersion := Version
	originalCommit := Commit
	originalBuildTime := BuildTime
	t.Cleanup(func() {
		Version = originalVersion
		Commit = originalCommit
		BuildTime = originalBuildTime
	})

	Version = "0.0.0-dev"
	Commit = "abc123def456"
	BuildTime = "2026-07-13T08:00:00Z"

	got := Current()

	if got.Version != Version {
		t.Fatalf("Version = %q, want %q", got.Version, Version)
	}
	if got.Commit != Commit {
		t.Fatalf("Commit = %q, want %q", got.Commit, Commit)
	}
	if got.BuildTime != BuildTime {
		t.Fatalf("BuildTime = %q, want %q", got.BuildTime, BuildTime)
	}
}
