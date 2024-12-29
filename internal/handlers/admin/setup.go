package admin

import (
	"ORDI/cmd/web"
	"ORDI/internal/utils"
	"net/http"

	"github.com/a-h/templ"
)

func (a *adminHandler) Setup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	email := r.URL.Query().Get(EMAIL_HEADER)
	if email == "" {
		http.Error(w, EmailRequiredError, http.StatusUnauthorized)
	}
	admin, err := a.adminRepository.FindByField(ctx, "email_id", email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if admin == nil {
		http.Error(w, "Admin has not been added to the team", http.StatusUnauthorized)
		return
	}
	// Admin exists

	if admin.Verified {
		// Admin has already setup the credentials. Take them to login
		templ.Handler(web.AdminLoginPage(utils.AdminLoginSubmit, false)).ServeHTTP(w, r)
		return
	}

	// Otherwise, take them to credentials page
	templ.Handler(web.AdminRegisterPage(utils.AdminRegisterSubmit)).ServeHTTP(w, r)
}
