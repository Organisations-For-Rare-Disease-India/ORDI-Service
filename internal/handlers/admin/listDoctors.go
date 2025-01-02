package admin

import (
	"ORDI/cmd/web"
	"net/http"

	"github.com/a-h/templ"
)

func (a *adminHandler) ListDoctors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// claims, err := token.ValidateJWT(w, r)
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	// We can use admin scope to find which doctors they have access to
	// Currently, listing all the doctors
	doctors, err := a.doctorRepository.FindAll(ctx)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Pass this list of doctors to be displayed
	templ.Handler(web.AdminDoctorsListViewPage(doctors)).ServeHTTP(w, r)
}
