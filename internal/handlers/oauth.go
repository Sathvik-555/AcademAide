package handlers

import (
	"academ_aide/internal/config"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{
	ClientID:     "", // Set in Init
	ClientSecret: "", // Set in Init
	RedirectURL:  "", // Set in Init
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

func InitOAuth() {
	googleOauthConfig.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	googleOauthConfig.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	googleOauthConfig.RedirectURL = os.Getenv("GOOGLE_REDIRECT_URL")
}

func GoogleLogin(c *gin.Context) {
	InitOAuth() // Ensure config is loaded
	oauthState := generateStateOauthCookie(c)
	u := googleOauthConfig.AuthCodeURL(oauthState)
	c.Redirect(http.StatusTemporaryRedirect, u)
}

func GoogleCallback(c *gin.Context) {
	InitOAuth()
	oauthState, _ := c.Cookie("oauthstate")

	if c.Query("state") != oauthState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}

	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Code exchange failed"})
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User info failed"})
		return
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Read body failed"})
		return
	}

	var googleUser map[string]interface{}
	json.Unmarshal(content, &googleUser)

	email := googleUser["email"].(string)

	// Check if student exists
	var studentID string
	err = config.PostgresDB.QueryRow("SELECT student_id FROM STUDENT WHERE s_email=$1", email).Scan(&studentID)

	if err == nil {
		// Student exists -> Login
		tokenString := generateToken(studentID, false)

		// Redirect to Frontend Callback with Login Mode
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/auth/callback?token="+tokenString+"&mode=login&student_id="+studentID)
	} else {
		// New Student -> Onboarding
		// Generate Temporary Token
		tempTokenString := generateToken(email, true) // Use Email as ID for now, marked as partial

		// Redirect to Frontend Callback with Signup Mode
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/auth/callback?token="+tempTokenString+"&mode=signup&email="+email)
	}
}

// CompleteRegistration Request
type CompleteRegistrationRequest struct {
	DeptID        string `json:"dept_id"`
	Semester      int    `json:"semester"`
	YearOfJoining int    `json:"year_of_joining"`
	PhoneNo       string `json:"phone_no"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
}

func CompleteRegistration(c *gin.Context) {
	// 1. Get Token from Header (MIDDLEWARE SHOULD HANDLE THIS usually, but for Onboarding logic we might need custom handling or use the AuthMiddleware but allow 'partial' tokens?)
	// For simplicity, let's assume the client sends the 'temp_token' in Authorization header.
	// We need to parse it here to get the 'email'.

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Parse Token manually to check 'partial' claim and get email
	// In a real app, Middleware would put claims in context.
	// ... (Simplified: We assume middleware passed, but we need to verify 'partial' claim if we want to be strict)

	// Let's rely on claims from Context if AuthMiddleware is used.
	// BUT AuthMiddleware might reject "partial" tokens if we don't adjust it.
	// For now, let's just parse it again or assume the `student_id` in context is actually `email` for partial tokens.

	// BETTER APPROACH: Just use the email from the claims.

	// NOTE: This requires AuthMiddleware to allow this flow.
	// Let's implement parsing logic here for safety if middleware isn't adapted yet.

	tokenString := authHeader[len("Bearer "):]
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	isPartial, _ := claims["partial"].(bool)
	if !isPartial {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not a registration token"})
		return
	}

	email := claims["student_id"].(string) // We stored email in student_id field for partial token

	var req CompleteRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Generate new Student ID
	newStudentID := "S" + time.Now().Format("20060102150405")

	// Insert into DB
	_, err := config.PostgresDB.Exec("INSERT INTO STUDENT (student_id, s_first_name, s_last_name, s_email, s_phone_no, semester, year_of_joining, dept_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		newStudentID, req.FirstName, req.LastName, email, req.PhoneNo, req.Semester, req.YearOfJoining, req.DeptID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Insert Failed: " + err.Error()})
		return
	}

	// Generate Final Token
	finalToken := generateToken(newStudentID, false)

	c.JSON(http.StatusOK, gin.H{"token": finalToken, "student_id": newStudentID})
}

func generateToken(id string, partial bool) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"student_id": id,
		"partial":    partial,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "supersecretkey"
	}

	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func generateStateOauthCookie(c *gin.Context) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	// Secure cookie (http-only)
	// Domain/Path should be set correctly for localhost
	c.SetCookie("oauthstate", state, 3600, "/", "localhost", false, true)

	return state
}
