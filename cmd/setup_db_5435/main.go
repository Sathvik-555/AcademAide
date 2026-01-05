package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Hardcoded DSN for port 5435 and dbname postgres
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5435 sslmode=disable"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB (is Docker container up?): %v", err)
	}
	fmt.Println("Connected to Docker Postgres (Port 5435).")

	// Read and execute files in order
	files := []string{
		"database/rag_setup.sql",
	}

	for _, file := range files {
		fmt.Printf("Executing %s...\n", file)
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read file %s: %v", file, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			log.Printf("Error executing %s: %v\n", file, err)
			// Don't exit fatal, maybe idempotent errors
		} else {
			fmt.Printf("Successfully executed %s\n", file)
		}
	}

	fmt.Println("Database setup complete.")
}
