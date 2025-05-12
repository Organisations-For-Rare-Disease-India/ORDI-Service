package repositories

import (
	"ORDI/internal/database"
	"ORDI/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type appointmentRepository struct {
	db database.Database
}

func NewAppointmentRepository(db database.Database) *appointmentRepository {
	if err := db.AutoMigrate(context.Background(), &models.Appointment{}); err != nil {
		panic("failed to migrate appointment database: " + err.Error())
	}
	return &appointmentRepository{db: db}
}

func (a *appointmentRepository) Save(ctx context.Context, appointment *models.Appointment) error {
	return a.db.Save(ctx, appointment)
}

func (a *appointmentRepository) FindByID(ctx context.Context, id uint) (*models.Appointment, error) {
	var ap *models.Appointment
	if err := a.db.FindByID(ctx, id, ap); err != nil {
		return nil, err
	}
	return ap, nil
}

func (a *appointmentRepository) Delete(ctx context.Context, patient *models.Appointment) error {
	return a.db.Delete(ctx, patient)
}

func (a *appointmentRepository) FindByField(ctx context.Context,
	field string, value any) (*models.Appointment, error) {
	var ap *models.Appointment
	if err := a.db.FindByField(ctx, ap, field, value); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ap, nil
}

func (a *appointmentRepository) FindAll(ctx context.Context) ([]models.Appointment, error) {
	var app []models.Appointment
	if err := a.db.FindAll(ctx, &app); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No doctor found
		}
		return nil, err
	}
	return app, nil
}
