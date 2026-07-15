package main

import (
	"context"
	"errors"

	"github.com/vvitem/item_all/internal/storage"
	"github.com/vvitem/item_all/internal/storage/sqlite"
)

type storageOpener interface {
	Open(context.Context, sqlite.Config) (*sqlite.Store, error)
}

type defaultStorageOpener struct{}

func (defaultStorageOpener) Open(ctx context.Context, config sqlite.Config) (*sqlite.Store, error) {
	return sqlite.Open(ctx, config)
}

type storageRuntime struct {
	store *sqlite.Store
	err   error
}

func bootstrapStorage(ctx context.Context, opener storageOpener, config sqlite.Config) storageRuntime {
	if opener == nil {
		return storageRuntime{err: storage.Wrap(storage.CodeOpenFailed, "bootstrap-storage", errors.New("storage opener is required"), false)}
	}
	store, err := opener.Open(ctx, config)
	if err != nil {
		return storageRuntime{err: err}
	}
	return storageRuntime{store: store}
}

func (r *storageRuntime) Close() error {
	if r == nil || r.store == nil {
		return nil
	}
	err := r.store.Close()
	r.store = nil
	return err
}
