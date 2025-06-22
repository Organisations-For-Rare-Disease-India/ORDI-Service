package calc

import (
	"ORDI/internal/handlers/token"
	"ORDI/internal/models"
	"ORDI/internal/repositories"
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

type yearMonthDate struct {
	year, month, date int
}

func FirstDay(t time.Time) (time.Time, error) {
	location, err := location()
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(
		t.Year(), t.Month(), 01, 0, 0, 0, 0, location), nil
}

func noOfDays(now time.Time) int {
	year, month, _ := now.Date()
	future := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC)
	noOfDays := future.Day()
	return noOfDays
}

func endTime(t time.Time) (time.Time, error) {
	totalDays := noOfDays(t)
	location, err := location()
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(
		t.Year(), t.Month(), totalDays, 23, 59, 59, 0, location), nil
}

func LastDay(t time.Time) (time.Time, error) {
	totalDays := noOfDays(t)
	location, err := location()
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(
		t.Year(), t.Month(), totalDays, 0, 0, 0, 0, location), nil
}

var monthLookUP = map[string]int{
	"January": 1, "Feburary": 2, "March": 3, "April": 4, "May": 5, "June": 6,
	"July": 7, "August": 8, "September": 9, "October": 10, "November": 11,
	"December": 12,
}

func IsInvalidYMD(ymd yearMonthDate) bool {
	return (ymd.year == 0 || ymd.month == 0 || ymd.date == 0)

}

func GetYearMonthDate(r *http.Request) (yearMonthDate, error) {
	ymd := yearMonthDate{}
	var err error
	y := r.PathValue("year")
	ymd.year, err = strToInt(y)
	if err != nil {
		return ymd, err
	}
	m := r.PathValue("month")
	if m != "" {
		ymd.month = monthLookUP[m]
	}
	d := r.PathValue("date")
	ymd.date, err = strToInt(d)
	if err != nil {
		return ymd, err
	}
	return ymd, err
}

func strToInt(in string) (int, error) {
	if in != "" {
		return strconv.Atoi(in)
	}
	return 0, nil
}

func DateTimeFormat(ymd yearMonthDate) (time.Time, error) {
	location, err := location()
	if err != nil {
		return time.Time{}, err
	}
	appointmentDate := time.Date(ymd.year,
		time.Month(ymd.month), ymd.date, 0, 0, 0, 0, location)
	return appointmentDate, nil
}

const LOCATION = `Asia/Calcutta`

func location() (*time.Location, error) {
	return time.LoadLocation(LOCATION)
}

func startTime(t time.Time) (time.Time, error) {
	location, err := location()
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(
		t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, location), nil

}

func AppointmentIterate(
	ctx context.Context,
	data []models.Appointment,
	dr repositories.Repository[models.Doctor],
	pr repositories.Repository[models.Patient]) ([]models.AppointmentData, []error) {
	if len(data) < 1 {
		return []models.AppointmentData{}, nil
	}
	ad := make([]models.AppointmentData, len(data))
	errList := make([]error, len(data))
	for i, v := range slices.All(data) {
		pdata, err := pr.FindByID(ctx, v.PatientID)
		if err != nil {
			errList[i] = err
		}
		ad[i] = models.AppointmentData{
			AppointmentID:   fmt.Sprintf("%d", v.ID),
			PatientName:     fmt.Sprintf("%s,%s", pdata.FirstName, pdata.LastName),
			DoctorName:      pdata.DoctorName,
			AppointmentDate: v.ApppointmentDate.Format(time.DateTime),
		}

	}
	return ad, errList
}

func ErrList(el []error) *strings.Builder {
	if len(el) < 1 {
		return nil
	}
	var errString *strings.Builder
	for v := range slices.Values(el) {
		if v != nil {
			fmt.Fprintf(errString, "%s,", v.Error())
		}

	}
	return errString
}

func GetLoggedinUser(r *http.Request,
	w http.ResponseWriter) (token.Claims, error) {
	claims, err := token.ValidateJWT(w, r)
	if err != nil {
		return token.Claims{}, err
	}
	return *claims, nil
}
