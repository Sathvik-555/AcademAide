package services

import (
	"academ_aide/internal/ai"
	"academ_aide/internal/config"
	"academ_aide/internal/models"
	"academ_aide/internal/repository"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type QuizService struct {
	Embedder *ai.Embedder
	Repo     *repository.CourseRepository
}

func NewQuizService() *QuizService {
	return &QuizService{
		Embedder: ai.NewEmbedder(),
		Repo:     repository.NewCourseRepository(config.PostgresDB),
	}
}

// GenerateQuiz generates a quiz for a given course based on its syllabus
func (s *QuizService) GenerateQuiz(courseID string, unit int) (*models.Quiz, error) {
	// 1. Fetch Syllabus Topics
	var rows *sql.Rows
	var err error

	if unit > 0 {
		rows, err = config.PostgresDB.Query("SELECT topic FROM SYLLABUS_UNIT WHERE course_id=$1 AND unit_no=$2", courseID, unit)
	} else {
		rows, err = config.PostgresDB.Query("SELECT topic FROM SYLLABUS_UNIT WHERE course_id=$1", courseID)
	}

	if err != nil {
		return nil, fmt.Errorf("fetching syllabus: %w", err)
	}
	defer rows.Close()

	var topics []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err == nil {
			topics = append(topics, t)
		}
	}

	if len(topics) == 0 {
		return nil, fmt.Errorf("no syllabus found for %s (Unit %d)", courseID, unit)
	}

	topicStr := strings.Join(topics, ", ")

	// 2. RAG Retrieval
	// Generate embedding for the broad topic context
	var query string
	if unit > 0 {
		query = fmt.Sprintf("Important concepts in %s Unit %d", courseID, unit)
	} else {
		query = fmt.Sprintf("Important concepts in %s: %s", courseID, topicStr)
	}

	embedding, err := s.Embedder.GenerateEmbedding(query)

	var contextText string
	if err != nil {
		log.Println("Quiz embedding failed, falling back to basic prompt:", err)
		contextText = "No course materials available."
	} else {
		// Fetch top 5 relevant chunks, filtering by unit if specified
		materials, err := s.Repo.SearchMaterials(context.Background(), embedding, 5, []string{courseID}, unit)
		if err != nil {
			log.Println("Quiz material search failed:", err)
		} else {
			var sb strings.Builder
			for _, m := range materials {
				sb.WriteString(fmt.Sprintf("---\nSource: %s\nUnit: %d\nContent: %s\n", m.SourceFile, m.UnitNo, m.Content))
			}
			contextText = sb.String()
		}
	}

	// 3. Prompt Ollama
	var unitContext string
	if unit > 0 {
		unitContext = fmt.Sprintf("focusing STRICTLY on Unit %d", unit)
	} else {
		unitContext = "covering the course topics"
	}

	prompt := fmt.Sprintf(`
You are a professor. Generate a quiz with 5 multiple-choice questions for the course %s, %s.

CRITICAL INSTRUCTION:
1. You MUST use ONLY the content provided below in "Context Materials" to generate the questions.
2. Do NOT use outside knowledge. If the provided materials are insufficient, do not make up facts.
3. For each question, the "reference" field MUST be the exact 'Source' filename provided in the Context Materials (e.g. "Unit1_Intro.pdf"). Do not hallucinate filenames.

Context Materials:
%s

Topics: %s

Return ONLY valid JSON in the following format, with no extra text:
{
  "questions": [
    {
      "id": 1,
      "text": "Question text here?",
      "options": ["Option A", "Option B", "Option C", "Option D"],
      "correct_option": 0,
      "reference": "Exact_Source_File_Name.pdf"
    }
  ]
}
`, courseID, unitContext, contextText, topicStr)

	jsonResp, err := s.callOllamaJSON(prompt)
	if err != nil {
		return nil, err
	}

	// 4. Parse Response
	var quizStructure struct {
		Questions []models.Question `json:"questions"`
	}
	if err := json.Unmarshal([]byte(jsonResp), &quizStructure); err != nil {
		log.Println("Ollama JSON Parse Error:", err)
		log.Println("Raw Response:", jsonResp)
		// Fallback attempt: try to find JSON in snippet if there's extra text
		start := strings.Index(jsonResp, "{")
		end := strings.LastIndex(jsonResp, "}")
		if start != -1 && end != -1 && end > start {
			jsonResp = jsonResp[start : end+1]
			if err := json.Unmarshal([]byte(jsonResp), &quizStructure); err != nil {
				return nil, fmt.Errorf("failed to parse AI response")
			}
		} else {
			return nil, fmt.Errorf("failed to parse AI response")
		}
	}

	// 5. Create and Save Quiz
	quiz := &models.Quiz{
		CourseID:  courseID,
		Topic:     "Generated from Course Materials",
		Questions: quizStructure.Questions,
		CreatedAt: time.Now(),
	}

	coll := config.MongoDB.Collection("quizzes")
	res, err := coll.InsertOne(context.Background(), quiz)
	if err != nil {
		return nil, fmt.Errorf("saving quiz: %w", err)
	}

	// Convert primitive.ObjectID to string if needed, but for now we rely on Mongo ID
	// If we needed the ID, we'd cast res.InsertedID
	_ = res

	return quiz, nil
}

func (s *QuizService) callOllamaJSON(prompt string) (string, error) {
	url := "http://localhost:11434/api/generate"
	reqBody := map[string]interface{}{
		"model":  "llama3.2",
		"prompt": prompt,
		"stream": false,
		"format": "json", // Force JSON mode
	}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("ollama connection failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var oResp struct {
		Response string `json:"response"`
	}
	json.Unmarshal(body, &oResp)

	return oResp.Response, nil
}
