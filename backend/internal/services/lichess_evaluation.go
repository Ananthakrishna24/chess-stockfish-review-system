package services

import (
	"math"
	"chess-backend/internal/models"
	"github.com/sirupsen/logrus"
)

// LichessEvaluationService implements the precise Lichess evaluation algorithms
// as described in their position evaluation system architecture
type LichessEvaluationService struct {
	// Lichess empirical constants derived from 75,000+ games of 2300+ rated players
	sigmoidCoefficient    float64 // -0.00368208 (Lichess research constant)
	maxCentipawns         int     // ±1000 centipawn cap for display normalization
	mateScoreThreshold    int     // Threshold for mate score detection (deprecated - now handled via mate conversion)
	smoothingEnabled      bool    // Enable smoothing algorithms
	cacheEnabled          bool    // Enable evaluation caching
	
	// Accuracy calculation constants
	accuracyBase          float64 // 103.1668 (Lichess accuracy formula base)
	accuracyExponent      float64 // -0.04354 (Lichess accuracy formula exponent)
	accuracyOffset        float64 // -3.1669 (Lichess accuracy formula offset)
}

// NewLichessEvaluationService creates a new Lichess evaluation service with exact parameters
func NewLichessEvaluationService() *LichessEvaluationService {
	return &LichessEvaluationService{
		// Core evaluation constants (exact Lichess values)
		sigmoidCoefficient: -0.00368208, // Derived from logistic regression of 75k+ games
		maxCentipawns:      1000,        // Cap at ±10 pawns for display consistency
		mateScoreThreshold: 3000,        // Threshold for mate detection
		smoothingEnabled:   true,        // Enable smoothing by default
		cacheEnabled:       true,        // Enable caching for performance
		
		// Accuracy calculation constants (exact Lichess formula)
		accuracyBase:     103.1668, // Empirically calibrated base
		accuracyExponent: -0.04354, // Exponential decay factor
		accuracyOffset:   -3.1669,  // Calibration offset
	}
}

// ConvertCentipawnsToWinProbability implements the exact Lichess win probability formula
// Formula: Win% = 50 + 50 * (2 / (1 + exp(-0.00368208 * centipawns)) - 1)
func (les *LichessEvaluationService) ConvertCentipawnsToWinProbability(centipawns int) float64 {
	// Step 1: Apply "ceiled" evaluation (cap extreme evaluations at ±1000 centipawns)
	// This happens BEFORE any other processing, as per Lichess WinPercent.scala
	ceiledCentipawns := les.ceilEvaluation(centipawns)
	
	// Step 2: Apply exact Lichess sigmoid formula to ceiled value
	// Win% = 50 + 50 * (2 / (1 + exp(-0.00368208 * centipawns)) - 1)
	cp := float64(ceiledCentipawns)
	
	// Prevent numerical overflow in exponential
	exponentInput := les.sigmoidCoefficient * cp
	if exponentInput > 700 { // e^700 is near float64 max
		return 0.999
	}
	if exponentInput < -700 {
		return 0.001
	}
	
	// Calculate the exact Lichess formula:
	// Win% = 50 + 50 * (2 / (1 + exp(-K * cp)) - 1)
	expTerm := math.Exp(exponentInput)
	innerTerm := 2.0 / (1.0 + expTerm) - 1.0
	winPercentage := 50.0 + 50.0*innerTerm
	
	// Convert percentage to probability (0.0 to 1.0)
	winProbability := winPercentage / 100.0
	
	// Ensure bounds
	if winProbability > 0.999 {
		winProbability = 0.999
	}
	if winProbability < 0.001 {
		winProbability = 0.001
	}
	
	logrus.Debugf("Lichess evaluation: %dcp -> %.1f%% -> %.4f probability", 
		centipawns, winPercentage, winProbability)
	
	return winProbability
}

