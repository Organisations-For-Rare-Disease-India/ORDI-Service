package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model           // Embed model for ID, CreatedAt, UpdatedAt, DeletedAt
	UserID     uint      `schema:"user_id" gorm:"column:user_id"`
	UserType   string       `schema:"user_type" gorm:"column:user_type"`
	Title      string    `schema:"title" gorm:"column:title;size:100"`
	UserEmail  string    `schema:"user_email" gorm:"column:user_email"`
	Message    string    `schema:"message" gorm:"column:message"`
	SentTime   time.Time `schema:"sent_time" gorm:"column:sent_time"`
	IsRead     bool      `schema:"is_read" gorm:"column:is_read"`
}

type ViewNotification struct {
	SentTime string
	Message  string
}
