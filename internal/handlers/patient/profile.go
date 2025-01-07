package patient

import (
	"ORDI/cmd/web"
	"ORDI/internal/handlers/token"
	"ORDI/internal/utils"
	"net/http"

	"github.com/a-h/templ"
)

func (p *patientHandler) Profile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, err := token.ValidateJWT(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	patient, err := p.patientRepository.FindByField(ctx, "email_id", claims.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	notificationCount, err := utils.GetNotificationCount(ctx, p.notificationRepository, claims.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	templ.Handler(web.PatientProfilePage(patient, notificationCount)).ServeHTTP(w, r)

}
