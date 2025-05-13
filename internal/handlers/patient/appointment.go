package patient

import (
	"ORDI/cmd/web"
	"ORDI/internal/models"
	"ORDI/internal/utils"
	"net/http"
	"strings"

	"github.com/a-h/templ"
)

func (s *patientHandler) Appointment(w http.ResponseWriter, r *http.Request) {
	if isAdminHost(r) {
		// add to patient calender
		if err := s.createAppointment(r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		templ.Handler(web.CalendarPage()).ServeHTTP(w, r)
	}
	// display only calender
	templ.Handler(web.CalendarPage()).ServeHTTP(w, r)
}

func isAdminHost(r *http.Request) bool {
	hostName := r.URL.Hostname()
	return strings.EqualFold(hostName, utils.InternalDomain)

}

func (s *patientHandler) createAppointment(r *http.Request) error {
	a := &models.Appointment{}
	return s.appointmentRepository.Save(r.Context(), a)
}
