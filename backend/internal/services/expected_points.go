package services

import (
	"math"
)

// ExpectedPointsService handles Expected Points calculations based on engine evaluations and player ratings
type ExpectedPointsService struct{}

// NewExpectedPointsService creates a new Expected Points service
func NewExpectedPointsService() *ExpectedPointsService {
	return &ExpectedPointsService{}
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
func (eps *ExpectedPointsService) getRatingFactor(rating int) float64 {
	// Base rating where factor = 1.0 (around intermediate level)
	baseRating := 1200.0
	
	// Rating effect scaling - each 200 rating points changes conversion by ~20%
	ratingScale := 400.0
	
	// Higher rating = higher factor (better conversion)
	// Lower rating = lower factor (worse conversion)
	factor := 1.0 + (float64(rating)-baseRating)/ratingScale
	
	// Clamp the factor to reasonable bounds
	// Beginners (800): ~0.8, Masters (2000): ~1.4, Super-GMs (2700+): ~1.75
	if factor < 0.5 {
		factor = 0.5
	}
	if factor > 2.0 {
		factor = 2.0
	}
	
	return factor
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
func (eps *ExpectedPointsService) GetAccuracyThresholds() map[string]float64 {
	return map[string]float64{
		"brilliant":   0.005, // Almost no EP loss for brilliant moves
		"great":       0.01,  // Very small EP loss
		"best":        0.00,  // No EP loss (best move)
		"excellent":   0.02,  // Small EP loss
		"good":        0.05,  // Minor EP loss
		"inaccuracy":  0.10,  // Noticeable EP loss
		"mistake":     0.20,  // Significant EP loss
		"blunder":     0.40,  // Very large EP loss
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