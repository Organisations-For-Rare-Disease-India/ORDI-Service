package server

import (
	"ORDI/cmd/web"
	"ORDI/internal/handlers/patient"
	"ORDI/internal/handlers/verification"
	"ORDI/internal/repositories"
	"encoding/json"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterPatientRoutes(r *chi.Mux, patientRepository repositories.Patient) {
	patientHandler := patient.NewPatientHandler(patient.PatientHandlerConfig{
		PatientRepo: patientRepository,
		Cache:       s.cache,
		Email:       s.email,
	})

	r.Get("/patient_signup", templ.Handler(web.PatientSignupPage()).ServeHTTP)
	r.Post("/patient_submit", patientHandler.Signup)
	r.Post("/patient_login", patientHandler.Login)
	r.Get("/patient_dashboard", templ.Handler(web.PatientDashboardPage()).ServeHTTP)
	r.Get("/appointments", patientHandler.Appointment)
	r.Get("/generate_captcha",patientHandler.GenerateCaptcha)
}

func (s *Server) RegisterVerificationRoutes(r *chi.Mux, patientRepository repositories.Patient) {
	verificationHandler := verification.NewVerificationHandler(verification.VerificationConfig{
		PatientRepo: patientRepository,
		Cache:       s.cache,
	})
	r.Get("/verify_patient", verificationHandler.VerifyPatient)
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", templ.Handler(web.HomePage()).ServeHTTP)

	r.Get("/health", s.healthHandler)

	// Static file serving
	fileServer := http.FileServer(http.FS(web.Files))
	r.Handle("/assets/*", fileServer)

	r.Get("/login", templ.Handler(web.LoginPage()).ServeHTTP)
	r.Get("/signup", templ.Handler(web.SignupPage()).ServeHTTP)
	r.Get("/signup_steps", templ.Handler(web.SignupStepsPage()).ServeHTTP)
	r.Get("/terms_and_conditions", templ.Handler(web.TermsAndConditionsPage()).ServeHTTP)

	// Patient specific handlers
	patientRepository := repositories.NewPatientRepository(s.db)
	s.RegisterPatientRoutes(r, patientRepository)

	// Verification handler

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
