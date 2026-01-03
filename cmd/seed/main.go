package main

import (
	"academ_aide/internal/config"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Initialize DBs
	config.InitDB()

	// 1. Insert Student
	studentID := "2024CS123"
	_, err := config.PostgresDB.Exec(`
		INSERT INTO STUDENT (student_id, s_first_name, s_last_name, s_email, s_phone_no, semester, year_of_joining, dept_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (student_id) DO NOTHING
	`, studentID, "Sathvik", "User", "sathvik@univ.edu", "9999999999", 3, 2024, "CS")

	if err != nil {
		log.Fatalf("Failed to insert student: %v", err)
	}
	fmt.Println("Student inserted (or already exists)")

	// 2. Insert Enrollment (for Timetable)
	// Enroll in CS101, CS102
	courses := []string{"CS101", "CS102"}
	for _, cid := range courses {
		_, err := config.PostgresDB.Exec(`
			INSERT INTO ENROLLS_IN (student_id, course_id, status)
			VALUES ($1, $2, 'Enrolled')
			ON CONFLICT (student_id, course_id) DO NOTHING
		`, studentID, cid)

		if err != nil {
			log.Printf("Failed to enroll in %s: %v", cid, err)
		} else {
			fmt.Printf("Enrolled in %s\n", cid)
		}
	}

	// 3. Ensure Schedule Exists for these courses (Already in schema.sql, but good to be safe if schema wasn't run fully)
	// We rely on existing seed data for schedules.

	fmt.Println("Seeding complete!")
}
