package verification

import (
	"ORDI/cmd/web"
	"ORDI/internal/cache"
	"ORDI/internal/email"
	"ORDI/internal/models"
	"ORDI/internal/repositories"
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/schema"
)

const (
	TokenHeader                = "token"
	TokenRequiredMessage       = "Token is required"
	InvalidRequestMessage      = "Invalid request"
	UnauthorizedRequestMessage = "Unauthorized request"
	InvalidCredentialsMessage  = "Invalid username or password"
)

var decoder = schema.NewDecoder()

type VerificationInterface interface {
	VerifyNewPatient(http.ResponseWriter, *http.Request)
	VerifyExistingPatient(http.ResponseWriter, *http.Request)
	CreateNewPassword(http.ResponseWriter, *http.Request)
	ForgotPassword(http.ResponseWriter, *http.Request)
}

type verificationHandler struct {
	patientRepository repositories.Repository[models.PatientInfo]
	cache             cache.Cache
	email             email.Email
}

type VerificationConfig struct {
	PatientRepo repositories.Repository[models.PatientInfo]
	Cache       cache.Cache
	EmailID     email.Email
}

func NewVerificationHandler(config VerificationConfig) VerificationInterface {
	return &verificationHandler{
		patientRepository: config.PatientRepo,
		cache:             config.Cache,
		email:             config.EmailID,
	}
}

func (s *verificationHandler) VerifyNewPatient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	email, err := s.verifyRequest(ctx, r, TokenHeader)
	if err != nil {
		http.Error(w, TokenRequiredMessage, http.StatusUnauthorized)
	}
	// User is verified

	// Display the login page
	templ.Handler(web.LoginPage()).ServeHTTP(w, r)

	// Mark the user as verified in the background
	go func() {
		patient, err := s.patientRepository.FindByField(ctx, "email_id", email)
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

func (s *verificationHandler) VerifyExistingPatient(w http.ResponseWriter, r *http.Request) {
	// Used incase patient has forgotten the password
	ctx := r.Context()

	_, err := s.verifyRequest(ctx, r, TokenHeader)
	if err != nil {
		http.Error(w, TokenRequiredMessage, http.StatusUnauthorized)
	}
	// User is verified

	// Display the Create new password page
	templ.Handler(web.CreateNewPasswordPage()).ServeHTTP(w, r)
}

// verifyRequest verifies that the token is valid, and returns the associated email with it
func (s *verificationHandler) verifyRequest(ctx context.Context, r *http.Request, header string) (string, error) {
	token := r.URL.Query().Get(header)
	if token == "" {
		return "", fmt.Errorf("invalid request")
	}

	email, err := s.cache.Get(ctx, token)
	if err != nil || email == "" {
		return "", fmt.Errorf("unauthorised request")
	}
	return email, nil
}
