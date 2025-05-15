package patient

import (
	"ORDI/cmd/web"
	"net/http"

	"github.com/a-h/templ"
)

func (s *patientHandler) Appointment(w http.ResponseWriter, r *http.Request) {
	// display only calender
	templ.Handler(web.CalendarPage()).ServeHTTP(w, r)
}
