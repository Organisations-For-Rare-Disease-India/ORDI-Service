package patient

import (
	"ORDI/internal/cache"
	"ORDI/internal/email"
	"ORDI/internal/repositories"
	"net/http"

	"github.com/mojocn/base64Captcha"
)

type PatientHandlerInterface interface {
	Signup(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
	Appointment(http.ResponseWriter, *http.Request)
	GenerateCaptcha(http.ResponseWriter, *http.Request)
}

type patientHandler struct {
	patientRepository repositories.Patient
	cache             cache.Cache
	email             email.Email
	captchaStore base64Captcha.Store
	captchaDriver  *base64Captcha.DriverDigit
}
// base64Captcha.NewDriverDigit(100, 240, 4, 0.7, 80)
type PatientHandlerConfig struct {
	PatientRepo repositories.Patient
	Cache       cache.Cache
	Email       email.Email
	CaptchaStore base64Captcha.Store
	CaptchaDriver  base64Captcha.DriverDigit
}

func NewPatientHandler(config PatientHandlerConfig) PatientHandlerInterface {
	return &patientHandler{
		patientRepository: config.PatientRepo,
		cache:             config.Cache,
		email:             config.Email,
		captchaStore: base64Captcha.DefaultMemStore,
		captchaDriver: base64Captcha.NewDriverDigit(100, 240, 4, 0.7, 80),
	}
}
