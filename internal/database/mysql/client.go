package mysql

import (
	"ORDI/internal/database"
	"context"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlService struct {
	db *gorm.DB
}

type SQLConfig struct {
	DatabaseName string
	Username     string
	Password     string
	Host         string
	Port         int
}

func NewDefaultSqlConnection() database.Database {
	sqlConfig := SQLConfig{
		DatabaseName: "ORDI",
		Username:     "root",
		Password:     "",
		Host:         "localhost",
		Port:         3306,
	}
	return NewMySqlConnection(sqlConfig)
}

func NewMySqlConnection(config SQLConfig) database.Database {
	// Constructing the data source name string
	rootDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
	)

	// Open connection to MySQL server
	rootDB, err := gorm.Open(mysql.Open(rootDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to MySQL server: ", err)
	}

	// Check if the database exists
	var exists int
	rootDB.Raw("SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", config.DatabaseName).Scan(&exists)

	if exists == 0 {
		// Create the database
		err = rootDB.Exec(fmt.Sprintf("CREATE DATABASE `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci", config.DatabaseName)).Error
		if err != nil {
			log.Fatal("Failed to create database: ", err)
		}
		log.Printf("Database %s created successfully", config.DatabaseName)
	} else {
		log.Printf("Database %s already exists", config.DatabaseName)
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DatabaseName,
	)

	// Opening a connection using gorm.
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}

	mysqlInstance := &mysqlService{
		db: db,
	}
	return mysqlInstance

}

func (s *mysqlService) Health() map[string]string {

	stats := make(map[string]string)
	sqlDB, err := s.db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
		stats["status"] = "down"
		stats["message"] = fmt.Sprintf("db down: %v", err)
		return stats
	}
	if err := sqlDB.Ping(); err != nil {
		stats["status"] = "down"
		stats["message"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	return stats

}

func (s *mysqlService) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (s *mysqlService) FindByID(ctx context.Context, id uint, entity interface{}) error {
	if err := s.db.WithContext(ctx).First(entity, id).Error; err != nil {
		return err
	}
	return nil
}

func (s *mysqlService) Delete(ctx context.Context, entity interface{}) error {
	if err := s.db.WithContext(ctx).Delete(entity).Error; err != nil {
		return err
	}
	return nil
}

func (s *mysqlService) Save(ctx context.Context, entity interface{}) error {
	if err := s.db.WithContext(ctx).Save(entity).Error; err != nil {
		return err
	}
	return nil
}

func (s *mysqlService) FindByField(ctx context.Context, entity interface{}, field string, value interface{}) error {
	if err := s.db.WithContext(ctx).Where(field+" = ?", value).First(entity).Error; err != nil {
		return err
	}
	return nil
}

func (s *mysqlService) AutoMigrate(ctx context.Context, entity interface{}) error {
	if err := s.db.WithContext(ctx).AutoMigrate(entity); err != nil {
		return err
	}
	return nil
}
