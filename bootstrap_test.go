package main

import (
	"context"
	"errors"
	"testing"

	"github.com/vvitem/item_all/internal/storage"
	"github.com/vvitem/item_all/internal/storage/sqlite"
)

func TestBootstrapStoragePreservesStartupFailure(t *testing.T) {
	startupErr := storage.Wrap(storage.CodePermissionDenied, "open", errors.New(`C:\Users\private\itemall.db`), true)
	runtime := bootstrapStorage(context.Background(), failingStorageOpener{err: startupErr}, sqlite.DefaultConfig(""))
	if runtime.store != nil {
		t.Fatal("partial Store retained")
	}
	if !errors.Is(runtime.err, startupErr) {
		t.Fatalf("error = %v", runtime.err)
	}
}

func TestStorageRuntimeCloseIsSafeWithoutStore(t *testing.T) {
	runtime := storageRuntime{}
	if err := runtime.Close(); err != nil {
		t.Fatal(err)
	}
}

type failingStorageOpener struct{ err error }

func (f failingStorageOpener) Open(context.Context, sqlite.Config) (*sqlite.Store, error) {
	return nil, f.err
}
