package verification

import (
	"ORDI/cmd/web"
	"ORDI/internal/cache"
	"ORDI/internal/repositories"
	"context"
	"net/http"

	"github.com/a-h/templ"
)

type VerificationInterface interface {
	Verify(http.ResponseWriter, *http.Request)
}

type verificationHandler struct {
	patientRepository repositories.Patient
	cache             cache.Cache
}

type VerificationConfig struct {
	PatientRepo repositories.Patient
	Cache       cache.Cache
}

func NewPatienVerificationtHandler(config VerificationConfig) VerificationInterface {
	return &verificationHandler{
		patientRepository: config.PatientRepo,
		cache:             config.Cache,
	}
}

func (s *verificationHandler) Verify(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	email, err := s.cache.Get(ctx, token)
	if err != nil || email == "" {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
	}

	// Display the login page
	go templ.Handler(web.LoginPage()).ServeHTTP(w, r)

	// Mark the user as verified
	patient, err := s.patientRepository.FindByField(ctx, "email", email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if patient != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	patient.Verified = true
	s.patientRepository.Save(context.Background(), patient)
}
