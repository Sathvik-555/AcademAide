package handlers

import (
	"academ_aide/internal/config"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TeacherHandler struct{}

func NewTeacherHandler() *TeacherHandler {
	return &TeacherHandler{}
}

// GetClassHealth godoc
// @Summary      Get Class Academic Health
// @Description  Returns attendance distribution and performance stats for a course
// @Tags         Teacher
// @Param        course_id query string true "Course ID"
// @Success      200 {object} map[string]interface{}
// @Router       /teacher/class-health [get]
func (h *TeacherHandler) GetClassHealth(c *gin.Context) {
	courseID := c.Query("course_id")
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "course_id is required"})
		return
	}

	// In a real app, query detailed attendance/grades.
	// For MVP without granular attendance table, we mock the distribution based on the course.
	// We verify the course exists first.
	var title string
	err := config.PostgresDB.QueryRow("SELECT title FROM COURSE WHERE course_id=$1", courseID).Scan(&title)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Error"})
		return
	}

	// Mock Data for Distribution
	healthData := gin.H{
		"course_id":   courseID,
		"title":       title,
		"attendance_distribution": gin.H{
			"90-100%": 15,
			"75-90%":  20,
			"60-75%":  5,
			"<60%":    2,
		},
		"performance_heatmap": gin.H{
			"Excellent (A/A+)": 10,
			"Good (B/B+)":      18,
			"Average (C/C+)":   10,
			"Poor (D/F)":       4,
		},
	}

	c.JSON(http.StatusOK, healthData)
}

// GetAtRiskStudents godoc
// @Summary      Get At-Risk Student Counts
// @Description  Returns aggregated counts of students at risk. Privacy safe.
// @Tags         Teacher
// @Param        course_id query string true "Course ID"
// @Success      200 {object} map[string]int
// @Router       /teacher/at-risk [get]
// Privacy note: This handler returns aggregated statistics only.
// Individual student identities are intentionally excluded.
func (h *TeacherHandler) GetAtRiskStudents(c *gin.Context) {
	courseID := c.Query("course_id")
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "course_id is required"})
		return
	}

	// Mock logic simulating risk analysis
	// In production, run complex queries on ENROLLS_IN and ATTENDANCE tables.
	response := gin.H{
		"high_risk_count":   3,  // < 60% attendance or failing grades
		"medium_risk_count": 8,  // 60-75% attendance or dropping grades
		"low_risk_count":    31, // Healthy
	}

	c.JSON(http.StatusOK, response)
}

// GetAlerts godoc
// @Summary      Get Early Warning Alerts
// @Description  Returns aggregated warnings for the teacher's courses
// @Tags         Teacher
// @Success      200 {object} []string
// @Router       /teacher/alerts [get]
func (h *TeacherHandler) GetAlerts(c *gin.Context) {
	userID := c.GetString("user_id") // Should be faculty_id

	// 1. Get Courses taught by this faculty
	rows, err := config.PostgresDB.Query("SELECT course_id FROM TEACHES WHERE faculty_id=$1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Error"})
		return
	}
	defer rows.Close()

	var alerts []string
	for rows.Next() {
		var cID string
		if err := rows.Scan(&cID); err == nil {
			// Generate mock alerts for each course
			alerts = append(alerts, "3 students in "+cID+" likely to fall below 75% attendance in 2 weeks")
		}
	}

	// Fallback if no courses or just to show something
	if len(alerts) == 0 {
		alerts = append(alerts, "No immediate alerts for your classes.")
	}

	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}
