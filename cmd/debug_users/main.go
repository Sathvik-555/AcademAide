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
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, defaults may fail if not localhost")
	}

	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=academ_aide port=5432 sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT student_id, s_first_name FROM STUDENT")
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	defer rows.Close()

	fmt.Println("--- Students in DB ---")
	count := 0
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: '%s', Name: '%s'\n", id, name)
		count++
	}
	fmt.Printf("Total: %d\n", count)
}
