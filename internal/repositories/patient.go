package repositories

import (
	"ORDI/internal/database"
	"ORDI/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type patientRepository struct {
	db database.Database
}

func NewPatientRepository(db database.Database) *patientRepository {
	// Create patient table if it doesn't already exist
	if err := db.AutoMigrate(context.Background(), &models.Patient{}); err != nil {
		// Panic if unable to create database
		panic("Failed to migrate patient database: " + err.Error())
	}
	return &patientRepository{
		db: db,
	}
}

func (r *patientRepository) Save(ctx context.Context, patient *models.Patient) error {
	return r.db.Save(ctx, patient)
}

func (r *patientRepository) FindByID(ctx context.Context, id uint) (*models.Patient, error) {
	var patient models.Patient
	if err := r.db.FindByID(ctx, id, &patient); err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *patientRepository) Delete(ctx context.Context, patient *models.Patient) error {
	return r.db.Delete(ctx, patient)
}

func (r *patientRepository) FindByField(ctx context.Context, field string, value interface{}) (*models.Patient, error) {
	var patient models.Patient
	if err := r.db.FindByField(ctx, &patient, field, value); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No patient found
		}
		return nil, err
	}
	return &patient, nil
}

func (r *patientRepository) FindAll(ctx context.Context) ([]models.Patient, error) {
	var patients []models.Patient
	if err := r.db.FindAll(ctx, &patients); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No doctor found
		}
		return nil, err
	}
	return patients, nil
}

func (r *patientRepository) FindAllWithPage(ctx context.Context) ([]models.Patient, error) {
	return []models.Patient{}, nil
}
