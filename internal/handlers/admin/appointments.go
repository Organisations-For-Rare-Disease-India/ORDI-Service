package admin

import (
	"ORDI/cmd/web"
	"ORDI/internal/utils"
	"net/http"
	"time"

	"github.com/a-h/templ"
)

func (a *adminHandler) Appointments(w http.ResponseWriter, r *http.Request) {
	templ.Handler(web.AdminAppointmentsPage(utils.AdminAppointments)).ServeHTTP(w, r)
}

func getMonthYear(now time.Time) (int, time.Month, int) {
	year, month, day := now.Date()
	return year, month, day

}

func noOfDays(now time.Time) int {
	year, month, _ := getMonthYear(now)
	future := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC)
	noOfDays := future.Day()
	return noOfDays
}
