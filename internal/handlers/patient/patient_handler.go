package patient

import (
	"ORDI/internal/cache"
	"ORDI/internal/email"
	"ORDI/internal/repositories"
	"net/http"
)

type PatientHandlerInterface interface {
	Signup(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
	Appointment(http.ResponseWriter, *http.Request)
}

type patientHandler struct {
	patientRepository repositories.Patient
	cache             cache.Cache
	email             email.Email
}

type PatientHandlerConfig struct {
	PatientRepo repositories.Patient
	Cache       cache.Cache
	Email       email.Email
}

func NewPatientHandler(config PatientHandlerConfig) PatientHandlerInterface {
	return &patientHandler{
		patientRepository: config.PatientRepo,
		cache:             config.Cache,
		email:             config.Email,
	}
}
