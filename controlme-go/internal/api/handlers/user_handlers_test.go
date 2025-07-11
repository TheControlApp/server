package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/thecontrolapp/controlme-go/internal/api/handlers"
	"github.com/thecontrolapp/controlme-go/internal/services"
)

type mockUserService struct{}

func (m *mockUserService) GetAllUsers() ([]services.CreateUserRequest, error) {
	return []services.CreateUserRequest{{LoginName: "test", ScreenName: "Test", Password: "pass"}}, nil
}

func TestGetUsersEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	userHandlers := handlers.NewUserHandlers(&mockUserService{})
	r.GET("/api/v1/users", userHandlers.GetUsers)

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
