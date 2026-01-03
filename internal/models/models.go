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

type StudyGroup struct {
	GroupID     int       `json:"group_id"`
	CourseID    string    `json:"course_id"`
	GroupName   string    `json:"group_name"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	MemberCount int       `json:"member_count,omitempty"`
}

type GroupMember struct {
	GroupID   int       `json:"group_id"`
	StudentID string    `json:"student_id"`
	JoinedAt  time.Time `json:"joined_at"`
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

type Quiz struct {
	ID        string     `bson:"_id,omitempty" json:"id"`
	CourseID  string     `bson:"course_id" json:"course_id"`
	Topic     string     `bson:"topic" json:"topic"`
	Questions []Question `bson:"questions" json:"questions"`
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
}

type Question struct {
	ID            int      `bson:"id" json:"id"`
	Text          string   `bson:"text" json:"text"`
	Options       []string `bson:"options" json:"options"`
	CorrectOption int      `bson:"correct_option" json:"correct_option"` // Index 0-3
}
