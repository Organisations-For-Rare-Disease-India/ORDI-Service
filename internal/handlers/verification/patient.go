package verification

import (
	"ORDI/cmd/web"
	"ORDI/internal/cache"
	"ORDI/internal/email"
	"ORDI/internal/models"
	"ORDI/internal/repositories"
	"ORDI/internal/utils"
	"context"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"golang.org/x/crypto/bcrypt"
)

type patientVerification struct {
	repository repositories.Repository[models.Patient]
	cache      cache.Cache
	email      email.Email
}

type PatientVerificationConfig struct {
	Repository repositories.Repository[models.Patient]
	Cache      cache.Cache
	Email      email.Email
}

func NewPatientVerificationHandler(config PatientVerificationConfig) Verification {
	return &patientVerification{
		repository: config.Repository,
		cache:      config.Cache,
		email:      config.Email,
	}
}

func (p *patientVerification) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var forgotPasswordDetails struct {
		Email string `schema:"email_id"`
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = decoder.Decode(&forgotPasswordDetails, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	patient, err := p.repository.FindByField(ctx, "email_id", forgotPasswordDetails.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if patient == nil {
		http.Error(w, "This user does not exist", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateVerificationtoken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = p.cache.Set(ctx, token, patient.Email, 15*time.Minute)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	htmlBody := utils.GenerateVerificationHTML(ctx, token, utils.PatientVerifyExisting, "You told us you forgot your password.")
	err = p.email.SendEmail(patient.Email, "Reset your Password", htmlBody, nil, "", "text/html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	templ.Handler(web.SubmitPage(emailMessage)).ServeHTTP(w, r)
}

func (p *patientVerification) CreateNewPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var createPasswordDetails struct {
		Email    string `schema:"email_id"`
		Password string `schema:"password"`
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = decoder.Decode(&createPasswordDetails, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Find patient from database
	patient, err := p.repository.FindByField(ctx, "email_id", createPasswordDetails.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if patient == nil {
		http.Error(w, "This user does not exist", http.StatusUnauthorized)
		return
	}

	// Update password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createPasswordDetails.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	patient.Password = string(hashedPassword)

	// Save updated details
	err = p.repository.Save(ctx, patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Take to login
	templ.Handler(web.LoginPage(utils.PatientLoginSubmit, utils.PatientForgotPasswordScreen, utils.PatientSignupSteps)).ServeHTTP(w, r)
}

func (p *patientVerification) VerifyNewUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.URL.Query().Get(TokenHeader)
	if token == "" {
		http.Error(w, TokenRequiredMessage, http.StatusUnauthorized)
	}

	email, err := p.cache.Get(ctx, token)
	if err != nil || email == "" {
		http.Error(w, TokenRequiredMessage, http.StatusUnauthorized)
	}
	if err != nil {
		http.Error(w, TokenRequiredMessage, http.StatusUnauthorized)
	}
	// User is verified

	// Display the login page
	templ.Handler(web.LoginPage(utils.PatientLoginSubmit, utils.PatientForgotPasswordScreen, utils.PatientSignupSteps)).ServeHTTP(w, r)

	// Mark the user as verified in the background
	go func() {
		patient, err := p.repository.FindByField(ctx, "email_id", email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if patient == nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		patient.Verified = true
		p.repository.Save(context.Background(), patient)
	}()
}

func (p *patientVerification) VerifyExistingUser(w http.ResponseWriter, r *http.Request) {
	// Used incase user has forgotten the password
	ctx := r.Context()

	token := r.URL.Query().Get(TokenHeader)
	if token == "" {
		http.Error(w, TokenRequiredMessage, http.StatusUnauthorized)
	}

	email, err := p.cache.Get(ctx, token)
	if err != nil || email == "" {
		http.Error(w, TokenRequiredMessage, http.StatusUnauthorized)
	}
	// User is verified

	// Display the Create new password page
	templ.Handler(web.CreateNewPasswordPage()).ServeHTTP(w, r)
}
