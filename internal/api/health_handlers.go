package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandlers struct{}

func newHealthHandlers() *HealthHandlers {
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
	c.JSON(http.StatusOK, HealthResponse{
		Status:  "ok",
		Message: "Server is running",
	})
}
