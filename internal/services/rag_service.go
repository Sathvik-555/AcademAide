package services

import (
	"academ_aide/internal/ai"
	"academ_aide/internal/config"
	"academ_aide/internal/models"
	"academ_aide/internal/repository"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RAGService struct {
	Embedder *ai.Embedder
	Repo     *repository.CourseRepository
}

func NewRAGService() *RAGService {
	return &RAGService{
		Embedder: ai.NewEmbedder(),
		Repo:     repository.NewCourseRepository(config.PostgresDB),
	}
}

// AnalyzeSentiment - Simple keyword-based or mock (since we are doing RAG)
func (s *RAGService) AnalyzeSentiment(message string) string {
	lower := strings.ToLower(message)
	if strings.Contains(lower, "bad") || strings.Contains(lower, "hate") || strings.Contains(lower, "fail") {
		return "negative"
	} else if strings.Contains(lower, "good") || strings.Contains(lower, "love") || strings.Contains(lower, "thanks") {
		return "positive"
	}
	return "neutral"
}

func (s *RAGService) getAgentPrompt(agentID string) string {
	switch agentID {
	case "socratic":
		return "You are a Socratic Tutor. Your goal is to guide the student to the answer by asking probing questions. Do not give the answer directly. Break down complex problems into smaller steps."
	case "code_reviewer":
		return "You are an expert Code Reviewer. Analyze the student's code for bugs, efficiency, and style. Provide constructive feedback and explain *why* something is an issue. Do not just fix it."
	case "research":
		return "You are a Research Assistant. Focus on providing academic context, summarizing key concepts, and suggesting related topics or papers. Be formal and precise."
	case "exam":
		return "You are an Exam Strategist. Focus on test-taking strategies, time management, and prioritizing questions. Help the student prepare effectively for exams."
	case "motivational":
		return "You are a Motivational Coach. Be encouraging, positive, and supportive. Help the student set believable goals and overcome anxiety or procrastination."
	default: // "general" or unknown
		return "You are AcademAide, an academic advisor."
	}
}

func (s *RAGService) ProcessChat(studentID, message, agentID string) (string, error) {
	ctx := context.Background()

	// 1. Check Cache (Semantic/Exact)
	hash := sha256.Sum256([]byte(message))
	cacheKey := "response:" + hex.EncodeToString(hash[:])
	cached, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		return cached, nil
	}

	// 2. Fetch Context (Postgres)
	// Profile
	var studentName, deptID string
	err = config.PostgresDB.QueryRow("SELECT s_first_name, dept_id FROM STUDENT WHERE student_id=$1", studentID).Scan(&studentName, &deptID)
	if err != nil {
		return "", fmt.Errorf("fetching profile: %w", err)
	}

	// Schedule/Timetable (Reuse query logic or call internal function)
	// For brevity, fetching just raw schedule summary text or similar.
	// We'll run a quick join to get context for today
	// Schedule/Timetable - Fetch Full Week
	scheduleRows, err := config.PostgresDB.Query(`
		SELECT c.title, sch.day_of_week, sch.start_time, sch.room_number 
		FROM SCHEDULE sch 
		JOIN COURSE c ON sch.course_id = c.course_id 
		JOIN ENROLLS_IN e ON c.course_id = e.course_id
		WHERE e.student_id=$1
		ORDER BY sch.day_of_week, sch.start_time`, studentID)

	if err != nil {
		log.Println("Error fetching schedule for RAG:", err)
	}

	var scheduleContextBuilder strings.Builder
	scheduleContextBuilder.WriteString("Weekly Schedule: ")
	if scheduleRows != nil {
		defer scheduleRows.Close()
		for scheduleRows.Next() {
			var title, day, start, room string
			scheduleRows.Scan(&title, &day, &start, &room)
			scheduleContextBuilder.WriteString(fmt.Sprintf("[%s: %s at %s in %s] ", day, title, start, room))
		}
	}
	scheduleContext := scheduleContextBuilder.String()

	// 2.5 Vector Search Context
	vectorContext := ""
	embedding, err := s.Embedder.GenerateEmbedding(message)
	if err != nil {
		log.Println("Embedding generation failed:", err)
	} else {
		courses, err := s.Repo.SearchByVector(ctx, embedding, 3) // Top 3 relevant courses
		if err != nil {
			log.Println("Vector search failed:", err)
		} else {
			var sb strings.Builder
			sb.WriteString("Relevant Courses: ")
			for _, c := range courses {
				sb.WriteString(fmt.Sprintf("[%s: %s - %s] ", c.CourseID, c.Title, c.Description))
			}
			vectorContext = sb.String()
		}
	}

	// 3. Fetch Last 5 Messages (Mongo)
	coll := config.MongoDB.Collection("ChatLogs")
	opts := options.Find().SetSort(bson.D{{"timestamp", -1}}).SetLimit(5)
	cursor, err := coll.Find(ctx, bson.M{"student_id": studentID}, opts)

	var history []models.ChatLog
	if err == nil {
		cursor.All(ctx, &history)
	}
	// Reverse history for prompt
	var historyText strings.Builder
	for i := len(history) - 1; i >= 0; i-- {
		sender := "User"
		if history[i].IsBot {
			sender = "Bot"
		}
		historyText.WriteString(fmt.Sprintf("%s: %s\n", sender, history[i].Message))
	}

	// 4. Sentiment
	sentiment := s.AnalyzeSentiment(message)

	// 5. Construct Prompt
	systemPrompt := s.getAgentPrompt(agentID)
	prompt := fmt.Sprintf(`
System: %s
Context: User is %s from Dept %s. %s. %s.
Sentiment: User seems %s.
History:
%s
User: %s
Assistant:`, systemPrompt, studentName, deptID, scheduleContext, vectorContext, sentiment, historyText.String(), message)

	// 6. Call Ollama
	response, err := s.callOllama(prompt)
	if err != nil {
		return "", err
	}

	// 7. Store in Mongo
	// User Msg
	userLog := models.ChatLog{
		StudentID: studentID,
		Message:   message,
		Intent:    "chat",
		Sentiment: sentiment,
		Timestamp: time.Now(),
		IsBot:     false,
	}
	coll.InsertOne(ctx, userLog)

	// Bot Msg
	botLog := models.ChatLog{
		StudentID: studentID,
		Message:   response,
		Intent:    "reply",
		Timestamp: time.Now(),
		IsBot:     true,
	}
	coll.InsertOne(ctx, botLog)

	// Update Context (Simple upsert)
	contextColl := config.MongoDB.Collection("ChatContext")
	contextColl.UpdateOne(ctx, bson.M{"student_id": studentID}, bson.M{
		"$set": bson.M{
			"last_topic":       "general", // would need logic to extract topic
			"emotion":          sentiment,
			"last_interaction": time.Now(),
		},
	}, options.Update().SetUpsert(true))

	// 8. Cache Response
	config.RedisClient.Set(ctx, cacheKey, response, 5*time.Minute)

	return response, nil
}

func (s *RAGService) ClearChatHistory(studentID string) error {
	ctx := context.Background()

	// 1. Delete Chat Logs
	logsColl := config.MongoDB.Collection("ChatLogs")
	_, err := logsColl.DeleteMany(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return fmt.Errorf("failed to delete logs: %w", err)
	}

	// 2. Delete/Reset Context
	contextColl := config.MongoDB.Collection("ChatContext")
	_, err = contextColl.DeleteMany(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return fmt.Errorf("failed to delete context: %w", err)
	}

	return nil
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func (s *RAGService) callOllama(prompt string) (string, error) {
	url := "http://localhost:11434/api/generate"
	reqBody := OllamaRequest{
		Model:  "llama3.2",
		Prompt: prompt,
		Stream: false,
	}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		// Fallback for dev/demo if Ollama not running
		log.Println("Ollama Unreachable:", err)
		return "Simulated AI Response: " + prompt, nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var oResp OllamaResponse
	json.Unmarshal(body, &oResp)

	if oResp.Response == "" {
		return "I'm having trouble thinking right now.", nil
	}
	return oResp.Response, nil
}
