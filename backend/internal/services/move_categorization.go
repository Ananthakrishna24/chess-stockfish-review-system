package services

import (
	"chess-backend/internal/models"
	"math"
	"strings"

	"github.com/notnil/chess"
)

// MoveCategorizer implements sophisticated move categorization based on Expected Points
type MoveCategorizer struct {
	expectedPointsService *ExpectedPointsService
	chessService          *ChessService
	openingService        *OpeningService
}

// NewMoveCategorizer creates a new move categorization service
func NewMoveCategorizer(eps *ExpectedPointsService, chess *ChessService, opening *OpeningService) *MoveCategorizer {
	return &MoveCategorizer{
		expectedPointsService: eps,
		chessService:          chess,
		openingService:        opening,
	}
}

// CategorizeMoveAdvanced categorizes a move using the sophisticated EP-based algorithm
func (mc *MoveCategorizer) CategorizeMoveAdvanced(moveData MoveCategoryData) models.MoveClassification {
	// Priority order as per algorithm specification
	
	// 1. Check for Book Move first (first ~10-15 moves from opening database)
	if mc.isBookMove(moveData) {
		return models.Book
	}
	
	// 2. Check for Brilliant Move (!!)
	if mc.isBrilliantMove(moveData) {
		return models.Brilliant
	}
	
	// 3. Check for Great Move (!)
	if mc.isGreatMove(moveData) {
		return models.Great
	}
	
	// 4. Check for Best Move
	if mc.isBestMove(moveData) {
		return models.Best
	}
	
	// 5-8. Classify based on EP loss thresholds
	return mc.classifyByAccuracy(moveData.EPLoss)
}

// MoveCategoryData contains all data needed for move categorization
type MoveCategoryData struct {
	MoveNumber        int
	Move              string
	SAN               string
	UCI               string
	IsWhiteToMove     bool
	BeforeEvaluation  *models.EngineEvaluation
	AfterEvaluation   *models.EngineEvaluation
	BestMove          string
	EPLoss            float64
	MoveAccuracy      float64
	MaterialBefore    models.MaterialValue
	MaterialAfter     models.MaterialValue
	Position          *chess.Position
	AlternativeMoves  []models.AlternativeMove
	OpeningName       string
	ECO               string
	IsPositionWinning bool // Whether position was already overwhelmingly winning
}

// isBookMove checks if the move is from opening theory
func (mc *MoveCategorizer) isBookMove(data MoveCategoryData) bool {
	// Check if we're in the opening phase (first 10-15 moves)
	if data.MoveNumber > 15 {
		return false
	}
	
	// Use opening service to check if this is a known theoretical move
	if mc.openingService != nil {
		// This would query the opening database
		// For now, we'll use a simple heuristic based on ECO codes and common openings
		return data.ECO != "" && data.MoveNumber <= 12
	}
	
	// Fallback: consider first 8 moves as potentially book moves if they're good
	return data.MoveNumber <= 8 && data.EPLoss <= 0.03
}

// isBrilliantMove checks for brilliant moves (!!)
func (mc *MoveCategorizer) isBrilliantMove(data MoveCategoryData) bool {
	// Brilliant move criteria:
	// 1. Very low EP loss (objectively strong)
	// 2. Involves a sound piece sacrifice
	// 3. Position wasn't already overwhelmingly winning
	
	// Must be a very good move first
	if data.EPLoss > 0.005 {
		return false
	}
	
	// Position shouldn't be already winning by more than 300 centipawns
	if data.IsPositionWinning {
		return false
	}
	
	// Check for piece sacrifice
	if !mc.involvesSacrifice(data.MaterialBefore, data.MaterialAfter) {
		return false
	}
	
	// The sacrifice must be sound (low EP loss proves it's objectively good)
	// Additional check: the move should be difficult to find (not the obvious best move in a simple position)
	return mc.isDifficultMove(data)
}

// isGreatMove checks for great moves (!)
func (mc *MoveCategorizer) isGreatMove(data MoveCategoryData) bool {
	// Great move: the ONLY move that doesn't significantly worsen the position
	// This requires analyzing alternative moves
	
	if data.EPLoss > 0.01 {
		return false
	}
	
	// Count how many alternative moves would have been much worse
	betterAlternatives := 0
	for _, alt := range data.AlternativeMoves {
		// Calculate EP loss for this alternative
		altEPLoss := mc.calculateAlternativeEPLoss(alt, data)
		if altEPLoss < data.EPLoss+0.05 { // Allow small margin
			betterAlternatives++
		}
	}
	
	// If there are very few good alternatives, this might be a great move
	return betterAlternatives <= 1 && len(data.AlternativeMoves) > 3
}

