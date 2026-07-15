package sqlite

import (
	"net/url"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestDefaultConfigUsesApprovedValues(t *testing.T) {
	path := filepath.Join(t.TempDir(), "itemall.db")
	got := DefaultConfig(path)
	if got.Path != path {
		t.Fatalf("Path = %q, want %q", got.Path, path)
	}
	if got.BusyTimeout != 5*time.Second {
		t.Fatalf("BusyTimeout = %s", got.BusyTimeout)
	}
	if got.MaxOpenConns != 4 || got.MaxIdleConns != 4 {
		t.Fatalf("pool = %d/%d", got.MaxOpenConns, got.MaxIdleConns)
	}
}

func TestBuildDSNUsesApprovedSQLiteOptions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "item all.db")
	dsn, err := buildDSN(DefaultConfig(path))
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Scheme != "file" {
		t.Fatalf("scheme = %q", parsed.Scheme)
	}
	query := parsed.Query()
	if query.Get("mode") != "rwc" {
		t.Fatalf("mode = %q", query.Get("mode"))
	}
	wantPragmas := []string{"journal_mode(WAL)", "foreign_keys(1)", "busy_timeout(5000)", "synchronous(NORMAL)"}
	if got := query["_pragma"]; !reflect.DeepEqual(got, wantPragmas) {
		t.Fatalf("_pragma = %#v, want %#v", got, wantPragmas)
	}
	for key, want := range map[string]string{"_txlock": "immediate", "_dqs": "false", "_time_format": "sqlite", "_timezone": "UTC"} {
		if got := query.Get(key); got != want {
			t.Fatalf("%s = %q, want %q", key, got, want)
		}
	}
}

func TestConfigValidateRejectsInvalidPoolAndTimeout(t *testing.T) {
	valid := DefaultConfig(filepath.Join(t.TempDir(), "itemall.db"))
	tests := map[string]Config{
		"zero timeout":      func() Config { c := valid; c.BusyTimeout = 0; return c }(),
		"negative timeout":  func() Config { c := valid; c.BusyTimeout = -time.Second; return c }(),
		"sub-millisecond":   func() Config { c := valid; c.BusyTimeout = time.Nanosecond; return c }(),
		"zero open":         func() Config { c := valid; c.MaxOpenConns = 0; return c }(),
		"negative idle":     func() Config { c := valid; c.MaxIdleConns = -1; return c }(),
		"idle greater open": func() Config { c := valid; c.MaxIdleConns = c.MaxOpenConns + 1; return c }(),
	}
	for name, config := range tests {
		t.Run(name, func(t *testing.T) {
			if err := config.validate(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestConfigValidateAllowsEmptyPathForPlatformDefault(t *testing.T) {
	if err := DefaultConfig("").validate(); err != nil {
		t.Fatal(err)
	}
}
