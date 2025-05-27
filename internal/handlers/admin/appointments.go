package admin

import (
	"ORDI/cmd/web"
	"ORDI/internal/models"
	"ORDI/internal/utils"
	"fmt"
	"net/http"
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
			AppointmentID:   int(v.ID),
			PatientName:     fmt.Sprintf("%d", v.PatientID),
			DoctorName:      fmt.Sprintf("%d", v.DoctorID),
			AppointmentDate: v.ApppointmentDate.Format(time.RFC3339),
		}

	}
	templ.Handler(web.AdminAppointmentsPage(utils.AdminAppointments, appointmentsData)).ServeHTTP(w, r)
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
