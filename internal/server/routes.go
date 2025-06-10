package server

import (
	"ORDI/cmd/web"
	"ORDI/internal/handlers/admin"
	"ORDI/internal/handlers/doctor"
	"ORDI/internal/handlers/masteradmin"
	"ORDI/internal/handlers/notification"
	"ORDI/internal/handlers/patient"
	"ORDI/internal/handlers/verification"
	"ORDI/internal/models"
	"ORDI/internal/repositories"
	"ORDI/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterPatientRoutes(
	r chi.Router,
	patientRepository repositories.Repository[models.Patient],
	appointmentRepo repositories.Repository[models.Appointment],
	notificationRepo repositories.Repository[models.Notification]) {
	patientHandler := patient.NewPatientHandler(patient.PatientHandlerConfig{
		PatientRepo:      patientRepository,
		NotificationRepo: notificationRepo,
		AppointmentRepo:  appointmentRepo,
		Cache:            s.cache,
		Email:            s.email,
	})

	r.Get(utils.PatientLoginScreen, templ.Handler(web.LoginPage(utils.PatientLoginSubmit, utils.PatientForgotPasswordScreen, utils.PatientSignupSteps)).ServeHTTP)
	r.Post(utils.PatientSignupSubmit, patientHandler.Signup)
	r.Post(utils.PatientLoginSubmit, patientHandler.Login)
	r.Get(utils.PatientDashboard, patientHandler.Dashboard)
	r.Get(utils.PatientAppointments, patientHandler.Appointment)
	r.Get(utils.PatientForgotPasswordScreen, templ.Handler(web.ForgotPasswordPage(utils.PatientForgotPasswordSubmit)).ServeHTTP)
	r.Get(utils.PatientSignupSteps, templ.Handler(web.SignupStepsPage(utils.CreatePatientSignupStepsMessage(), utils.PatientSignupForm)).ServeHTTP)
	r.Get(utils.PatientSignupForm, templ.Handler(web.PatientSignupFormPage(utils.PatientSignupSubmit)).ServeHTTP)
	r.Get(utils.PatientProfile, patientHandler.Profile)

	patientVerificationHandler := verification.NewPatientVerificationHandler(verification.PatientVerificationConfig{
		Repository: patientRepository,
		Cache:      s.cache,
		Email:      s.email,
	})
	r.Get(utils.PatientVerifyNew, patientVerificationHandler.VerifyNewUser)
	r.Get(utils.PatientVerifyExisting, patientVerificationHandler.VerifyExistingUser)
	r.Post(utils.PatientNewPassword, patientVerificationHandler.CreateNewPassword)
	r.Post(utils.PatientForgotPasswordSubmit, patientVerificationHandler.ForgotPassword)
	r.Get(utils.GenerateCaptcha, patientHandler.GenerateCaptcha)
	r.Post(utils.VerifyCaptcha, patientHandler.VerifyCaptcha)

}

func (s *Server) RegisterCommonRoutes(r chi.Router, notificationRepository repositories.Repository[models.Notification]) {
	config := notification.NotificationHandlerConfig{
		NotificationRepository: notificationRepository,
	}
	notificationHandler := notification.NewNotificationHandler(config)
	r.Get(utils.Notifications, notificationHandler.ShowNotifications)
}

func (s *Server) RegisterDoctorRoutes(r chi.Router, doctorRepository repositories.Repository[models.Doctor], notificationRepository repositories.Repository[models.Notification]) {
	doctorHandler := doctor.NewDoctorHandler(doctor.DoctorHandlerConfig{
		DoctorRepo:       doctorRepository,
		NotificationRepo: notificationRepository,
		Cache:            s.cache,
		Email:            s.email,
	})

	r.Get(utils.DoctorLoginScreen, templ.Handler(web.LoginPage(utils.DoctorLoginSubmit, utils.DoctorForgotPasswordScreen, utils.DoctorSignupSteps)).ServeHTTP)
	r.Post(utils.DoctorSignupSubmit, doctorHandler.Signup)
	r.Post(utils.DoctorLoginSubmit, doctorHandler.Login)
	r.Get(utils.DoctorDashboard, doctorHandler.Dashboard)
	r.Get(utils.DoctorAppointments, doctorHandler.Appointment)
	r.Get(utils.DoctorForgotPasswordScreen, templ.Handler(web.ForgotPasswordPage(utils.DoctorForgotPasswordSubmit)).ServeHTTP)
	r.Get(utils.DoctorSignupSteps, templ.Handler(web.SignupStepsPage(utils.CreateDoctorSignupStepsMessage(), utils.DoctorSignupForm)).ServeHTTP)
	r.Get(utils.DoctorSignupForm, templ.Handler(web.DoctorSignupFormPage(utils.DoctorSignupSubmit)).ServeHTTP)
	r.Get(utils.DoctorProfile, doctorHandler.Profile)

	doctorVerificationHandler := verification.NewDoctorVerificationHandler(verification.DoctorVerificationConfig{
		Repository: doctorRepository,
		Cache:      s.cache,
		Email:      s.email,
	})
	r.Get(utils.DoctorVerifyNew, doctorVerificationHandler.VerifyNewUser)
	r.Get(utils.DoctorVerifyExisting, doctorVerificationHandler.VerifyExistingUser)
	r.Post(utils.DoctorNewPassword, doctorVerificationHandler.CreateNewPassword)
	r.Post(utils.DoctorForgotPasswordSubmit, doctorVerificationHandler.ForgotPassword)

}

