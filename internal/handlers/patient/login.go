package patient

import (
	"ORDI/cmd/web"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/mojocn/base64Captcha"
	"golang.org/x/crypto/bcrypt"
)

func (s *patientHandler) Login(w http.ResponseWriter, r *http.Request) {
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
		log.Fatalf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("%+v ::::::::; loginDetails", loginDetails)
	captchaID , err := s.cache.Get(ctx,"1")
	if err != nil {
		log.Fatalf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	isValid := base64Captcha.DefaultMemStore.Verify(captchaID, loginDetails.Captcha, false)

	if isValid {
		log.Printf("%+v", isValid)
	} else {
		log.Printf("%+v", isValid)
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

func (s *patientHandler) GenerateCaptcha(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, captchaImage, err := s.generateCaptcha()
	err = s.cache.Set(ctx,"1",id,time.Duration(10)*time.Second)
	if err != nil {
		log.Fatalf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("%s", captchaImage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")

	json.NewEncoder(w).Encode(map[string]string{
		"captchaImage": captchaImage,
	})

}

func (s *patientHandler) generateCaptcha() (string, string, error) {
	captcha := base64Captcha.NewCaptcha(s.captchaDriver, s.captchaStore)
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		return "", "", err
	}
	return id, b64s, nil
}
