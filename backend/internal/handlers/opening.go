package handlers

import (
	"net/http"
	"strings"

	"chess-backend/internal/models"
	"chess-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// OpeningHandler handles opening database HTTP requests
type OpeningHandler struct {
	openingService *services.OpeningService
}

// NewOpeningHandler creates a new opening handler
func NewOpeningHandler(openingService *services.OpeningService) *OpeningHandler {
	return &OpeningHandler{
		openingService: openingService,
	}
}

// SearchOpenings searches opening database by various criteria
// GET /api/openings/search
func (h *OpeningHandler) SearchOpenings(c *gin.Context) {
	// Get query parameters
	ecoCode := c.Query("eco")
	fen := c.Query("fen")
	movesStr := c.Query("moves")
	name := c.Query("name")

	var results []models.OpeningInfo
	var err error

	// Determine search type and execute
	switch {
	case ecoCode != "":
		// Search by ECO code
		result, searchErr := h.openingService.SearchByECO(ecoCode)
		if searchErr != nil {
			err = searchErr
		} else {
			results = []models.OpeningInfo{*result}
		}

	case fen != "":
		// Search by FEN position
		result, searchErr := h.openingService.SearchByPosition(fen)
		if searchErr != nil {
			err = searchErr
		} else {
			results = []models.OpeningInfo{*result}
		}

	case movesStr != "":
		// Search by move sequence
		moves := strings.Fields(movesStr)
		result, searchErr := h.openingService.SearchByMoves(moves)
		if searchErr != nil {
			err = searchErr
		} else {
			results = []models.OpeningInfo{*result}
		}

	case name != "":
		// Search by opening name
		searchResults, searchErr := h.openingService.GetOpeningByName(name)
		if searchErr != nil {
			err = searchErr
		} else {
			results = make([]models.OpeningInfo, len(searchResults))
			for i, opening := range searchResults {
				results[i] = *opening
			}
		}

	default:
		// No search criteria provided - return error
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Search criteria required: eco, fen, moves, or name",
		})
		return
	}

	// Handle search errors
	if err != nil {
		logrus.Errorf("Opening search failed: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Opening not found",
			"details": err.Error(),
		})
		return
	}

	// Return successful results
	response := models.OpeningSearchResponse{
		Results: results,
		Count:   len(results),
	}

	c.JSON(http.StatusOK, response)
}

// GetOpeningByECO retrieves specific opening by ECO code
// GET /api/openings/:eco
func (h *OpeningHandler) GetOpeningByECO(c *gin.Context) {
	ecoCode := c.Param("eco")
	if ecoCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ECO code is required",
		})
		return
	}

	opening, err := h.openingService.SearchByECO(ecoCode)
	if err != nil {
		logrus.Errorf("Failed to find opening %s: %v", ecoCode, err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Opening not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, opening)
}

// GetAllOpenings returns all openings in the database
// GET /api/openings
func (h *OpeningHandler) GetAllOpenings(c *gin.Context) {
	openings := h.openingService.GetAllOpenings()
	
	response := models.OpeningSearchResponse{
		Results: openings,
		Count:   len(openings),
	}

	c.JSON(http.StatusOK, response)
}

// GetECOCategories returns openings grouped by ECO categories
// GET /api/openings/categories
func (h *OpeningHandler) GetECOCategories(c *gin.Context) {
	categories := h.openingService.GetECOCategories()
	
	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"total": len(categories),
	})
} 