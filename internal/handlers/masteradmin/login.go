package masteradmin

import (
	"ORDI/internal/handlers/token"
	"ORDI/internal/utils"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type LoginDetails struct {
	Email    string `schema:"email_id"`
	Password string `schema:"password"`
}

func (a *masterAdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var credentials LoginDetails

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
	admin, err := a.masterAdminRepository.FindByField(ctx, "email_id", credentials.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if admin == nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	cookie, err := token.CreateTokenCookie(credentials.Email)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If the user logs in with correct credentials, this handler will set a cookie on the client side with JWT value.
	// Once a cookie is set on client, it is sent along with every request henceforth.
	http.SetCookie(w, cookie)

	// Take admin to dashboard if credentials match
	http.Redirect(w, r, utils.AdminCreate, http.StatusSeeOther)
}
