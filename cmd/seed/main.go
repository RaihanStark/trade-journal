package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/raihanstark/trade-journal/internal/seed"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get database URL
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Connect to database
	dbConn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Test database connection
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to database")

	// Seed the random number generator
	gofakeit.Seed(time.Now().UnixNano())

	// Create seeder instance
	seeder := seed.NewSeeder(dbConn)
	ctx := context.Background()

	// Clear existing data (optional - comment out if you want to keep existing data)
	if err := seeder.ClearData(ctx); err != nil {
		log.Fatalf("Failed to clear data: %v", err)
	}

	// Run seeding with default configuration
	config := seed.DefaultConfig()
	if err := seeder.Run(ctx, config); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("\nTest credentials:")
	log.Println("Email: test1@example.com, Password: password123")
	log.Println("Email: test2@example.com, Password: password123")
	log.Println("Email: test3@example.com, Password: password123")
}
