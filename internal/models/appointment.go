package models

import "gorm.io/gorm"

type Appointment struct {
	gorm.Model
	PatientID           int64  `schema:"patiend_id" gorm:"column:patient_id"`
	DoctorID            int64  `schema:"doctor_id" gorm:"column:doctor_id"`
	Notes               string `schema:"notes" gorm:"column:notes"`
	Remarks             string `schema:"remarks" gorm:"column:remarks"`
	PreAppointmentNotes string `schema:"pre_appointment_notes" gorm:"column:pre_appointment_notes"`
}

type TestRecommended struct {
	gorm.Model
	AppointmentID int64  `schema:"appointment_id" gorm:"column:appointment_id"`
	Name          string `schema:"name" gorm:"column:name"`
	Description   string `schema:"description" gorm:"column:description"`
	Status        string `schema:"status" grom:"column:status"`
}
