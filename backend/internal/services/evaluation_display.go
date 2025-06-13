package services

import (
	"math"
	"chess-backend/internal/models"
	"github.com/sirupsen/logrus"
)

// EvaluationDisplayService handles stable evaluation display using research-based normalization
type EvaluationDisplayService struct {
	// Lichess-style sigmoid parameters (based on real game data)
	sigmoidK          float64 // Scaling factor: 0.00368208 (Lichess research)
	maxDisplayEval    int     // Maximum centipawn value to display (cap extreme values)
	smoothingFactor   float64 // Smoothing factor for move-to-move transitions
	winThreshold      float64 // Win probability threshold for "winning" positions
	equalityThreshold float64 // Win probability threshold for "equal" positions
}

// NewEvaluationDisplayService creates a new evaluation display service
func NewEvaluationDisplayService() *EvaluationDisplayService {
	return &EvaluationDisplayService{
		sigmoidK:          0.00368208, // Lichess empirical constant
		maxDisplayEval:    1000,       // Cap at ±10 pawns for display
		smoothingFactor:   0.15,       // 15% smoothing between moves
		winThreshold:      0.75,       // 75% = clearly winning
		equalityThreshold: 0.45,       // 45-55% = roughly equal
	}
}

// NormalizeForDisplay converts raw centipawn evaluation to stable display format
func (eds *EvaluationDisplayService) NormalizeForDisplay(
	rawCentipawns int, 
	isWhiteToMove bool, 
	previousDisplay *models.DisplayEvaluation,
) *models.DisplayEvaluation {
	
	// Step 1: Cap extreme evaluations to prevent UI chaos
	cappedEval := eds.capEvaluation(rawCentipawns)
	
	// Step 2: Convert to win probability using Lichess-style sigmoid
	winProb := eds.centipawnsToWinProbability(cappedEval, isWhiteToMove)
	
	// Step 3: Apply smoothing if we have previous evaluation
	if previousDisplay != nil {
		winProb = eds.applySmoothingTransition(winProb, previousDisplay.WinProbability)
	}
	
	// Step 4: Create evaluation bar value (-1 to +1 for UI)
	evalBar := eds.winProbabilityToEvalBar(winProb)
	
	// Step 5: Determine position assessment
	assessment := eds.assessPosition(winProb)
	
	// Step 6: Check if evaluation is stable
	isStable := eds.isEvaluationStable(winProb, previousDisplay)
	
	logrus.Debugf("Evaluation normalized: %dcp -> %.3f win prob -> %.2f bar", 
		rawCentipawns, winProb, evalBar)
	
	return &models.DisplayEvaluation{
		WinProbability:     winProb,
		DisplayScore:       cappedEval,
		EvaluationBar:      evalBar,
		PositionAssessment: assessment,
		IsStable:           isStable,
	}
}

// centipawnsToWinProbability converts centipawns to win probability using Lichess formula
func (eds *EvaluationDisplayService) centipawnsToWinProbability(centipawns int, isWhiteToMove bool) float64 {
	// Normalize for current player perspective
	normalizedCP := float64(centipawns)
	if !isWhiteToMove {
		normalizedCP = -normalizedCP
	}
	
	// Lichess formula: Win% = 50 + 50 * (2 / (1 + exp(-K * centipawns)) - 1)
	// Simplified: Win% = 1 / (1 + exp(-K * centipawns))
	sigmoidInput := eds.sigmoidK * normalizedCP
	
	// Prevent overflow in exponential
	if sigmoidInput > 10 {
		return 0.9999
	}
	if sigmoidInput < -10 {
		return 0.0001
	}
	
	winProbability := 1.0 / (1.0 + math.Exp(-sigmoidInput))
	
	// Clamp to reasonable bounds
	if winProbability > 0.999 {
		winProbability = 0.999
	}
	if winProbability < 0.001 {
		winProbability = 0.001
	}
	
	return winProbability
}

// capEvaluation caps extreme evaluations to prevent UI volatility
func (eds *EvaluationDisplayService) capEvaluation(centipawns int) int {
	if centipawns > eds.maxDisplayEval {
		return eds.maxDisplayEval
	}
	if centipawns < -eds.maxDisplayEval {
		return -eds.maxDisplayEval
	}
	return centipawns
}

// applySmoothingTransition smooths evaluation changes between moves
func (eds *EvaluationDisplayService) applySmoothingTransition(newWinProb, oldWinProb float64) float64 {
	// Apply exponential smoothing: new = α * current + (1-α) * previous
	smoothed := eds.smoothingFactor*newWinProb + (1-eds.smoothingFactor)*oldWinProb
	
	// But don't smooth dramatic changes (blunders, tactics)
	diff := math.Abs(newWinProb - oldWinProb)
	if diff > 0.3 { // Major evaluation swing - don't smooth as much
		smoothed = 0.7*newWinProb + 0.3*oldWinProb
	}
	
	return smoothed
}

