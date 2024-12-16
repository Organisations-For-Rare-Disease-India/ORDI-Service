package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"

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
	envConfig, err := LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Add the DB, email and cache configs from environment and initialize them
	dbConfig := mysql.SQLConfig{
		DatabaseName: envConfig.DBName,
		Username:     envConfig.DBUser,
		Password:     envConfig.DBPass,
		Host:         envConfig.DBHost,
		Port:         envConfig.DBPort,
	}

	emailConfig := emailSender.EmailConfig{
		SMTPHost:     envConfig.SMTPHost,
		SMTPPort:     envConfig.SMTPPort,
		SMTPUsername: envConfig.SMTPUser,
		SMTPPassword: envConfig.SMTPPass,
		FromAddress:  envConfig.EmailFromAddress,
	}

	cacheConfig := redisClient.RedisConfig{
		ADDR:     envConfig.RedisAddr,
		Password: envConfig.RedisPwd,
	}

	NewServer := &Server{
		url:   envConfig.Url,
		port:  envConfig.Port,
		db:    mysql.NewMySqlConnection(dbConfig),
		email: emailSender.NewEmailSender(&emailConfig),
		cache: redisClient.NewRedisClient(cacheConfig),
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

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config EnvConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
