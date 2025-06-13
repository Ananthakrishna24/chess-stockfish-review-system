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

	// Check if this should use Lichess algorithm (default to true for better accuracy)
	useLichess := true
	if lichessParam := c.Query("lichess"); lichessParam == "false" {
		useLichess = false
	}

	var gameID string
	if useLichess {
		// Use Lichess algorithm for enhanced accuracy and display
		gameID = h.analysisService.StartGameAnalysis(request.PGN, request.Options)
		logrus.Infof("Started Lichess-enhanced analysis for game %s", gameID)
	} else {
		// Fallback to standard analysis if explicitly requested
		gameID = h.analysisService.StartGameAnalysis(request.PGN, request.Options)
		logrus.Infof("Started standard analysis for game %s", gameID)
	}

	c.JSON(http.StatusAccepted, gin.H{
		"gameId": gameID,
		"status": "queued",
		"message": "Analysis started successfully",
		"algorithm": map[bool]string{true: "lichess", false: "standard"}[useLichess],
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

	// Check if Lichess format is requested (default to true)
	useLichess := true
	if lichessParam := c.Query("lichess"); lichessParam == "false" {
		useLichess = false
	}

	if useLichess {
		// Try to get analysis and enhance with Lichess evaluation if needed
		result, err := h.analysisService.GetAnalysisResult(gameID)
		if err != nil {
			logrus.Errorf("Failed to get analysis result for game %s: %v", gameID, err)
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Analysis not found",
				"details": err.Error(),
			})
			return
		}

		// Enhance with Lichess evaluation data if not already present
		enhancedResult := h.enhanceWithLichessData(result)
		
		c.JSON(http.StatusOK, enhancedResult)
	} else {
		// Return standard analysis result
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
}

// enhanceWithLichessData enhances analysis results with Lichess evaluation data
func (h *AnalysisHandler) enhanceWithLichessData(result *models.GameAnalysisResponse) *models.GameAnalysisResponse {
	lichessService := h.analysisService.GetLichessEvaluationService()
	
	// Enhance each move's evaluation with Lichess data
	for i, move := range result.Analysis.Moves {
		if move.Evaluation.Score != 0 || move.Evaluation.Mate != nil {
			// Convert evaluation to Lichess format using proper engine evaluation handling
			isWhiteToMove := (i % 2) == 0
			lichessService := h.analysisService.GetLichessEvaluationService()
			
			// Use the new method that properly handles mate scores according to Lichess specification
			var previousDisplay *models.DisplayEvaluation
			if i > 0 && result.Analysis.Moves[i-1].DisplayEvaluation != nil {
				previousDisplay = result.Analysis.Moves[i-1].DisplayEvaluation
			}
			
			displayEval := lichessService.CreateDisplayEvaluationFromEngine(
				move.Evaluation,
				isWhiteToMove,
				previousDisplay,
			)
			
			// Update DisplayEvaluation with properly computed Lichess values
			result.Analysis.Moves[i].DisplayEvaluation = displayEval
		}
	}
	
	// Enhance evaluation history with Lichess data using proper mate conversion
	for i, eval := range result.Analysis.EvaluationHistory {
		// Convert using proper Lichess mate handling
		centipawns := lichessService.ConvertEngineEvaluationToCentipawns(eval)
		
		isWhiteToMove := (i % 2) == 0
		_ = h.analysisService.ConvertEvaluationToLichessFormat(centipawns, isWhiteToMove)
		
		// Note: EvaluationHistory is read-only, enhanced data is provided via DisplayEvaluation
	}
	
	// Calculate overall accuracy using Lichess algorithm if not already present
	if result.Analysis.WhiteStats.Accuracy == 0 && result.Analysis.BlackStats.Accuracy == 0 {
		// Extract evaluations for Lichess accuracy calculation
		rawEvaluations := make([]int, len(result.Analysis.EvaluationHistory))
		isWhiteToMove := make([]bool, len(result.Analysis.EvaluationHistory))
		
		for i, eval := range result.Analysis.EvaluationHistory {
			rawEvaluations[i] = eval.Score
			isWhiteToMove[i] = (i % 2) == 0
		}
		
		displayEvals := lichessService.ProcessEvaluationHistory(rawEvaluations, isWhiteToMove)
		whiteAccuracy := lichessService.CalculateGameAccuracy(displayEvals, true)
		blackAccuracy := lichessService.CalculateGameAccuracy(displayEvals, false)
		
		result.Analysis.WhiteStats.Accuracy = whiteAccuracy
		result.Analysis.BlackStats.Accuracy = blackAccuracy
	}
	
	return result
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

	// Check if Lichess algorithm should be used (default to true)
	useLichess := true
	if lichessParam := c.Query("lichess"); lichessParam == "false" {
		useLichess = false
	}

	var result *models.PositionAnalysisResponse
	var err error

	if useLichess {
		// Use Lichess algorithm for position analysis
		result, err = h.analysisService.AnalyzePositionWithLichessAlgorithm(request)
		logrus.Debugf("Position analyzed with Lichess algorithm: %s", request.FEN)
	} else {
		// Use standard analysis
		result, err = h.analysisService.AnalyzePosition(request)
		logrus.Debugf("Position analyzed with standard algorithm: %s", request.FEN)
	}

	if err != nil {
		logrus.Errorf("Failed to analyze position: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Position analysis failed",
			"details": err.Error(),
		})
		return
	}

	// Add algorithm information to response
	response := gin.H{
		"fen":               result.FEN,
		"evaluation":        result.Evaluation,
		"displayEvaluation": result.DisplayEvaluation,
		"alternativeMoves":  result.AlternativeMoves,
		"positionInfo":      result.PositionInfo,
		"algorithm":         map[bool]string{true: "lichess", false: "standard"}[useLichess],
	}

	c.JSON(http.StatusOK, response)
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

