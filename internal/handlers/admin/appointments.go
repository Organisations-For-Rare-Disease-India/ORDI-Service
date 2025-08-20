package admin

import (
	"ORDI/cmd/web"
	"ORDI/internal/models"
	"ORDI/internal/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
)

func (a *adminHandler) Appointments(w http.ResponseWriter, r *http.Request) {
	appointments, err := a.appointmentRepository.FindAllWithPage(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	appointmentsData := make([]models.AppointmentData, len(appointments))
	for i, v := range appointments {
		// get patient
		patientFromStore, err := a.patientRepository.FindByID(r.Context(), v.PatientID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		appointmentsData[i] = models.AppointmentData{
			AppointmentID:   fmt.Sprintf("%d", v.ID),
			PatientName:     fmt.Sprintf("%s,%s", patientFromStore.FirstName, patientFromStore.LastName),
			DoctorName:      patientFromStore.DoctorName,
			AppointmentDate: v.ApppointmentDate.Format(time.DateTime),
		}

	}
	templ.Handler(web.AdminAppointmentsPage(utils.AdminAppointments, appointmentsData)).ServeHTTP(w, r)
}

// updates existing appointment
func (a *adminHandler) PutAppointment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, fmt.Errorf("method:%s not allowed", r.Method).Error(), http.StatusBadRequest)
		return
	}
	// get data from req
	if err := r.ParseForm(); err != nil {
		http.Error(w,
			fmt.Errorf("error reading form data:%v", err).Error(),
			http.StatusBadRequest)
		return
	}
	formValues := r.Form
	aid := formValues.Get("appointmentID")
	at := formValues.Get("appointmentDate")
	appointmentID, err := strconv.Atoi(aid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	appointmentFromStore, err := a.appointmentRepository.FindByID(r.Context(), uint(appointmentID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	updatedTime, err := time.Parse(time.DateTime, at)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	appointmentFromStore.ApppointmentDate = updatedTime
	appointmentFromStore.UpdatedAt = time.Now()
	if err := a.appointmentRepository.Save(r.Context(),
		appointmentFromStore); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// get patient
	patientFromStore, err := a.patientRepository.FindByID(r.Context(),
		appointmentFromStore.PatientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	appointmentData := models.AppointmentData{
		AppointmentID:   fmt.Sprintf("%d", appointmentID),
		AppointmentDate: updatedTime.Format(time.DateTime),
		PatientName:     fmt.Sprintf("%s,%s", patientFromStore.FirstName, patientFromStore.LastName),
		DoctorName:      patientFromStore.DoctorName,
	}
	templ.Handler(web.GetAppointmentView(&appointmentData)).ServeHTTP(w, r)
}

func (a *adminHandler) GetAppointmentID() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, fmt.Errorf("empty value received for id").Error(),
				http.StatusBadRequest)
			return
		}
		appointmentID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		appointment, err := a.appointmentRepository.FindByID(r.Context(), uint(appointmentID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// get patient
		patientFromStore, err := a.patientRepository.FindByID(r.Context(),
			appointment.PatientID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		a := &models.AppointmentData{
			AppointmentID:   fmt.Sprintf("%d", appointment.ID),
			PatientName:     fmt.Sprintf("%s,%s", patientFromStore.FirstName, patientFromStore.LastName),
			DoctorName:      patientFromStore.DoctorName,
			AppointmentDate: appointment.ApppointmentDate.Format(time.DateTime),
		}

		templ.Handler(web.GetAppointmentByID(*a)).ServeHTTP(w, r)

	})
}

func (a *adminHandler) GetAppointmentIDView() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, fmt.Errorf("empty value received for id").Error(),
				http.StatusBadRequest)
			return
		}
		appointmentID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		appointment, err := a.appointmentRepository.FindByID(r.Context(), uint(appointmentID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// get patient
		patientFromStore, err := a.patientRepository.FindByID(r.Context(),
			appointment.PatientID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		a := &models.AppointmentData{
			AppointmentID: fmt.Sprintf("%d", appointment.ID),
			PatientName: fmt.Sprintf("%s,%s", patientFromStore.FirstName,
				patientFromStore.LastName),
			DoctorName:      patientFromStore.DoctorName,
			AppointmentDate: appointment.ApppointmentDate.Format(time.DateTime),
		}

		templ.Handler(web.GetAppointmentView(a)).ServeHTTP(w, r)

	})
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
