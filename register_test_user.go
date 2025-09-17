package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const serverURL = "http://localhost:8080"

type RegisterRequest struct {
	ScreenName string `json:"screen_name"`
	LoginName  string `json:"login_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

func main() {
	fmt.Println("Registering test user...")

	registerReq := RegisterRequest{
		ScreenName: "Test User",
		LoginName:  "test_user",
		Email:      "test@example.com",
		Password:   "test_password",
	}

	jsonData, _ := json.Marshal(registerReq)
	resp, err := http.Post(serverURL+"/api/v1/auth/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Registration request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		fmt.Println("✅ User registered successfully!")
	} else {
		fmt.Printf("Registration status: %d\n", resp.StatusCode)
		// User might already exist, which is fine for testing
	}
}
