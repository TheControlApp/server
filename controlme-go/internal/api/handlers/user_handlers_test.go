package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"workspace/server/controlme-go/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func TestGetUsersEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/users", handlers.GetUsers)

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
