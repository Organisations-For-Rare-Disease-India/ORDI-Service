package repositories

import (
	"ORDI/internal/database"
	"ORDI/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type doctorRepository struct {
	db database.Database
}

func NewDoctorRepository(db database.Database) *doctorRepository {
	// Create doctor table if it doesn't already exist
	if err := db.AutoMigrate(context.Background(), &models.Doctor{}); err != nil {
		// Panic if unable to create database
		panic("Failed to migrate doctor database: " + err.Error())
	}
	return &doctorRepository{
		db: db,
	}
}

func (r *doctorRepository) Save(ctx context.Context, doctor *models.Doctor) error {
	return r.db.Save(ctx, doctor)
}

func (r *doctorRepository) FindByID(ctx context.Context, id uint) (*models.Doctor, error) {
	var doctor models.Doctor
	if err := r.db.FindByID(ctx, id, &doctor); err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (r *doctorRepository) Delete(ctx context.Context, doctor *models.Doctor) error {
	return r.db.Delete(ctx, doctor)
}

func (r *doctorRepository) FindByField(ctx context.Context, field string, value interface{}) (*models.Doctor, error) {
	var doctor models.Doctor
	if err := r.db.FindByField(ctx, &doctor, field, value); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No doctor found
		}
		return nil, err
	}
	return &doctor, nil
}
