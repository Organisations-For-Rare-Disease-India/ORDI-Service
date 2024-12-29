package utils

import (
	"ORDI/internal/models"
	"bytes"
	"fmt"
)

// Convert Patient to a PDF
func PatientToPDF(patient models.Patient) (*bytes.Buffer, error) {

	builder := NewPDFBuilder()
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

func DoctorToPDF(doctor models.Doctor) (*bytes.Buffer, error) {

	builder := NewPDFBuilder()
	builder.AddTitle("Patient Information").
		AddField("Name", fmt.Sprintf("%s %s", doctor.FirstName, doctor.LastName)).
		AddField("Gender", doctor.Gender).
		AddField("ABHA ID", fmt.Sprintf("%d", doctor.ABHAID)).
		AddField("Email", doctor.Email).
		AddField("Address", fmt.Sprintf("%s, %s, %s, %s, %s",
			doctor.StreetAddress, doctor.City, doctor.Region, doctor.Country, doctor.PostalCode))

	return builder.Output()
}

func AdminToPDF(admin models.Admin) (*bytes.Buffer, error) {

	builder := NewPDFBuilder()
	builder.AddTitle("Your Information").
		AddField("Name", fmt.Sprintf("%s %s", admin.FirstName, admin.LastName)).
		AddField("Email", admin.Email).
		AddField("Address", fmt.Sprintf("%s, %s, %s, %s, %s",
			admin.StreetAddress, admin.City, admin.Region, admin.Country, admin.PostalCode))

	return builder.Output()

}
