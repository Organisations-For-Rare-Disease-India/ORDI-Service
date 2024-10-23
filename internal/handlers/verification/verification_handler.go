package verification

import (
	"ORDI/cmd/web"
	"ORDI/internal/cache"
	"ORDI/internal/repositories"
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

type VerificationInterface interface {
	VerifyPatient(http.ResponseWriter, *http.Request)
}

type verificationHandler struct {
	patientRepository repositories.Patient
	cache             cache.Cache
}

type VerificationConfig struct {
	PatientRepo repositories.Patient
	Cache       cache.Cache
}

func NewVerificationHandler(config VerificationConfig) VerificationInterface {
	return &verificationHandler{
		patientRepository: config.PatientRepo,
		cache:             config.Cache,
	}
}

func (s *verificationHandler) VerifyPatient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	email, err := s.verifyRequest(ctx, r)
	if err != nil {
		http.Error(w, "Token is required", http.StatusUnauthorized)
	}
	// User is verified

	// Display the login page
	templ.Handler(web.LoginPage()).ServeHTTP(w, r)

	// Mark the user as verified in the background
	go func() {
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
	}()
}

// verifyRequest verifies that the token is valid, and returns the associate email with it
func (s *verificationHandler) verifyRequest(ctx context.Context, r *http.Request) (string, error) {
	token := r.URL.Query().Get("token")
	if token == "" {
		return "", fmt.Errorf("invalid request")
	}

	email, err := s.cache.Get(ctx, token)
	if err != nil || email == "" {
		return "", fmt.Errorf("unauthorised request")
	}
	return email, nil
}
