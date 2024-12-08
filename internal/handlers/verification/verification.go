package verification

import (
	"ORDI/internal/messages"
	"net/http"

	"github.com/gorilla/schema"
)

type Verification interface {
	VerifyNewUser(http.ResponseWriter, *http.Request)
	VerifyExistingUser(http.ResponseWriter, *http.Request)
	CreateNewPassword(http.ResponseWriter, *http.Request)
	ForgotPassword(http.ResponseWriter, *http.Request)
}

var emailMessage messages.SubmitMessage = messages.SubmitMessage{
	Title:   "Forgot Password Request sent",
	Message: "A verification email has been sent to your email address. Please click on the link provided to proceed ahead.",
}

var decoder = schema.NewDecoder()

const (
	TokenHeader                = "token"
	TokenRequiredMessage       = "Token is required"
	InvalidRequestMessage      = "Invalid request"
	UnauthorizedRequestMessage = "Unauthorized request"
	InvalidCredentialsMessage  = "Invalid username or password"
)