// CalculateAccuracy implements the exact Lichess accuracy formula
// Formula: Accuracy% = 103.1668 * exp(-0.04354 * (winPercentBefore - winPercentAfter)) - 3.1669
func (les *LichessEvaluationService) CalculateAccuracy(winPercentBefore, winPercentAfter float64) float64 {
	// Convert probabilities to percentages for calculation
	percentBefore := winPercentBefore * 100.0
	percentAfter := winPercentAfter * 100.0
	
	// Calculate win percentage change (loss if positive)
	winPercentChange := percentBefore - percentAfter
	
	// Apply Lichess accuracy formula
	// Accuracy% = 103.1668 * exp(-0.04354 * change) - 3.1669
	exponentInput := les.accuracyExponent * winPercentChange
	
	// Prevent numerical issues
	if exponentInput > 700 {
		return 0.0 // Very large loss = 0% accuracy
	}
	if exponentInput < -700 {
		return 100.0 // Very large gain = 100% accuracy (capped)
	}
	
	accuracy := les.accuracyBase*math.Exp(exponentInput) + les.accuracyOffset
	
	// Clamp to 0-100 range
	if accuracy < 0 {
		accuracy = 0
	}
	if accuracy > 100 {
		accuracy = 100
	}
	
	logrus.Debugf("Lichess accuracy: %.1f%% -> %.1f%% = %.1f%% accuracy", 
		percentBefore, percentAfter, accuracy)
	
	return accuracy
}

// CreateDisplayEvaluationFromEngine creates a display evaluation from EngineEvaluation using exact Lichess algorithms
// This is the preferred method as it handles both centipawns and mate scores correctly
func (les *LichessEvaluationService) CreateDisplayEvaluationFromEngine(
	engineEval models.EngineEvaluation,
	isWhiteToMove bool,
	previousDisplay *models.DisplayEvaluation,
) *models.DisplayEvaluation {
	// Convert engine evaluation to centipawns using Lichess mate conversion if needed
	rawCentipawns := les.ConvertEngineEvaluationToCentipawns(engineEval)
	
	// Use the standard pipeline
	return les.CreateDisplayEvaluation(rawCentipawns, isWhiteToMove, previousDisplay)
}

// CreateDisplayEvaluation creates a display evaluation using Lichess algorithms
// Note: This method assumes rawCentipawns is already converted from mate scores if applicable
func (les *LichessEvaluationService) CreateDisplayEvaluation(
	rawCentipawns int, 
	isWhiteToMove bool, 
	previousDisplay *models.DisplayEvaluation,
) *models.DisplayEvaluation {
	
	// Step 1: Convert to win probability using Lichess formula (from white's perspective)
	// Note: Mate scores should already be converted to centipawns via ConvertEngineEvaluationToCentipawns
	baseWinProb := les.ConvertCentipawnsToWinProbability(rawCentipawns)
	
	// Step 2: Apply smoothing if enabled and previous evaluation exists
	smoothedWinProb := baseWinProb
	if les.smoothingEnabled && previousDisplay != nil {
		smoothedWinProb = les.applySmoothingTransition(baseWinProb, previousDisplay.WinProbability)
	}
	
	// Step 3: Create evaluation bar value using improved visual mapping
	evalBar := les.convertToEvaluationBar(smoothedWinProb)
	
	// Step 4: Determine position assessment using Lichess thresholds
	assessment := les.assessPosition(smoothedWinProb)
	
	// Step 5: Check evaluation stability
	isStable := les.isEvaluationStable(smoothedWinProb, previousDisplay)
	
	// Step 6: Apply more aggressive smoothing to display score
	displayScore := rawCentipawns
	if les.smoothingEnabled && previousDisplay != nil {
		// Much more aggressive smoothing for display score
		prevScore := float64(previousDisplay.DisplayScore)
		currentScore := float64(rawCentipawns)
		smoothFactor := 0.15 // Only 15% of new value, 85% of old (very conservative)
		
		// Less smoothing only for very large changes (major tactics)
		change := math.Abs(currentScore - prevScore)
		if change > 200 { // Very large change - reduce smoothing
			smoothFactor = 0.35
		} else if change > 50 { // Medium change - moderate smoothing  
			smoothFactor = 0.25
		} else { // Small changes - maximum smoothing
			smoothFactor = 0.1
		}
		
		smoothedScore := smoothFactor*currentScore + (1-smoothFactor)*prevScore
		displayScore = int(math.Round(smoothedScore))
	}
	
	// Cap the display score
	cappedScore := les.capEvaluation(displayScore)
	
	return &models.DisplayEvaluation{
		WinProbability:     smoothedWinProb,
		DisplayScore:       cappedScore,
		EvaluationBar:      evalBar,
		PositionAssessment: assessment,
		IsStable:           isStable,
	}
}

