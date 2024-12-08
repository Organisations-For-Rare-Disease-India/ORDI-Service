package server

import (
	"ORDI/cmd/web"
	"ORDI/internal/handlers/doctor"
	"ORDI/internal/handlers/patient"
	"ORDI/internal/handlers/verification"
	"ORDI/internal/models"
	"ORDI/internal/repositories"
	"ORDI/internal/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterPatientRoutes(r *chi.Mux, patientRepository repositories.Repository[models.Patient]) {
	patientHandler := patient.NewPatientHandler(patient.PatientHandlerConfig{
		PatientRepo: patientRepository,
		Cache:       s.cache,
		Email:       s.email,
	})

	r.Get(utils.PatientLoginScreen, templ.Handler(web.LoginPage(utils.PatientLoginSubmit, utils.PatientForgotPasswordScreen, utils.PatientSignupSteps)).ServeHTTP)
	r.Post(utils.PatientSignupSubmit, patientHandler.Signup)
	r.Post(utils.PatientLoginSubmit, patientHandler.Login)
	r.Get(utils.PatientDashboard, patientHandler.Dashboard)
	r.Get(utils.PatientAppointments, patientHandler.Appointment)
	r.Get(utils.PatientForgotPasswordScreen, templ.Handler(web.ForgotPasswordPage(utils.PatientForgotPasswordSubmit)).ServeHTTP)
	r.Get(utils.PatientSignupSteps, templ.Handler(web.SignupStepsPage(utils.CreatePatientSignupStepsMessage(), utils.PatientSignupForm)).ServeHTTP)
	r.Get(utils.PatientSignupForm, templ.Handler(web.PatientSignupFormPage(utils.PatientSignupSubmit)).ServeHTTP)
	r.Get(utils.PatientProfile, patientHandler.Profile)

	patientVerificationHandler := verification.NewPatientVerificationHandler(verification.PatientVerificationConfig{
		Repository: patientRepository,
		Cache:      s.cache,
		Email:      s.email,
	})
	r.Get(utils.PatientVerifyNew, patientVerificationHandler.VerifyNewUser)
	r.Get(utils.PatientVerifyExisting, patientVerificationHandler.VerifyExistingUser)
	r.Post(utils.PatientNewPassword, patientVerificationHandler.CreateNewPassword)
	r.Post(utils.PatientForgotPasswordSubmit, patientVerificationHandler.ForgotPassword)

}

func (s *Server) RegisterDoctorRoutes(r *chi.Mux, doctorRepository repositories.Repository[models.Doctor]) {
	doctorHandler := doctor.NewDoctorHandler(doctor.DoctorHandlerConfig{
		DoctorRepo: doctorRepository,
		Cache:      s.cache,
		Email:      s.email,
	})

	r.Get(utils.DoctorLoginScreen, templ.Handler(web.LoginPage(utils.DoctorLoginSubmit, utils.DoctorForgotPasswordScreen, utils.DoctorSignupSteps)).ServeHTTP)
	r.Post(utils.DoctorSignupSubmit, doctorHandler.Signup)
	r.Post(utils.DoctorLoginSubmit, doctorHandler.Login)
	r.Get(utils.DoctorDashboard, doctorHandler.Dashboard)
	r.Get(utils.DoctorAppointments, doctorHandler.Appointment)
	r.Get(utils.DoctorForgotPasswordScreen, templ.Handler(web.ForgotPasswordPage(utils.DoctorForgotPasswordSubmit)).ServeHTTP)
	r.Get(utils.DoctorSignupSteps, templ.Handler(web.SignupStepsPage(utils.CreatePatientSignupStepsMessage(), utils.DoctorSignupForm)).ServeHTTP)
	r.Get(utils.DoctorSignupForm, templ.Handler(web.DoctorSignupFormPage(utils.DoctorSignupSubmit)).ServeHTTP)
	r.Get(utils.DoctorProfile, doctorHandler.Profile)

	doctorVerificationHandler := verification.NewDoctorVerificationHandler(verification.DoctorVerificationConfig{
		Repository: doctorRepository,
		Cache:      s.cache,
		Email:      s.email,
	})
	r.Get(utils.DoctorVerifyNew, doctorVerificationHandler.VerifyNewUser)
	r.Get(utils.DoctorVerifyExisting, doctorVerificationHandler.VerifyExistingUser)
	r.Post(utils.DoctorNewPassword, doctorVerificationHandler.CreateNewPassword)
	r.Post(utils.DoctorForgotPasswordSubmit, doctorVerificationHandler.ForgotPassword)

}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", templ.Handler(web.HomePage()).ServeHTTP)

	r.Get("/health", s.healthHandler)

	// Static file serving
	fileServer := http.FileServer(http.FS(web.Files))
	r.Handle("/assets/*", fileServer)

	r.Get(utils.HomeLogin, templ.Handler(web.ChooseRolePage(utils.DoctorLoginScreen, utils.PatientLoginScreen)).ServeHTTP)
	r.Get(utils.HomeSignup, templ.Handler(web.ChooseRolePage(utils.DoctorSignupSteps, utils.PatientSignupSteps)).ServeHTTP)

	// Patient specific handlers
	patientRepository := repositories.NewPatientRepository(s.db)
	s.RegisterPatientRoutes(r, patientRepository)

	// Doctor specific handlers
	doctorRepository := repositories.NewDoctorRepository(s.db)
	s.RegisterDoctorRoutes(r, doctorRepository)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
