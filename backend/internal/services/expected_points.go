package services

import (
	"chess-backend/internal/models"
	"math"
)

// ExpectedPointsService handles Expected Points calculations based on engine evaluations and player ratings
type ExpectedPointsService struct {
	calibrationService *CalibrationService
	thresholds         map[models.RatingBucket]models.EPThresholds
}

// NewExpectedPointsService creates a new Expected Points service
func NewExpectedPointsService() *ExpectedPointsService {
	return &ExpectedPointsService{
		thresholds: make(map[models.RatingBucket]models.EPThresholds),
	}
}

// SetCalibrationService sets the calibration service for dynamic thresholds
func (eps *ExpectedPointsService) SetCalibrationService(cs *CalibrationService) {
	eps.calibrationService = cs
	// Load thresholds from calibration service
	if thresholds, err := cs.LoadThresholds(); err == nil {
		eps.thresholds = thresholds
	}
}

// CalculateExpectedPoints calculates win probability (0-1) based on engine evaluation and player rating
// Uses a sigmoid function: EP = 1 / (1 + exp(-evaluation * rating_factor))
func (eps *ExpectedPointsService) CalculateExpectedPoints(evaluationCentipawns int, playerRating int) float64 {
	// Handle mate scores
	if evaluationCentipawns >= 3000 {
		return 0.999 // Essentially winning
	}
	if evaluationCentipawns <= -3000 {
		return 0.001 // Essentially losing
	}
	
	// Rating factor: stronger players convert advantages better
	// Base rating factor around 1200 (intermediate level)
	// Higher rating = better conversion, lower rating = worse conversion
	ratingFactor := eps.getRatingFactor(playerRating)
	
	// Convert centipawns to a normalized evaluation
	// 100 centipawns = 1.0 evaluation unit
	normalizedEval := float64(evaluationCentipawns) / 100.0
	
	// Apply rating factor - stronger players need smaller advantages to win
	adjustedEval := normalizedEval * ratingFactor
	
	// Sigmoid function: 1 / (1 + exp(-x))
	// This gives us win probability between 0 and 1
	expectedPoints := 1.0 / (1.0 + math.Exp(-adjustedEval))
	
	return expectedPoints
}

// getRatingFactor calculates how well a player converts advantages based on their rating
// Updated to match the EP algorithm specification
func (eps *ExpectedPointsService) getRatingFactor(rating int) float64 {
	// Linear or logistic mapping from rating→[0.5–2.0] as per algorithm spec
	// clamp((rating - 1200)/2000 + 1.0, 0.5, 2.0)
	f := (float64(rating)-1200.0)/2000.0 + 1.0
	
	// Clamp to algorithm-specified bounds
	if f < 0.5 {
		f = 0.5
	}
	if f > 2.0 {
		f = 2.0
	}
	
	return f
}

// CalculateExpectedPointsLoss calculates the EP loss from one position to another
func (eps *ExpectedPointsService) CalculateExpectedPointsLoss(beforeEval, afterEval int, playerRating int) float64 {
	epBefore := eps.CalculateExpectedPoints(beforeEval, playerRating)
	epAfter := eps.CalculateExpectedPoints(afterEval, playerRating)
	
	// EP loss is the difference (positive = loss, negative = gain)
	return epBefore - epAfter
}

// CalculateMoveAccuracy calculates accuracy for a single move based on EP loss
func (eps *ExpectedPointsService) CalculateMoveAccuracy(epLoss float64) float64 {
	// Accuracy = (1 - ep_loss) * 100
	// Clamp to 0-100 range
	accuracy := (1.0 - epLoss) * 100.0
	
	if accuracy < 0 {
		accuracy = 0
	}
	if accuracy > 100 {
		accuracy = 100
	}
	
	return accuracy
}

// GetAccuracyThresholds returns EP loss thresholds for move classifications
// Updated to match Chess.com's exact thresholds
func (eps *ExpectedPointsService) GetAccuracyThresholds() map[string]float64 {
	return map[string]float64{
		"brilliant":   0.005, // Almost no EP loss for brilliant moves (requires sacrifice)
		"great":       0.01,  // Very small EP loss for critical moves
		"best":        0.00,  // No EP loss (best move)
		"excellent":   0.02,  // Chess.com: ≤0.02 (tiny losses)
		"good":        0.05,  // Chess.com: 0.02–0.05 lost
		"inaccuracy":  0.10,  // Chess.com: 0.05–0.10 lost
		"mistake":     0.20,  // Chess.com: 0.10–0.20 lost
		"blunder":     1.00,  // Chess.com: >0.20 lost (effectively no upper limit)
	}
}

