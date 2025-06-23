package token

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

const JWT_TOKEN_HEADER = "token"

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

type Claims struct {
	UserId uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func ValidateJWT(w http.ResponseWriter, r *http.Request) (*Claims, error) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie(JWT_TOKEN_HEADER)
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return nil, err
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return nil, err
		}
		w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, err
	}
	return claims, nil
}
