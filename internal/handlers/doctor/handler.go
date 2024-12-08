package doctor

import (
	"ORDI/internal/cache"
	"ORDI/internal/email"
	"ORDI/internal/models"
	"ORDI/internal/repositories"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

const JWT_TOKEN_HEADER = "token"

type Doctor interface {
	Signup(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
	Appointment(http.ResponseWriter, *http.Request)
	Dashboard(http.ResponseWriter, *http.Request)
	Profile(http.ResponseWriter, *http.Request)
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

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
