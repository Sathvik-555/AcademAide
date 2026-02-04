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

// GetMyCourses godoc
// @Summary      Get Faculty Courses
// @Description  Returns list of courses taught by the logged-in faculty
// @Tags         Teacher
// @Router       /teacher/courses [get]
func (h *TeacherHandler) GetMyCourses(c *gin.Context) {
	facultyID := c.GetString("user_id") // Authenticated Faculty ID
	rows, err := config.PostgresDB.Query(`
		SELECT c.course_id, c.title, t.section_name 
		FROM TEACHES t
		JOIN COURSE c ON t.course_id = c.course_id
		WHERE t.faculty_id=$1
	`, facultyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Error"})
		return
	}
	defer rows.Close()

	var courses []gin.H
	for rows.Next() {
		var id, title, section string
		if err := rows.Scan(&id, &title, &section); err == nil {
			courses = append(courses, gin.H{"course_id": id, "title": title, "section": section})
		}
	}
	c.JSON(http.StatusOK, courses)
}

// GetEnrolledStudents godoc
// @Summary      Get Students List
// @Description  Returns enrolled students for a course
// @Tags         Teacher
// @Param        course_id query string true "Course ID"
// @Router       /teacher/students [get]
func (h *TeacherHandler) GetEnrolledStudents(c *gin.Context) {
	courseID := c.Query("course_id")
	rows, err := config.PostgresDB.Query(`
		SELECT s.student_id, s.s_first_name, s.s_last_name, s.s_email
		FROM ENROLLS_IN e
		JOIN STUDENT s ON e.student_id = s.student_id
		WHERE e.course_id=$1
	`, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Error"})
		return
	}
	defer rows.Close()

	var students []gin.H
	for rows.Next() {
		var id, fname, lname, email string
		if err := rows.Scan(&id, &fname, &lname, &email); err == nil {
			students = append(students, gin.H{
				"student_id": id,
				"name":       fname + " " + lname,
				"email":      email,
			})
		}
	}
	c.JSON(http.StatusOK, students)
}

func (h *TeacherHandler) GetClassHealth(c *gin.Context) {
	courseID := c.Query("course_id")
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "course_id is required"})
		return
	}

	var title string
	err := config.PostgresDB.QueryRow("SELECT title FROM COURSE WHERE course_id=$1", courseID).Scan(&title)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Error"})
		return
	}

	// Dynamic Grades
	rows, err := config.PostgresDB.Query("SELECT grade FROM ENROLLS_IN WHERE course_id=$1 AND grade IS NOT NULL", courseID)
	perfMap := map[string]int{"Excellent (A/A+)": 0, "Good (B/B+)": 0, "Average (C/C+)": 0, "Poor (D/F)": 0}

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var g string
			if err := rows.Scan(&g); err == nil {
				switch g {
				case "O", "A+", "A":
					perfMap["Excellent (A/A+)"]++
				case "B+", "B":
					perfMap["Good (B/B+)"]++
				case "C+", "C":
					perfMap["Average (C/C+)"]++
				case "D", "F":
					perfMap["Poor (D/F)"]++
				}
			}
		}
	}

	healthData := gin.H{
		"course_id": courseID,
		"title":     title,
		// Attendance still mocked as we lack attendance table
		"attendance_distribution": gin.H{"90-100%": 12, "75-90%": 8, "60-75%": 3, "<60%": 1},
		"performance_heatmap":     perfMap,
	}
	c.JSON(http.StatusOK, healthData)
}

func (h *TeacherHandler) GetAtRiskStudents(c *gin.Context) {
	courseID := c.Query("course_id")
	// Dynamic Risk Count (Based on D/F grades)
	var highRisk int
	config.PostgresDB.QueryRow("SELECT COUNT(*) FROM ENROLLS_IN WHERE course_id=$1 AND grade IN ('D', 'F')", courseID).Scan(&highRisk)

	c.JSON(http.StatusOK, gin.H{
		"high_risk_count":   highRisk,
		"medium_risk_count": 2,  // Mocked
		"low_risk_count":    20, // Mocked
	})
}

