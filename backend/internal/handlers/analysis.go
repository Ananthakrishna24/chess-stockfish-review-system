package handlers

import (
	"net/http"

	"chess-backend/internal/models"
	"chess-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AnalysisHandler handles analysis-related HTTP requests
type AnalysisHandler struct {
	analysisService *services.AnalysisService
}

// NewAnalysisHandler creates a new analysis handler
func NewAnalysisHandler(analysisService *services.AnalysisService) *AnalysisHandler {
	return &AnalysisHandler{
		analysisService: analysisService,
	}
}

// AnalyzeGame starts game analysis
// POST /api/games/analyze
func (h *AnalysisHandler) AnalyzeGame(c *gin.Context) {
	var request models.AnalyzeGameRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Errorf("Invalid game analysis request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate PGN is not empty
	if request.PGN == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "PGN cannot be empty",
		})
		return
	}

	// Start analysis
	gameID := h.analysisService.StartGameAnalysis(request.PGN, request.Options)

	logrus.Infof("Started analysis for game %s", gameID)

	c.JSON(http.StatusAccepted, gin.H{
		"gameId": gameID,
		"status": "queued",
		"message": "Analysis started successfully",
	})
}

// GetAnalysis retrieves completed analysis results
// GET /api/games/analyze/:gameId
func (h *AnalysisHandler) GetAnalysis(c *gin.Context) {
	gameID := c.Param("gameId")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Game ID is required",
		})
		return
	}

	result, err := h.analysisService.GetAnalysisResult(gameID)
	if err != nil {
		logrus.Errorf("Failed to get analysis result for game %s: %v", gameID, err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Analysis not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetProgress retrieves analysis progress
// GET /api/games/analyze/:gameId/progress
func (h *AnalysisHandler) GetProgress(c *gin.Context) {
	gameID := c.Param("gameId")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Game ID is required",
		})
		return
	}

	progress, err := h.analysisService.GetAnalysisProgress(gameID)
	if err != nil {
		logrus.Errorf("Failed to get analysis progress for game %s: %v", gameID, err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Analysis job not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// AnalyzePosition analyzes a single position
// POST /api/positions/analyze
func (h *AnalysisHandler) AnalyzePosition(c *gin.Context) {
	var request models.AnalyzePositionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Errorf("Invalid position analysis request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate FEN is not empty
	if request.FEN == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "FEN cannot be empty",
		})
		return
	}

	// Validate depth and time limits
	if request.Depth < 0 || request.Depth > 24 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Depth must be between 0 and 24",
		})
		return
	}

	if request.TimeLimit < 0 || request.TimeLimit > 30000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Time limit must be between 0 and 30000ms",
		})
		return
	}

	// Analyze position
	result, err := h.analysisService.AnalyzePosition(request)
	if err != nil {
		logrus.Errorf("Failed to analyze position: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Position analysis failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetEngineConfig returns current engine configuration
// GET /api/engine/config
func (h *AnalysisHandler) GetEngineConfig(c *gin.Context) {
	config := h.analysisService.GetEngineConfig()
	c.JSON(http.StatusOK, config)
}

// UpdateEngineConfig updates engine configuration
// POST /api/engine/config
func (h *AnalysisHandler) UpdateEngineConfig(c *gin.Context) {
	var request models.UpdateEngineConfigRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Errorf("Invalid engine config update request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Update configuration
	if err := h.analysisService.UpdateEngineConfig(request); err != nil {
		logrus.Errorf("Failed to update engine configuration: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Configuration update failed",
			"details": err.Error(),
		})
		return
	}

	// Return updated configuration
	config := h.analysisService.GetEngineConfig()
	c.JSON(http.StatusOK, gin.H{
		"message": "Configuration updated successfully",
		"config": config,
	})
} 