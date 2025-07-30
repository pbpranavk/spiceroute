package main

import (
	"log"
	"os"

	"spiceroute/pkg/database"
)

func main() {
	// Check if DB_DSN is set
	if os.Getenv("DB_DSN") == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	// Initialize database connection
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Database migrations completed successfully!")
}
