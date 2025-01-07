package admin

import (
	"ORDI/cmd/web"
	"ORDI/internal/handlers/token"
	"ORDI/internal/utils"
	"net/http"

	"github.com/a-h/templ"
)

func (a *adminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, err := token.ValidateJWT(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	admin, err := a.adminRepository.FindByField(ctx, "email_id", claims.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	notificationCount, err := utils.GetNotificationCount(ctx, a.notificationRepository, claims.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// We can later use more admin information to define scopes
	templ.Handler(web.AdminDashboardPage(admin.FirstName, utils.AdminProfile, utils.AdminViewDoctorList, utils.AdminViewPatientList, notificationCount)).ServeHTTP(w, r)
}
