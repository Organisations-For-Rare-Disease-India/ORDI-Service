package models

import "gorm.io/gorm"

type PatientInfo struct {
	PersonalInfo `schema:",inline" gorm:"embedded"`
	DoctorInfo   `schema:",inline" gorm:"embedded"`
	SiblingInfo  `schema:",inline" gorm:"embedded"`
}

type PersonalInfo struct {
	gorm.Model              // Embed model for ID, CreatedAt, UpdatedAt, DeletedAt
	FirstName        string `schema:"first_name" gorm:"column:first_name"`
	LastName         string `schema:"last_name" gorm:"column:last_name"`
	Gender           string `schema:"gender" gorm:"column:gender"`
	FatherName       string `schema:"father_name" gorm:"column:father_name"`
	FatherOccupation string `schema:"father_occupation" gorm:"column:father_occupation"`
	ABHAID           int    `schema:"abha_id" gorm:"column:abha_id"`
	Email            string `schema:"email_id" gorm:"column:email_id"`
	Password         string `schema:"password" gorm:"column:password"`
	Country          string `schema:"country" gorm:"column:country"`
	StreetAddress    string `schema:"street_address" gorm:"column:street_address"`
	City             string `schema:"city" gorm:"column:city"`
	Region           string `schema:"region" gorm:"column:region"`
	PostalCode       string `schema:"postal_code" gorm:"column:postal_code"`
}

type DoctorInfo struct {
	DoctorName    string `schema:"doctor_name" gorm:"column:doctor_name"`
	HospitalName  string `schema:"hospital_name" gorm:"column:hospital_name"`
	DoctorAddress string `schema:"doctor_address" gorm:"column:doctor_address"`
	DoctorEmail   string `schema:"doctor_email_id" gorm:"column:doctor_email_id"`
	DoctorContact string `schema:"doctor_contact" gorm:"column:doctor_contact"`
	DoctorRemarks string `schema:"doctor_remarks" gorm:"column:doctor_remarks"`
	DiseaseInfo   `schema:",inline"`
}

type SiblingInfo struct {
	HasBrother            bool `schema:"has_brother" gorm:"column:has_brother"`
	HasSister             bool `schema:"has_sister" gorm:"column:has_sister"`
	SiblingHasRareDisease bool `schema:"sibling_has_rare_disease" gorm:"column:sibling_has_rare_disease"`
}

type DiseaseInfo struct {
	DiseaseName     string `schema:"disease_name" gorm:"column:disease_name"`
	DiseaseSymptoms string `schema:"disease_symptoms" gorm:"column:disease_symptoms"`
}
