package repositories

import (
	"ORDI/internal/database"
	"ORDI/internal/models"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type masterAdminRepository struct {
	db database.Database
}

func NewMasterAdminRepository(db database.Database) *masterAdminRepository {
	// Create admin table if it doesn't already exist
	if err := db.AutoMigrate(context.Background(), &models.MasterAdmin{}); err != nil {
		// Panic if unable to create database
		panic("Failed to migrate doctor database: " + err.Error())
	}
	return &masterAdminRepository{
		db: db,
	}
}

func (r *masterAdminRepository) Save(ctx context.Context, admin *models.MasterAdmin) error {
	return r.db.Save(ctx, admin)
}

func (r *masterAdminRepository) FindByID(ctx context.Context, id uint) (*models.MasterAdmin, error) {
	var admin models.MasterAdmin
	if err := r.db.FindByID(ctx, id, &admin); err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *masterAdminRepository) Delete(ctx context.Context, admin *models.MasterAdmin) error {
	return r.db.Delete(ctx, admin)
}

func (r *masterAdminRepository) FindByField(ctx context.Context, field string, value interface{}) (*models.MasterAdmin, error) {
	var admin models.MasterAdmin
	if err := r.db.FindByField(ctx, &admin, field, value); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No admin found
		}
		return nil, err
	}
	return &admin, nil
}

func (r *masterAdminRepository) FindAll(ctx context.Context) ([]models.MasterAdmin, error) {
	var masterAdmin []models.MasterAdmin
	if err := r.db.FindAll(ctx, &masterAdmin); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No doctor found
		}
		return nil, err
	}
	return masterAdmin, nil
}

func (r *masterAdminRepository) FindAllByField(ctx context.Context, field string, value interface{}) ([]models.MasterAdmin, error) {
	var admin []models.MasterAdmin
	if err := r.db.FindByField(ctx, &admin, field, value); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No Admin Found
		}
		return nil, err
	}
	return admin, nil
}

func (r *masterAdminRepository) FindAllWithPage(ctx context.Context) ([]models.MasterAdmin, error) {
	return []models.MasterAdmin{}, nil
}

func (r *masterAdminRepository) FilterByDate(ctx context.Context, idField string,
	idValue uint, filterField string,
	filterFieldValue time.Time) ([]models.MasterAdmin, error) {
	return []models.MasterAdmin{}, nil
}

func (r *masterAdminRepository) FilterBetweenDates(ctx context.Context,
	idField string,
	idValue uint, field string, start, end time.Time) ([]models.MasterAdmin, error) {
	return nil, nil
}
