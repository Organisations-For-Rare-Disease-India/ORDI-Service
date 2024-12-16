package doctor

import (
	"ORDI/cmd/web"
	"net/http"

	"github.com/a-h/templ"
)

func (s *doctorHandler) Appointment(w http.ResponseWriter, r *http.Request) {
	// TODO: Admin should be able to add appointments to Patient and Doctor's calendar
	templ.Handler(web.CalendarPage()).ServeHTTP(w, r)
}
