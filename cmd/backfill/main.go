package main

import (
	"academ_aide/internal/ai"
	"academ_aide/internal/config"
	"academ_aide/internal/repository"
	"context"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Initialize DB
	config.InitDB()
	db := config.PostgresDB

	repo := repository.NewCourseRepository(db)
	embedder := ai.NewEmbedder()

	log.Println("Starting backfill for course embeddings...")

	courses, err := repo.GetAllCoursesWithoutEmbeddings(context.Background())
	if err != nil {
		log.Fatalf("Failed to fetch courses: %v", err)
	}

	log.Printf("Found %d courses needing embeddings", len(courses))

	for _, course := range courses {
		log.Printf("Processing course: %s", course.Title)

		if course.Description == "" {
			log.Printf("Skipping course %s (no description)", course.CourseID)
			continue
		}

		embedding, err := embedder.GenerateEmbedding(course.Description)
		if err != nil {
			log.Printf("Failed to generate embedding for %s: %v", course.CourseID, err)
			continue
		}

		if err := repo.UpdateEmbedding(context.Background(), course.CourseID, embedding); err != nil {
			log.Printf("Failed to update embedding for %s: %v", course.CourseID, err)
			continue
		}

		log.Printf("Successfully updated embedding for %s", course.CourseID)
	}

	log.Println("Backfill completed.")
}
