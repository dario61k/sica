package main

import (
	"log"
	"os"
	"sica/internal"
	"sica/internal/database"
	"sica/internal/models"
	"sica/internal/repositories"
	"sica/pkg/bcrypt"
)

func main() {

	database.GetDB()
	database.AutoMigrate()

	createAdmin()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Println("⚠️  PORT variable not defined, using default:", port)
	}

	server := internal.SetupRoutes()
	log.Println("🚀 Server running on port:", port)
	if err := server.Run(":" + port); err != nil {
		log.Fatal("❌ Error starting server:", err)
	}

}

func createAdmin() {

	r := repositories.NewAuthRepository()
	_, err := r.Get(1)
	if err == nil {
		return
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		log.Fatal("ADMIN_PASSWORD environment variable not set")
	}

	hashedPass, err := bcrypt.HashPassword(adminPassword)
	if err != nil {
		log.Fatal("error hashing admin password: ", err)
	}

	admin := models.Auth{
		ID:       1,
		Password: hashedPass,
	}

	if _, err := r.Create(&admin); err != nil {
		log.Fatal("error creating admin: ", err)
	}

}
