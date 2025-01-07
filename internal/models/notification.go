package models

import "gorm.io/gorm"

type Notification struct {
	gorm.Model        // Embed model for ID, CreatedAt, UpdatedAt, DeletedAt
	UserEmail  string `schema:"user_email" gorm:"column:user_email"`
	Message    string `schema:"message" gorm:"column:message"`
	SentTime   string `schema:"sent_time" gorm:"column:sent_time"`
	IsRead     bool   `schema:"is_read" gorm:"column:is_read"`
}
