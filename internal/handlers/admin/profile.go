package admin

import (
	"ORDI/cmd/web"
	"ORDI/internal/handlers/token"
	"net/http"

	"github.com/a-h/templ"
)

func (a *adminHandler) Profile(w http.ResponseWriter, r *http.Request) {
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
	templ.Handler(web.AdminProfilePage(admin)).ServeHTTP(w, r)
}
