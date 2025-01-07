package repositories

import (
	"ORDI/internal/database"
	"ORDI/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type notificationRepository struct {
	db database.Database
}

func NewNotificationRepository(db database.Database) *notificationRepository {
	// Create admin table if it doesn't already exist
	if err := db.AutoMigrate(context.Background(), &models.Notification{}); err != nil {
		// Panic if unable to create database
		panic("Failed to migrate admin database: " + err.Error())
	}
	return &notificationRepository{
		db: db,
	}
}

func (r *notificationRepository) Save(ctx context.Context, notification *models.Notification) error {
	return r.db.Save(ctx, notification)
}

func (r *notificationRepository) FindByID(ctx context.Context, id uint) (*models.Notification, error) {
	var notification models.Notification
	if err := r.db.FindByID(ctx, id, &notification); err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) Delete(ctx context.Context, notification *models.Notification) error {
	return r.db.Delete(ctx, notification)
}

func (r *notificationRepository) FindAllByField(ctx context.Context, field string, value interface{}) ([]models.Notification, error) {
	var notification []models.Notification
	if err := r.db.FindAllByField(ctx, &notification, field, value); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No record found
		}
		return nil, err
	}
	return notification, nil
}

func (r *notificationRepository) FindByField(ctx context.Context, field string, value interface{}) (*models.Notification, error) {
	var notification models.Notification
	if err := r.db.FindByField(ctx, &notification, field, value); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No record found
		}
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) FindAll(ctx context.Context) ([]models.Notification, error) {
	var notification []models.Notification
	if err := r.db.FindAll(ctx, &notification); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No doctor found
		}
		return nil, err
	}
	return notification, nil
}
