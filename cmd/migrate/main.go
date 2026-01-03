package main

import (
	"academ_aide/internal/config"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	config.InitDB()
	db := config.PostgresDB

	// Read schema file
	schemaBytes, err := os.ReadFile("database/schema.sql")
	if err != nil {
		log.Fatal("Failed to read database/schema.sql:", err)
	}

	log.Println("Applying schema...")
	// Execute schema
	_, err = db.Exec(string(schemaBytes))
	if err != nil {
		log.Fatal("Failed to apply schema:", err)
	}

	log.Println("Schema applied successfully!")
}
