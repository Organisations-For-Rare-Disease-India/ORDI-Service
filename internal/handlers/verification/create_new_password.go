package verification

import (
	"ORDI/cmd/web"
	"net/http"

	"github.com/a-h/templ"
	"golang.org/x/crypto/bcrypt"
)

func (s *verificationHandler) CreateNewPassword(w http.ResponseWriter, r *http.Request) {
	// Accept password
	// Update password
	// Take to login page

	ctx := r.Context()

	var createPasswordDetails struct {
		Email    string `schema:"email_id"`
		Password string `schema:"password"`
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = decoder.Decode(&createPasswordDetails, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Find patient from database
	patient, err := s.patientRepository.FindByField(ctx, "email_id", createPasswordDetails.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if patient == nil {
		http.Error(w, "This user does not exist", http.StatusUnauthorized)
		return
	}

	// Update password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createPasswordDetails.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	patient.Password = string(hashedPassword)

	// Save updated details
	err = s.patientRepository.Save(ctx, patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Take to login
	templ.Handler(web.LoginPage()).ServeHTTP(w, r)

}
