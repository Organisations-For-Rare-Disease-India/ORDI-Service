package admin

import (
	"ORDI/cmd/web"
	"net/http"

	"github.com/a-h/templ"
)

func (a *adminHandler) Appointments(w http.ResponseWriter, r *http.Request) {
	// TODO: make this edit using username/patientname
	templ.Handler(web.CalendarPage()).ServeHTTP(w, r)
}
