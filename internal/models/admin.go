package models

import "gorm.io/gorm"

type Admin struct {
	AdminInfo `schema:",inline" gorm:"embedded"`
	Verified  bool `gorm:"column:verified"`
}

type AdminInfo struct {
	gorm.Model           // Embed model for ID, CreatedAt, UpdatedAt, DeletedAt
	FirstName     string `schema:"first_name" gorm:"column:first_name"`
	LastName      string `schema:"last_name" gorm:"column:last_name"`
	Email         string `schema:"email_id" gorm:"column:email_id"`
	Password      string `schema:"password" gorm:"column:password"`
	Gender        string `schema:"gender" gorm:"column:gender"`
	Country       string `schema:"country" gorm:"column:country"`
	StreetAddress string `schema:"street_address" gorm:"column:street_address"`
	City          string `schema:"city" gorm:"column:city"`
	Region        string `schema:"region" gorm:"column:region"`
	PostalCode    string `schema:"postal_code" gorm:"column:postal_code"`
	Notes         string `schema:"notes" gorm:"column: notes"`
}
