package patient

import (
	"ORDI/cmd/web"
	"ORDI/internal/utils"
	"net/http"
	"strings"

	"github.com/a-h/templ"
)

func (s *patientHandler) Appointment(w http.ResponseWriter, r *http.Request) {
	// TODO: Admin should be able to add appointments to Patient and Doctor's calendar
	if isAdminHost(r) {
		// add to patient calender
		return templ.Handler(web.CalendarPage()).ServeHTTP(w, r)
	}
	// display only calender
	return templ.Handler(web.CalendarPage()).ServeHTTP(w, r)
}

func isAdminHost(r *http.Request) bool {
	hostName := r.URL.Hostname()
	return strings.EqualFold(hostName, utils.InternalDomain)

}
