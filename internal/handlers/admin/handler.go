package admin

import (
	"ORDI/internal/cache"
	"ORDI/internal/email"
	"ORDI/internal/models"
	"ORDI/internal/repositories"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/schema"
)

const JWT_TOKEN_HEADER = "token"
const EMAIL_HEADER = "email"
const EmailRequiredError = "Incorrect URL"

var decoder = schema.NewDecoder()

// MasterAdmin Creates the Admin
// Admin would then Register their profile
type Admin interface {
	Login(http.ResponseWriter, *http.Request)
	Dashboard(http.ResponseWriter, *http.Request)
	Profile(http.ResponseWriter, *http.Request)
	Register(http.ResponseWriter, *http.Request)
	Setup(http.ResponseWriter, *http.Request)
	ListDoctors(http.ResponseWriter, *http.Request)
	ListPatients(http.ResponseWriter, *http.Request)
}

type adminHandler struct {
	adminRepository       repositories.Repository[models.Admin]
	patientRepository     repositories.Repository[models.Patient]
	doctorRepository      repositories.Repository[models.Doctor]
	appointmentRepository repositories.Repository[models.Appointment]
	email                 email.Email
}

type AdminHandlerConfig struct {
	AdminRepo       repositories.Repository[models.Admin]
	PatientRepo     repositories.Repository[models.Patient]
	DoctorRepo      repositories.Repository[models.Doctor]
	AppointmentRepo repositories.Repository[models.Appointment]
	Cache           cache.Cache
	Email           email.Email
}

func NewAdminHandler(config AdminHandlerConfig) Admin {
	return &adminHandler{
		adminRepository:       config.AdminRepo,
		patientRepository:     config.PatientRepo,
		doctorRepository:      config.DoctorRepo,
		appointmentRepository: config.AppointmentRepo,
		email:                 config.Email,
	}
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
