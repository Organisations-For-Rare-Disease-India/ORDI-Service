package patient

import (
	"ORDI/internal/handlers/token"
	"ORDI/internal/utils"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/mojocn/base64Captcha"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("ORDI_CaptchaStore"))

type LoginDetails struct {
	Email    string `schema:"email_id"`
	Password string `schema:"password"`
	Captcha  string `schema:"captcha"`
}

func (p *patientHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var credentials LoginDetails

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Decode the form data into the Patient struct
	err = decoder.Decode(&credentials, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Find patient from database
	patient, err := p.patientRepository.FindByField(ctx, "email_id", credentials.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if patient == nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(patient.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	cookie, err := token.CreateTokenCookie(patient.ID, credentials.Email)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If the user logs in with correct credentials, this handler will set a cookie on the client side with JWT value.
	// Once a cookie is set on client, it is sent along with every request henceforth.
	http.SetCookie(w, cookie)

	// Take patient to dashboard if credentials match
	http.Redirect(w, r, utils.PatientDashboard, http.StatusSeeOther)
}

func (s *patientHandler) VerifyCaptcha(w http.ResponseWriter, r *http.Request) {

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
	session, err := store.Get(r, "captcha-session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	// Get the stored CAPTCHA ID
	storedID, ok := session.Values["captcha_id"].(string)
	if !ok || storedID == "" {
		http.Error(w, "CAPTCHA not found in session", http.StatusBadRequest)
		return
	}
	userInput := loginDetails.Captcha

	isValid := s.captchaStore.Verify(storedID, userInput, true)
	if !isValid {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<div id="error-message" class="mb-6 text-sm text-red-600 font-medium">Wrong CAPTCHA</div>`))
		return
	}
}

func (s *patientHandler) GenerateCaptcha(w http.ResponseWriter, r *http.Request) {

	// Generate the CAPTCHA
	id, captchaImage, err := s.generateCaptchaImage()
	if err != nil {
		http.Error(w, "Failed to generate CAPTCHA", http.StatusInternalServerError)
		return
	}

	// Save the CAPTCHA ID in a session
	session, err := store.Get(r, "captcha-session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}
	session.Values["captcha_id"] = id

	// Save the session (write it to the client as a cookie)
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<img id="captcha-image" src="%s" alt="CAPTCHA">`, captchaImage)))
}

func (s *patientHandler) generateCaptchaImage() (string, string, error) {
	captcha := base64Captcha.NewCaptcha(s.captchaDriver, s.captchaStore)
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		return "", "", err
	}
	return id, b64s, nil
}