// isBestMove checks if the played move matches the engine's best move
func (mc *MoveCategorizer) isBestMove(data MoveCategoryData) bool {
	// Direct comparison with engine's recommended move
	return strings.EqualFold(data.UCI, data.BestMove) || strings.EqualFold(data.Move, data.BestMove)
}

// classifyByAccuracy classifies moves based on EP loss thresholds
func (mc *MoveCategorizer) classifyByAccuracy(epLoss float64) models.MoveClassification {
	thresholds := mc.expectedPointsService.GetAccuracyThresholds()
	
	switch {
	case epLoss <= thresholds["excellent"]:
		return models.Excellent
	case epLoss <= thresholds["good"]:
		return models.Good
	case epLoss <= thresholds["inaccuracy"]:
		return models.Inaccuracy
	case epLoss <= thresholds["mistake"]:
		return models.Mistake
	default:
		return models.Blunder
	}
}

// involvesSacrifice checks if the move involves a piece sacrifice
func (mc *MoveCategorizer) involvesSacrifice(before, after models.MaterialValue) bool {
	// Calculate material change
	materialChange := before.Total - after.Total
	
	// A sacrifice is when we lose significant material (more than a pawn)
	// But we need to account for captures too
	return materialChange >= 300 // About 3 pawns worth of material
}

// isDifficultMove checks if a move is non-obvious or difficult to find
func (mc *MoveCategorizer) isDifficultMove(data MoveCategoryData) bool {
	// Heuristics for difficulty:
	// 1. Involves a sacrifice
	// 2. Not the first choice in simple positions
	// 3. Creates complex tactical ideas
	
	// For now, we'll use a simple heuristic
	// In the future, this could analyze move complexity, tactical patterns, etc.
	return mc.involvesSacrifice(data.MaterialBefore, data.MaterialAfter)
}

// calculateAlternativeEPLoss calculates EP loss for an alternative move
func (mc *MoveCategorizer) calculateAlternativeEPLoss(alt models.AlternativeMove, data MoveCategoryData) float64 {
	// This would require evaluating what the position would be after the alternative move
	// For now, we'll use the evaluation difference as a proxy
	
	if data.BeforeEvaluation == nil {
		return 0.0
	}
	
	// Normalize the evaluation for the current player
	playerRating := 1500 // Default rating if not available
	
	beforeEval := mc.expectedPointsService.NormalizeEvaluationForPlayer(
		data.BeforeEvaluation.Score, data.IsWhiteToMove)
	afterEval := mc.expectedPointsService.NormalizeEvaluationForPlayer(
		alt.Evaluation.Score, data.IsWhiteToMove)
	
	return mc.expectedPointsService.CalculateExpectedPointsLoss(beforeEval, afterEval, playerRating)
}

// CalculateMaterialValue calculates the total material value for a position
func (mc *MoveCategorizer) CalculateMaterialValue(position *chess.Position, color chess.Color) models.MaterialValue {
	if position == nil {
		return models.MaterialValue{}
	}
	
	material := models.MaterialValue{}
	
	// Material values in centipawns
	values := map[chess.PieceType]int{
		chess.Pawn:   100,
		chess.Knight: 320,
		chess.Bishop: 330,
		chess.Rook:   500,
		chess.Queen:  900,
		chess.King:   0, // King has no material value
	}
	
	// Count pieces for the specified color
	for square := chess.A1; square <= chess.H8; square++ {
		piece := position.Board().Piece(square)
		if piece != chess.NoPiece && piece.Color() == color {
			switch piece.Type() {
			case chess.Pawn:
				material.Pawns++
			case chess.Knight:
				material.Knights++
			case chess.Bishop:
				material.Bishops++
			case chess.Rook:
				material.Rooks++
			case chess.Queen:
				material.Queens++
			}
			material.Total += values[piece.Type()]
		}
	}
	
	return material
}

// IsPositionWinning determines if a position is already overwhelmingly winning
func (mc *MoveCategorizer) IsPositionWinning(evaluation *models.EngineEvaluation, isWhiteToMove bool) bool {
	if evaluation == nil {
		return false
	}
	
	// Handle mate scores
	if evaluation.Mate != nil {
		mateScore := *evaluation.Mate
		if isWhiteToMove {
			return mateScore > 0 // White is winning
		} else {
			return mateScore < 0 // Black is winning
		}
	}
	
	// Consider position winning if advantage is > 300 centipawns
	winningThreshold := 300
	
	if isWhiteToMove {
		return evaluation.Score > winningThreshold
	} else {
		return evaluation.Score < -winningThreshold
	}
}

