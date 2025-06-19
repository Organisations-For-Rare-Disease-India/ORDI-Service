package repositories

import (
	"context"
	"time"
)

type Repository[T any] interface {
	Save(ctx context.Context, entity *T) error

	FindByID(ctx context.Context, id uint) (*T, error)

	Delete(ctx context.Context, entity *T) error

	FindByField(ctx context.Context, field string, value any) (*T, error)

	FindAll(ctx context.Context) ([]T, error)

	FindAllWithPage(ctx context.Context) ([]T, error)

	FindAllByField(ctx context.Context, field string, value interface{}) ([]T, error)

	FilterBetweenDates(ctx context.Context, idField string, idValue uint,
		field string, start, end time.Time) ([]T, error)

	FilterByDate(ctx context.Context, idField string, idValue uint,
		filterField string, filterFieldValue time.Time) ([]T, error)
}
