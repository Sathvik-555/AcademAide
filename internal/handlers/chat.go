package handlers

import (
	"academ_aide/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	StudentID string `json:"student_id"`
	FacultyID string `json:"faculty_id"`
	Role      string `json:"role"` // "student" or "teacher"
	Message   string `json:"message"`
	AgentID   string `json:"agent_id"` // Optional, defaults to "general" if empty
}

func ChatHandler(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	rag := services.NewRAGService()

	// Determine UserID and Role defaulting
	userID := req.StudentID
	role := "student"
	if req.Role == "teacher" {
		userID = req.FacultyID
		role = "teacher"
	}

	// If AgentID is missing, it will default to "general" in getAgentPrompt
	response, err := rag.ProcessChat(userID, role, req.Message, req.AgentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI Processing Failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": response,
		"user_id":  userID,
		"role":     role,
	})
}

func ClearChatHandler(c *gin.Context) {
	studentID := c.Query("student_id")
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_id is required"})
		return
	}

	rag := services.NewRAGService()
	if err := rag.ClearChatHistory(studentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat history cleared"})
}
