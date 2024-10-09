package server

import (
	"ORDI/cmd/web"
	"encoding/json"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", templ.Handler(web.HomePage()).ServeHTTP)

	r.Get("/health", s.healthHandler)

	fileServer := http.FileServer(http.FS(web.Files))
	r.Handle("/assets/*", fileServer)
	r.Get("/login", templ.Handler(web.LoginPage()).ServeHTTP)
	r.Get("/signup", templ.Handler(web.SignupPage()).ServeHTTP)
	r.Get("/patient_signup", templ.Handler(web.PatientSignupPage()).ServeHTTP)
	r.Get("/signup_steps", templ.Handler(web.SignupStepsPage()).ServeHTTP)
	r.Post("/patient_submit", s.PatientSignupFormHandler)
	r.Post("/patient_login", s.PatientLoginHandler)
	r.Get("/patient_dashboard", templ.Handler(web.PatientDashboardPage()).ServeHTTP)
	r.Get("/verify", s.EmailVerificationHandler)
	r.Get("/appointments", s.AppointmentHandler)
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