// GetStudentDetails godoc
// @Summary      Get Student Details
// @Description  Returns detailed info for a student including grade and risk status
// @Tags         Teacher
// @Param        student_id query string true "Student ID"
// @Param        course_id query string true "Course ID"
// @Router       /teacher/student-details [get]
func (h *TeacherHandler) GetStudentDetails(c *gin.Context) {
	studentID := c.Query("student_id")
	courseID := c.Query("course_id")

	var fname, lname, email string
	err := config.PostgresDB.QueryRow("SELECT s_first_name, s_last_name, s_email FROM STUDENT WHERE student_id=$1", studentID).Scan(&fname, &lname, &email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	var grade sql.NullString
	config.PostgresDB.QueryRow("SELECT grade FROM ENROLLS_IN WHERE student_id=$1 AND course_id=$2", studentID, courseID).Scan(&grade)

	currentGrade := "N/A"
	if grade.Valid {
		currentGrade = grade.String
	}

	riskStatus := "No Risk"
	if currentGrade == "D" || currentGrade == "F" {
		riskStatus = "High Risk"
	}

	c.JSON(http.StatusOK, gin.H{
		"student_id":    studentID,
		"name":          fname + " " + lname,
		"email":         email,
		"course_id":     courseID,
		"current_grade": currentGrade,
		"risk_status":   riskStatus,
		"attendance":    "85%",                          // Mocked
		"last_active":   "2 days ago",                   // Mocked
		"other_courses": []string{"CS354TA", "IS353IA"}, // Mocked
	})
}

// PostAnnouncement godoc
// @Summary      Post Announcement
// @Description  Creates a new class announcement
// @Tags         Teacher
// @Router       /teacher/announce [post]
func (h *TeacherHandler) PostAnnouncement(c *gin.Context) {
	var req struct {
		CourseID string `json:"course_id"`
		Content  string `json:"content"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	facultyID := c.GetString("user_id")

	_, err := config.PostgresDB.Exec("INSERT INTO ANNOUNCEMENT (faculty_id, course_id, content) VALUES ($1, $2, $3)", facultyID, req.CourseID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post announcement"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Announcement posted successfully"})
}

// GetAlerts (Unchanged mostly, but validated)
func (h *TeacherHandler) GetAlerts(c *gin.Context) {
	userID := c.GetString("user_id")
	rows, err := config.PostgresDB.Query("SELECT course_id FROM TEACHES WHERE faculty_id=$1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Error"})
		return
	}
	defer rows.Close()
	var alerts []string
	for rows.Next() {
		var cID string
		rows.Scan(&cID)
		// Mock logic but tied to course
		alerts = append(alerts, "Update for "+cID+": Check student progress.")
	}
	// Check for recent announcements
	var count int
	config.PostgresDB.QueryRow("SELECT COUNT(*) FROM ANNOUNCEMENT WHERE faculty_id=$1 AND created_at > NOW() - INTERVAL '1 day'", userID).Scan(&count)
	if count > 0 {
		alerts = append(alerts, "You posted "+string(rune(count))+" announcements recently.")
	}

	if len(alerts) == 0 {
		alerts = append(alerts, "No immediate alerts.")
	}
	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}

// GetProfile godoc
// @Summary      Get Faculty Profile
// @Description  Returns profile information for the logged-in faculty
// @Tags         Teacher
// @Router       /teacher/profile [get]
func (h *TeacherHandler) GetProfile(c *gin.Context) {
	facultyID := c.GetString("user_id")

	var f struct {
		FacultyID string `json:"faculty_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		PhoneNo   string `json:"phone_no"`
	}

	err := config.PostgresDB.QueryRow(`
		SELECT faculty_id, f_first_name, f_last_name, f_email, f_phone_no
		FROM FACULTY WHERE faculty_id=$1
	`, facultyID).Scan(&f.FacultyID, &f.FirstName, &f.LastName, &f.Email, &f.PhoneNo)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Faculty not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Error"})
		return
	}

	// Get Affiliated Departments
	rows, err := config.PostgresDB.Query(`
		SELECT DISTINCT d.dept_name
		FROM TEACHES t
		JOIN COURSE c ON t.course_id = c.course_id
		JOIN DEPARTMENT d ON c.dept_id = d.dept_id
		WHERE t.faculty_id=$1
	`, facultyID)

	var departments []string
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var dName string
			if err := rows.Scan(&dName); err == nil {
				departments = append(departments, dName)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"faculty_id":  f.FacultyID,
		"first_name":  f.FirstName,
		"last_name":   f.LastName,
		"email":       f.Email,
		"phone_no":    f.PhoneNo,
		"departments": departments,
	})
}
