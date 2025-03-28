// main.go
package main

import (
	"cal-blog-service/models"
	"cal-blog-service/routes"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func initDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db
	log.Println("Database connected successfully")
}

func runMigrations() {
	err := DB.AutoMigrate(&models.BlogPost{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
	log.Println("Database migrations applied successfully")
}

func main() {
	initDatabase()
	runMigrations()

	// Setup router from the routes package
	r := routes.SetupRouter(DB)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Default port if not specified
	}

	// Start the server
	log.Printf("Server starting on port %s...", port)
	err := r.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
