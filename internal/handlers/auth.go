package handlers

import (
	"academ_aide/internal/config"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Simple Login - In production use proper Password Hashing
type LoginRequest struct {
	StudentID string `json:"student_id"`
	Password  string `json:"password"`
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Verify Student Exists in Postgres
	var exists bool
	err := config.PostgresDB.QueryRow("SELECT EXISTS(SELECT 1 FROM STUDENT WHERE student_id=$1)", req.StudentID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Error"})
		return
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Student not found"})
		return
	}

	// TODO: Verify Password (req.Password) here against DB hash
	// For now, we assume if student exists, and password is provided (not empty), it's valid for this lab.
	if req.Password == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password required"})
		return
	}

	// Generate Real JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"student_id": req.StudentID,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "supersecretkey"
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Store in Redis (Session tracking can still be useful, or just rely on JWT statelessness)
	// We'll keep Redis for "last login" or similar if needed, or just strict session management.
	// The original code stored "token" in Redis. We can still do that.
	err = config.RedisClient.Set(context.Background(), "session:"+req.StudentID, tokenString, 24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "student_id": req.StudentID})
}