// ProcessEvaluationHistory processes a sequence of evaluations with Lichess smoothing
func (les *LichessEvaluationService) ProcessEvaluationHistory(
	rawEvaluations []int, 
	isWhiteToMove []bool,
) []*models.DisplayEvaluation {
	
	if len(rawEvaluations) != len(isWhiteToMove) {
		logrus.Errorf("Evaluation arrays length mismatch: %d vs %d", 
			len(rawEvaluations), len(isWhiteToMove))
		return nil
	}
	
	displayEvaluations := make([]*models.DisplayEvaluation, len(rawEvaluations))
	
	// Process each evaluation - all evaluations should be from white's perspective
	for i, rawEval := range rawEvaluations {
		var previous *models.DisplayEvaluation
		if i > 0 {
			previous = displayEvaluations[i-1]
		}
		
		// All evaluations are treated as white's perspective
		// The perspective handling is done at the frontend level for display
		displayEvaluations[i] = les.CreateDisplayEvaluation(
			rawEval, 
			true, // Always treat as white's perspective for consistency
			previous,
		)
	}
	
	// Apply windowing system smoothing if enabled
	if les.smoothingEnabled && len(displayEvaluations) > 4 {
		displayEvaluations = les.applyWindowingSmoothing(displayEvaluations)
	}
	
	return displayEvaluations
}

// CalculateGameAccuracy calculates overall game accuracy using Lichess method
func (les *LichessEvaluationService) CalculateGameAccuracy(
	evaluationHistory []*models.DisplayEvaluation,
	isWhitePlayer bool,
) float64 {
	
	if len(evaluationHistory) < 2 {
		return 100.0 // No moves to evaluate
	}
	
	totalAccuracy := 0.0
	moveCount := 0
	
	// Calculate accuracy for each move (skip first position)
	for i := 1; i < len(evaluationHistory); i++ {
		// Only evaluate this player's moves
		isPlayerMove := (i%2 == 1) == isWhitePlayer
		if !isPlayerMove {
			continue
		}
		
		beforeWinProb := evaluationHistory[i-1].WinProbability
		afterWinProb := evaluationHistory[i].WinProbability
		
		// Adjust perspective for current player
		if !isWhitePlayer {
			beforeWinProb = 1.0 - beforeWinProb
			afterWinProb = 1.0 - afterWinProb
		}
		
		moveAccuracy := les.CalculateAccuracy(beforeWinProb, afterWinProb)
		totalAccuracy += moveAccuracy
		moveCount++
	}
	
	if moveCount == 0 {
		return 100.0
	}
	
	overallAccuracy := totalAccuracy / float64(moveCount)
	logrus.Debugf("Game accuracy calculated: %.1f%% over %d moves", 
		overallAccuracy, moveCount)
	
	return overallAccuracy
}

// ceilEvaluation caps extreme evaluations at ±1000 centipawns (Lichess "ceiled" operation)
// This is the exact equivalent of cp.ceiled in Lichess WinPercent.scala
func (les *LichessEvaluationService) ceilEvaluation(centipawns int) int {
	if centipawns > les.maxCentipawns {
		return les.maxCentipawns
	}
	if centipawns < -les.maxCentipawns {
		return -les.maxCentipawns
	}
	return centipawns
}

