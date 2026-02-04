package handlers

import (
	"academ_aide/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	aiService *services.AIService
}

func NewAIHandler() *AIHandler {
	return &AIHandler{
		aiService: services.NewAIService(),
	}
}

// GetInsights godoc
// @Summary      Get AI Insights for a student
// @Description  Returns risk analysis and suggestions
// @Tags         AI
// @Param        student_id query string true "Student ID"
// @Success      200 {object} services.AIInsightsResponse
// @Router       /ai/insights [get]
func (h *AIHandler) GetInsights(c *gin.Context) {
	studentID := c.Query("student_id")
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_id is required"})
		return
	}

	insights, err := h.aiService.GetStudentInsights(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze data"})
		return
	}

	c.JSON(http.StatusOK, insights)
}

type WhatIfRequest struct {
	StudentID     string `json:"student_id"`
	MissedClasses int    `json:"missed_classes"`
}

func (h *AIHandler) CalculateWhatIf(c *gin.Context) {
	var req WhatIfRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result, err := h.aiService.CalculateWhatIf(req.StudentID, req.MissedClasses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Simulation failed"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// AnalyzeQuiz godoc
// @Summary      Analyze Quiz Performance
// @Description  Returns weak areas and study priorities based on quiz results
// @Tags         AI
// @Accept       json
// @Produce      json
// @Param        request body services.QuizSubmission true "Quiz Results"
// @Success      200 {object} services.QuizAnalysisResponse
// @Router       /ai/quiz-analysis [post]
func (h *AIHandler) AnalyzeQuiz(c *gin.Context) {
	var req services.QuizSubmission
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	analysis, err := h.aiService.AnalyzeQuizPerformance(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Analysis failed"})
		return
	}

	c.JSON(http.StatusOK, analysis)
}
