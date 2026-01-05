package main

import (
	"academ_aide/internal/config"
	"academ_aide/internal/handlers"
	"academ_aide/internal/middleware"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Initialize DBs
	config.InitDB()

	// Setup Router
	r := gin.Default()

	// Apply CORS Middleware
	r.Use(middleware.CORSMiddleware())

	// Routes
	r.POST("/login", handlers.LoginHandler)

	studentGroup := r.Group("/student")
	studentGroup.Use(middleware.AuthMiddleware())
	{
		studentGroup.GET("/profile", handlers.GetStudentProfile)
		studentGroup.GET("/timetable", handlers.GetStudentTimetable)
		studentGroup.GET("/resources", handlers.GetStudentResources)
		studentGroup.GET("/announcements", handlers.GetAnnouncements)
	}

	chatGroup := r.Group("/chat")
	chatGroup.Use(middleware.AuthMiddleware())
	{
		chatGroup.POST("/message", handlers.ChatHandler)
		chatGroup.DELETE("/history", handlers.ClearChatHandler)
	}

	// OAuth Routes
	r.GET("/auth/google/login", handlers.GoogleLogin)
	r.GET("/auth/google/callback", handlers.GoogleCallback)
	r.POST("/auth/complete-registration", handlers.CompleteRegistration)

	// Feature: AI Quizzes
	r.POST("/quiz/generate", middleware.AuthMiddleware(), handlers.GenerateQuiz)

	// Feature: and Study Groups
	groupRoutes := r.Group("/groups")
	groupRoutes.Use(middleware.AuthMiddleware())
	{
		groupRoutes.GET("/peers", handlers.FindPeers)
		groupRoutes.POST("/create", handlers.CreateGroup)
		groupRoutes.GET("/list", handlers.ListGroups)

	}

	// Feature: AI Academic Intelligence
	aiHandler := handlers.NewAIHandler()
	aiGroup := r.Group("/ai")
	// aiGroup.Use(middleware.AuthMiddleware()) // Optional: Enable auth if needed
	{
		aiGroup.GET("/insights", aiHandler.GetInsights)
		aiGroup.POST("/what-if", aiHandler.CalculateWhatIf)
		aiGroup.POST("/quiz-analysis", aiHandler.AnalyzeQuiz)
	}

	// Feature: Teacher Dashboard
	teacherHandler := handlers.NewTeacherHandler()
	teacherGroup := r.Group("/teacher")
	teacherGroup.Use(middleware.AuthMiddleware())
	teacherGroup.Use(middleware.RoleMiddleware("teacher"))
	{
		teacherGroup.GET("/class-health", teacherHandler.GetClassHealth)
		teacherGroup.GET("/at-risk", teacherHandler.GetAtRiskStudents)
		teacherGroup.GET("/alerts", teacherHandler.GetAlerts)
		teacherGroup.GET("/courses", teacherHandler.GetMyCourses)
		teacherGroup.GET("/students", teacherHandler.GetEnrolledStudents)
		teacherGroup.GET("/student-details", teacherHandler.GetStudentDetails)
		teacherGroup.POST("/announce", teacherHandler.PostAnnouncement)
	}

	// Start Server
	log.Println("Server executing on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server start failed: ", err)
	}
}
