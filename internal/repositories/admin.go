package repositories

import (
	"ORDI/internal/database"
	"ORDI/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type adminRepository struct {
	db database.Database
}

func NewAdminRepository(db database.Database) *adminRepository {
	// Create admin table if it doesn't already exist
	if err := db.AutoMigrate(context.Background(), &models.Admin{}); err != nil {
		// Panic if unable to create database
		panic("Failed to migrate admin database: " + err.Error())
	}
	return &adminRepository{
		db: db,
	}
}

func (r *adminRepository) Save(ctx context.Context, admin *models.Admin) error {
	return r.db.Save(ctx, admin)
}

func (r *adminRepository) FindByID(ctx context.Context, id uint) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.FindByID(ctx, id, &admin); err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) Delete(ctx context.Context, admin *models.Admin) error {
	return r.db.Delete(ctx, admin)
}

func (r *adminRepository) FindByField(ctx context.Context, field string, value any) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.FindByField(ctx, &admin, field, value); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No admin found
		}
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) FindAll(ctx context.Context) ([]models.Admin, error) {
	var admin []models.Admin
	if err := r.db.FindAll(ctx, &admin); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No doctor found
		}
		return nil, err
	}
	return admin, nil
}

func (r *adminRepository) FindAllWithPage(ctx context.Context) ([]models.Admin, error) {
	return []models.Admin{}, nil
}
