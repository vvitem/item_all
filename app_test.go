package main

import (
	"strings"
	"testing"

	"github.com/vvitem/item_all/internal/buildinfo"
)

func TestGetAppInfoReturnsSafeReadOnlyMetadata(t *testing.T) {
	originalVersion := buildinfo.Version
	originalCommit := buildinfo.Commit
	originalBuildTime := buildinfo.BuildTime
	t.Cleanup(func() {
		buildinfo.Version = originalVersion
		buildinfo.Commit = originalCommit
		buildinfo.BuildTime = originalBuildTime
	})

	buildinfo.Version = "0.0.0-dev"
	buildinfo.Commit = "abc123def456"
	buildinfo.BuildTime = "2026-07-13T08:00:00Z"

	got := NewApp().GetAppInfo()

	if got.Name != "ItemAll" {
		t.Fatalf("Name = %q, want ItemAll", got.Name)
	}
	if got.Tagline != "Local-first SafeOps" {
		t.Fatalf("Tagline = %q", got.Tagline)
	}
	if got.Version != "0.0.0-dev" {
		t.Fatalf("Version = %q", got.Version)
	}
	if got.Commit != "abc123def456" {
		t.Fatalf("Commit = %q", got.Commit)
	}
	if got.BuildTime != "2026-07-13T08:00:00Z" {
		t.Fatalf("BuildTime = %q", got.BuildTime)
	}
	if got.Status != "ready" {
		t.Fatalf("Status = %q", got.Status)
	}
	if !strings.Contains(got.Runtime, "/") {
		t.Fatalf("Runtime = %q, want Go version and OS/arch", got.Runtime)
	}
}
