package services

import (
	"academ_aide/internal/config"
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

// Data Structures

type StudentRisk struct {
	Type     string `json:"type"`     // "Attendance", "Grades"
	Severity string `json:"severity"` // "High", "Medium", "Low"
	Message  string `json:"message"`
	Subject  string `json:"subject,omitempty"`
}

type Suggestion struct {
	Suggestion string `json:"suggestion"`
	Reason     string `json:"reason"`
}

type AIInsightsResponse struct {
	Risks       []StudentRisk `json:"risks"`
	Suggestions []Suggestion  `json:"suggestions"`
}

type WhatIfScenario struct {
	PercentageDrop      float64 `json:"percentage_drop"`
	ProjectedAttendance float64 `json:"projected_attendance"`
	RiskLevel           string  `json:"risk_level"`
}

// Service Interface regarding AI capabilities
type AIService struct {
	db *sql.DB
}

func NewAIService() *AIService {
	return &AIService{
		db: config.PostgresDB,
	}
}

// Constants for Heuristic Rules
const (
	AttendanceThresholdCritical = 75.0
	AttendanceThresholdWarning  = 85.0 // Stricter for demo
	GradeDropThreshold          = 5.0
)

// GetStudentInsights returns risks and suggestions for a student
func (s *AIService) GetStudentInsights(studentID string) (*AIInsightsResponse, error) {
	var risks []StudentRisk
	var suggestions []Suggestion

	// 1. Fetch Enrolled Courses & Grades from DB
	rows, err := s.db.Query(`
		SELECT c.title, e.grade, e.status, e.course_id
		FROM ENROLLS_IN e
		JOIN COURSE c ON e.course_id = c.course_id
		WHERE e.student_id = $1
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var title, gradeStr, status, courseID string
		var grade sql.NullString
		if err := rows.Scan(&title, &grade, &status, &courseID); err != nil {
			continue
		}
		if grade.Valid {
			gradeStr = grade.String
		}

		// --- A. Grade Analysis (REAL) ---
		// Map grade to risk
		if status == "Enrolled" || status == "Completed" {
			switch gradeStr {
			case "F", "D", "E":
				risks = append(risks, StudentRisk{
					Type:     "Grades",
					Severity: "High",
					Message:  fmt.Sprintf("Critical performance (Grade: %s) in %s", gradeStr, title),
					Subject:  title,
				})
				suggestions = append(suggestions, Suggestion{
					Suggestion: fmt.Sprintf("Schedule remedial session for %s", title),
					Reason:     "Current grade puts you at risk of failing or academic probation.",
				})
			case "C", "C+":
				risks = append(risks, StudentRisk{
					Type:     "Grades",
					Severity: "Medium",
					Message:  fmt.Sprintf("Average performance (Grade: %s) in %s", gradeStr, title),
					Subject:  title,
				})
				suggestions = append(suggestions, Suggestion{
					Suggestion: fmt.Sprintf("Review core concepts in %s", title),
					Reason:     "Grade is average; improving understanding now can boost final score.",
				})
			}
		}

		// --- B. Attendance Analysis (SIMULATED / DETERMINISTIC) ---
		// We simulate attendance based on a hash of (StudentID + CourseID) so it's consistent for the user but varies by course.
		// Hash -> 0-100
		h := sha256.New()
		h.Write([]byte(studentID + courseID + "attendance_salt"))
		sum := h.Sum(nil)
		// Use first byte to determine rough percentage (50-100 range likely)
		// To make it interesting:
		// S1001 (Senior) -> usually good.
		// Freshman -> maybe mixed.
		// Let's just use the hash modulo 101.
		// To ensure we have some "Good" and some "Bad", we map 0-255 to 40-100.
		// val = 40 + (byte / 255) * 60
		val := int(sum[0])
		attendancePct := 60.0 + (float64(val)/255.0)*40.0 // Range 60% to 100%

		// Override for specific Test Users to ensure we see specific UI states
		if strings.Contains(studentID, "TEST_STRESSED") && (courseID == "CS102" || courseID == "CS103") {
			attendancePct = 65.0 // Bad attendance
		}
		if strings.Contains(studentID, "TEST_FRESHMAN") && courseID == "CS101" {
			attendancePct = 78.0 // Warning
		}
		if strings.Contains(studentID, "TEST_SENIOR") {
			attendancePct = 92.0 + (float64(val)/255.0)*6.0 // 92-98%
		}

		if attendancePct < AttendanceThresholdCritical {
			risks = append(risks, StudentRisk{
				Type:     "Attendance",
				Severity: "High",
				Message:  fmt.Sprintf("Attendance CRITICAL (%.1f%%) in %s", attendancePct, title),
				Subject:  title,
			})
			suggestions = append(suggestions, Suggestion{
				Suggestion: fmt.Sprintf("Meet %s instructor immediately", title),
				Reason:     "Attendance is below 75%; you may be debarred from exams.",
			})

		} else if attendancePct < AttendanceThresholdWarning {
			risks = append(risks, StudentRisk{
				Type:     "Attendance",
				Severity: "Medium",
				Message:  fmt.Sprintf("Low Attendance (%.1f%%) in %s", attendancePct, title),
				Subject:  title,
			})
			suggestions = append(suggestions, Suggestion{
				Suggestion: fmt.Sprintf("Attend next all classes for %s", title),
				Reason:     "Attendance is borderline. Missing more classes will trigger critical status.",
			})
		}
	}

	if len(risks) == 0 {
		suggestions = append(suggestions, Suggestion{
			Suggestion: "Maintain current study schedule",
			Reason:     "Your academic health is green! All grades and attendance metrics are satisfactory.",
		})
	}

	return &AIInsightsResponse{
		Risks:       risks,
		Suggestions: suggestions,
	}, nil
}

// CalculateWhatIf simulates attendance scenarios
func (s *AIService) CalculateWhatIf(studentID string, missedClasses int) (*WhatIfScenario, error) {
	// Mock Base State (Same as frontend default)
	currentTotal := 40.0
	currentAttended := 34.0 // 85% initially

	currentPct := (currentAttended / currentTotal) * 100

	// Simulation
	newTotal := currentTotal + float64(missedClasses)
	// newAttended remains same as we are missing classes
	newPct := (currentAttended / newTotal) * 100

	risk := "Low"
	if newPct < AttendanceThresholdCritical {
		risk = "High"
	} else if newPct < AttendanceThresholdWarning {
		risk = "Medium"
	}

	return &WhatIfScenario{
		ProjectedAttendance: newPct,
		PercentageDrop:      currentPct - newPct,
		RiskLevel:           risk,
	}, nil
}

// Quiz Analysis Structures
type QuizSubmission struct {
	CourseID       string              `json:"course_id"`
	WrongQuestions []QuizWrongQuestion `json:"wrong_questions"`
	TotalQuestions int                 `json:"total_questions"`
	Score          int                 `json:"score"`
}

type QuizWrongQuestion struct {
	QuestionText  string `json:"question_text"`
	CorrectAnswer string `json:"correct_answer"`
	UserAnswer    string `json:"user_answer"`
	Reference     string `json:"reference"`
}

type QuizAnalysisResponse struct {
	WeakAreas       []string            `json:"weak_areas"`
	StudyPriorities []StudyPriorityItem `json:"study_priorities"`
}

type StudyPriorityItem struct {
	Topic    string `json:"topic"`
	Priority string `json:"priority"` // High, Medium, Low
	Reason   string `json:"reason"`
}

// AnalyzeQuizPerformance generates insights based on quiz results
func (s *AIService) AnalyzeQuizPerformance(sub QuizSubmission) (*QuizAnalysisResponse, error) {
	if len(sub.WrongQuestions) == 0 {
		return &QuizAnalysisResponse{
			WeakAreas: []string{},
			StudyPriorities: []StudyPriorityItem{
				{Topic: "General Review", Priority: "Low", Reason: "Perfection! Just review the course summary."},
			},
		}, nil
	}

	// Build Context for AI
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Course: %s\n", sub.CourseID))
	sb.WriteString("Incorrectly Answered Questions:\n")
	for i, q := range sub.WrongQuestions {
		sb.WriteString(fmt.Sprintf("%d. Question: %s\n   Correct: %s\n   User Answered: %s\n   Source: %s\n",
			i+1, q.QuestionText, q.CorrectAnswer, q.UserAnswer, q.Reference))
	}

	prompt := fmt.Sprintf(`
You are an academic tutor. A student just took a quiz for course %s and got the following questions wrong.
Analyze their mistakes to identify specific weak topics.
Give a prioritized study list.

%s

Return ONLY valid JSON in this format:
{
  "weak_areas": ["Topic A", "Topic B"],
  "study_priorities": [
    { "topic": "Specific Concept", "priority": "High", "reason": "Explanation based on mistakes" },
    { "topic": "Another Concept", "priority": "Medium", "reason": "..." }
  ]
}
`, sub.CourseID, sb.String())

	// Call Ollama
	jsonResp, err := s.callOllamaJSON(prompt)
	if err != nil {
		return nil, err
	}

	var analysis QuizAnalysisResponse
	if err := json.Unmarshal([]byte(jsonResp), &analysis); err != nil {
		// Try fallback parsing if markdown code blocks exist
		clean := strings.ReplaceAll(jsonResp, "```json", "")
		clean = strings.ReplaceAll(clean, "```", "")
		if err2 := json.Unmarshal([]byte(clean), &analysis); err2 != nil {
			return nil, fmt.Errorf("failed to parse AI analysis: %v", err)
		}
	}

	// Sort by Priority: High > Medium > Low
	priorityMap := map[string]int{
		"High":   3,
		"Medium": 2,
		"Low":    1,
	}

	sort.Slice(analysis.StudyPriorities, func(i, j int) bool {
		p1 := priorityMap[analysis.StudyPriorities[i].Priority]
		p2 := priorityMap[analysis.StudyPriorities[j].Priority]
		return p1 > p2 // Descending order
	})

	return &analysis, nil
}

// Helper: Call Ollama (Duplicated to avoid circular deps or refactoring overhead for now)
func (s *AIService) callOllamaJSON(prompt string) (string, error) {
	url := "http://localhost:11434/api/generate"
	reqBody := map[string]interface{}{
		"model":  "llama3.2",
		"prompt": prompt,
		"stream": false,
		"format": "json",
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
	// Simple unmarshal
	json.Unmarshal(body, &oResp)

	return oResp.Response, nil
}
