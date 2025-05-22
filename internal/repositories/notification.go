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
	if err := db.AutoMigrate(context.Background(), &models.Notification{}); err != nil {
		// Panic when defining Notification Schema on the Backend
		panic("Failed to migrate notification table schema " + err.Error())
	}
	return &notificationRepository{
		db: db,
	}
}

// Creating a notification on the backend
func (n *notificationRepository) Save(ctx context.Context, notification *models.Notification) error {
	return n.db.Save(ctx, notification)
}

// Finding notification by ID
func (n *notificationRepository) FindByID(ctx context.Context, id uint) (*models.Notification, error) {
	var notification models.Notification
	if err := n.db.FindByID(ctx, id, &notification); err != nil {
		return nil, err
	}
	return &notification, nil
}

// Deleting particular notification
func (n *notificationRepository) Delete(ctx context.Context, notification *models.Notification) error {
	return n.db.Delete(ctx, notification)
}


func (n *notificationRepository) FindByField(ctx context.Context, field string, value interface{}) (*models.Notification, error) {
	var notification models.Notification
	if err := n.db.FindByField(ctx, &notification, field, value); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No admin found
		}
		return nil, err
	}
	return &notification, nil
}

func (n *notificationRepository) FindAll(ctx context.Context) ([]models.Notification, error) {
	var notifications []models.Notification
	if err := n.db.FindAll(ctx, &notifications); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No doctor found
		}
		return nil, err
	}
	return notifications, nil
}

func (n *notificationRepository) FindAllByField(ctx context.Context, field string, value interface{}) ([]models.Notification, error) {
	var notifications []models.Notification
	if err := n.db.FindAllByField(ctx, &notifications, field, value); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No Admin Found 
		}
		return nil, err
	}
	return notifications, nil
}
