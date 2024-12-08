package doctor

import (
	"ORDI/cmd/web"
	"ORDI/internal/handlers/token"
	"net/http"

	"github.com/a-h/templ"
)

func (d *doctorHandler) Profile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, err := token.ValidateJWT(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	doctor, err := d.doctorRepository.FindByField(ctx, "email_id", claims.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	templ.Handler(web.DoctorProfilePage(doctor)).ServeHTTP(w, r)
}