// capEvaluation caps extreme evaluations at ±1000 centipawns (legacy method)
func (les *LichessEvaluationService) capEvaluation(centipawns int) int {
	return les.ceilEvaluation(centipawns)
}

// IsMateScore detects if evaluation represents a mate score (public method)
func (les *LichessEvaluationService) IsMateScore(centipawns int) bool {
	absEval := centipawns
	if absEval < 0 {
		absEval = -absEval
	}
	return absEval >= les.mateScoreThreshold
}

// isMateScore detects if evaluation represents a mate score (private method)
func (les *LichessEvaluationService) isMateScore(centipawns int) bool {
	return les.IsMateScore(centipawns)
}

// convertMateToEquivalentCentipawns converts mate-in-N to equivalent centipawns using Lichess formula
// Formula: cp = 100*(21 - min(10, N)) where N is mate distance
func (les *LichessEvaluationService) convertMateToEquivalentCentipawns(mateInN int) int {
	// Extract absolute mate distance
	mateDistance := mateInN
	if mateDistance < 0 {
		mateDistance = -mateDistance
	}
	
	// Apply Lichess mate conversion formula: cp = 100*(21 - min(10, N))
	// Cap mate distance at 10 moves as per Lichess specification
	clampedDistance := mateDistance
	if clampedDistance > 10 {
		clampedDistance = 10
	}
	
	equivalentCp := 100 * (21 - clampedDistance)
	
	// Restore sign: positive for mate-for, negative for mate-against
	if mateInN < 0 {
		equivalentCp = -equivalentCp
	}
	
	logrus.Debugf("Converted mate %d to %d centipawns using Lichess formula", mateInN, equivalentCp)
	return equivalentCp
}

// ConvertEngineEvaluationToCentipawns converts an EngineEvaluation to centipawns using Lichess algorithms
// This handles both regular centipawn scores and mate scores according to Lichess specification
func (les *LichessEvaluationService) ConvertEngineEvaluationToCentipawns(eval models.EngineEvaluation) int {
	// If it's a mate score, convert using Lichess mate formula
	if eval.Mate != nil {
		return les.convertMateToEquivalentCentipawns(*eval.Mate)
	}
	
	// Otherwise, use the regular centipawn score
	return eval.Score
}

// CalculateStandardDeviation calculates the standard deviation of win probabilities
// Used for Lichess volatility weighting in windowing system
func (les *LichessEvaluationService) CalculateStandardDeviation(winProbs []float64) float64 {
	if len(winProbs) <= 1 {
		return 0.5 // Minimum volatility as per Lichess bounds
	}
	
	// Calculate mean
	mean := 0.0
	for _, prob := range winProbs {
		mean += prob
	}
	mean /= float64(len(winProbs))
	
	// Calculate variance
	variance := 0.0
	for _, prob := range winProbs {
		diff := prob - mean
		variance += diff * diff
	}
	variance /= float64(len(winProbs))
	
	// Calculate standard deviation and apply Lichess bounds
	stdDev := math.Sqrt(variance)
	
	// Apply Lichess bounds: atLeast(0.5).atMost(12)
	if stdDev < 0.5 {
		stdDev = 0.5
	}
	if stdDev > 12.0 {
		stdDev = 12.0
	}
	
	return stdDev
}

