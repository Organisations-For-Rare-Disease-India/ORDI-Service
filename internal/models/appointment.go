package models

import (
	"time"

	"gorm.io/gorm"
)

type Appointment struct {
	gorm.Model
	PatientID           uint              `schema:"patiend_id" gorm:"column:patient_id;type bigint unsigned"`
	DoctorID            uint              `schema:"doctor_id" gorm:"column:doctor_id;type bigint unsigned"`
	Notes               string            `schema:"notes" gorm:"column:notes"`
	Remarks             string            `schema:"remarks" gorm:"column:remarks"`
	PreAppointmentNotes string            `schema:"pre_appointment_notes" gorm:"column:pre_appointment_notes"`
	RecommendedTests    []RecommendedTest `gorm:"foreignKey:AppointmentID"`
	ApppointmentDate    time.Time         `schema:"appointment_date" gorm:"column:appointment_date"`
}

type RecommendedTest struct {
	gorm.Model
	AppointmentID uint   `schema:"appointment_id" gorm:"type:bigint unsigned"`
	Name          string `schema:"name" gorm:"column:name"`
	Description   string `schema:"description" gorm:"column:description"`
	Status        string `schema:"status" grom:"column:status"`
}

type AppointmentData struct {
	PatientName string
	DoctorName  string
	Appointment
}
