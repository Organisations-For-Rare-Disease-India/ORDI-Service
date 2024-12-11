package doctor

import (
	"ORDI/cmd/web"
	"net/http"

	"github.com/a-h/templ"
	"golang.org/x/crypto/bcrypt"
)

func (d *doctorHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var loginDetails struct {
		Email    string `schema:"email_id"`
		Password string `schema:"password"`
		Captcha  string `schema:"captcha"`
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

	// Find doctor from database
	doctor, err := d.doctorRepository.FindByField(ctx, "email_id", loginDetails.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if doctor == nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(doctor.Password), []byte(loginDetails.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Show dashboard if  password matches
	templ.Handler(web.PatientDashboardPage()).ServeHTTP(w, r)
}
