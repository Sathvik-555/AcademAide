package handlers

import (
	"academ_aide/internal/config"
	"academ_aide/internal/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FindPeersRequest - GET param ?course_id=...

func FindPeers(c *gin.Context) {
	currentStudentID, _ := c.Get("student_id")
	courseID := c.Query("course_id")

	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "course_id is required"})
		return
	}

	// Find students enrolled in the same course
	rows, err := config.PostgresDB.Query(`
		SELECT s.student_id, s.s_first_name, s.s_last_name, s.s_email 
		FROM ENROLLS_IN e
		JOIN STUDENT s ON e.student_id = s.student_id
		WHERE e.course_id = $1 AND s.student_id != $2
	`, courseID, currentStudentID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch peers"})
		return
	}
	defer rows.Close()

	var peers []models.Student
	for rows.Next() {
		var p models.Student
		if err := rows.Scan(&p.StudentID, &p.FirstName, &p.LastName, &p.Email); err == nil {
			peers = append(peers, p)
		}
	}

	c.JSON(http.StatusOK, peers)
}

type CreateGroupRequest struct {
	CourseID    string `json:"course_id"`
	GroupName   string `json:"group_name"`
	Description string `json:"description"`
}

func CreateGroup(c *gin.Context) {
	studentIDVal, _ := c.Get("student_id")
	studentID := studentIDVal.(string)

	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	tx, err := config.PostgresDB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Error"})
		return
	}

	// 1. Create Group
	var groupID int
	err = tx.QueryRow(`
		INSERT INTO STUDY_GROUP (course_id, group_name, description, created_by)
		VALUES ($1, $2, $3, $4) RETURNING group_id
	`, req.CourseID, req.GroupName, req.Description, studentID).Scan(&groupID)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group: " + err.Error()})
		return
	}

	// 2. Add Creator as Member
	_, err = tx.Exec(`
		INSERT INTO GROUP_MEMBER (group_id, student_id) VALUES ($1, $2)
	`, groupID, studentID)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"group_id": groupID, "message": "Study Group Created"})
}

func ListGroups(c *gin.Context) {
	courseID := c.Query("course_id")

	query := `
		SELECT g.group_id, g.course_id, g.group_name, g.description, g.created_by, g.created_at,
		(SELECT COUNT(*) FROM GROUP_MEMBER gm WHERE gm.group_id = g.group_id) as member_count
		FROM STUDY_GROUP g
	`
	var rows *sql.Rows
	var err error

	if courseID != "" {
		query += " WHERE g.course_id = $1"
		rows, err = config.PostgresDB.Query(query, courseID)
	} else {
		rows, err = config.PostgresDB.Query(query)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var groups []models.StudyGroup
	for rows.Next() {
		var g models.StudyGroup
		if err := rows.Scan(&g.GroupID, &g.CourseID, &g.GroupName, &g.Description, &g.CreatedBy, &g.CreatedAt, &g.MemberCount); err == nil {
			groups = append(groups, g)
		}
	}
	c.JSON(http.StatusOK, groups)
}