// GetChessComThresholds returns Chess.com's exact classification thresholds
func (eps *ExpectedPointsService) GetChessComThresholds(playerRating int) map[string]float64 {
	// Chess.com's published thresholds
	base := map[string]float64{
		"excellent":   0.02,  // ≤0.02 (tiny losses up to 0.02)
		"good":        0.05,  // 0.02–0.05 expected-points lost
		"inaccuracy":  0.10,  // 0.05–0.10 lost
		"mistake":     0.20,  // 0.10–0.20 lost
		"blunder":     1.00,  // >0.20 lost (no upper limit)
	}
	
	// Apply rating-based scaling
	ratingFactor := eps.getChessComRatingFactor(playerRating)
	
	scaled := make(map[string]float64)
	for key, value := range base {
		if key == "blunder" {
			scaled[key] = value // Blunder threshold doesn't scale
		} else {
			scaled[key] = value * ratingFactor
		}
	}
	
	return scaled
}

// getChessComRatingFactor returns rating scaling factor for Chess.com thresholds
func (eps *ExpectedPointsService) getChessComRatingFactor(rating int) float64 {
	// Chess.com scales thresholds by rating
	// Higher rated players have stricter standards
	// Rating 1200 = 1.0 (baseline)
	// Rating 2400 = 0.8 (20% stricter)
	// Rating 800 = 1.2 (20% more lenient)
	
	baseFactor := 1.0
	ratingDiff := float64(rating - 1200)
	scalingRate := -0.0001 // -0.01% per rating point
	
	factor := baseFactor + (ratingDiff * scalingRate)
	
	// Clamp between 0.6 and 1.4 to prevent extreme values
	if factor < 0.6 {
		factor = 0.6
	}
	if factor > 1.4 {
		factor = 1.4
	}
	
	return factor
}

// GetDynamicThresholds returns rating-specific EP loss thresholds
func (eps *ExpectedPointsService) GetDynamicThresholds(rating int) models.EPThresholds {
	bucket := eps.getRatingBucket(rating)
	
	if thresholds, exists := eps.thresholds[bucket]; exists {
		return thresholds
	}
	
	// Fallback to default thresholds if not calibrated
	return eps.getDefaultThresholds()[bucket]
}

// getRatingBucket determines which rating bucket a player belongs to
func (eps *ExpectedPointsService) getRatingBucket(rating int) models.RatingBucket {
	switch {
	case rating >= 2001:
		return models.RatingBucket2001Plus
	case rating >= 1601:
		return models.RatingBucket1601to2000
	case rating >= 1201:
		return models.RatingBucket1201to1600
	default:
		return models.RatingBucket800to1200
	}
}

// getDefaultThresholds provides fallback thresholds when calibration data is not available
func (eps *ExpectedPointsService) getDefaultThresholds() map[models.RatingBucket]models.EPThresholds {
	return map[models.RatingBucket]models.EPThresholds{
		models.RatingBucket800to1200: {
			P1: 0.002, P5: 0.008, P10: 0.015, P25: 0.040, P50: 0.080, P75: 0.150, P90: 0.250,
		},
		models.RatingBucket1201to1600: {
			P1: 0.001, P5: 0.005, P10: 0.012, P25: 0.030, P50: 0.060, P75: 0.120, P90: 0.200,
		},
		models.RatingBucket1601to2000: {
			P1: 0.001, P5: 0.003, P10: 0.008, P25: 0.020, P50: 0.045, P75: 0.090, P90: 0.150,
		},
		models.RatingBucket2001Plus: {
			P1: 0.000, P5: 0.002, P10: 0.005, P25: 0.015, P50: 0.035, P75: 0.070, P90: 0.120,
		},
	}
}

// NormalizeEvaluationForPlayer normalizes evaluation from the current player's perspective
// Stockfish evaluations are from White's perspective, so we need to flip for Black
func (eps *ExpectedPointsService) NormalizeEvaluationForPlayer(evaluation int, isWhiteToMove bool) int {
	if isWhiteToMove {
		return evaluation
	}
	// Flip evaluation for Black's perspective
	return -evaluation
}

// CalculateWinProbabilityChange calculates the change in win probability between two positions
func (eps *ExpectedPointsService) CalculateWinProbabilityChange(beforeEval, afterEval int, playerRating int, isWhiteToMove bool) float64 {
	// Normalize evaluations for current player
	normalizedBefore := eps.NormalizeEvaluationForPlayer(beforeEval, isWhiteToMove)
	normalizedAfter := eps.NormalizeEvaluationForPlayer(afterEval, isWhiteToMove)
	
	// Calculate EP change
	return eps.CalculateExpectedPointsLoss(normalizedBefore, normalizedAfter, playerRating)
} 