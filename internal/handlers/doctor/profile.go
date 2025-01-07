package doctor

import (
	"ORDI/cmd/web"
	"ORDI/internal/handlers/token"
	"ORDI/internal/utils"
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
	notificationCount, err := utils.GetNotificationCount(ctx, d.notificationRepository, claims.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	templ.Handler(web.DoctorProfilePage(doctor, notificationCount)).ServeHTTP(w, r)
}
