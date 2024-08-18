package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"ORDI/internal/database"
	"ORDI/internal/database/mysql"
	"ORDI/internal/storage"
	"ORDI/internal/storage/s3"
)

type Server struct {
	port int
	db   database.Database
	s3   storage.Storage
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
		db:   mysql.NewMySqlConnection(),
		s3:   s3.NewS3ServiceConnection(),
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
