package models

type PatientInfo struct {
	PersonalInfo `form:",inline"`
	DoctorInfo   `form:",inline"`
	SiblingInfo  `form:",inline"`
}

type PersonalInfo struct {
	FirstName        string `form:"first_name"`
	LastName         string `form:"last_name"`
	Gender           string `form:"gender"`
	FatherName       string `form:"father_name"`
	FatherOccupation string `form:"father_occupation"`
	ABHAID           int    `form:"abha_id"`
	Email            string `form:"email_id"`
	Country          string `form:"country"`
	StreetAddress    string `form:"street_address"`
	City             string `form:"city"`
	Region           string `form:"region"`
	PostalCode       string `form:"postal_code"`
}

type DoctorInfo struct {
	DoctorName    string `form:"doctor_name"`
	HospitalName  string `form:"hospital_name"`
	DoctorAddress string `form:"doctor_address"`
	DoctorEmail   string `form:"doctor_email"`
	DoctorContact string `form:"doctor_contact"`
}

type SiblingInfo struct {
	HasBrother            bool `form:"has_brother"`
	HasSister             bool `form:"has_sister"`
	SiblingHasRareDisease bool `form:"sibling_has_rare_disease"`
}
