package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5435 sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := `
	CREATE TABLE IF NOT EXISTS COURSE_MATERIAL_CHUNK (
		chunk_id SERIAL PRIMARY KEY,
		course_id VARCHAR(10) NOT NULL,
		unit_no INTEGER NOT NULL,       
		content_text TEXT NOT NULL,     
		embedding vector(768),          
		source_file VARCHAR(255),       
		chunk_index INTEGER,            
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_material_course FOREIGN KEY (course_id) REFERENCES COURSE(course_id)
	);
	`
	fmt.Println("Creating COURSE_MATERIAL_CHUNK table...")
	_, err = db.Exec(query)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Println("Table created successfully!")
	}
}
