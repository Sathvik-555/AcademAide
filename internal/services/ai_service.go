package services

import (
	"bytes"
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
	// In the future, this would hold references to repositories
}

func NewAIService() *AIService {
	return &AIService{}
}

// Constants for Heuristic Rules
const (
	AttendanceThresholdCritical = 75.0
	AttendanceThresholdWarning  = 80.0
	GradeDropThreshold          = 5.0
)

// GetStudentInsights returns risks and suggestions for a student
func (s *AIService) GetStudentInsights(studentID string) (*AIInsightsResponse, error) {
	// MOCK DATA FETCHING
	// In a real implementation, we would call s.repo.GetAttendance(studentID) etc.
	// For MVP, we reconstruct the mock data used in the frontend prototype.

	attendanceData := []struct {
		Subject  string
		Attended float64
		Total    float64
	}{
		{"DBMS", 28, 40}, // 70%
		{"OS", 35, 40},   // 87.5%
		{"DAA", 32, 40},  // 80%
	}

	gradesData := []struct {
		Subject  string
		Current  float64 // Percentage equivalent
		Previous float64
	}{
		{"DBMS", 65, 80}, // Big drop
		{"OS", 85, 82},
	}

	// 1. Analyze Risks
	var risks []StudentRisk

	// Check Attendance
	for _, sub := range attendanceData {
		percentage := (sub.Attended / sub.Total) * 100
		if percentage < AttendanceThresholdCritical {
			risks = append(risks, StudentRisk{
				Type:     "Attendance",
				Severity: "High",
				Message:  fmt.Sprintf("Attendance below 75%% in %s", sub.Subject),
				Subject:  sub.Subject,
			})
		} else if percentage < AttendanceThresholdWarning {
			risks = append(risks, StudentRisk{
				Type:     "Attendance",
				Severity: "Medium",
				Message:  fmt.Sprintf("Attendance risk in %s (Current: %.1f%%)", sub.Subject, percentage),
				Subject:  sub.Subject,
			})
		}
	}

	// Check Grades
	for _, sub := range gradesData {
		if sub.Previous-sub.Current > GradeDropThreshold {
			risks = append(risks, StudentRisk{
				Type:     "Grades",
				Severity: "Medium",
				Message:  fmt.Sprintf("Performance drop detected in %s", sub.Subject),
				Subject:  sub.Subject,
			})
		}
	}

	// 2. Generate Suggestions
	var suggestions []Suggestion
	for _, r := range risks {
		if r.Type == "Attendance" {
			if r.Severity == "High" {
				suggestions = append(suggestions, Suggestion{
					Suggestion: fmt.Sprintf("Consider meeting the course instructor for %s", r.Subject),
					Reason:     fmt.Sprintf("Attendance in %s is critical (<75%%). Immediate action prevents debarment.", r.Subject),
				})
			} else {
				suggestions = append(suggestions, Suggestion{
					Suggestion: fmt.Sprintf("Attend next 3 classes for %s", r.Subject),
					Reason:     "Current attendance trend shows a potential drop below safe limits in upcoming weeks.",
				})
			}
		} else if r.Type == "Grades" {
			suggestions = append(suggestions, Suggestion{
				Suggestion: fmt.Sprintf("Review recent topics in %s with a study group", r.Subject),
				Reason:     fmt.Sprintf("Recent scores indicate a %s drop in performance compared to previous assessments.", r.Severity), // Severity is "Medium" usually
			})
		}
	}

	if len(suggestions) == 0 {
		suggestions = append(suggestions, Suggestion{
			Suggestion: "Maintain current study schedule",
			Reason:     "All metrics are within healthy ranges. Consistency is key.",
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