// applySmoothingTransition applies exponential smoothing between evaluations
func (les *LichessEvaluationService) applySmoothingTransition(newWinProb, oldWinProb float64) float64 {
	// Use much more aggressive smoothing to reduce volatility
	change := math.Abs(newWinProb - oldWinProb)
	
	// Very aggressive smoothing for small changes (opening moves, minor positional changes)
	smoothingFactor := 0.15 // Default 15% smoothing (very aggressive)
	
	if change > 0.15 {       // Major swing - preserve important tactical changes
		smoothingFactor = 0.4
	} else if change > 0.05 { // Medium changes - moderate smoothing
		smoothingFactor = 0.25
	} else {                 // Small changes - maximum smoothing
		smoothingFactor = 0.1  // Only 10% of new value for tiny changes
	}
	
	return smoothingFactor*newWinProb + (1-smoothingFactor)*oldWinProb
}

// convertToEvaluationBar converts win probability to evaluation bar value (-1 to +1)
func (les *LichessEvaluationService) convertToEvaluationBar(winProb float64) float64 {
	// Ensure exactly 0.0 for equal positions
	if math.Abs(winProb-0.5) < 0.001 {
		return 0.0
	}
	
	// Apply much more conservative scaling - only significant advantages should be visible
	if winProb > 0.5 {
		// Winning side: use quartic root for very conservative visual scaling
		advantage := (winProb - 0.5) * 2.0 // 0.0 to 1.0
		// Apply threshold - only advantages >5% win probability get visual representation
		if advantage < 0.1 { // Less than 55% win probability
			return advantage * 0.3 // Very minimal visual change
		}
		return math.Pow(advantage, 0.6) * 0.6  // Quartic-ish root, scaled to max 0.6
	} else {
		// Losing side
		disadvantage := (0.5 - winProb) * 2.0 // 0.0 to 1.0
		if disadvantage < 0.1 { // Less than 45% win probability 
			return -disadvantage * 0.3 // Very minimal visual change
		}
		return -math.Pow(disadvantage, 0.6) * 0.6 // Negative quartic-ish root
	}
}

// assessPosition provides human-readable position assessment
func (les *LichessEvaluationService) assessPosition(winProb float64) string {
	switch {
	case winProb >= 0.90:
		return "winning"
	case winProb >= 0.75:
		return "much_better"
	case winProb >= 0.60:
		return "slightly_better"
	case winProb >= 0.40 && winProb <= 0.60:
		return "equal"
	case winProb >= 0.25:
		return "slightly_worse"
	case winProb >= 0.10:
		return "much_worse"
	default:
		return "losing"
	}
}

// isEvaluationStable checks if evaluation has stabilized
func (les *LichessEvaluationService) isEvaluationStable(currentWinProb float64, previous *models.DisplayEvaluation) bool {
	if previous == nil {
		return false
	}
	
	change := math.Abs(currentWinProb - previous.WinProbability)
	return change < 0.05 // 5% threshold for stability
}

// applyWindowingSmoothing applies windowing system smoothing exactly as described in Lichess AccuracyPercent.scala
// Uses sliding window analysis with volatility weighting based on standard deviation
func (les *LichessEvaluationService) applyWindowingSmoothing(evaluations []*models.DisplayEvaluation) []*models.DisplayEvaluation {
	gameLength := len(evaluations)
	
	// Calculate dynamic window size: (cps.size / 10).atLeast(2).atMost(8)
	windowSize := gameLength / 10
	if windowSize < 2 {
		windowSize = 2
	}
	if windowSize > 8 {
		windowSize = 8
	}
	
	// Ensure we have enough positions for windowing
	if gameLength < windowSize*2+1 {
		logrus.Debugf("Game too short (%d moves) for windowing with size %d", gameLength, windowSize)
		return evaluations
	}
	
	// Extract win probabilities for processing
	allWinPercents := make([]float64, len(evaluations))
	for i, eval := range evaluations {
		allWinPercents[i] = eval.WinProbability
	}
	
	// Create windows and calculate volatility weights as per Lichess specification
	// val windows = List.fill(windowSize.atMost(allWinPercentValues.size) - 2)(allWinPercentValues.take(windowSize))
	maxWindows := len(allWinPercents) - 2
	if windowSize < maxWindows {
		maxWindows = windowSize
	}
	
	volatilityWeights := make([]float64, maxWindows)
	for i := 0; i < maxWindows; i++ {
		// Take windowSize elements starting from position i
		windowEnd := i + windowSize
		if windowEnd > len(allWinPercents) {
			windowEnd = len(allWinPercents)
		}
		
		windowValues := allWinPercents[i:windowEnd]
		volatilityWeights[i] = les.CalculateStandardDeviation(windowValues)
	}
	
	smoothed := make([]*models.DisplayEvaluation, len(evaluations))
	copy(smoothed, evaluations) // Start with original values
	
	// Apply volatility-weighted smoothing to positions that have sufficient context
	for i := windowSize; i < len(evaluations)-windowSize; i++ {
		smoothed[i] = les.calculateVolatilityWeightedAverage(evaluations, allWinPercents, i, windowSize, volatilityWeights)
	}
	
	logrus.Debugf("Applied Lichess windowing smoothing with window size %d and %d volatility weights", windowSize, len(volatilityWeights))
	return smoothed
}

