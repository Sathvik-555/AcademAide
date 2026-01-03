package handlers

import (
	"academ_aide/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GenerateQuizRequest struct {
	CourseID string `json:"course_id" binding:"required"`
}

func GenerateQuiz(c *gin.Context) {
	var req GenerateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	svc := services.NewQuizService()
	quiz, err := svc.GenerateQuiz(req.CourseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quiz)
}
