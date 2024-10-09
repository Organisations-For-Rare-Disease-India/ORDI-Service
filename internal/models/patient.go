package models

type PatientInfo struct {
	PersonalInfo `schema:",inline"`
	DoctorInfo   `schema:",inline"`
	SiblingInfo  `schema:",inline"`
}

type PersonalInfo struct {
	FirstName        string `schema:"first_name"`
	LastName         string `schema:"last_name"`
	Gender           string `schema:"gender"`
	FatherName       string `schema:"father_name"`
	FatherOccupation string `schema:"father_occupation"`
	ABHAID           int    `schema:"abha_id"`
	Email            string `schema:"email_id"`
	Password         string `schema:"password`
	Country          string `schema:"country"`
	StreetAddress    string `schema:"street_address"`
	City             string `schema:"city"`
	Region           string `schema:"region"`
	PostalCode       string `schema:"postal_code"`
}

type DoctorInfo struct {
	DoctorName    string `schema:"doctor_name"`
	HospitalName  string `schema:"hospital_name"`
	DoctorAddress string `schema:"doctor_address"`
	DoctorEmail   string `schema:"doctor_email_id"`
	DoctorContact string `schema:"doctor_contact"`
	DoctorRemarks string `schema:"doctor_remarks"`
	DiseaseInfo   `schema:",inline"`
}

type SiblingInfo struct {
	HasBrother            bool `schema:"has_brother"`
	HasSister             bool `schema:"has_sister"`
	SiblingHasRareDisease bool `schema:"sibling_has_rare_disease"`
}

type DiseaseInfo struct {
	DiseaseName     string `schema:"disease_name"`
	DiseaseSymptoms string `schema:"disease_symptoms"`
}
