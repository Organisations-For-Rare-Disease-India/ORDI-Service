package patient

import (
	"ORDI/cmd/web"
	"net/http"

	"github.com/a-h/templ"
	"golang.org/x/crypto/bcrypt"
)

func (s *patientHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var loginDetails struct {
		Email    string `schema:"email"`
		Password string `schema:"password"`
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = decoder.Decode(&loginDetails, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Find patient from database
	patient, err := s.patientRepository.FindByField(ctx, "email", loginDetails.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if patient != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(patient.Password), []byte(loginDetails.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Show dashboard if  password matches
	templ.Handler(web.PatientDashboardPage()).ServeHTTP(w, r)
}
