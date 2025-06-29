package admin

import (
	"ORDI/cmd/web"
	"ORDI/internal/utils"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func (a *adminHandler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Errorf("method:%s not allowed", r.Method).Error(), http.StatusBadRequest)
		return
	}
	templ.Handler(web.AdminCreateAppointment(utils.AdminAppointmentCreate)).ServeHTTP(w, r)
}
