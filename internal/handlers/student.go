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
