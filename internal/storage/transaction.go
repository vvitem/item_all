package storage

import "context"

// Transactor owns transaction lifecycle and exposes only repository-safe operations.
type Transactor interface {
	WithTx(ctx context.Context, fn func(context.Context, Querier) error) error
}
