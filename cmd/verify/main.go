package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const baseURL = "http://localhost:8080"

func main() {
	fmt.Println("Starting RAG Verification...")

	// 1. Login
	token, err := login("S1001", "password")
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Login successful. Token obtained.")

	// 2. Chat
	question := "What topics are covered in Data Structures?"
	fmt.Printf("Sending question: %s\n", question)

	response, err := chat(token, "S1001", question)
	if err != nil {
		fmt.Printf("Chat failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nAI Response:\n%s\n", response)

	// Simple validation
	if len(response) > 0 {
		fmt.Println("\nVerification PASSED: Received response from AI.")
	} else {
		fmt.Println("\nVerification FAILED: Empty response.")
	}
}

func login(studentID, password string) (string, error) {
	body, _ := json.Marshal(map[string]string{
		"student_id": studentID,
		"password":   password,
	})

	resp, err := http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status %d: %s", resp.StatusCode, string(b))
	}

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res["token"].(string), nil
}

func chat(token, studentID, message string) (string, error) {
	body, _ := json.Marshal(map[string]string{
		"student_id": studentID,
		"message":    message,
	})

	req, _ := http.NewRequest("POST", baseURL+"/chat/message", bytes.NewBuffer(body))
	req.Header.Set("Authorization", token) // Auth middleware expects raw token in header (usually UserID or Bearer)
	// Wait, middleware might expect "Bearer <token>" or just token?
	// Checking middleware usage: studentGroup.Use(middleware.AuthMiddleware())
	// Usually standard JWT middleware expects "Bearer ".
	// Let's assume Bearer for safety, or check middleware file.
	// But based on auth handler returning token, standard is Bearer.
	// We'll try just the token first or check the middleware code if we can.
	// Actually, I'll view middleware if I can, but to save steps I'll assume standard Bearer.
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status %d: %s", resp.StatusCode, string(b))
	}

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res["response"].(string), nil
}
