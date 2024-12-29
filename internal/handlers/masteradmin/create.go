package masteradmin

import (
	"ORDI/cmd/web"
	"ORDI/internal/messages"
	"ORDI/internal/models"
	"ORDI/internal/utils"
	"fmt"

	"net/http"

	"github.com/a-h/templ"
	"golang.org/x/crypto/bcrypt"
)

func (a *masterAdminHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var admin models.Admin

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Decode the form data into the Admin struct
	err = decoder.Decode(&admin, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate html body for verification mail
	htmlBody := utils.GenerateWelcomeAdminHTML(ctx, utils.AdminSetCredentials, admin.Email)
	// Prepare attachement to send to ORDI and new admin
	attachementName := fmt.Sprintf("%s.pdf", admin.FirstName)
	pdfBuffer, err := utils.AdminToPDF(admin)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = a.email.SendEmail(admin.Email, "ORDI Onboarding Email", htmlBody, pdfBuffer, attachementName, "text/html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// NOTE : This is a temporary password, it would be changed the admin
	adminPassword := "temporary"
	//Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	admin.Password = string(hashedPassword)

	// Store admin on database
	admin.Verified = false // Mark verified as false
	a.adminRepository.Save(ctx, &admin)

	// Render the submit page with a success message
	// Render the template with the success flag
	templ.Handler(web.SubmitPage(messages.SubmitMessage{
		Title:   "Created new admin",
		Message: "An email has been sent to the new admin. Please request them to check their account and onboard.",
	})).ServeHTTP(w, r)
}