// calculateVolatilityWeightedAverage calculates weighted average using Lichess volatility weighting
// This implements the exact Lichess algorithm with standard deviation-based weights
func (les *LichessEvaluationService) calculateVolatilityWeightedAverage(
	evaluations []*models.DisplayEvaluation,
	allWinPercents []float64,
	center int,
	windowSize int,
	volatilityWeights []float64,
) *models.DisplayEvaluation {
	
	// Calculate volatility-weighted average of win probabilities within window
	totalWeight := 0.0
	weightedSum := 0.0
	
	for i := center - windowSize; i <= center + windowSize; i++ {
		if i < 0 || i >= len(evaluations) {
			continue
		}
		
		// Use volatility weight based on standard deviation
		// Higher volatility = lower weight for smoothing
		volatilityIndex := i - windowSize
		if volatilityIndex < 0 {
			volatilityIndex = 0
		}
		if volatilityIndex >= len(volatilityWeights) {
			volatilityIndex = len(volatilityWeights) - 1
		}
		
		// Weight is inversely proportional to volatility (higher volatility = lower weight)
		volatility := volatilityWeights[volatilityIndex]
		weight := 1.0 / volatility // Higher volatility gets lower weight
		
		totalWeight += weight
		weightedSum += weight * evaluations[i].WinProbability
	}
	
	smoothedWinProb := weightedSum / totalWeight
	
	// Create new display evaluation with smoothed probability
	original := evaluations[center]
	return &models.DisplayEvaluation{
		WinProbability:     smoothedWinProb,
		DisplayScore:       original.DisplayScore,
		EvaluationBar:      les.convertToEvaluationBar(smoothedWinProb),
		PositionAssessment: les.assessPosition(smoothedWinProb),
		IsStable:           true, // Smoothed values are considered stable
	}
}

// GetLichessConstants returns the empirical constants used in calculations
func (les *LichessEvaluationService) GetLichessConstants() map[string]float64 {
	return map[string]float64{
		"sigmoidCoefficient": les.sigmoidCoefficient,
		"accuracyBase":       les.accuracyBase,
		"accuracyExponent":   les.accuracyExponent,
		"accuracyOffset":     les.accuracyOffset,
		"maxCentipawns":      float64(les.maxCentipawns),
	}
}

// EnableSmoothing enables or disables smoothing algorithms
func (les *LichessEvaluationService) EnableSmoothing(enabled bool) {
	les.smoothingEnabled = enabled
	logrus.Infof("Lichess evaluation smoothing %s", 
		map[bool]string{true: "enabled", false: "disabled"}[enabled])
}

// EnableCaching enables or disables evaluation caching
func (les *LichessEvaluationService) EnableCaching(enabled bool) {
	les.cacheEnabled = enabled
	logrus.Infof("Lichess evaluation caching %s", 
		map[bool]string{true: "enabled", false: "disabled"}[enabled])
} 