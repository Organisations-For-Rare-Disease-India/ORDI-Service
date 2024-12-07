package doctor

import (
	"ORDI/internal/cache"
	"ORDI/internal/email"
	"ORDI/internal/models"
	"ORDI/internal/repositories"
	"net/http"
)

type Doctor interface {
	Signup(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
	Appointment(http.ResponseWriter, *http.Request)
}

type doctorHandler struct {
	doctorRepository repositories.Repository[models.Doctor]
	cache            cache.Cache
	email            email.Email
}

type DoctorHandlerConfig struct {
	DoctorRepo repositories.Repository[models.Doctor]
	Cache      cache.Cache
	Email      email.Email
}

func NewDoctorHandler(config DoctorHandlerConfig) Doctor {
	return &doctorHandler{
		doctorRepository: config.DoctorRepo,
		cache:            config.Cache,
		email:            config.Email,
	}
}
