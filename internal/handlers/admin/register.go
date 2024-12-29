package admin

import (
	"ORDI/cmd/web"
	"ORDI/internal/utils"
	"net/http"

	"github.com/a-h/templ"
	"golang.org/x/crypto/bcrypt"
)

type RegisterDetails struct {
	Email    string `schema:"email_id"`
	Password string `schema:"password"`
}

func (a *adminHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var credentials RegisterDetails

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Decode the form data into credentials
	err = decoder.Decode(&credentials, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Find admin from database
	admin, err := a.adminRepository.FindByField(ctx, "email_id", credentials.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if admin == nil {
		http.Error(w, "You have not been added to the system", http.StatusUnauthorized)
		return
	}

	// Use the password provided by admin for further login
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	admin.Password = string(hashedPassword)
	// Admin has setup their credentials, hence now verified.
	admin.Verified = true

	// Save updated details
	err = a.adminRepository.Save(ctx, admin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Once the admin sets up the password
	// Take them to login
	templ.Handler(web.AdminLoginPage(utils.AdminLoginSubmit, false)).ServeHTTP(w, r)
}
