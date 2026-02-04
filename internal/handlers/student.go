package handlers

import (
	"academ_aide/internal/config"
	"academ_aide/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func GetStudentProfile(c *gin.Context) {
	// Extract ID from Token via Context
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	studentID := val.(string)

	ctx := context.Background()
	cacheKey := "student_profile_v2:" + studentID

	// 1. Check Redis Cache
	valCached, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache Hit
		var s models.Student
		if err := json.Unmarshal([]byte(valCached), &s); err == nil {
			fmt.Println("Cache Hit for", studentID) // Logging for verification
			c.JSON(http.StatusOK, s)
			return
		}
	} else if err != redis.Nil {
		fmt.Println("Redis error:", err)
	}

	// 2. Cache Miss - Query Postgres
	fmt.Println("Cache Miss for", studentID) // Logging for verification
	var s models.Student
	err = config.PostgresDB.QueryRow(`
		SELECT student_id, s_first_name, s_last_name, s_email, s_phone_no, semester, year_of_joining, dept_id 
		FROM STUDENT WHERE student_id=$1`, studentID).Scan(
		&s.StudentID, &s.FirstName, &s.LastName, &s.Email, &s.PhoneNo, &s.Semester, &s.YearOfJoining, &s.DeptID,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch Courses Enrolled Count
	var count int
	err = config.PostgresDB.QueryRow(`
		SELECT COUNT(*) FROM ENROLLS_IN WHERE student_id=$1 AND status='Enrolled'
	`, studentID).Scan(&count)
	if err == nil {
		s.CoursesEnrolled = count
	} else {
		fmt.Println("Error fetching course count:", err)
	}

	// Calculate CGPA based on grades
	// Assumption: 10-point scale. A=10, B=8, C=6, D=4, E=2, F=0.
	rows, err := config.PostgresDB.Query(`
		SELECT e.grade, c.credits
		FROM ENROLLS_IN e
		JOIN COURSE c ON e.course_id = c.course_id
		WHERE e.student_id=$1 AND e.grade IS NOT NULL
	`, studentID)
	if err == nil {
		defer rows.Close()
		totalCredits := 0
		totalPoints := 0.0
		for rows.Next() {
			var grade string
			var credits int
			if err := rows.Scan(&grade, &credits); err != nil {
				continue
			}
			points := 0.0
			switch grade {
			case "O", "A+":
				points = 10.0
			case "A":
				points = 9.0
			case "B+":
				points = 8.0
			case "B":
				points = 7.0
			case "C+":
				points = 6.0
			case "C":
				points = 5.0
			case "D":
				points = 4.0
			case "F":
				points = 0.0
			default:
				points = 0.0 // Fail or unknown
			}
			totalPoints += points * float64(credits)
			totalCredits += credits
		}
		if totalCredits > 0 {
			s.CGPA = float64(int((totalPoints/float64(totalCredits))*100)) / 100 // Round to 2 decimal places
		}
	} else {
		fmt.Println("Error fetching grades:", err)
	}

	// Fetch Next Class (Dashboard Feature)
	currentTime := time.Now()
	dayOfWeek := currentTime.Weekday().String()
	currentTimeStr := currentTime.Format("15:04:00") // 24h format

	var nextTitle, nextStart string
	// Find the first class strictly after current time today
	schedQuery := `
		SELECT c.title, TO_CHAR(sch.start_time, 'HH24:MI')
		FROM SCHEDULE sch 
		JOIN ENROLLS_IN e ON sch.course_id = e.course_id
		JOIN COURSE c ON sch.course_id = c.course_id
		WHERE e.student_id=$1 AND sch.day_of_week=$2 AND sch.start_time > $3
		ORDER BY sch.start_time ASC
		LIMIT 1
	`
	err = config.PostgresDB.QueryRow(schedQuery, studentID, dayOfWeek, currentTimeStr).Scan(&nextTitle, &nextStart)
	if err == nil {
		s.NextClass = nextTitle
		s.NextClassTime = nextStart
	} else {
		// Fallback: Check if there's an ongoing class? Or just say "None"
		// If NO class upcoming today, we could check for *ongoing* class or set blank.
		// For simplicity, let's leave it empty/default if no future class today.
		s.NextClass = "No Upcoming Classes"
		s.NextClassTime = "Today"
	}

	// 3. Store result in Redis (5 Minutes TTL)
	if jsonBytes, err := json.Marshal(s); err == nil {
		config.RedisClient.Set(ctx, cacheKey, jsonBytes, 5*time.Minute)
	}

	c.JSON(http.StatusOK, s)
}

