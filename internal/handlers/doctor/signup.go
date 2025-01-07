package doctor

import (
	"ORDI/cmd/web"
	"ORDI/internal/messages"
	"ORDI/internal/models"
	"ORDI/internal/utils"
	"fmt"
	"time"

	"net/http"

	"github.com/a-h/templ"
	"golang.org/x/crypto/bcrypt"
)

func (d *doctorHandler) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var doctor models.Doctor

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Decode the form data into the Patient struct
	err = decoder.Decode(&doctor, r.PostForm)
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
	d.cache.Set(ctx, token, doctor.Email, 15*time.Minute)

	// Generate html body for verification mail
	htmlBody := utils.GenerateVerificationHTML(ctx, token, utils.DoctorVerifyNew, "Thank you for Registering to ORDI!")
	d.notificationRepository.Save(ctx, &models.Notification{
		UserEmail: doctor.Email,
		Message:   "Thank you for Registering to ORDI!",
		SentTime:  time.Now().String(),
		IsRead:    false,
	})
	err = d.email.SendEmail(doctor.Email, "ORDI Email Verification", htmlBody, nil, "", "text/html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(doctor.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	doctor.Password = string(hashedPassword)

	// Store patient on database
	doctor.Verified = false // Mark verified as false
	err = d.doctorRepository.Save(ctx, &doctor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Prepare attachement to send to ORDI
	pdfBuffer, err := utils.DoctorToPDF(doctor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Email details
	subject := "New Doctor registration"
	bodyType := "text/plain"
	body := "A new doctor has signed up. Please find the attached PDF with the details"
	to := "jhavedantamay@gmail.com"
	attachementName := fmt.Sprintf("%s.pdf", doctor.FirstName)
	err = d.email.SendEmail(to, subject, body, pdfBuffer, attachementName, bodyType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the submit page with a success message
	// Render the template with the success flag
	templ.Handler(web.SubmitPage(messages.SubmitMessage{
		Title:   "Successfully uploaded",
		Message: "A verification email has been sent to your email address. Please check your inbox to verify your account.",
	})).ServeHTTP(w, r)

}
