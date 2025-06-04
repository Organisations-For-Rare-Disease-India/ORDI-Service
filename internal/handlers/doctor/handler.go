package doctor

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

var decoder = schema.NewDecoder()

type Doctor interface {
	Signup(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
	Appointment(http.ResponseWriter, *http.Request)
	Dashboard(http.ResponseWriter, *http.Request)
	Profile(http.ResponseWriter, *http.Request)
}

type doctorHandler struct {
	doctorRepository       repositories.Repository[models.Doctor]
	appointmentRepository  repositories.Repository[models.Appointment]
	notificationRepository repositories.Repository[models.Notification]
	cache                  cache.Cache
	email                  email.Email
}

type DoctorHandlerConfig struct {
	DoctorRepo       repositories.Repository[models.Doctor]
	AppointmentRepo  repositories.Repository[models.Appointment]
	NotificationRepo repositories.Repository[models.Notification]
	Cache            cache.Cache
	Email            email.Email
}

func NewDoctorHandler(config DoctorHandlerConfig) Doctor {
	return &doctorHandler{
		doctorRepository:       config.DoctorRepo,
		appointmentRepository:  config.AppointmentRepo,
		notificationRepository: config.NotificationRepo,
		cache:                  config.Cache,
		email:                  config.Email,
	}
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
