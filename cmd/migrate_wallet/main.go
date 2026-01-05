package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it, proceeding with defaults")
	}

	// Connect to DB using the same logic as config (simplified)
	dsn := "host=localhost user=postgres password=postgres dbname=academ_aide port=5432 sslmode=disable"
	// Check if env var override exists
	if envDSN := os.Getenv("POSTGRES_DSN"); envDSN != "" {
		dsn = envDSN
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}
	fmt.Println("Connected to DB.")

	// Read SQL file
	// We'll just hardcode the query here to avoid file path issues in execution from random generic directories
	query := `
	ALTER TABLE STUDENT ADD COLUMN IF NOT EXISTS wallet_address VARCHAR(42) UNIQUE;
	ALTER TABLE STUDENT ADD COLUMN IF NOT EXISTS encrypted_private_key TEXT;
	`

	fmt.Println("Executing migration...")
	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migration successful: wallet_address and encrypted_private_key columns added.")
}
