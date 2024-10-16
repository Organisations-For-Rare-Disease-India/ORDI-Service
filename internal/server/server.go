package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"ORDI/internal/cache"
	"ORDI/internal/cache/redisClient"
	"ORDI/internal/database"
	"ORDI/internal/database/mysql"
	"ORDI/internal/email"
	"ORDI/internal/email/emailSender"
)

type Server struct {
	url   string
	port  int
	db    database.Database
	email email.Email
	cache cache.Cache
}

func NewServer() *http.Server {
	baseURL := os.Getenv("BASE_URL")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		url:   baseURL,
		port:  port,
		db:    mysql.NewDefaultSqlConnection(),
		email: emailSender.NewDefaultEmailSender(),
		cache: redisClient.NewDefaultRedisClient(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