// CalculateAccuracyScores calculates overall accuracy for both players
func (mc *MoveCategorizer) CalculateAccuracyScores(moves []models.MoveAnalysis, excludeBookMoves bool) (whiteAccuracy, blackAccuracy float64) {
	var whiteTotal, blackTotal float64
	var whiteCount, blackCount int
	
	for _, move := range moves {
		// Skip book moves if requested
		if excludeBookMoves && move.IsBookMove {
			continue
		}
		
		// Determine if this is a white or black move
		isWhiteMove := move.MoveNumber%2 == 1
		
		if isWhiteMove {
			whiteTotal += move.MoveAccuracy
			whiteCount++
		} else {
			blackTotal += move.MoveAccuracy
			blackCount++
		}
	}
	
	// Calculate averages
	if whiteCount > 0 {
		whiteAccuracy = whiteTotal / float64(whiteCount)
	}
	if blackCount > 0 {
		blackAccuracy = blackTotal / float64(blackCount)
	}
	
	return whiteAccuracy, blackAccuracy
}

// EnhancedMoveClassification provides additional context for move classifications
type EnhancedMoveClassification struct {
	Classification models.MoveClassification `json:"classification"`
	Reason         string                    `json:"reason"`
	Confidence     float64                   `json:"confidence"`
	Context        ClassificationContext     `json:"context"`
}

// ClassificationContext provides context about why a move was classified a certain way
type ClassificationContext struct {
	IsSacrifice       bool    `json:"isSacrifice"`
	MaterialChange    int     `json:"materialChange"`
	EPLoss            float64 `json:"epLoss"`
	PositionComplexity string  `json:"positionComplexity"`
	AlternativeCount  int     `json:"alternativeCount"`
}

// ClassifyMoveWithContext provides detailed classification with reasoning
func (mc *MoveCategorizer) ClassifyMoveWithContext(data MoveCategoryData) EnhancedMoveClassification {
	classification := mc.CategorizeMoveAdvanced(data)
	
	// Generate reasoning and context
	context := ClassificationContext{
		IsSacrifice:       mc.involvesSacrifice(data.MaterialBefore, data.MaterialAfter),
		MaterialChange:    data.MaterialBefore.Total - data.MaterialAfter.Total,
		EPLoss:            data.EPLoss,
		AlternativeCount:  len(data.AlternativeMoves),
	}
	
	reason := mc.generateClassificationReason(classification, context)
	confidence := mc.calculateClassificationConfidence(classification, data)
	
	return EnhancedMoveClassification{
		Classification: classification,
		Reason:         reason,
		Confidence:     confidence,
		Context:        context,
	}
}

// generateClassificationReason generates human-readable reason for classification
func (mc *MoveCategorizer) generateClassificationReason(classification models.MoveClassification, context ClassificationContext) string {
	switch classification {
	case models.Brilliant:
		if context.IsSacrifice {
			return "Brilliant sacrificial move that objectively improves the position"
		}
		return "Brilliant move that finds the best continuation in a complex position"
	case models.Great:
		return "Great move - the only good option in a difficult position"
	case models.Best:
		return "Best move according to the engine"
	case models.Excellent:
		return "Excellent move with minimal loss of advantage"
	case models.Good:
		return "Good move maintaining the position"
	case models.Book:
		return "Theoretical opening move"
	case models.Inaccuracy:
		return "Inaccuracy that slightly worsens the position"
	case models.Mistake:
		return "Mistake that significantly worsens the position"
	case models.Blunder:
		return "Blunder that severely damages the position"
	default:
		return "Move analyzed"
	}
}

// calculateClassificationConfidence calculates confidence in the classification
func (mc *MoveCategorizer) calculateClassificationConfidence(classification models.MoveClassification, data MoveCategoryData) float64 {
	// Base confidence on various factors
	confidence := 0.8 // Base confidence
	
	// Higher confidence for clear-cut cases
	if classification == models.Best && data.EPLoss < 0.001 {
		confidence = 0.95
	}
	
	if classification == models.Blunder && data.EPLoss > 0.5 {
		confidence = 0.95
	}
	
	// Lower confidence for borderline cases
	thresholds := mc.expectedPointsService.GetAccuracyThresholds()
	for _, threshold := range thresholds {
		if math.Abs(data.EPLoss-threshold) < 0.01 {
			confidence = 0.6 // Borderline case
			break
		}
	}
	
	return confidence
} 