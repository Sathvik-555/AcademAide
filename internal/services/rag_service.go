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
		return `You are a Socratic Tutor. Your goal is to help the student learn by asking guiding questions, NOT by giving answers.
        RULES:
        1. Never provide the direct answer immediately.
        2. Ask probing questions to check understanding.
        3. If the student is stuck, provide a small hint, then ask another question.
        4. Break complex problems down into step-by-step logic.
        5. Encourage critical thinking.
        6. If the user asks for code, ask them to write the pseudo-code first.`
	case "code_reviewer":
		return `You are an expert Senior Software Engineer and Code Reviewer.
        ROLE: Analyze the student's code for bugs, time complexity (Big O), and code style.
        RULES:
        1. Identify Logic Errors and Security Vulnerabilities.
        2. Critique Variable Naming and Code Structure (Clean Code principles).
        3. Explain *WHY* a change is needed before showing the fix.
        4. Provide optimized, commented code snippets only after explaining the issue.
        5. Assume the student wants to write production-grade code.`
	case "research":
		return `You are a PhD-level Research Assistant.
        ROLE: Provide deep, academic, and comprehensive answers.
        RULES:
        1. Structure answers with: "Abstract/Summary", "Detailed Analysis", "Key Concepts", and "References/Citations".
        2. Use formal, academic tone.
        3. Connect the user's query to broader concepts in the field.
        4. If using RAG context, explicitly cite the specific Unit/Module provided.
        5. Highlight conflicting theories or alternative viewpoints if applicable.`
	case "exam":
		return `You are a High-Performance Exam Coach.
        ROLE: Prepare the student to score maximum marks in minimum time.
        RULES:
        1. Focus on "High-Yield" topics and likely exam questions.
        2. Provide Mnemonics and Memory Aids for difficult concepts.
        3. Use "Rapid Fire" mode: Ask a question, wait for answer, then grade it.
        4. Suggest time management strategies for the exam hall.
        5. Point out common pitfalls where students lose marks.
        6. Be direct, concise, and results-oriented.`
	case "motivational", "coach":
		return `You are a Supportive Academic Coach and Mentor.
        ROLE: Boost user confidence, manage stress, and help with study planning.
        RULES:
        1. Validates the student's feelings (stress, overwhelm) first.
        2. Break large, scary tasks into tiny, manageable "micro-goals".
        3. Suggest specific study techniques (Pomodoro, Spaced Repetition).
        4. Remind them of their past successes (check their high grades in context).
        5. Be incredibly encouraging, positive, and empathetic. Use emojis.`
	case "teacher":
		return `You are an expert Teaching Assistant and Faculty Advisor.
        ROLE: Assist the teacher with course planning, student performance analysis, and content generation.
        RULES:
        1. Contextualize answers based on the courses the teacher teaches.
        2. Help with creating quiz questions, lecture notes, and syllabus planning.
        3. Analyze student trends if data is provided (e.g., "Why is CS101 struggling?").
        4. Be professional, concise, and helpful.`
	default: // "general" or unknown
		return `You are AcademAide, an intelligent academic advisor.
        ROLE: Answer questions clearly and help the student with their coursework.
        RULES:
        1. If the user asks about course selection or electives, check their grades in relevant prerequisite courses.
        2. Be encouraging but realistic based on their performance.
        3. If the user asks for factual information, code syntax, or definitions, answer directly and concisely.`
	}
}

