package server

import (
	"ORDI/cmd/web"
	"net/http"

	"github.com/a-h/templ"
)

func (s *Server) PatientLoginHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Verify email and password
	// TODO: Grab name and details from database
	templ.Handler(web.PatientDashboardPage()).ServeHTTP(w, r)
}
