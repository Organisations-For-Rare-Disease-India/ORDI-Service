package server

import (
	"ORDI/cmd/web"
	"context"
	"net/http"

	"github.com/a-h/templ"
)

func (s *Server) EmailVerificationHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	email, err := s.cache.Get(ctx, token)
	if err != nil || email == "" {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
	}

	// TODO: Mark user as verified on the database

	templ.Handler(web.LoginPage()).ServeHTTP(w, r)
}
