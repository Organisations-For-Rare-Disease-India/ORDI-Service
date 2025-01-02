package models

import "gorm.io/gorm"

type MasterAdmin struct {
	gorm.Model        // Embed model for ID, CreatedAt, UpdatedAt, DeletedAt
	Email      string `schema:"email_id" gorm:"column:email_id"`
	Password   string `schema:"password" gorm:"column:password"`
}
