package handlers

import (
	"academ_aide/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	StudentID string `json:"student_id"`
	Message   string `json:"message"`
}

func ChatHandler(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	rag := services.NewRAGService()
	response, err := rag.ProcessChat(req.StudentID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI Processing Failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": response,
		"student_id": req.StudentID,
	})
}
