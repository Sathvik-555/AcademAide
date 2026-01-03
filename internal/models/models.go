package models

import "time"

// PostgreSQL Entities

type Student struct {
	StudentID       string  `json:"student_id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Email           string  `json:"email"`
	PhoneNo         string  `json:"phone_no"`
	Semester        int     `json:"semester"`
	YearOfJoining   int     `json:"year_of_joining"`
	DeptID          string  `json:"dept_id"`
	CoursesEnrolled int     `json:"courses_enrolled,omitempty"`
	CGPA            float64 `json:"cgpa,omitempty"`
}

type ScheduleItem struct {
	CourseID    string `json:"course_id"`
	Title       string `json:"title"`
	SectionName string `json:"section_name"`
	DayOfWeek   string `json:"day_of_week"`
	StartTime   string `json:"start_time"` // Returning as string for JSON simplicity
	EndTime     string `json:"end_time"`
	RoomNumber  string `json:"room_number"`
}

// MongoDB Entities

type ChatLog struct {
	StudentID string    `bson:"student_id" json:"student_id"`
	Message   string    `bson:"message" json:"message"`
	Intent    string    `bson:"intent" json:"intent"`
	Sentiment string    `bson:"sentiment" json:"sentiment"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	IsBot     bool      `bson:"is_bot" json:"is_bot"` // To distinguish user vs bot
}

type ChatContext struct {
	StudentID       string    `bson:"student_id" json:"student_id"`
	LastTopic       string    `bson:"last_topic" json:"last_topic"`
	Emotion         string    `bson:"emotion" json:"emotion"`
	LastInteraction time.Time `bson:"last_interaction" json:"last_interaction"`
}
