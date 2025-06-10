package patient

import (
	"ORDI/internal/cache"
	"ORDI/internal/email"
	"ORDI/internal/models"
	"ORDI/internal/repositories"
	"net/http"

	"github.com/mojocn/base64Captcha"
)

type Patient interface {
	Signup(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
	Appointment(http.ResponseWriter, *http.Request)
	GenerateCaptcha(http.ResponseWriter, *http.Request)
	VerifyCaptcha(http.ResponseWriter, *http.Request)
	Dashboard(http.ResponseWriter, *http.Request)
	Profile(http.ResponseWriter, *http.Request)
	GetAppointments(w http.ResponseWriter, r *http.Request)
}

type patientHandler struct {
	patientRepository      repositories.Repository[models.Patient]
	appointmentRepository  repositories.Repository[models.Appointment]
	notificationRepository repositories.Repository[models.Notification]
	cache                  cache.Cache
	email                  email.Email
	captchaStore           base64Captcha.Store
	captchaDriver          *base64Captcha.DriverDigit
}
type PatientHandlerConfig struct {
	PatientRepo      repositories.Repository[models.Patient]
	NotificationRepo repositories.Repository[models.Notification]
	AppointmentRepo  repositories.Repository[models.Appointment]
	Cache            cache.Cache
	Email            email.Email
	CaptchaStore     base64Captcha.Store
	CaptchaDriver    base64Captcha.DriverDigit
}

func NewPatientHandler(config PatientHandlerConfig) Patient {
	return &patientHandler{
		patientRepository:      config.PatientRepo,
		notificationRepository: config.NotificationRepo,
		appointmentRepository:  config.AppointmentRepo,
		cache:                  config.Cache,
		email:                  config.Email,
		captchaStore:           base64Captcha.DefaultMemStore,
		captchaDriver:          base64Captcha.NewDriverDigit(50, 120, 4, 0.7, 80),
	}
}
