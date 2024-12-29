package masteradmin

import (
	"ORDI/internal/cache"
	"ORDI/internal/email"
	"ORDI/internal/models"
	"ORDI/internal/repositories"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/schema"
)

const JWT_TOKEN_HEADER = "token"
const EMAIL_HEADER = "email"
const EmailRequiredError = "Incorrect URL"

var decoder = schema.NewDecoder()

// MasterAdmin Creates the Admin
// Admin would then Register their profile
type MasterAdmin interface {
	Login(http.ResponseWriter, *http.Request)
	Create(http.ResponseWriter, *http.Request)
}

type masterAdminHandler struct {
	adminRepository       repositories.Repository[models.Admin]
	masterAdminRepository repositories.Repository[models.MasterAdmin]
	email                 email.Email
}

type MasterAdminHandlerConfig struct {
	AdminRepo       repositories.Repository[models.Admin]
	MasterAdminRepo repositories.Repository[models.MasterAdmin]
	Cache           cache.Cache
	Email           email.Email
}

func NewMasterAdminHandler(config MasterAdminHandlerConfig) MasterAdmin {
	return &masterAdminHandler{
		adminRepository:       config.AdminRepo,
		masterAdminRepository: config.MasterAdminRepo,
		email:                 config.Email,
	}
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
