package doctor

import (
	"ORDI/cmd/web"
	"ORDI/internal/handlers/calc"
	"ORDI/internal/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
)

// GetMonthlyAppointment fetches all appointments for logged in patient
// by the current month
func (s *doctorHandler) GetMonthlyAppointment(
	w http.ResponseWriter, r *http.Request) {
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
	startDate, err := calc.FirstDay(time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	lastDate, err := calc.LastDay(time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	appointmentsFromStore, err := s.appointmentRepository.FilterBetweenDates(
		r.Context(),
		"doctor_id", doctorFromStore.ID,
		"appointment_date", startDate, lastDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// get patient data
	appointmentsData, errs := calc.AppointmentIterate(r.Context(),
		appointmentsFromStore,
		s.doctorRepository, s.patientRepository)

	if estr := calc.ErrList(errs); estr != nil {
		http.Error(w, estr.String(), http.StatusInternalServerError)
		return
	}
	templ.Handler(web.GetDoctorAppointments(
		utils.DoctorAppointments,
		appointmentsData)).ServeHTTP(w, r)
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
	ymd, err := calc.GetYearMonthDate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if calc.IsInvalidYMD(ymd) {
		// 3. get monthly appointments
		http.Error(w, fmt.Errorf(
			"invalid url param recieved for year month and date").
			Error(), http.StatusBadRequest)
		return
	}
	appointmentDate, err := calc.DateTimeFormat(ymd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// get appointments for selected date
	appointmentsFromStore, err := s.appointmentRepository.
		FilterByDate(r.Context(), "doctor_id",
			doctorFromStore.ID, "appointment_date", appointmentDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	appointmentsData, errs := calc.AppointmentIterate(r.Context(),
		appointmentsFromStore,
		s.doctorRepository, s.patientRepository)
	if estr := calc.ErrList(errs); estr != nil {
		http.Error(w, estr.String(), http.StatusInternalServerError)
		return
	}

	templ.Handler(web.GetDoctorAppointmentByDate(
		utils.DoctorAppointmentByID,
		appointmentsData)).ServeHTTP(w, r)
}
