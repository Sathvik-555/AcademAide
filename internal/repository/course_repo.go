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
	sb.WriteString("[") // pgvector uses [0.1,0.2]
	for i, v := range vec {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("%f", v))
	}
	sb.WriteString("]")
	return sb.String()
}

type CourseMaterial struct {
	Content    string
	CourseID   string
	UnitNo     int
	Score      float64
	SourceFile string
}

// SearchMaterials finds the most relevant material chunks using the custom cosine_similarity function
func (r *CourseRepository) SearchMaterials(ctx context.Context, embedding []float32, limit int, courseIDFilter []string, unitFilter int) ([]CourseMaterial, error) {
	// Format as PGVector string: [0.1, 0.2, ...]
	vectorStr := vecToString(embedding)

	// Build Query
	var query string
	var args []interface{}

	if len(courseIDFilter) > 0 {
		// PostgreSQL ANY operator for array filtering
		// We need to convert the slice to a postgres array string literal if we were binding it as a single string,
		// but using pq.Array or similar is better. However, since we might not have lib/pq,
		// we can simpler use "course_id = ANY($3)" and pass the slice directly if the driver supports it.
		// Standard lib/pq supports passing []string for []text or []varchar.
		// Let's assume standard driver behavior.

		if unitFilter > 0 {
			// Filter by Course IDs AND Unit
			query = `
				SELECT content_text, course_id, unit_no, 1 - (embedding <=> $1::vector) as score, source_file
				FROM COURSE_MATERIAL_CHUNK
				WHERE course_id = ANY($3::text[]) AND unit_no = $4
				ORDER BY embedding <=> $1::vector ASC
				LIMIT $2
			`
			args = []interface{}{vectorStr, limit, (courseIDFilter), unitFilter}
		} else {
			// Filter by Course IDs only
			query = `
				SELECT content_text, course_id, unit_no, 1 - (embedding <=> $1::vector) as score, source_file
				FROM COURSE_MATERIAL_CHUNK
				WHERE course_id = ANY($3::text[])
				ORDER BY embedding <=> $1::vector ASC
				LIMIT $2
			`
			args = []interface{}{vectorStr, limit, (courseIDFilter)}
		}
	} else {
		// Global Search (No Course Filter)
		query = `
			SELECT content_text, course_id, unit_no, 1 - (embedding <=> $1::vector) as score, source_file
			FROM COURSE_MATERIAL_CHUNK
			ORDER BY embedding <=> $1::vector ASC
			LIMIT $2
		`
		args = []interface{}{vectorStr, limit}
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var materials []CourseMaterial
	for rows.Next() {
		var m CourseMaterial
		var sourceFile sql.NullString
		if err := rows.Scan(&m.Content, &m.CourseID, &m.UnitNo, &m.Score, &sourceFile); err != nil {
			return nil, err
		}
		if sourceFile.Valid {
			m.SourceFile = sourceFile.String
		}
		materials = append(materials, m)
	}
	return materials, nil
}
