package models

import "errors"

type Role int

const (
	DoctorType = iota
	PatientType
	AdminType
)

func (r Role) String() string {
	switch r {
	case DoctorType:
		return "doctor"
	case PatientType:
		return "patient"
	case AdminType:
		return "admin"
	default:
		return "unknown"
	}
}

func ParseRole(s string) (Role, error) {
	switch s {
	case "Doctor":
		return DoctorType, nil
	case "Patient":
		return PatientType, nil
	case "Admin":
		return AdminType, nil
	default:
		return 0, errors.New("invalid role")
	}
}
