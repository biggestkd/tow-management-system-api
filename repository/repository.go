package repository

import (
	"context"
)

// Repository defines a generic CRUD interface for any model type T.
// Implementations should handle persistence for their respective data source (e.g., MongoDB, DynamoDB).
type Repository[T any] interface {
	Create(ctx context.Context, item *T) error
	Find(ctx context.Context, filter *T) ([]*T, error)
	Update(ctx context.Context, id string, updateData *T) error
	Delete(ctx context.Context, id string) error
}
