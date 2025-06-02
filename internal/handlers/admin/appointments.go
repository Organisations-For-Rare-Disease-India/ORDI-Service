package admin

import (
	"ORDI/cmd/web"
	"ORDI/internal/models"
	"ORDI/internal/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
)

func (a *adminHandler) Appointments(w http.ResponseWriter, r *http.Request) {
	appointments, err := a.appointmentRepository.FindAllWithPage(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	appointmentsData := make([]models.AppointmentData, len(appointments))
	for i, v := range appointments {
		appointmentsData[i] = models.AppointmentData{
			AppointmentID: int(v.ID),
			// TODO: get name from patient and doctor table
			PatientName:     fmt.Sprintf("%d", v.PatientID),
			DoctorName:      fmt.Sprintf("%d", v.DoctorID),
			AppointmentDate: v.ApppointmentDate.Format(time.RFC3339),
		}

	}
	templ.Handler(web.AdminAppointmentsPage(utils.AdminAppointments, appointmentsData)).ServeHTTP(w, r)
}

func (a *adminHandler) GetAppointmentID() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, fmt.Errorf("empty value received for id").Error(),
				http.StatusBadRequest)
			return
		}
		appointmentID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		appointment, err := a.appointmentRepository.FindByID(r.Context(), uint(appointmentID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		a := &models.AppointmentData{
			AppointmentID:   int(appointment.ID),
			PatientName:     fmt.Sprintf("%d", appointment.PatientID),
			DoctorName:      fmt.Sprintf("%d", appointment.DoctorID),
			AppointmentDate: appointment.ApppointmentDate.Format(time.RFC3339),
		}

		templ.Handler(web.GetAppointmentByID(a)).ServeHTTP(w, r)

	})
}

func getMonthYear(now time.Time) (int, time.Month, int) {
	year, month, day := now.Date()
	return year, month, day

}

func noOfDays(now time.Time) int {
	year, month, _ := getMonthYear(now)
	future := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC)
	noOfDays := future.Day()
	return noOfDays
}
