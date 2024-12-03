package patient

import (
	"ORDI/cmd/web"
	"ORDI/internal/messages"
	"ORDI/internal/models"
	"ORDI/internal/utils"
	"bytes"
	"fmt"
	"time"

	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/schema"
	"golang.org/x/crypto/bcrypt"
)

var decoder = schema.NewDecoder()

func (s *patientHandler) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var patient models.PatientInfo

	// Render the submit page with a success message
	// Render the template with the success flag
	go templ.Handler(web.SubmitPage(messages.SubmitMessage{
		Title:   "Successfully uploaded",
		Message: "A verification email has been sent to your email address. Please check your inbox to verify your account.",
	})).ServeHTTP(w, r)

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Decode the form data into the Patient struct
	err = decoder.Decode(&patient, r.PostForm)
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
	go s.cache.Set(ctx, token, patient.Email, 15*time.Minute)

	// Generate html body for verification mail
	htmlBody := utils.GenerateVerificationHTML(ctx, token, "verify_patient", "Thank you for Registering to ORDI!")
	s.email.SendEmail(patient.Email, "ORDI Email Verification", htmlBody, nil, "", "text/html")

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(patient.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	patient.Password = string(hashedPassword)

	// Store patient on database
	patient.Verified = false // Mark verified as false
	s.patientRepository.Save(ctx, &patient)

	// Prepare attachement to send to ORDI
	pdfBuffer, err := patientToPDF(patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Email details
	subject := "New Patient registration"
	bodyType := "text/plain"
	body := "A new patient has signed up. Please find the attached PDF with the details"
	to := "jhavedantamay@gmail.com"
	attachementName := fmt.Sprintf("%s.pdf", patient.FirstName)
	err = s.email.SendEmail(to, subject, body, pdfBuffer, attachementName, bodyType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Convert Patient to a PDF
func patientToPDF(patient models.PatientInfo) (*bytes.Buffer, error) {

	builder := utils.NewPDFBuilder()
	builder.AddTitle("Patient Information").
		AddField("Name", fmt.Sprintf("%s %s", patient.FirstName, patient.LastName)).
		AddField("Gender", patient.Gender).
		AddField("Father's Name", patient.FatherName).
		AddField("Father's Occupation", patient.FatherOccupation).
		AddField("ABHA ID", fmt.Sprintf("%d", patient.ABHAID)).
		AddField("Email", patient.Email).
		AddField("Address", fmt.Sprintf("%s, %s, %s, %s, %s",
			patient.StreetAddress, patient.City, patient.Region, patient.Country, patient.PostalCode)).
		AddSectionTitle("Doctor Information").
		AddField("Doctor Name", patient.DoctorName).
		AddField("Hospital Name", patient.HospitalName).
		AddField("Doctor Address", patient.DoctorAddress).
		AddField("Doctor Email", patient.DoctorEmail).
		AddField("Doctor Contact", patient.DoctorContact).
		AddField("Doctor Remarks", patient.DoctorRemarks).
		AddSectionTitle("Disease Information").
		AddField("Disease Name", patient.DiseaseName).
		AddField("Symptoms", patient.DiseaseSymptoms).
		AddSectionTitle("Sibling Information").
		AddField("Has Brother", fmt.Sprintf("%t", patient.HasBrother)).
		AddField("Has Sister", fmt.Sprintf("%t", patient.HasSister)).
		AddField("Sibling Has Rare Disease", fmt.Sprintf("%t", patient.SiblingHasRareDisease))

	return builder.Output()
}
