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

type doctorVerification struct {
	repository repositories.Repository[models.Doctor]
	cache      cache.Cache
	email      email.Email
}

type DoctorVerificationConfig struct {
	Repository repositories.Repository[models.Doctor]
	Cache      cache.Cache
	Email      email.Email
}

func NewDoctorVerificationHandler(config DoctorVerificationConfig) Verification {
	return &doctorVerification{
		repository: config.Repository,
		cache:      config.Cache,
		email:      config.Email,
	}
}

func (d *doctorVerification) ForgotPassword(w http.ResponseWriter, r *http.Request) {
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

	doctor, err := d.repository.FindByField(ctx, "email_id", forgotPasswordDetails.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if doctor == nil {
		http.Error(w, "This user does not exist", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateVerificationtoken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = d.cache.Set(ctx, token, doctor.Email, 15*time.Minute)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	htmlBody := utils.GenerateVerificationHTML(ctx, token, utils.DoctorVerifyExisting, "You told us you forgot your password.")
	err = d.email.SendEmail(doctor.Email, "Reset your Password", htmlBody, nil, "", "text/html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	templ.Handler(web.SubmitPage(emailMessage)).ServeHTTP(w, r)
}

func (d *doctorVerification) CreateNewPassword(w http.ResponseWriter, r *http.Request) {
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

	// Find doctor from database
	doctor, err := d.repository.FindByField(ctx, "email_id", createPasswordDetails.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if doctor == nil {
		http.Error(w, "This user does not exist", http.StatusUnauthorized)
		return
	}

	// Update password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createPasswordDetails.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	doctor.Password = string(hashedPassword)

	// Save updated details
	err = d.repository.Save(ctx, doctor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Take to login
	templ.Handler(web.LoginPage(utils.DoctorLoginSubmit, utils.DoctorForgotPasswordScreen, utils.DoctorSignupSteps)).ServeHTTP(w, r)
}

func (d *doctorVerification) VerifyNewUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.URL.Query().Get(TokenHeader)
	if token == "" {
		http.Error(w, TokenRequiredMessage, http.StatusUnauthorized)
	}

	email, err := d.cache.Get(ctx, token)
	if err != nil || email == "" {
		http.Error(w, TokenRequiredMessage, http.StatusUnauthorized)
	}
	if err != nil {
		http.Error(w, TokenRequiredMessage, http.StatusUnauthorized)
	}
	// User is verified

	// Display the login page
	templ.Handler(web.LoginPage(utils.DoctorLoginSubmit, utils.DoctorForgotPasswordScreen, utils.DoctorSignupSteps)).ServeHTTP(w, r)

	// Mark the user as verified in the background
	go func() {
		doctor, err := d.repository.FindByField(ctx, "email_id", email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if doctor == nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		doctor.Verified = true
		d.repository.Save(context.Background(), doctor)
	}()
}

func (d *doctorVerification) VerifyExistingUser(w http.ResponseWriter, r *http.Request) {
	// Used incase user has forgotten the password
	ctx := r.Context()
	token := r.URL.Query().Get(TokenHeader)
	if token == "" {
		http.Error(w, TokenRequiredMessage, http.StatusUnauthorized)
	}

	email, err := d.cache.Get(ctx, token)
	if err != nil || email == "" {
		http.Error(w, TokenRequiredMessage, http.StatusUnauthorized)
	}
	// User is verified

	// Display the Create new password page
	templ.Handler(web.CreateNewPasswordPage(utils.DoctorNewPassword)).ServeHTTP(w, r)
}