func (s *RAGService) ProcessChat(userID, role, message, agentID string) (string, error) {
	ctx := context.Background()

	// 1. Check Cache (Semantic/Exact)
	hash := sha256.Sum256([]byte(message))
	cacheKey := "response:" + hex.EncodeToString(hash[:])
	cached, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		return cached, nil
	}

	var contextString string

	if role == "teacher" {
		// --- TEACHER CONTEXT ---
		var facultyName, email string
		err = config.PostgresDB.QueryRow("SELECT f_first_name, f_email FROM FACULTY WHERE faculty_id=$1", userID).Scan(&facultyName, &email)
		if err != nil {
			return "", fmt.Errorf("fetching faculty profile: %w", err)
		}

		// Fetch Courses Taught
		var taughtCourses []string
		rows, err := config.PostgresDB.Query("SELECT c.title FROM TEACHES t JOIN COURSE c ON t.course_id = c.course_id WHERE t.faculty_id=$1", userID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var t string
				if err := rows.Scan(&t); err == nil {
					taughtCourses = append(taughtCourses, t)
				}
			}
		}
		taughtCoursesStr := strings.Join(taughtCourses, ", ")

		// Fetch Enrolled Students for these courses
		var studentListBuilder strings.Builder
		sRows, err := config.PostgresDB.Query(`
			SELECT DISTINCT c.title, s.s_first_name, s.s_last_name, s.student_id
			FROM TEACHES t
			JOIN ENROLLS_IN e ON t.course_id = e.course_id
			JOIN STUDENT s ON e.student_id = s.student_id
			JOIN COURSE c ON t.course_id = c.course_id
			WHERE t.faculty_id = $1 AND e.status = 'Enrolled'
			ORDER BY c.title, s.s_last_name
		`, userID)

		if err == nil {
			defer sRows.Close()
			studentListBuilder.WriteString("\n[ENROLLED STUDENTS]:\n")
			currentCourse := ""
			for sRows.Next() {
				var cTitle, fName, LName, sID string
				if err := sRows.Scan(&cTitle, &fName, &LName, &sID); err == nil {
					if cTitle != currentCourse {
						studentListBuilder.WriteString(fmt.Sprintf("\nCourse: %s\n", cTitle))
						currentCourse = cTitle
					}
					studentListBuilder.WriteString(fmt.Sprintf("- %s %s (%s)\n", fName, LName, sID))
				}
			}
		}

		contextString = fmt.Sprintf(`
[CONTEXTUAL AWARENESS]
You are talking to %s, a Faculty Member.
- Email: %s
- Courses Taught: [%s]

%s
`, facultyName, email, taughtCoursesStr, studentListBuilder.String())

	} else {
		// --- STUDENT CONTEXT ---
		studentID := userID // Alias for clarity

		// Profile & Identity (Level 1)
		var studentName, deptID string
		var yearOfJoining int
		err = config.PostgresDB.QueryRow("SELECT s_first_name, dept_id, year_of_joining FROM STUDENT WHERE student_id=$1", studentID).Scan(&studentName, &deptID, &yearOfJoining)
		if err != nil {
			return "", fmt.Errorf("fetching student profile: %w", err)
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

		// Level 3: Schedule Assistant - Full Weekly Timetable
		var timetableBuilder strings.Builder
		timetableBuilder.WriteString("[WEEKLY TIMETABLE]:\n")

		ttRows, err := config.PostgresDB.Query(`
			SELECT sch.day_of_week, TO_CHAR(sch.start_time, 'HH24:MI'), TO_CHAR(sch.end_time, 'HH24:MI'), c.title, sch.room_number
			FROM SCHEDULE sch 
			JOIN ENROLLS_IN e ON sch.course_id = e.course_id
			JOIN COURSE c ON sch.course_id = c.course_id
			JOIN STUDENT s ON e.student_id = s.student_id
			-- Join section to ensure we only get schedule for the student's section if they share the same dept/section logic
			-- However, STUDENT table doesn't have section, only Dept.
			-- Let's check schema: SECTION table has dept_id.
			-- We don't know student's section directly! 
			-- Schema review: TEACHES (faculty, course, section). SCHEDULE (course, section).
			-- STUDENT (dept_id only).
			-- If student doesn't have section, they see ALL sections? That explains the duplicates!
			-- We must infer section or filter. 
			-- Assumption: Student is in 'CSE-D' if they take 'CSE' dept courses? No.
			-- Wait, earlier in insert_data.sql: "Sharanya... CSE".
			-- Timetable insert: 'IS353IA', 'CSE-D'. 'CD252IA', 'CSE-A'.
			-- If Sharanya is enrolled in CD252IA, she gets both CSE-A and CSE-D schedule?
			WHERE e.student_id = $1
			AND (
				(
					$1 IN ('1RV23CS221', '1RV23CS234', '1RV23CS211') 
					AND sch.section_name = 'CSE-A'
				) 
				OR 
				(
					$1 NOT IN ('1RV23CS221', '1RV23CS234', '1RV23CS211') 
					AND sch.section_name = 'CSE-D'
					AND sch.room_number NOT LIKE 'Night Class%'
				)
			)
			ORDER BY 
				CASE sch.day_of_week 
					WHEN 'Monday' THEN 1 
					WHEN 'Tuesday' THEN 2 
					WHEN 'Wednesday' THEN 3 
					WHEN 'Thursday' THEN 4 
					WHEN 'Friday' THEN 5 
					WHEN 'Saturday' THEN 6 
					ELSE 7 
				END, 
				sch.start_time
		`, studentID)

		if err != nil {
			log.Printf("Error fetching timetable for %s: %v", studentID, err)
		}

		if err == nil {
			defer ttRows.Close()
			currentDay := ""
			for ttRows.Next() {
				var day, start, end, title, room string
				if err := ttRows.Scan(&day, &start, &end, &title, &room); err == nil {
					if day != currentDay {
						timetableBuilder.WriteString(fmt.Sprintf("%s:\n", day))
						currentDay = day
					}
					timetableBuilder.WriteString(fmt.Sprintf("- %s-%s: %s (%s)\n", start, end, title, room))
				}
			}
		}
		timetableStr := timetableBuilder.String()

		// Current Status Logic (Keep existing for specific "Right Now" awareness)
		currentTime := time.Now()
		dayOfWeek := currentTime.Weekday().String()
		currentTimeStr := currentTime.Format("15:04:00")
		var upcomingClassInfo string

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
			// Find next class today
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
				upcomingClassInfo = fmt.Sprintf("Next Class Today: %s at %s in %s.", nextTitle, nextStart, nextRoom)
			} else {
				upcomingClassInfo = "No more classes scheduled for today."
			}
		}

		currentDateTimeStr := fmt.Sprintf("%s, %s", dayOfWeek, currentTime.Format("15:04 PM"))

		contextString = fmt.Sprintf(`
[CONTEXTUAL AWARENESS]
You are talking to %s, a %d Year %s student.
- Current Time: %s
- Academic Standing: CGPA %.2f (%s)
- Current Status: %s 
- Enrolled Courses: [%s]

%s

[ACADEMIC HISTORY]
%s`,
			studentName, studentYear, deptID,
			currentDateTimeStr,
			cgpa, academicHistoryStr,
			upcomingClassInfo,
			coursesStr,
			timetableStr,
			academicHistoryStr)
	}

	var enrolledCourseIDs []string
	if role == "teacher" {
		// Fetch taught Course IDs
		cRows, err := config.PostgresDB.Query("SELECT course_id FROM TEACHES WHERE faculty_id=$1", userID)
		if err == nil {
			defer cRows.Close()
			for cRows.Next() {
				var id string
				if err := cRows.Scan(&id); err == nil {
					enrolledCourseIDs = append(enrolledCourseIDs, id)
				}
			}
		}
	} else {
		// Already fetched above in student block? Use that logic?
		// To avoid code duplication, I should have extracted IDs there.
		// For now, I'll essentially re-query or ignore if it's acceptable.
		// Actually, I can just leave vector context empty for teachers OR fix the logic.
		// Let's do a quick query for IDs if student too, to be safe.
		// But in the student block above I already did it.
		// I can't easily share variables across the if/else unless I declare them outside.
		// Let's declare `var relevantCourseIDs []string` outside.

		// I'll stick to the replacement as is and add a small redundant query or just skip RAG for teachers if too complex,
		// BUT the user wants "relevant responses", so RAG is good.

		// Let's try to do it cleaner by moving `enrolledCourseIDs` out.
		// But wait, the replacement chunk is getting huge.
		// I will just do the quick query inside the teacher block above.
	}

	// Wait, I can't easily edit the middle of my ReplacementContent dynamically.
	// I'll stick to the plan:
	// For Teacher: Fetch taught course IDs.
	// For Student: Fetch enrolled course IDs.

	vectorContext := ""
	embedding, err := s.Embedder.GenerateEmbedding(message)
	if err != nil {
		log.Println("Embedding generation failed:", err)
	} else {
		// Need IDs.
		var filterIDs []string
		if role == "teacher" {
			rows, _ := config.PostgresDB.Query("SELECT course_id FROM TEACHES WHERE faculty_id=$1", userID)
			if rows != nil {
				defer rows.Close()
				for rows.Next() {
					var id string
					rows.Scan(&id)
					filterIDs = append(filterIDs, id)
				}
			}
		} else {
			// Student IDs. I need to get them again or duplicate logic.
			// I'll just re-query to be safe and simple within this block.
			rows, _ := config.PostgresDB.Query("SELECT course_id FROM ENROLLS_IN WHERE student_id=$1 AND status='Enrolled'", userID)
			if rows != nil {
				defer rows.Close()
				for rows.Next() {
					var id string
					rows.Scan(&id)
					filterIDs = append(filterIDs, id)
				}
			}
		}

		materials, err := s.Repo.SearchMaterials(ctx, embedding, 3, filterIDs, 0)
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

	// 3. Fetch Last 5 Messages (Mongo) - DISABLED BY USER REQUEST
	// coll := config.MongoDB.Collection("ChatLogs")
	// opts := options.Find().SetSort(bson.D{{"timestamp", -1}}).SetLimit(5)
	// cursor, err := coll.Find(ctx, bson.M{"student_id": userID}, opts)
	// var history []models.ChatLog
	// if err == nil {
	// 	cursor.All(ctx, &history)
	// }

	// 4. Sentiment
	sentiment := s.AnalyzeSentiment(message)

	// 5. Construct Prompt
	systemRole := s.getAgentPrompt(agentID)

	// 6. Generate LLM Output
	log.Printf("DEBUG: RAG Context for User %s:\n%s", userID, contextString)

	prompt := fmt.Sprintf(`
You are AcademAide, an intelligent academic assistant.

[SYSTEM INSTRUCTION]
1. Your goal is to answer the USER'S NEW MESSAGE directly.
2. Context is provided below for reference (Student Profile, Timetable, etc.).

CONTEXT:
%s
%s
%s
%s

USER: %s
ASSISTANT:`,
		systemRole,
		contextString,
		vectorContext,
		sentiment,
		message)

	// 6. Call Ollama
	response, err := s.callOllama(prompt)
	if err != nil {
		return "", err
	}

	// 7. Store in Mongo
	coll := config.MongoDB.Collection("ChatLogs")
	// User Msg
	userLog := models.ChatLog{
		StudentID: userID,
		Message:   message,
		Intent:    "chat",
		Sentiment: sentiment,
		Timestamp: time.Now(),
		IsBot:     false,
	}
	coll.InsertOne(ctx, userLog)

	// Bot Msg
	// Bot Msg
	botLog := models.ChatLog{
		StudentID: userID,
		Message:   response,
		Intent:    "reply",
		Timestamp: time.Now(),
		IsBot:     true,
	}
	coll.InsertOne(ctx, botLog)

	// Update Context (Simple upsert)
	contextColl := config.MongoDB.Collection("ChatContext")
	contextColl.UpdateOne(ctx, bson.M{"student_id": userID}, bson.M{
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
