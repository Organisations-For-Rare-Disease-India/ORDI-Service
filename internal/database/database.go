package database

import (
	"context"
)

// Database represents a service that interacts with a database.
type Database interface {
	// Save would insert or update the entity on the database
	Save(ctx context.Context, entity interface{}) error

	// FindByID retrieves a record of type entity with primary key id
	FindByID(ctx context.Context, id uint, entity interface{}) error

	// Delete removes the specified entity from the database
	// The primary id is extracted from the entity
	Delete(ctx context.Context, entity interface{}) error

	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// FindByField retrieves the first record with field and provided value
	FindByField(ctx context.Context, entity interface{}, field string, value interface{}) error

	// FindAllByField retrieves all the records with field and provided value
	FindAllByField(ctx context.Context, entity interface{}, field string, value interface{}) error

	// AutoMigrate migrates database schema to match the struct definitions
	AutoMigrate(ctx context.Context, entity interface{}) error

	// FindAll finds all instances of type entity
	FindAll(ctx context.Context, entity interface{}) error

	FindAllWithPage(ctx context.Context, paginate Paginate, entity any) error
}

type Paginate struct {
	Offset     int   `json:"offset,omitempty"`
	Page       int   `json:"page,omitempty"`
	Limit      int   `json:"limit,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int64 `json:"total_pages,omitempty"`
}
