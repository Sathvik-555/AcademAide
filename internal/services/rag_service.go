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
	// Profile & Identity (Level 1)
	var studentName, deptID string
	var yearOfJoining int
	err = config.PostgresDB.QueryRow("SELECT s_first_name, dept_id, year_of_joining FROM STUDENT WHERE student_id=$1", studentID).Scan(&studentName, &deptID, &yearOfJoining)
	if err != nil {
		return "", fmt.Errorf("fetching profile: %w", err)
	}

	// Level 1: Calculate Year
	currentYearVal := time.Now().Year() // Dynamic based on system time
	studentYear := currentYearVal - yearOfJoining
	if studentYear <= 0 {
		studentYear = 1
	} // Fallback

	// Level 4: RAG Filtering - Fetch Enrolled Course IDs
	var enrolledCourseIDs []string
	var courseTitles []string
	cRows, err := config.PostgresDB.Query("SELECT c.course_id, c.title FROM ENROLLS_IN e JOIN COURSE c ON e.course_id = c.course_id WHERE e.student_id=$1 AND e.status='Enrolled'", studentID)
	if err == nil {
		defer cRows.Close()
		for cRows.Next() {
			var id, t string
			if err := cRows.Scan(&id, &t); err == nil {
				enrolledCourseIDs = append(enrolledCourseIDs, id)
				courseTitles = append(courseTitles, t)
			}
		}
	}
	coursesStr := strings.Join(courseTitles, ", ")

	// Level 2: Grade-Based Advisor - Academic History
	var academicHistoryBuilder strings.Builder
	var cgpa float64
	gRows, err := config.PostgresDB.Query(`
		SELECT c.title, e.grade, c.credits, e.course_id
		FROM ENROLLS_IN e
		JOIN COURSE c ON e.course_id = c.course_id
		WHERE e.student_id=$1 AND e.grade IS NOT NULL
	`, studentID)

	totalCredits := 0
	totalPoints := 0.0
	academicHistoryBuilder.WriteString("Academic History: [")

	if err == nil {
		defer gRows.Close()
		first := true
		for gRows.Next() {
			var title, grade, courseID string
			var credits int
			if err := gRows.Scan(&title, &grade, &credits, &courseID); err == nil {
				if !first {
					academicHistoryBuilder.WriteString(", ")
				}
				academicHistoryBuilder.WriteString(fmt.Sprintf("%s: %s", title, grade))
				first = false

				points := 0.0
				switch grade {
				case "O", "A+":
					points = 10.0
				case "A":
					points = 9.0
				case "B+":
					points = 8.0
				case "B":
					points = 7.0
				case "C+":
					points = 6.0
				case "C":
					points = 5.0
				case "D":
					points = 4.0
				default:
					points = 0.0
				}
				totalPoints += points * float64(credits)
				totalCredits += credits
			}
		}
	}
	academicHistoryBuilder.WriteString("]")
	if totalCredits > 0 {
		cgpa = float64(int((totalPoints/float64(totalCredits))*100)) / 100
	}
	academicHistoryStr := academicHistoryBuilder.String()

	// Level 3: Schedule Assistant - Next Class Logic
	// Logic: Find the FIRST class where (day = today AND start_time > now) OR (day > today)
	// For simplicity, we'll just check today's remaining classes first.
	currentTime := time.Now()
	// NOTE: User metadata says 2026-01-05T22:50:15+05:30 which is a Monday.
	// We need to parse correctly.
	dayOfWeek := currentTime.Weekday().String()      // "Monday"
	currentTimeStr := currentTime.Format("15:04:00") // "22:50:00"

	var upcomingClassInfo string

	// Query for today's remaining classes
	// Query for ONGOING class (where user should be RIGHT NOW)
	currentQuery := `
		SELECT c.title, sch.room_number, sch.end_time
		FROM SCHEDULE sch 
		JOIN ENROLLS_IN e ON sch.course_id = e.course_id
		JOIN COURSE c ON sch.course_id = c.course_id
		WHERE e.student_id=$1 AND sch.day_of_week=$2 AND sch.start_time <= $3 AND sch.end_time > $3
		LIMIT 1
	`
	var ongoingTitle, ongoingRoom, endTime string
	err = config.PostgresDB.QueryRow(currentQuery, studentID, dayOfWeek, currentTimeStr).Scan(&ongoingTitle, &ongoingRoom, &endTime)

	if err == nil {
		upcomingClassInfo = fmt.Sprintf("HAPPENING NOW: You should be in %s (Room %s) until %s.", ongoingTitle, ongoingRoom, endTime)
	} else {
		// If no ongoing class, find NEXT upcoming class today
		schedQuery := `
			SELECT c.title, sch.start_time, sch.room_number 
			FROM SCHEDULE sch 
			JOIN ENROLLS_IN e ON sch.course_id = e.course_id
			JOIN COURSE c ON sch.course_id = c.course_id
			WHERE e.student_id=$1 AND sch.day_of_week=$2 AND sch.start_time > $3
			ORDER BY sch.start_time ASC
			LIMIT 1
		`
		var nextTitle, nextStart, nextRoom string
		err = config.PostgresDB.QueryRow(schedQuery, studentID, dayOfWeek, currentTimeStr).Scan(&nextTitle, &nextStart, &nextRoom)

		if err == nil {
			upcomingClassInfo = fmt.Sprintf("Next Class: %s at %s in %s (Today).", nextTitle, nextStart, nextRoom)
		} else {
			upcomingClassInfo = "No more classes scheduled for today."
		}
	}

	currentDateTimeStr := fmt.Sprintf("%s, %s", dayOfWeek, currentTime.Format("15:04 PM")) // "Monday, 22:50 PM"

	// 2.5 Vector Search Context (Materials) with Level 4 Filtering
	vectorContext := ""
	embedding, err := s.Embedder.GenerateEmbedding(message)
	if err != nil {
		log.Println("Embedding generation failed:", err)
	} else {
		// Level 4: RAG Filtering - Pass enrolledCourseIDs
		materials, err := s.Repo.SearchMaterials(ctx, embedding, 3, enrolledCourseIDs, 0)
		if err != nil {
			log.Println("Material search failed:", err)
		} else {
			var sb strings.Builder
			sb.WriteString("Relevant Materials:\n")
			for _, m := range materials {
				sb.WriteString(fmt.Sprintf("- [Course %s Unit %d] %s\n", m.CourseID, m.UnitNo, m.Content))
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

	// Injecting all 4 Levels into Context
	prompt := fmt.Sprintf(`
System: %s
IMPORTANT: You have access to the user's real-time data below. Use it to answer questions about Schedule, Grades, and Profile DIRECTLY. Do not be vague.
Answer based primarily on the provided Context.
- If the user asks about a specific COURSE or SYLLABUS not in 'Enrolled Courses', politely inform them it's not in their enrollment.
- If the user asks about a GENERAL CONCEPT (e.g., 'pointers', 'algorithms'), you may answer using general knowledge if context is missing, but ADAPT COMPLEXITY to their Year:
- 1st/2nd Year: Use simple language and analogies.
- 3rd/4th Year: Use technical terms, depth, and industry context.
CGPA Context (10-point scale):
- 9.0+: Outstanding.
- 7.5-9.0: Good/Very Good.
- 6.0-7.5: Average.
- < 6.0: Low/Concerning. Be supportive but acknowledge they need to improve.
You are talking to %s, a %d Year %s student.
Context:
- Current Time: %s
- %s
- Enrolled Courses: [%s]
- CGPA: %.2f
- %s
- Materials: %s
- Sentiment: User seems %s.

History:
%s
User: %s
Assistant:`,
		systemPrompt,
		studentName, studentYear, deptID, // Level 1
		currentDateTimeStr, // Level 3
		upcomingClassInfo,  // Level 3
		coursesStr,
		cgpa,
		academicHistoryStr, // Level 2
		vectorContext,      // Level 4 (Filtered)
		sentiment,
		historyText.String(),
		message)

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
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
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
		Options: map[string]interface{}{
			"stop": []string{"User:", "System:", "Assistant:"},
		},
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