func (s *Server) RegisterAdminRoutes(r chi.Router, adminRepository repositories.Repository[models.Admin],
	patientRepository repositories.Repository[models.Patient],
	doctorRepository repositories.Repository[models.Doctor],
	appointmentRepository repositories.Repository[models.Appointment],
	notificationRepository repositories.Repository[models.Notification],
) {
	adminHandler := admin.NewAdminHandler(admin.AdminHandlerConfig{
		AdminRepo:        adminRepository,
		PatientRepo:      patientRepository,
		DoctorRepo:       doctorRepository,
		NotificationRepo: notificationRepository,
		AppointmentRepo:  appointmentRepository,
		Cache:            s.cache,
		Email:            s.email,
	})

	r.Get(utils.AdminLoginScreen, templ.Handler(web.AdminLoginPage(utils.AdminLoginSubmit, false)).ServeHTTP)
	r.Get(utils.MasterAdminLoginScreen, templ.Handler(web.AdminLoginPage(utils.MasterAdminLoginSubmit, true)).ServeHTTP)
	r.Post(utils.AdminLoginSubmit, adminHandler.Login)
	r.Get(utils.DoctorProfile, adminHandler.Profile)
	r.Get(utils.AdminCreate, templ.Handler(web.AdminCreationFormPage(utils.AdminCreateSubmit)).ServeHTTP)
	r.Get(utils.AdminLoginScreen, templ.Handler(web.AdminLoginPage(utils.AdminLoginSubmit, false)).ServeHTTP)
	r.Get(utils.AdminSetCredentials, adminHandler.Setup)
	r.Post(utils.AdminRegisterSubmit, adminHandler.Register)
	r.Get(utils.AdminDashboard, adminHandler.Dashboard)
	r.Get(utils.AdminViewDoctorList, adminHandler.ListDoctors)
	r.Get(utils.AdminViewPatientList, adminHandler.ListPatients)
	r.Get(utils.AdminAppointments, adminHandler.Appointments)
	r.Get(fmt.Sprintf("%s/{id}/edit", utils.AdminAppointmentByID), adminHandler.GetAppointmentID())
	r.Get(fmt.Sprintf("%s/{id}", utils.AdminAppointmentByID), adminHandler.GetAppointmentIDView())
	r.Put(utils.AdminAppointments, adminHandler.PutAppointment)
	r.Get("/", templ.Handler(web.AdminHomePage(utils.AdminLoginScreen, utils.MasterAdminLoginScreen)).ServeHTTP)
}

func (s *Server) RegisterMasterAdminRoutes(r chi.Router, adminRepository repositories.Repository[models.Admin],
	masterAdminRepository repositories.Repository[models.MasterAdmin],
) {
	masterAdminHandler := masteradmin.NewMasterAdminHandler(masteradmin.MasterAdminHandlerConfig{
		AdminRepo:       adminRepository,
		MasterAdminRepo: masterAdminRepository,
		Cache:           s.cache,
		Email:           s.email,
	})

	r.Get(utils.MasterAdminLoginScreen, templ.Handler(web.AdminLoginPage(utils.MasterAdminLoginSubmit, true)).ServeHTTP)
	r.Post(utils.MasterAdminLoginSubmit, masterAdminHandler.Login)
	r.Post(utils.AdminCreateSubmit, masterAdminHandler.Create)
}
func (s *Server) RegisterRoutes() http.Handler {
	mainRouter := chi.NewRouter()

	// Add all middleware first
	mainRouter.Use(middleware.Logger)

	// Initialise all the repositories
	adminRepository := repositories.NewAdminRepository(s.db)
	masterAdminRepository := repositories.NewMasterAdminRepository(s.db)
	patientRepository := repositories.NewPatientRepository(s.db)
	doctorRepository := repositories.NewDoctorRepository(s.db)
	appointmentRepository := repositories.NewAppointmentRepository(s.db)
	notificationRepository := repositories.NewNotificationRepository(s.db)

	mainRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip domain check for assets and health endpoint
			if strings.HasPrefix(r.URL.Path, "/assets/") || r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			domain := extractDomain(r.Host)
			// if domain == utils.InternalDomain {
			if domain == "localhost" {
				// Create admin router
				adminRouter := chi.NewRouter()
				s.RegisterAdminRoutes(
					adminRouter, adminRepository, patientRepository,
					doctorRepository, appointmentRepository, notificationRepository)
				s.RegisterMasterAdminRoutes(adminRouter, adminRepository, masterAdminRepository)
				adminRouter.ServeHTTP(w, r)
				return
			}

			// Create public router
			publicRouter := chi.NewRouter()
			publicRouter.Get("/", templ.Handler(web.HomePage()).ServeHTTP)
			publicRouter.Get(utils.HomeLogin, templ.Handler(web.ChooseRolePage(utils.DoctorLoginScreen, utils.PatientLoginScreen)).ServeHTTP)
			publicRouter.Get(utils.HomeSignup, templ.Handler(web.ChooseRolePage(utils.DoctorSignupSteps, utils.PatientSignupSteps)).ServeHTTP)

			// Patient specific handlers
			s.RegisterPatientRoutes(publicRouter, patientRepository, appointmentRepository, notificationRepository)

			// Doctor specific handlers
			s.RegisterDoctorRoutes(publicRouter, doctorRepository, notificationRepository)

			// Public Routes handlers
			s.RegisterCommonRoutes(publicRouter, notificationRepository)

			publicRouter.ServeHTTP(w, r)
		})
	})

	// After all middleware, add the universal routes
	fileServer := http.FileServer(http.FS(web.Files))
	mainRouter.Handle("/assets/*", fileServer)
	mainRouter.Get("/health", s.healthHandler)

	return mainRouter
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

// extractDomain extracts the domain part from the full host (with or without port)
func extractDomain(host string) string {
	// Split the host by ':' to handle the case with port
	hostParts := strings.Split(host, ":")
	return hostParts[0] // Return only the domain part
}
