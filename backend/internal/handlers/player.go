package handlers

import (
	"net/http"
	"strconv"

	"chess-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// PlayerHandler handles player statistics HTTP requests
type PlayerHandler struct {
	playerService *services.PlayerService
}

// NewPlayerHandler creates a new player handler
func NewPlayerHandler(playerService *services.PlayerService) *PlayerHandler {
	return &PlayerHandler{
		playerService: playerService,
	}
}

// GetPlayerStatistics retrieves comprehensive player statistics
// GET /api/stats/player/:playername
func (h *PlayerHandler) GetPlayerStatistics(c *gin.Context) {
	playerName := c.Param("playername")
	if playerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Player name is required",
		})
		return
	}

	stats, err := h.playerService.GetPlayerStatistics(playerName)
	if err != nil {
		logrus.Errorf("Failed to get player statistics for %s: %v", playerName, err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Player not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetAllPlayers returns a list of all tracked players
// GET /api/stats/players
func (h *PlayerHandler) GetAllPlayers(c *gin.Context) {
	players := h.playerService.GetAllPlayers()
	
	c.JSON(http.StatusOK, gin.H{
		"players": players,
		"count": len(players),
	})
}

// GetTopPlayers returns players ranked by average accuracy
// GET /api/stats/leaderboard
func (h *PlayerHandler) GetTopPlayers(c *gin.Context) {
	// Get limit from query parameter (default to 10)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	
	if limit > 100 {
		limit = 100 // Cap at 100 to prevent abuse
	}

	rankings := h.playerService.GetTopPlayers(limit)
	
	c.JSON(http.StatusOK, gin.H{
		"rankings": rankings,
		"count": len(rankings),
		"limit": limit,
	})
}

// GetPlayerGames returns game history for a specific player
// GET /api/stats/player/:playername/games
func (h *PlayerHandler) GetPlayerGames(c *gin.Context) {
	playerName := c.Param("playername")
	if playerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Player name is required",
		})
		return
	}

	// Get player statistics (which includes recent games)
	stats, err := h.playerService.GetPlayerStatistics(playerName)
	if err != nil {
		logrus.Errorf("Failed to get player games for %s: %v", playerName, err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Player not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"playerName": stats.PlayerName,
		"games": stats.RecentGames,
		"totalGames": stats.GamesAnalyzed,
	})
} 