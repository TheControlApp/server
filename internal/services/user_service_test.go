package services

import (
	"testing"
)

func TestCreateUserRequestValidation(t *testing.T) {
	req := CreateUserRequest{
		LoginName:  "testuser",
		ScreenName: "Test User",
		Password:   "password123",
	}
	user, err := CreateUser(req)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if user.LoginName != req.LoginName {
		t.Errorf("LoginName mismatch: got %s, want %s", user.LoginName, req.LoginName)
	}
	if user.ScreenName != req.ScreenName {
		t.Errorf("ScreenName mismatch: got %s, want %s", user.ScreenName, req.ScreenName)
	}
}
