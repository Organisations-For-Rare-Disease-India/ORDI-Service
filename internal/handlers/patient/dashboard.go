package patient

import (
	"ORDI/cmd/web"
	"ORDI/internal/handlers/token"
	"ORDI/internal/utils"
	"net/http"

	"github.com/a-h/templ"
)

func (p *patientHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
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

	// TODO: Add logic here to show current status of patient
	templ.Handler(web.PatientDashboardPage(patient.FirstName,
		utils.PatientProfile)).ServeHTTP(w, r)
}
