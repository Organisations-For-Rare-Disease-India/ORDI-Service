package patient

import (
	"ORDI/cmd/web"
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/mojocn/base64Captcha"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("ORDI_CaptchaStore"))

func (s *patientHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var loginDetails struct {
		Email    string `schema:"email_id"`
		Password string `schema:"password"`
		Captcha  string `schema:"captcha"`
	}

	// Find patient from database
	patient, err := s.patientRepository.FindByField(ctx, "email_id", loginDetails.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if patient == nil {
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
	log.Printf("%+v",loginDetails)

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
	log.Print(userInput)
	log.Print(storedID)
	log.Print(isValid)
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