// winProbabilityToEvalBar converts win probability to evaluation bar value
func (eds *EvaluationDisplayService) winProbabilityToEvalBar(winProb float64) float64 {
	// Convert 0-1 probability to -1 to +1 evaluation bar
	// 0.5 (equal) -> 0.0, 1.0 (winning) -> +1.0, 0.0 (losing) -> -1.0
	
	// Use a non-linear transformation for better visual representation
	if winProb > 0.5 {
		// Winning side: compress high probabilities
		advantage := (winProb - 0.5) * 2.0 // 0.0 to 1.0
		return math.Sqrt(advantage)        // Square root for compression
	} else {
		// Losing side: compress low probabilities
		disadvantage := (0.5 - winProb) * 2.0 // 0.0 to 1.0
		return -math.Sqrt(disadvantage)       // Negative square root
	}
}

// assessPosition provides human-readable position assessment
func (eds *EvaluationDisplayService) assessPosition(winProb float64) string {
	switch {
	case winProb >= 0.90:
		return "winning"
	case winProb >= eds.winThreshold: // 0.75
		return "much_better"
	case winProb >= 0.60:
		return "slightly_better"
	case winProb >= eds.equalityThreshold && winProb <= (1.0-eds.equalityThreshold): // 0.45-0.55
		return "equal"
	case winProb >= 0.40:
		return "slightly_worse"
	case winProb >= 0.25:
		return "much_worse"
	default:
		return "losing"
	}
}

// isEvaluationStable checks if the evaluation has stabilized
func (eds *EvaluationDisplayService) isEvaluationStable(currentWinProb float64, previous *models.DisplayEvaluation) bool {
	if previous == nil {
		return false
	}
	
	// Consider stable if change is less than 5%
	change := math.Abs(currentWinProb - previous.WinProbability)
	return change < 0.05
}

// GetEvaluationTrend analyzes evaluation trend over multiple moves
func (eds *EvaluationDisplayService) GetEvaluationTrend(evaluations []*models.DisplayEvaluation) string {
	if len(evaluations) < 3 {
		return "insufficient_data"
	}
	
	recent := evaluations[len(evaluations)-3:]
	
	// Check if generally improving, declining, or stable
	trend := 0.0
	for i := 1; i < len(recent); i++ {
		trend += recent[i].WinProbability - recent[i-1].WinProbability
	}
	
	switch {
	case trend > 0.1:
		return "improving"
	case trend < -0.1:
		return "declining"
	default:
		return "stable"
	}
}

// CreateEvaluationHistory creates a stable evaluation history for game analysis
func (eds *EvaluationDisplayService) CreateEvaluationHistory(rawEvaluations []int, isWhiteToMove []bool) []*models.DisplayEvaluation {
	if len(rawEvaluations) != len(isWhiteToMove) {
		logrus.Error("Mismatched evaluation and color arrays")
		return nil
	}
	
	history := make([]*models.DisplayEvaluation, len(rawEvaluations))
	var previous *models.DisplayEvaluation
	
	for i, rawEval := range rawEvaluations {
		display := eds.NormalizeForDisplay(rawEval, isWhiteToMove[i], previous)
		history[i] = display
		previous = display
	}
	
	return history
}

// GetRecommendedDisplaySettings returns optimal display settings for different use cases
func (eds *EvaluationDisplayService) GetRecommendedDisplaySettings(useCase string) map[string]interface{} {
	settings := make(map[string]interface{})
	
	switch useCase {
	case "live_game":
		settings["smoothing_factor"] = 0.25 // More smoothing for live games
		settings["update_frequency"] = 500  // Update every 500ms
		settings["show_raw_centipawns"] = false
		
	case "game_analysis":
		settings["smoothing_factor"] = 0.10 // Less smoothing for analysis
		settings["update_frequency"] = 100  // Immediate updates
		settings["show_raw_centipawns"] = true
		
	case "beginner_mode":
		settings["smoothing_factor"] = 0.35 // Heavy smoothing
		settings["update_frequency"] = 1000 // Slow updates
		settings["show_raw_centipawns"] = false
		settings["simplified_labels"] = true
		
	default:
		settings["smoothing_factor"] = eds.smoothingFactor
		settings["update_frequency"] = 250
		settings["show_raw_centipawns"] = true
	}
	
	return settings
} 