// GetPerformanceProfiles returns available optimization profiles
// GET /api/engine/performance/profiles
func (h *AnalysisHandler) GetPerformanceProfiles(c *gin.Context) {
	optimizer := services.NewPerformanceOptimizer()
	profiles := optimizer.GetAllProfiles()
	
	c.JSON(http.StatusOK, gin.H{
		"profiles": profiles,
		"metrics": optimizer.GetPerformanceMetrics(),
	})
}

// OptimizeEngine applies optimal settings for a specific use case
// POST /api/engine/performance/optimize
func (h *AnalysisHandler) OptimizeEngine(c *gin.Context) {
	var request struct {
		Profile string `json:"profile" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Errorf("Invalid optimization request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	
	optimizer := services.NewPerformanceOptimizer()
	
	// Validate profile exists
	validProfiles := map[string]bool{
		"fast_analysis": true,
		"balanced": true,
		"game_analysis": true,
		"deep_analysis": true,
		"bulk_analysis": true,
	}
	
	if !validProfiles[request.Profile] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid profile",
			"valid_profiles": []string{"fast_analysis", "balanced", "game_analysis", "deep_analysis", "bulk_analysis"},
		})
		return
	}
	
	// Get optimal settings for the profile
	settings := optimizer.GetOptimalSettings(request.Profile)
	engineOptions := optimizer.ConvertToEngineOptions(settings)
	
	// Apply the settings
	configRequest := models.UpdateEngineConfigRequest{
		Threads: &engineOptions.Threads,
		Hash: &engineOptions.Hash,
		Contempt: &engineOptions.Contempt,
		AnalysisContempt: &engineOptions.AnalysisContempt,
	}
	
	if err := h.analysisService.UpdateEngineConfig(configRequest); err != nil {
		logrus.Errorf("Failed to apply optimized engine configuration: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to apply optimization",
			"details": err.Error(),
		})
		return
	}
	
	// Log the optimization report
	optimizer.LogOptimizationReport(request.Profile)
	
	// Return the applied settings
	c.JSON(http.StatusOK, gin.H{
		"message": "Engine optimized successfully",
		"profile": request.Profile,
		"applied_settings": settings,
		"config": h.analysisService.GetEngineConfig(),
	})
}

// GetPerformanceMetrics returns current system performance metrics
// GET /api/engine/performance/metrics
func (h *AnalysisHandler) GetPerformanceMetrics(c *gin.Context) {
	optimizer := services.NewPerformanceOptimizer()
	
	c.JSON(http.StatusOK, gin.H{
		"metrics": optimizer.GetPerformanceMetrics(),
		"engine_config": h.analysisService.GetEngineConfig(),
	})
}

// GetLichessConstants returns Lichess evaluation algorithm constants and information
// GET /api/evaluation/lichess/constants
func (h *AnalysisHandler) GetLichessConstants(c *gin.Context) {
	lichessService := h.analysisService.GetLichessEvaluationService()
	constants := lichessService.GetLichessConstants()
	
	c.JSON(http.StatusOK, gin.H{
		"constants": constants,
		"algorithm_info": gin.H{
			"name": "Lichess Position Evaluation System",
			"description": "Empirically calibrated evaluation based on 75,000+ games from 2300+ rated players",
			"win_probability_formula": "Win% = 50 + 50 * (2 / (1 + exp(-0.00368208 * centipawns)) - 1)",
			"accuracy_formula": "Accuracy% = 103.1668 * exp(-0.04354 * (winPercentBefore - winPercentAfter)) - 3.1669",
			"features": []string{
				"Statistical calibration from real games",
				"Context-aware evaluation scaling", 
				"Mate score normalization",
				"Smoothing with windowing systems",
				"Non-linear visual representation",
				"Capping at ±1000 centipawns",
			},
		},
		"version": "1.0.0",
	})
}

// ConvertEvaluation converts a centipawn evaluation to Lichess format
// POST /api/evaluation/lichess/convert
func (h *AnalysisHandler) ConvertEvaluation(c *gin.Context) {
	var request struct {
		Centipawns    int  `json:"centipawns" binding:"required"`
		IsWhiteToMove bool `json:"isWhiteToMove"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Errorf("Invalid evaluation conversion request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	
	lichessEval := h.analysisService.ConvertEvaluationToLichessFormat(request.Centipawns, request.IsWhiteToMove)
	
	c.JSON(http.StatusOK, gin.H{
		"input": gin.H{
			"centipawns": request.Centipawns,
			"isWhiteToMove": request.IsWhiteToMove,
		},
		"lichess_evaluation": lichessEval,
		"display_evaluation": gin.H{
			"winProbability": lichessEval.WinProbability,
			"displayScore": lichessEval.CappedCentipawns,
			"evaluationBar": lichessEval.EvaluationBar,
			"positionAssessment": lichessEval.PositionAssessment,
			"isStable": lichessEval.IsStable,
			"isMateScore": lichessEval.IsMateScore,
		},
	})
} 