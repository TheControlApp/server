package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thecontrolapp/server/internal/api/responses"
)

type HealthHandlers struct{}

func NewHealthHandlers() *HealthHandlers {
	return &HealthHandlers{}
}

// HealthCheck performs a health check
// HealthCheck godoc
//
//	@Summary		Health check
//	@Description	Check if the server is running and responsive
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	responses.HealthResponse
//	@Failure		500	{object}	responses.ErrorResponse
//	@Router			/health [get]
func (h *HealthHandlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, responses.HealthResponse{
		Status:  "ok",
		Message: "Server is running",
	})
}
