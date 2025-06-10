package patient

import (
	"ORDI/cmd/web"
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

// TODO: enable to get patient email from cookie set during login time
func (s *patientHandler) GetAppointments(w http.ResponseWriter, r *http.Request) {
	// claims, err := token.ValidateJWT(w, r)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusUnauthorized)
	// 	return
	// }
	// TODO: change to get from cookie
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
