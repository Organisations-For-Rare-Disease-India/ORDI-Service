package token

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateTokenCookie(userId uint, email string) (*http.Cookie, error) {
	// Expiration time of JWT token is set as 30 mins
	expirationTime := time.Now().Add(30 * time.Minute)

	// Create the JWT claims, which includes the patient email and expiration time
	claims := &Claims{
		UserId: userId,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}
	return &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	}, nil
}
