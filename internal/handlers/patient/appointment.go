package patient

import (
	"ORDI/cmd/web"
	"ORDI/internal/handlers/token"
	"ORDI/internal/models"
	"ORDI/internal/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
)

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
	appointments, err := s.appointmentRepository.FilterBetweenDates(
		r.Context(),
		"patient_id", patientFromStore.ID,
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
		utils.PatientAppointments,
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

const LOCATION = `Asia/Calcutta`

func location() (*time.Location, error) {
	return time.LoadLocation(LOCATION)
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
// NOTES: when patient visit appointments page, it shows list of appointment
// for current month along with calendar on right to select a date .
// When a date is selected from calendar it displays only the appointments
// for the selected date. It fetches the patient details from the token initialized
// during logged in
func (s *patientHandler) GetAppointments(w http.ResponseWriter, r *http.Request) {
	// 1. get logged in user from the token
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
	// 2. get url path parameters from request url
	ymd, err := getYearMonthDate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ymd.year == 0 || ymd.month == 0 || ymd.date == 0 {
		// 3. get monthly appointments
		s.GetMonthlyAppointment(w, r)
		return
	}
	appointmentDate, err := dateTimeFormat(ymd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// get appointments for selected date
	appointmentsFromStore, err := s.appointmentRepository.
		FilterByDate(r.Context(), "patient_id",
			patientFromStore.ID, "appointment_date", appointmentDate)
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
	templ.Handler(web.GetPatientAppointments(
		utils.PatientAppointments,
		appointmentsData)).ServeHTTP(w, r)
}

type PatientAppointmentsData struct {
	DoctorName      string `json:"doctor_name"`
	AppointmentDate string `json:"appointment_date"`
}

type yearMonthDate struct {
	year, month, date int
}

func getYearMonthDate(r *http.Request) (yearMonthDate, error) {
	ymd := yearMonthDate{}
	var err error
	y := r.PathValue("year")
	ymd.year, err = strToInt(y)
	if err != nil {
		return ymd, err
	}
	m := r.PathValue("month")
	ymd.month, err = strToInt(m)
	if err != nil {
		return ymd, err
	}
	d := r.PathValue("date")
	ymd.date, err = strToInt(d)
	if err != nil {
		return ymd, err
	}
	return ymd, err
}

func strToInt(in string) (int, error) {
	if in != "" {
		return strconv.Atoi(in)
	}
	return 0, nil
}

func dateTimeFormat(ymd yearMonthDate) (time.Time, error) {
	location, err := location()
	if err != nil {
		return time.Time{}, err
	}
	appointmentDate := time.Date(ymd.year,
		time.Month(ymd.month), ymd.date, 0, 0, 0, 0, location)
	return appointmentDate, nil
}
