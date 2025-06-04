package patient

import (
	"ORDI/cmd/web"
	"ORDI/internal/messages"
	"ORDI/internal/models"
	"ORDI/internal/utils"
	"fmt"
	"time"

	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/schema"
	"golang.org/x/crypto/bcrypt"
)

var decoder = schema.NewDecoder()

func (s *patientHandler) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var patient models.Patient

	// Render the submit page with a success message

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Decode the form data into the Patient struct
	err = decoder.Decode(&patient, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate token
	token, err := utils.GenerateVerificationtoken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store token in cache
	go s.cache.Set(ctx, token, patient.Email, 15*time.Minute)

	// Generate html body for verification mail
	htmlBody := utils.GenerateVerificationHTML(ctx, token, utils.PatientVerifyNew, "Thank you for Registering to ORDI!")
	s.email.SendEmail(patient.Email, "ORDI Email Verification", htmlBody, nil, "", "text/html")

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(patient.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	patient.Password = string(hashedPassword)

	// Store patient on database
	patient.Verified = false // Mark verified as false
	err = s.patientRepository.Save(ctx, &patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare attachement to send to ORDI
	pdfBuffer, err := utils.PatientToPDF(patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Email details
	subject := "New Patient registration"
	bodyType := "text/plain"
	body := "A new patient has signed up. Please find the attached PDF with the details"
	to := "jhavedantamay@gmail.com"
	attachementName := fmt.Sprintf("%s.pdf", patient.FirstName)
	err = s.email.SendEmail(to, subject, body, pdfBuffer, attachementName, bodyType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.notificationRepository.Save(ctx, &models.Notification{
		UserID:    patient.ID,
		UserType:  models.Role.String(models.PatientType),
		Title:     "Welcome",
		UserEmail: patient.Email,
		Message:   "Thank you for registration to ORDI!",
		IsRead:    false,
		SentTime:  time.Now(),
	})

	// Render the template with the success flag
	templ.Handler(web.SubmitPage(messages.SubmitMessage{
		Title:   "Successfully uploaded",
		Message: "A verification email has been sent to your email address. Please check your inbox to verify your account.",
	})).ServeHTTP(w, r)
}
