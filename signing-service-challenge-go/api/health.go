package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// Health evaluates the health of the service and writes a standardized response.
func (s *Server) Health(c *gin.Context) {
	health := HealthResponse{
		Status:  "pass",
		Version: "v0",
	}

	c.JSON(http.StatusOK, Response{Data: health})
}
