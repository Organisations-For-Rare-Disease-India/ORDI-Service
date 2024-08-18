package server

import (
	"ORDI/cmd/web"
	"ORDI/internal/models"
	"context"

	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func (s *Server) PatientSignupFormHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var patient models.PatientInfo

	// Render the submit page with a success message
	// Render the template with the success flag
	templ.Handler(web.PatientSubmitPage()).ServeHTTP(w, r)

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

	// Push the data to the database
	err = s.AddPatientToDataBase(ctx, patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) AddPatientToDataBase(ctx context.Context, patient models.PatientInfo) error {
	// Add the patient to the database

	// SQL query to insert the patient into the database
	err := s.db.Insert(ctx, `
		INSERT INTO patient_info (
			first_name, last_name, gender, father_name, father_occupation,
			abha_id, email, country, street_address, city, region, postal_code,
			doctor_name, hospital_name, doctor_address, doctor_email, doctor_contact,
			doctor_remarks, has_brother, has_sister, sibling_has_rare_disease
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		patient.FirstName, patient.LastName, patient.Gender, patient.FatherName, patient.FatherOccupation,
		patient.ABHAID, patient.Email, patient.Country, patient.StreetAddress, patient.City, patient.Region, patient.PostalCode,
		patient.DoctorName, patient.HospitalName, patient.DoctorAddress, patient.DoctorEmail, patient.DoctorContact,
		patient.DoctorRemarks, patient.HasBrother, patient.HasSister, patient.SiblingHasRareDisease,
	)

	return err
}
