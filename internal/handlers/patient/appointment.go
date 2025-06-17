package patient

import (
	"ORDI/cmd/web"
	"ORDI/internal/handlers/token"
	"ORDI/internal/models"
	"ORDI/internal/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
)

func (s *patientHandler) Appointment(w http.ResponseWriter, r *http.Request) {
	// display only calender
	templ.Handler(web.CalendarPage()).ServeHTTP(w, r)
}

func (s *patientHandler) GetMonthlyAppointment(w http.ResponseWriter, r *http.Request) {
	// claims, err := getLoggedinUser(r, w)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusUnauthorized)
	// 	return
	// }
	// TODO: change to get from cookie
	// email := claims.Email
	email := "email.id@domain.com"
	patientFromStore, err := s.patientRepository.FindByField(r.Context(),
		"email_id", email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	startDate, err := firstDay(time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	lastDate, err := lastDay(time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	appointments, err := s.appointmentRepository.FilterByDate(
		r.Context(),
		"patiend_id", patientFromStore.ID,
		"appointment_date", startDate, lastDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appointmentsData := make([]models.AppointmentData, len(appointments))
	for i, appointment := range appointments {
		appointmentsData[i] = models.AppointmentData{
			AppointmentID:   fmt.Sprintf("%d", appointment.ID),
			DoctorName:      patientFromStore.DoctorName,
			AppointmentDate: appointment.ApppointmentDate.Format(time.DateTime),
		}

	}
	templ.Handler(web.GetPatientAppointments(
		utils.GetPatientAppointments,
		appointmentsData)).ServeHTTP(w, r)
}

func getLoggedinUser(r *http.Request,
	w http.ResponseWriter) (token.Claims, error) {
	claims, err := token.ValidateJWT(w, r)
	if err != nil {
		return token.Claims{}, err
	}
	return *claims, nil
}

func location() (*time.Location, error) {
	return time.LoadLocation(`Asia/Calcutta`)
}

func startTime(t time.Time) (time.Time, error) {
	location, err := location()
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(
		t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, location), nil

}

func getMonthYear(now time.Time) (int, time.Month, int) {
	year, month, day := now.Date()
	return year, month, day

}

func noOfDays(now time.Time) int {
	year, month, _ := now.Date()
	future := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC)
	noOfDays := future.Day()
	return noOfDays
}
func endTime(t time.Time) (time.Time, error) {
	totalDays := noOfDays(t)
	location, err := location()
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(
		t.Year(), t.Month(), totalDays, 23, 59, 59, 0, location), nil
}

func firstDay(t time.Time) (time.Time, error) {
	location, err := location()
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(
		t.Year(), t.Month(), 01, 0, 0, 0, 0, location), nil
}

func lastDay(t time.Time) (time.Time, error) {
	totalDays := noOfDays(t)
	location, err := location()
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(
		t.Year(), t.Month(), totalDays, 0, 0, 0, 0, location), nil
}

// TODO: enable to get patient email from cookie set during login time
func (s *patientHandler) GetAppointments(w http.ResponseWriter, r *http.Request) {
	// claims, err := getLoggedinUser(r, w)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusUnauthorized)
	// 	return
	// }
	// TODO: change to get from cookie
	// email := claims.Email
	email := "email.id@domain.com"
	patientFromStore, err := s.patientRepository.FindByField(r.Context(),
		"email_id", email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	appointmentsFromStore, err := s.appointmentRepository.FindAllByField(r.Context(), "patient_id", patientFromStore.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	appointmentsData := make([]models.AppointmentData, len(appointmentsFromStore))
	for i, appointment := range appointmentsFromStore {
		appointmentsData[i] = models.AppointmentData{
			AppointmentID:   fmt.Sprintf("%d", appointment.ID),
			DoctorName:      patientFromStore.DoctorName,
			AppointmentDate: appointment.ApppointmentDate.Format(time.DateTime),
		}

	}
	for _, v := range appointmentsData {
		fmt.Printf("appointmentsData:%#v\n", v)

	}
	templ.Handler(web.GetPatientAppointments(
		utils.GetPatientAppointments,
		appointmentsData)).ServeHTTP(w, r)
}

type PatientAppointmentsData struct {
	DoctorName      string `json:"doctor_name"`
	AppointmentDate string `json:"appointment_date"`
}