func GetStudentTimetable(c *gin.Context) {
	rawID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	studentID := rawID.(string)

	ctx := context.Background()
	cacheKey := "timetable:" + studentID

	// Check Redis Cache
	val, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache Hit
		var schedule []models.ScheduleItem
		json.Unmarshal([]byte(val), &schedule)
		c.JSON(http.StatusOK, gin.H{"source": "cache", "data": schedule})
		return
	} else if err != redis.Nil {
		// Redis error, log it but continue to DB
		fmt.Println("Redis error:", err)
	}

	// Cache Miss - Query Postgres
	// JOIN Schedule -> Course, Schedule -> Section? Already in schedule.
	// Logic: Student -> Enrolls_In -> Course -> Schedule
	query := `
		SELECT c.course_id, c.title, sch.section_name, sch.day_of_week, 
		       TO_CHAR(sch.start_time, 'HH24:MI'), TO_CHAR(sch.end_time, 'HH24:MI'), sch.room_number
		FROM ENROLLS_IN e
		JOIN COURSE c ON e.course_id = c.course_id
		JOIN SCHEDULE sch ON c.course_id = sch.course_id 
		-- Note: Ideally join on section too if student is assigned to section explicitly. 
		-- Schema has Enrolls_In(student, course), but Section is not explicitly linked to Student in Enrolls_in?
		-- Schema: ENROLLS_IN (student_id, course_id). SECTION (section, dept).
		-- Assuming student attends ALL scheduled slots for the course OR we need section assignment.
		-- For simplicity, assuming 1 section per course or returning all.
		-- Wait, Schedule has section_name. We need to know which section the student is in.
		-- Missing 'section_name' in ENROLLS_IN. 
		-- Assuming for this task, we return ALL schedule entries for the enrolled course.
		WHERE e.student_id = $1
		ORDER BY sch.day_of_week, sch.start_time
	`
	// Correction: "JOIN SCHEDULE sch ON c.course_id = sch.course_id" might duplicate if multiple sections.
	// Without section in Enrolls_in, we'll return all. User Context: "Follow specification EXACTLY".
	// Spec doesn't put Section in Enrolls_In. So we fetch all.

	rows, err := config.PostgresDB.Query(query, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var schedule []models.ScheduleItem
	for rows.Next() {
		var i models.ScheduleItem
		if err := rows.Scan(&i.CourseID, &i.Title, &i.SectionName, &i.DayOfWeek, &i.StartTime, &i.EndTime, &i.RoomNumber); err != nil {
			continue
		}
		schedule = append(schedule, i)
	}

	// Store in Redis (1 Hour TTL)
	jsonBytes, _ := json.Marshal(schedule)
	config.RedisClient.Set(ctx, cacheKey, jsonBytes, time.Hour)

	c.JSON(http.StatusOK, gin.H{"source": "database", "data": schedule})
}

func GetStudentResources(c *gin.Context) {
	rawID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	studentID := rawID.(string)

	query := `
		SELECT r.resource_id, r.title, r.description, r.type, r.course_id, r.link
		FROM RESOURCE r
		JOIN ENROLLS_IN e ON r.course_id = e.course_id
		WHERE e.student_id = $1
	`
	rows, err := config.PostgresDB.Query(query, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var resources []models.Resource
	for rows.Next() {
		var r models.Resource
		var desc, link sql.NullString
		if err := rows.Scan(&r.ResourceID, &r.Title, &desc, &r.Type, &r.CourseID, &link); err != nil {
			continue
		}
		if desc.Valid {
			r.Description = desc.String
		}
		if link.Valid {
			r.Link = link.String
		}
		resources = append(resources, r)
	}

	c.JSON(http.StatusOK, resources)
}

// GetAnnouncements godoc
// @Summary      Get Course Announcements
// @Description  Returns enrolled course announcements
// @Tags         Student
// @Router       /student/announcements [get]
func GetAnnouncements(c *gin.Context) {
	rawID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	studentID := rawID.(string)

	query := `
		SELECT a.content, a.created_at, c.title
		FROM ANNOUNCEMENT a
		JOIN ENROLLS_IN e ON a.course_id = e.course_id
		JOIN COURSE c ON a.course_id = c.course_id
		WHERE e.student_id = $1
		ORDER BY a.created_at DESC
	`
	rows, err := config.PostgresDB.Query(query, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	announcements := make([]gin.H, 0)
	for rows.Next() {
		var content, courseTitle string
		var createdAt time.Time
		if err := rows.Scan(&content, &createdAt, &courseTitle); err == nil {
			announcements = append(announcements, gin.H{
				"course":  courseTitle,
				"content": content,
				"date":    createdAt.Format("Jan 02, 15:04"),
			})
		}
	}
	c.JSON(http.StatusOK, announcements)
}

func GetStudentCourses(c *gin.Context) {
	rawID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	studentID := rawID.(string)

	query := `
		SELECT c.course_id, c.title
		FROM ENROLLS_IN e
		JOIN COURSE c ON e.course_id = c.course_id
		WHERE e.student_id = $1 AND e.status = 'Enrolled'
	`
	rows, err := config.PostgresDB.Query(query, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type Course struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	var courses []Course
	for rows.Next() {
		var c Course
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			continue
		}
		courses = append(courses, c)
	}

	c.JSON(http.StatusOK, courses)
}

func GetTeachers(c *gin.Context) {
	rawID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	studentID := rawID.(string)

	// Fetch teachers for courses the student is enrolled in
	query := `
		SELECT DISTINCT f.f_first_name, f.f_last_name, f.f_email, c.title, c.course_id
		FROM TEACHES t
		JOIN FACULTY f ON t.faculty_id = f.faculty_id
		JOIN COURSE c ON t.course_id = c.course_id
		JOIN ENROLLS_IN e ON t.course_id = e.course_id
		WHERE e.student_id = $1 AND e.status = 'Enrolled'
		ORDER BY c.title
	`

	rows, err := config.PostgresDB.Query(query, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type TeacherInfo struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Course   string `json:"course"`
		CourseID string `json:"course_id"`
	}

	var teachers []TeacherInfo
	for rows.Next() {
		var t TeacherInfo
		var fFirst, fLast string
		if err := rows.Scan(&fFirst, &fLast, &t.Email, &t.Course, &t.CourseID); err != nil {
			continue
		}
		t.Name = fFirst + " " + fLast
		teachers = append(teachers, t)
	}

	c.JSON(http.StatusOK, teachers)
}
