package admin

import (
	"ORDI/cmd/web"
	"net/http"

	"github.com/a-h/templ"
)

func (a *adminHandler) ListPatients(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// claims, err := token.ValidateJWT(w, r)
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	// We can use admin scope to find which patients they have access to
	// Currently, listing all the patients
	patients, err := a.patientRepository.FindAll(ctx)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Pass this list of doctors to be displayed
	templ.Handler(web.AdminPatientsListViewPage(patients)).ServeHTTP(w, r)
}
