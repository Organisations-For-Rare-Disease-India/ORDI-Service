package main

import (
	"ORDI/internal/database/mysql"
	"ORDI/internal/models"
	"ORDI/internal/repositories"
	"ORDI/internal/server"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Initialize the database
	envConfig, err := server.LoadConfig(".")
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

	db := mysql.NewMySqlConnection(dbConfig)
	defer db.Close()
	repo := repositories.NewMasterAdminRepository(db)

	ctx := context.Background()

	reader := bufio.NewReader(os.Stdin)
	for {
		// Take email input
		fmt.Print("Enter email: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)

		// Check if email already exists
		existingAdmin, err := repo.FindByField(ctx, "email_id", email)
		if err != nil {
			log.Fatalf("Error checking email: %v", err)
		}
		if existingAdmin != nil {
			fmt.Println("Admin with this email already exists. Try again.")
			continue
		}

		// Take password input
		fmt.Print("Enter password: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Error hashing password: %v", err)
		}

		// Save to the database
		admin := &models.MasterAdmin{
			Email:    email,
			Password: string(hashedPassword),
		}
		if err := repo.Save(ctx, admin); err != nil {
			log.Fatalf("Error saving admin: %v", err)
		}

		fmt.Println("Admin added successfully!")

		// Ask if the user wants to add another admin
		fmt.Print("Do you want to add another admin? (y/n): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		if strings.ToLower(choice) != "y" {

			break
		}
	}
}
