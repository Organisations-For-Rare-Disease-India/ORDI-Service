package doctor

import (
	"ORDI/cmd/web"
	"ORDI/internal/handlers/token"
	"ORDI/internal/utils"
	"net/http"

	"github.com/a-h/templ"
)

func (d *doctorHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
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

	templ.Handler(web.DoctorDashboardPage(doctor.FirstName, utils.DoctorProfile)).ServeHTTP(w, r)
}
