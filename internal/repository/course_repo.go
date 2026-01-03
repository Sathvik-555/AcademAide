package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type CourseRepository struct {
	DB *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{DB: db}
}

type Course struct {
	CourseID    string
	Title       string
	Description string
	Credits     int
	DeptID      string
}

// UpdateEmbedding updates the embedding for a specific course
func (r *CourseRepository) UpdateEmbedding(ctx context.Context, courseID string, embedding []float32) error {
	// Convert embedding to string format "[0.1, 0.2, ...]" for pgvector
	vectorStr := vecToString(embedding)
	query := `UPDATE COURSE SET embedding = $1 WHERE course_id = $2`
	_, err := r.DB.ExecContext(ctx, query, vectorStr, courseID)
	return err
}

// SearchByVector finds courses closest to the query vector
func (r *CourseRepository) SearchByVector(ctx context.Context, embedding []float32, limit int) ([]Course, error) {
	vectorStr := vecToString(embedding)
	query := `
		SELECT course_id, title, description, credits, dept_id
		FROM COURSE
		ORDER BY embedding <=> $1
		LIMIT $2
	`
	rows, err := r.DB.QueryContext(ctx, query, vectorStr, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var c Course
		var desc sql.NullString // Description might be null if not backfilled fully
		if err := rows.Scan(&c.CourseID, &c.Title, &desc, &c.Credits, &c.DeptID); err != nil {
			return nil, err
		}
		if desc.Valid {
			c.Description = desc.String
		}
		courses = append(courses, c)
	}
	return courses, nil
}

// GetAllCoursesWithoutEmbeddings returns courses that possess a description but no embedding
func (r *CourseRepository) GetAllCoursesWithoutEmbeddings(ctx context.Context) ([]Course, error) {
	query := `SELECT course_id, title, description, credits, dept_id FROM COURSE WHERE embedding IS NULL AND description IS NOT NULL`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var c Course
		var desc sql.NullString
		if err := rows.Scan(&c.CourseID, &c.Title, &desc, &c.Credits, &c.DeptID); err != nil {
			return nil, err
		}
		c.Description = desc.String
		courses = append(courses, c)
	}
	return courses, nil
}

// Helper to format float slice to vector string
func vecToString(vec []float32) string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, v := range vec {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("%f", v))
	}
	sb.WriteString("]")
	return sb.String()
}
