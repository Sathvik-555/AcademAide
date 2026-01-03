package services

import (
	"academ_aide/internal/config"
	"academ_aide/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type QuizService struct{}

func NewQuizService() *QuizService {
	return &QuizService{}
}

// GenerateQuiz generates a quiz for a given course based on its syllabus
func (s *QuizService) GenerateQuiz(courseID string) (*models.Quiz, error) {
	// 1. Fetch Syllabus Topics
	rows, err := config.PostgresDB.Query("SELECT topic FROM SYLLABUS_UNIT WHERE course_id=$1", courseID)
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
		return nil, fmt.Errorf("no syllabus found for course %s", courseID)
	}

	topicStr := strings.Join(topics, ", ")

	// 2. Prompt Ollama
	prompt := fmt.Sprintf(`
You are a professor. Generate a quiz with 5 multiple-choice questions for the course %s based on these topics: %s.
Return ONLY valid JSON in the following format, with no extra text:
{
  "questions": [
    {
      "id": 1,
      "text": "Question text here?",
      "options": ["Option A", "Option B", "Option C", "Option D"],
      "correct_option": 0
    }
  ]
}
`, courseID, topicStr)

	jsonResp, err := s.callOllamaJSON(prompt)
	if err != nil {
		return nil, err
	}

	// 3. Parse Response
	var quizStructure struct {
		Questions []models.Question `json:"questions"`
	}
	if err := json.Unmarshal([]byte(jsonResp), &quizStructure); err != nil {
		log.Println("Ollama JSON Parse Error:", err)
		log.Println("Raw Response:", jsonResp)
		return nil, fmt.Errorf("failed to parse AI response")
	}

	// 4. Create and Save Quiz
	quiz := &models.Quiz{
		CourseID:  courseID,
		Topic:     "General Syllabus",
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
