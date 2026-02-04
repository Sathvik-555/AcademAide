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
	ID       string `json:"id"` // Unified ID field
	Password string `json:"password"`
	Role     string `json:"role"` // "student" or "teacher"
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Role != "student" && req.Role != "teacher" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Must be 'student' or 'teacher'"})
		return
	}

	var exists bool

	if req.Role == "student" {
		// Verify Student Exists
		err := config.PostgresDB.QueryRow("SELECT EXISTS(SELECT 1 FROM STUDENT WHERE student_id=$1)", req.ID).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Error"})
			return
		}
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Student not found"})
			return
		}
	} else if req.Role == "teacher" {
		// Verify Faculty Exists
		err := config.PostgresDB.QueryRow("SELECT EXISTS(SELECT 1 FROM FACULTY WHERE faculty_id=$1)", req.ID).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Error"})
			return
		}
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Faculty not found"})
			return
		}
		// No wallet logic for teachers
	}

	// Mock Password Verification (Common)
	if req.Password == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password required"})
		return
	}

	// Generate Real JWT
	claims := jwt.MapClaims{
		"user_id": req.ID,
		"role":    req.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	// Add wallet to claims only if student
	if req.Role == "student" {

		claims["student_id"] = req.ID // Backwards compatibility if needed
	} else {
		claims["faculty_id"] = req.ID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "supersecretkey"
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Store in Redis (Session)
	err = config.RedisClient.Set(context.Background(), "session:"+req.ID, tokenString, 24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis Error"})
		return
	}

	response := gin.H{
		"token":   tokenString,
		"user_id": req.ID,
		"role":    req.Role,
	}

	c.JSON(http.StatusOK, response)
}
