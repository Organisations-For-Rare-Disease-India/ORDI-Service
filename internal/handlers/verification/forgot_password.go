package verification

import (
	"ORDI/cmd/web"
	"ORDI/internal/messages"
	"ORDI/internal/utils"
	"net/http"
	"time"

	"github.com/a-h/templ"
)

func (s *verificationHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	// Receive an email from form
	// Verify that this email is an existing user
	// Generate token
	// Send email
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

	// Find patient from database
	patient, err := s.patientRepository.FindByField(ctx, "email_id", forgotPasswordDetails.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if patient == nil {
		http.Error(w, "This user does not exist", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := utils.GenerateVerificationtoken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store token in cache
	err = s.cache.Set(ctx, token, patient.Email, 15*time.Minute)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate html body for verification mail
	htmlBody := utils.GenerateVerificationHTML(ctx, token, "verify_existing_patient", "You told us you forgot your password.")
	err = s.email.SendEmail(patient.Email, "Reset your Password", htmlBody, nil, "", "text/html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Display the message that email has been sent
	templ.Handler(web.SubmitPage(messages.SubmitMessage{
		Title:   "Forgot Password Request sent",
		Message: "A verification email has been sent to your email address. Please click on the link provided to proceed ahead.",
	})).ServeHTTP(w, r)
}
