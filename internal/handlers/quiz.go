package handlers

import (
	"academ_aide/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GenerateQuizRequest struct {
	CourseID     string `json:"course_id" binding:"required"`
	Unit         int    `json:"unit"`          // Optional
	NumQuestions int    `json:"num_questions"` // Optional, Default 5
}

func GenerateQuiz(c *gin.Context) {
	var req GenerateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Default to 5 questions if not specified
	if req.NumQuestions <= 0 {
		req.NumQuestions = 5
	}

	svc := services.NewQuizService()
	quiz, err := svc.GenerateQuiz(req.CourseID, req.Unit, req.NumQuestions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quiz)
}
