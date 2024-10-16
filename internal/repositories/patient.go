package repositories

import (
	"ORDI/internal/database"
	"ORDI/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type Patient interface {
	Save(ctx context.Context, patient *models.PatientInfo) error

	FindByID(ctx context.Context, id uint) (*models.PatientInfo, error)

	Delete(ctx context.Context, patient *models.PatientInfo) error

	FindByField(ctx context.Context, field string, value interface{}) (*models.PatientInfo, error)
}

type patientRepository struct {
	db database.Database
}

func NewPatientRepository(db database.Database) *patientRepository {
	return &patientRepository{
		db: db,
	}
}

func (r *patientRepository) Save(ctx context.Context, patient *models.PatientInfo) error {
	return r.db.Save(ctx, patient)
}

func (r *patientRepository) FindByID(ctx context.Context, id uint) (*models.PatientInfo, error) {
	var patient models.PatientInfo
	if err := r.db.FindByID(ctx, id, &patient); err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *patientRepository) Delete(ctx context.Context, patient *models.PatientInfo) error {
	return r.db.Delete(ctx, patient)
}

func (r *patientRepository) FindByField(ctx context.Context, field string, value interface{}) (*models.PatientInfo, error) {
	var patient models.PatientInfo
	if err := r.db.FindByField(ctx, &patient, field, value); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No patient found
		}
		return nil, err
	}
	return &patient, nil
}
