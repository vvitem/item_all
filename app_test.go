package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/vvitem/item_all/internal/buildinfo"
	"github.com/vvitem/item_all/internal/storage"
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

	got := NewApp(nil).GetAppInfo()

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
	if got.Error != nil {
		t.Fatalf("Error = %+v", got.Error)
	}
	if !strings.Contains(got.Runtime, "/") {
		t.Fatalf("Runtime = %q, want Go version and OS/arch", got.Runtime)
	}
}

func TestGetAppInfoReturnsSafeStorageFailure(t *testing.T) {
	startupErr := storage.Wrap(
		storage.CodePermissionDenied,
		"open",
		errors.New(`C:\Users\private\itemall.db: SQLITE_CANTOPEN`),
		true,
	)
	got := NewApp(startupErr).GetAppInfo()
	if got.Status != appStatusStorageError {
		t.Fatalf("Status = %q", got.Status)
	}
	if got.Error == nil {
		t.Fatal("Error = nil")
	}
	if got.Error.Code != string(storage.CodePermissionDenied) || !got.Error.Retryable {
		t.Fatalf("Error = %+v", got.Error)
	}
	if strings.Contains(got.Error.SafeMessage, "Users") || strings.Contains(got.Error.SafeMessage, "SQLITE") || strings.Contains(got.Error.SafeMessage, "itemall.db") {
		t.Fatalf("unsafe message = %q", got.Error.SafeMessage)
	}
}
