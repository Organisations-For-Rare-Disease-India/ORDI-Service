package models

import "gorm.io/gorm"

type Doctor struct {
	DoctorPersonalInfo `schema:",inline" gorm:"embedded"`
	DoctorHospitalInfo `schema:",inline" gorm:"embedded"`
	Verified           bool `gorm:"column:verified"` // Add this line
}

type DoctorPersonalInfo struct {
	gorm.Model           // Embed model for ID, CreatedAt, UpdatedAt, DeletedAt
	FirstName     string `schema:"first_name" gorm:"column:first_name"`
	LastName      string `schema:"last_name" gorm:"column:last_name"`
	Gender        string `schema:"gender" gorm:"column:gender"`
	ABHAID        int    `schema:"abha_id" gorm:"column:abha_id"`
	Email         string `schema:"email_id" gorm:"column:email_id"`
	Password      string `schema:"password" gorm:"column:password"`
	Country       string `schema:"country" gorm:"column:country"`
	StreetAddress string `schema:"street_address" gorm:"column:street_address"`
	City          string `schema:"city" gorm:"column:city"`
	Region        string `schema:"region" gorm:"column:region"`
	PostalCode    string `schema:"postal_code" gorm:"column:postal_code"`
}

type DoctorHospitalInfo struct {
	HospitalCountry       string `schema:"hospital_country" gorm:"column:hospital_country"`
	HospitalStreetAddress string `schema:"hospital_street_address" gorm:"column:hospital_street_address"`
	HospitalCity          string `schema:"hospital_city" gorm:"column:hospital_city"`
	HospitalRegion        string `schema:"hospital_region" gorm:"column:hospital_region"`
	HospitalPostalCode    string `schema:"hospital_postal_code" gorm:"column:hospital_postal_code"`
}
