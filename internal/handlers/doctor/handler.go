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
	doctorRepository repositories.Repository[models.Doctor]
	notificationRepository repositories.Repository[models.Notification]
	cache            cache.Cache
	email            email.Email
}

type DoctorHandlerConfig struct {
	DoctorRepo repositories.Repository[models.Doctor]
	NotficationRepo repositories.Repository[models.Notification]
	Cache      cache.Cache
	Email      email.Email
}

func NewDoctorHandler(config DoctorHandlerConfig) Doctor {
	return &doctorHandler{
		doctorRepository: config.DoctorRepo,
		notificationRepository: config.NotficationRepo,
		cache:            config.Cache,
		email:            config.Email,
	}
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
