package repositories

import (
	"context"
)

type Repository[T any] interface {
	Save(ctx context.Context, entity *T) error

	FindByID(ctx context.Context, id uint) (*T, error)

	Delete(ctx context.Context, entity *T) error

	FindByField(ctx context.Context, field string, value interface{}) (*T, error)

	FindAll(ctx context.Context) ([]T, error)

	FindAllWithPage(ctx context.Context) ([]T, error)
}
