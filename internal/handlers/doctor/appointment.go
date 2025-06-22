package doctor

import (
	"ORDI/cmd/web"
	"ORDI/internal/handlers/token"
	"ORDI/internal/models"
	"ORDI/internal/utils"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
)

// GetMonthlyAppointment fetches all appointments for logged in patient
// by the current month
func (s *doctorHandler) GetMonthlyAppointment(w http.ResponseWriter, r *http.Request) {
	// precheck
	if y := r.PathValue("year"); y != "" {
		http.Error(w,
			fmt.Errorf("invalid request year month date path param received").
				Error(), http.StatusBadRequest)
		return
	}
	// claims, err := getLoggedinUser(r, w)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusUnauthorized)
	// 	return
	// }
	// TODO: change to get from cookie
	// email := claims.Email
	email := "doctor.email@domain.com"
	doctorFromStore, err := s.doctorRepository.FindByField(r.Context(),
		"email_id", email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if doctorFromStore == nil {
		http.Error(w,
			fmt.Errorf("doctor with email:%q not found", email).Error(), http.StatusBadRequest)
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
		"patient_id", doctorFromStore.ID,
		"appointment_date", startDate, lastDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// get patient data
	errors := make([]error, len(appointments))
	appointmentsData := make([]models.AppointmentData, len(appointments))
	for i, appointment := range appointments {
		patientFromStore, err := s.patientRepository.FindByID(r.Context(), appointment.PatientID)
		if err != nil {
			errors[i] = err
			continue
		}
		appointmentsData[i] = models.AppointmentData{
			AppointmentID:   fmt.Sprintf("%d", appointment.ID),
			PatientName:     fmt.Sprintf("%s,%s", patientFromStore.FirstName, patientFromStore.LastName),
			AppointmentDate: appointment.ApppointmentDate.Format(time.DateTime),
		}

	}
	var errString strings.Builder
	for v := range slices.Values(errors) {
		if v != nil {
			fmt.Fprintf(&errString, "%s,", v.Error())
		}

	}
	if errString.Len() > 0 {
		http.Error(w, errString.String(), http.StatusInternalServerError)
		return
	}

	templ.Handler(web.GetDoctorAppointments(
		utils.DoctorAppointments,
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
func (s *doctorHandler) GetAppointmentByDate(
	w http.ResponseWriter, r *http.Request) {
	// 1. get logged in user from the token
	// claims, err := getLoggedinUser(r, w)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusUnauthorized)
	// 	return
	// }
	// TODO: change to get from cookie
	// email := claims.Email
	email := "doctor.email@domain.com"
	doctorFromStore, err := s.doctorRepository.FindByField(r.Context(),
		"email_id", email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if doctorFromStore == nil {
		http.Error(w,
			fmt.Errorf("doctor with email:%q not found", email).Error(), http.StatusBadRequest)
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
		http.Error(w, fmt.Errorf(
			"invalid url param recieved for year month and date").
			Error(), http.StatusBadRequest)
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
			doctorFromStore.ID, "appointment_date", appointmentDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	errors := make([]error, len(appointmentsFromStore))
	appointmentsData := make([]models.AppointmentData, len(appointmentsFromStore))
	for i, appointment := range appointmentsFromStore {
		patientFromStore, err := s.patientRepository.FindByID(r.Context(), appointment.PatientID)
		if err != nil {
			errors[i] = err
			continue
		}
		appointmentsData[i] = models.AppointmentData{
			AppointmentID: fmt.Sprintf("%d", appointment.ID),
			PatientName: fmt.Sprintf("%s,%s", patientFromStore.FirstName,
				patientFromStore.LastName),
			AppointmentDate: appointment.ApppointmentDate.Format(time.DateTime),
		}
	}

	var errString strings.Builder
	for v := range slices.Values(errors) {
		if v != nil {
			fmt.Fprintf(&errString, "%s,", v.Error())
		}

	}
	if errString.Len() > 0 {
		http.Error(w, errString.String(), http.StatusInternalServerError)
		return
	}

	templ.Handler(web.GetDoctorAppointmentByDate(
		utils.DoctorAppointmentByID,
		appointmentsData)).ServeHTTP(w, r)
}

type DoctorAppointmentData struct {
	DoctorName      string `json:"doctor_name"`
	AppointmentDate string `json:"appointment_date"`
}

type yearMonthDate struct {
	year, month, date int
}

var monthLookUP = map[string]int{
	"January": 1, "Feburary": 2, "March": 3, "April": 4, "May": 5, "June": 6,
	"July": 7, "August": 8, "September": 9, "October": 10, "November": 11,
	"December": 12,
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
	if m != "" {
		ymd.month = monthLookUP[m]
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
