package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/thecontrolapp/controlme-go/internal/auth"
	"github.com/thecontrolapp/controlme-go/internal/config"
	"github.com/thecontrolapp/controlme-go/internal/services"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	return db
}

func TestLegacyHandlers_AppCommand_WrongVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	db := setupTestDB(t)

	// Create services
	authService := auth.NewAuthService("test-key", "test-secret", 0)
	userService := services.NewUserService(db, authService)
	cmdService := services.NewCommandService(db)
	legacyService := services.NewLegacyService(db, authService)
	cfg := &config.Config{}

	// Create handler
	handler := NewLegacyHandlers(db, userService, cmdService, legacyService, authService, cfg)

	// Setup router
	router := gin.New()
	router.GET("/AppCommand.aspx", handler.AppCommand)

	// Test wrong version
	req, err := http.NewRequest("GET", "/AppCommand.aspx?vrs=001&usernm=test&pwd=test&cmd=Outstanding", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Wrong version.")
}

func TestLegacyHandlers_AppCommand_MissingParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	db := setupTestDB(t)

	// Create services
	authService := auth.NewAuthService("test-key", "test-secret", 0)
	userService := services.NewUserService(db, authService)
	cmdService := services.NewCommandService(db)
	legacyService := services.NewLegacyService(db, authService)
	cfg := &config.Config{}

	// Create handler
	handler := NewLegacyHandlers(db, userService, cmdService, legacyService, authService, cfg)

	// Setup router
	router := gin.New()
	router.GET("/AppCommand.aspx", handler.AppCommand)

	// Test missing parameters
	req, err := http.NewRequest("GET", "/AppCommand.aspx?vrs=012", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	// Should contain error message due to missing parameters
	assert.Contains(t, rr.Body.String(), "<!DOCTYPE html>")
}

func TestLegacyHandlers_GetContent_WrongVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	db := setupTestDB(t)

	// Create services
	authService := auth.NewAuthService("test-key", "test-secret", 0)
	userService := services.NewUserService(db, authService)
	cmdService := services.NewCommandService(db)
	legacyService := services.NewLegacyService(db, authService)
	cfg := &config.Config{}

	// Create handler
	handler := NewLegacyHandlers(db, userService, cmdService, legacyService, authService, cfg)

	// Setup router
	router := gin.New()
	router.GET("/GetContent.aspx", handler.GetContent)

	// Test wrong version
	req, err := http.NewRequest("GET", "/GetContent.aspx?vrs=001&usernm=test&pwd=test", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "<!DOCTYPE html>")
	// Should have empty labels when version is wrong
	assert.Contains(t, rr.Body.String(), `<asp:Label ID="SenderId" runat="server"></asp:Label>`)
}

func TestLegacyHandlers_GetCount_WrongVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	db := setupTestDB(t)

	// Create services
	authService := auth.NewAuthService("test-key", "test-secret", 0)
	userService := services.NewUserService(db, authService)
	cmdService := services.NewCommandService(db)
	legacyService := services.NewLegacyService(db, authService)
	cfg := &config.Config{}

	// Create handler
	handler := NewLegacyHandlers(db, userService, cmdService, legacyService, authService, cfg)

	// Setup router
	router := gin.New()
	router.GET("/GetCount.aspx", handler.GetCount)

	// Test wrong version
	req, err := http.NewRequest("GET", "/GetCount.aspx?vrs=001&usernm=test&pwd=test", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "<!DOCTYPE html>")
	// Should have empty labels when version is wrong
	assert.Contains(t, rr.Body.String(), `<asp:Label ID="result" runat="server"></asp:Label>`)
}

func TestLegacyHandlers_ProcessComplete_MissingParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	db := setupTestDB(t)

	// Create services
	authService := auth.NewAuthService("test-key", "test-secret", 0)
	userService := services.NewUserService(db, authService)
	cmdService := services.NewCommandService(db)
	legacyService := services.NewLegacyService(db, authService)
	cfg := &config.Config{}

	// Create handler
	handler := NewLegacyHandlers(db, userService, cmdService, legacyService, authService, cfg)

	// Setup router
	router := gin.New()
	router.POST("/ProcessComplete.aspx", handler.ProcessComplete)

	// Test missing parameters
	req, err := http.NewRequest("POST", "/ProcessComplete.aspx", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Missing required parameters")
}

func TestLegacyHandlers_DeleteOut_MissingParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	db := setupTestDB(t)

	// Create services
	authService := auth.NewAuthService("test-key", "test-secret", 0)
	userService := services.NewUserService(db, authService)
	cmdService := services.NewCommandService(db)
	legacyService := services.NewLegacyService(db, authService)
	cfg := &config.Config{}

	// Create handler
	handler := NewLegacyHandlers(db, userService, cmdService, legacyService, authService, cfg)

	// Setup router
	router := gin.New()
	router.POST("/DeleteOut.aspx", handler.DeleteOut)

	// Test missing parameters
	req, err := http.NewRequest("POST", "/DeleteOut.aspx", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Missing required parameters")
}

func TestLegacyHandlers_GetOptions_MissingParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	db := setupTestDB(t)

	// Create services
	authService := auth.NewAuthService("test-key", "test-secret", 0)
	userService := services.NewUserService(db, authService)
	cmdService := services.NewCommandService(db)
	legacyService := services.NewLegacyService(db, authService)
	cfg := &config.Config{}

	// Create handler
	handler := NewLegacyHandlers(db, userService, cmdService, legacyService, authService, cfg)

	// Setup router
	router := gin.New()
	router.GET("/GetOptions.aspx", handler.GetOptions)

	// Test missing parameters
	req, err := http.NewRequest("GET", "/GetOptions.aspx", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Missing required parameters")
}
