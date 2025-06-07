package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health and status endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health returns basic health status
// GET /api/health
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "chess-analysis-backend",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC(),
		"uptime":    time.Since(startTime).Seconds(),
	})
}

// Stats returns system statistics
// GET /api/stats
func (h *HealthHandler) Stats(c *gin.Context) {
	// This would typically include more detailed metrics
	// For now, return basic stats
	c.JSON(http.StatusOK, gin.H{
		"service":        "chess-analysis-backend",
		"version":        "1.0.0",
		"uptime_seconds": time.Since(startTime).Seconds(),
		"timestamp":      time.Now().UTC(),
		"endpoints": gin.H{
			"games_analyze":     "/api/games/analyze",
			"positions_analyze": "/api/positions/analyze",
			"engine_config":     "/api/engine/config",
			"health":            "/api/health",
		},
	})
}

// Global variable to track startup time
var startTime = time.Now() 