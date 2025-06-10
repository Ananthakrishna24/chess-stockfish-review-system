package services

import (
	"chess-backend/internal/models"
	"strings"

	"github.com/notnil/chess"
)

// MoveCategorizer implements Chess.com's move categorization algorithm
type MoveCategorizer struct {
	expectedPointsService *ImprovedExpectedPointsService
	chessService          *ChessService
	openingService        *OpeningService
}

// NewMoveCategorizer creates a new move categorization service
func NewMoveCategorizer(eps *ExpectedPointsService, chess *ChessService, opening *OpeningService) *MoveCategorizer {
	// Use the improved expected points service instead
	improvedEPS := NewImprovedExpectedPointsService()
	return &MoveCategorizer{
		expectedPointsService: improvedEPS,
		chessService:          chess,
		openingService:        opening,
	}
}

// ChessComMoveClassifier implements Chess.com's exact move classification algorithm
type ChessComMoveClassifier struct {
	expectedPointsService *ImprovedExpectedPointsService
}

// NewChessComMoveClassifier creates a new Chess.com-style move classifier
func NewChessComMoveClassifier(eps *ExpectedPointsService) *ChessComMoveClassifier {
	// Use the improved expected points service instead
	improvedEPS := NewImprovedExpectedPointsService()
	return &ChessComMoveClassifier{
		expectedPointsService: improvedEPS,
	}
}

// ClassifyMove implements Chess.com's move classification algorithm exactly
func (c *ChessComMoveClassifier) ClassifyMove(data MoveCategoryData) models.MoveClassification {
	// Step 1: Check for Book moves first (opening theory)
	if c.isBookMove(data) {
		return models.Book
	}
	
	// Step 2: Calculate expected points loss using Chess.com's method
	epLoss := c.calculateExpectedPointsLoss(data)
	
	// Step 3: Check for special moves (Brilliant) before thresholds
	if c.isBrilliantMove(data, epLoss) {
		return models.Brilliant
	}
	
	// Step 4: Apply Chess.com's exact thresholds for classification
	return c.classifyByExpectedPointsLoss(epLoss, data.PlayerRating)
}

// calculateExpectedPointsLoss calculates EP loss using Chess.com's method
func (c *ChessComMoveClassifier) calculateExpectedPointsLoss(data MoveCategoryData) float64 {
	if data.BeforeEvaluation == nil || data.AfterEvaluation == nil {
		return 0.0
	}
	
	// Get player rating (default to 1500 if not provided)
	playerRating := data.PlayerRating
	if playerRating == 0 {
		playerRating = 1500
	}
	
	// Normalize evaluations for the current player
	beforeEval := c.expectedPointsService.NormalizeEvaluationForPlayer(
		data.BeforeEvaluation.Score, data.IsWhiteToMove)
	afterEval := c.expectedPointsService.NormalizeEvaluationForPlayer(
		data.AfterEvaluation.Score, data.IsWhiteToMove)
	
	// Convert to expected points (win probability 0.00-1.00)
	epBefore := c.expectedPointsService.CalculateExpectedPoints(beforeEval, playerRating)
	epAfter := c.expectedPointsService.CalculateExpectedPoints(afterEval, playerRating)
	
	// EP loss = reduction in win probability
	epLoss := epBefore - epAfter
	
	// Ensure non-negative (we're only interested in losses)
	if epLoss < 0 {
		epLoss = 0
	}
	
	return epLoss
}

// classifyByExpectedPointsLoss applies Chess.com's exact thresholds
func (c *ChessComMoveClassifier) classifyByExpectedPointsLoss(epLoss float64, playerRating int) models.MoveClassification {
	// Chess.com's exact thresholds (rating-independent for now):
	// Best/Excellent: ≤0.02
	// Good: 0.02–0.05  
	// Inaccuracy: 0.05–0.10
	// Mistake: 0.10–0.20
	// Blunder: >0.20
	
	// Apply rating-based scaling if needed
	thresholds := c.getChessComThresholds(playerRating)
	
	switch {
	case epLoss == 0.0:
		return models.Best // No EP loss = engine's top move
	case epLoss <= thresholds.Excellent:
		return models.Excellent // Tiny losses up to 0.02
	case epLoss <= thresholds.Good:
		return models.Good // 0.02–0.05 lost
	case epLoss <= thresholds.Inaccuracy:
		return models.Inaccuracy // 0.05–0.10 lost
	case epLoss <= thresholds.Mistake:
		return models.Mistake // 0.10–0.20 lost
	default:
		return models.Blunder // >0.20 lost
	}
}

// ChessComThresholds represents Chess.com's classification thresholds
type ChessComThresholds struct {
	Excellent   float64
	Good        float64
	Inaccuracy  float64
	Mistake     float64
	Blunder     float64
}

// getChessComThresholds returns Chess.com's exact thresholds
func (c *ChessComMoveClassifier) getChessComThresholds(playerRating int) ChessComThresholds {
	// Chess.com's published thresholds
	base := ChessComThresholds{
		Excellent:  0.02,  // ≤0.02 (tiny losses)
		Good:       0.05,  // 0.02–0.05 lost
		Inaccuracy: 0.10,  // 0.05–0.10 lost
		Mistake:    0.20,  // 0.10–0.20 lost
		Blunder:    1.00,  // >0.20 lost (effectively infinite)
	}
	
	// Chess.com scales thresholds by rating
	// Higher rated players have stricter thresholds
	ratingFactor := c.getRatingScalingFactor(playerRating)
	
	return ChessComThresholds{
		Excellent:  base.Excellent * ratingFactor,
		Good:       base.Good * ratingFactor,
		Inaccuracy: base.Inaccuracy * ratingFactor,
		Mistake:    base.Mistake * ratingFactor,
		Blunder:    base.Blunder, // Blunder threshold doesn't scale
	}
}

// getRatingScalingFactor scales thresholds based on player rating
func (c *ChessComMoveClassifier) getRatingScalingFactor(rating int) float64 {
	// Higher rated players have stricter standards
	// Rating 1200 = 1.0 (baseline)
	// Rating 2400 = 0.8 (20% stricter)
	// Rating 800 = 1.2 (20% more lenient)
	
	baseFactor := 1.0
	ratingDiff := float64(rating - 1200)
	scalingRate := -0.0001 // -0.01% per rating point
	
	factor := baseFactor + (ratingDiff * scalingRate)
	
	// Clamp between 0.6 and 1.4
	if factor < 0.6 {
		factor = 0.6
	}
	if factor > 1.4 {
		factor = 1.4
	}
	
	return factor
}

// isBookMove checks if the move is from opening theory
func (c *ChessComMoveClassifier) isBookMove(data MoveCategoryData) bool {
	// Book moves are treated specially in opening theory
	// Chess.com marks well-known opening moves as Book
	
	// Must be in opening phase (first 12-15 moves typically)
	if data.MoveNumber > 15 {
		return false
	}
	
	// If we have ECO classification, likely book move
	if data.ECO != "" && data.MoveNumber <= 12 {
		return true
	}
	
	// If opening service is available, check against database
	if c.expectedPointsService != nil && data.OpeningName != "" && data.MoveNumber <= 10 {
		return true
	}
	
	// Conservative fallback: very early moves with minimal EP loss
	if data.MoveNumber <= 6 {
		epLoss := c.calculateExpectedPointsLoss(data)
		return epLoss <= 0.01 // Very small tolerance for book moves
	}
	
	return false
}

// isBrilliantMove implements Chess.com's brilliant move detection
func (c *ChessComMoveClassifier) isBrilliantMove(data MoveCategoryData, epLoss float64) bool {
	// Chess.com's criteria for Brilliant moves:
	// 1. Must be a good piece sacrifice
	// 2. Must be among the engine's best moves (minimal EP loss)
	// 3. Should not be in a bad position after the move
	// 4. Should not already be completely winning
	
	// Must have very low EP loss (engine agrees it's best or nearly best)
	if epLoss > 0.02 {
		return false
	}
	
	// Must involve a piece sacrifice
	if !c.involvesPieceSacrifice(data) {
		return false
	}
	
	// Position shouldn't be completely winning already
	if c.isPositionAlreadyWinning(data.BeforeEvaluation, data.IsWhiteToMove) {
		return false
	}
	
	// Position shouldn't be bad after the move
	if c.isPositionBad(data.AfterEvaluation, data.IsWhiteToMove) {
		return false
	}
	
	return true
}

// involvesPieceSacrifice checks if the move involves a piece sacrifice
func (c *ChessComMoveClassifier) involvesPieceSacrifice(data MoveCategoryData) bool {
	// Calculate material change
	materialBefore := data.MaterialBefore.Total
	materialAfter := data.MaterialAfter.Total
	
	// A piece sacrifice means losing significant material (not just trading)
	// Chess.com is "more generous in defining a piece sacrifice" for lower rated players
	
	materialLoss := materialBefore - materialAfter
	
	// Rating-based sacrifice threshold
	sacrificeThreshold := c.getSacrificeThreshold(data.PlayerRating)
	
	return materialLoss >= sacrificeThreshold
}

// getSacrificeThreshold returns the material loss threshold for sacrifice detection
func (c *ChessComMoveClassifier) getSacrificeThreshold(rating int) int {
	// Lower rated players: more generous (250 centipawns = 2.5 pawns)
	// Higher rated players: stricter (350 centipawns = 3.5 pawns)
	
	baseThreshold := 300 // 3 pawns
	
	if rating < 1200 {
		return 250 // More generous for beginners
	} else if rating > 2000 {
		return 350 // Stricter for masters
	}
	
	return baseThreshold
}

// isPositionAlreadyWinning checks if position was already overwhelmingly winning
func (c *ChessComMoveClassifier) isPositionAlreadyWinning(eval *models.EngineEvaluation, isWhiteToMove bool) bool {
	if eval == nil {
		return false
	}
	
	// Check for mate scores
	if eval.Mate != nil {
		return true
	}
	
	// Normalize evaluation for current player
	playerEval := eval.Score
	if !isWhiteToMove {
		playerEval = -playerEval
	}
	
	// Position is "completely winning" if advantage > +500 centipawns (roughly 5 pawns)
	return playerEval > 500
}

// isPositionBad checks if position is bad after the move
func (c *ChessComMoveClassifier) isPositionBad(eval *models.EngineEvaluation, isWhiteToMove bool) bool {
	if eval == nil {
		return false
	}
	
	// Check for being mated
	if eval.Mate != nil && *eval.Mate < 0 {
		return true
	}
	
	// Normalize evaluation for current player
	playerEval := eval.Score
	if !isWhiteToMove {
		playerEval = -playerEval
	}
	
	// Position is "bad" if disadvantage > -200 centipawns (roughly 2 pawns)
	return playerEval < -200
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
	PlayerRating      int  // Player's rating for dynamic threshold calculation
	IsBestMove        bool // Whether this move matches the engine's best move
	IsSacrifice       bool // Whether this move involves a material sacrifice
}

// CategorizeMoveAdvanced categorizes a move using Chess.com's algorithm
func (mc *MoveCategorizer) CategorizeMoveAdvanced(moveData MoveCategoryData) models.MoveClassification {
	classifier := &ChessComMoveClassifier{expectedPointsService: mc.expectedPointsService}
	return classifier.ClassifyMove(moveData)
}

// Legacy methods for backward compatibility

// isBookMove checks if the move is from opening theory
func (mc *MoveCategorizer) isBookMove(data MoveCategoryData) bool {
	classifier := &ChessComMoveClassifier{expectedPointsService: mc.expectedPointsService}
	return classifier.isBookMove(data)
}

// isBrilliantMove checks for brilliant moves (!!)
func (mc *MoveCategorizer) isBrilliantMove(data MoveCategoryData) bool {
	classifier := &ChessComMoveClassifier{expectedPointsService: mc.expectedPointsService}
	epLoss := classifier.calculateExpectedPointsLoss(data)
	return classifier.isBrilliantMove(data, epLoss)
}

// isGreatMove checks for great moves (!)
func (mc *MoveCategorizer) isGreatMove(data MoveCategoryData) bool {
	// Great moves are "critical to the outcome" - turning losing into draw/win
	// This requires more complex analysis than simple EP loss
	// For now, we'll use a simplified implementation
	
	epLoss := mc.expectedPointsService.CalculateExpectedPointsLoss(
		data.BeforeEvaluation.Score, data.AfterEvaluation.Score, data.PlayerRating)
	
	// Must be a very good move
	if epLoss > 0.01 {
		return false
	}
	
	// Check if this move significantly improved a difficult position
	return mc.isGameChangingMove(data)
}

// isGameChangingMove checks if the move was critical to the game outcome
func (mc *MoveCategorizer) isGameChangingMove(data MoveCategoryData) bool {
	// Simplified heuristic: move that improves a losing position significantly
	if data.BeforeEvaluation == nil || data.AfterEvaluation == nil {
		return false
	}
	
	beforeEval := data.BeforeEvaluation.Score
	afterEval := data.AfterEvaluation.Score
	
	// Normalize for current player
	if !data.IsWhiteToMove {
		beforeEval = -beforeEval
		afterEval = -afterEval
	}
	
	// Was losing before (< -100) but improved significantly
	wasLosing := beforeEval < -100
	improvement := afterEval - beforeEval
	
	return wasLosing && improvement > 200 // Improved by 2+ pawns
}

// isBestMove checks if the played move matches the engine's best move
func (mc *MoveCategorizer) isBestMove(data MoveCategoryData) bool {
	return strings.EqualFold(data.UCI, data.BestMove) || strings.EqualFold(data.Move, data.BestMove)
}

// classifyByAccuracy classifies moves based on EP loss thresholds
func (mc *MoveCategorizer) classifyByAccuracy(epLoss float64) models.MoveClassification {
	classifier := &ChessComMoveClassifier{expectedPointsService: mc.expectedPointsService}
	return classifier.classifyByExpectedPointsLoss(epLoss, 1500) // Default rating
}

// involvesSacrifice checks if the move involves a piece sacrifice
func (mc *MoveCategorizer) involvesSacrifice(before, after models.MaterialValue) bool {
	materialChange := before.Total - after.Total
	return materialChange >= 300 // About 3 pawns worth of material
}

// isDifficultMove checks if a move is non-obvious or difficult to find
func (mc *MoveCategorizer) isDifficultMove(data MoveCategoryData) bool {
	return mc.involvesSacrifice(data.MaterialBefore, data.MaterialAfter)
}

// calculateAlternativeEPLoss calculates EP loss for an alternative move
func (mc *MoveCategorizer) calculateAlternativeEPLoss(alt models.AlternativeMove, data MoveCategoryData) float64 {
	if data.BeforeEvaluation == nil {
		return 0.0
	}
	
	playerRating := data.PlayerRating
	if playerRating == 0 {
		playerRating = 1500
	}
	
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
	
	// Count pieces for the specified color
	board := position.Board()
	for sq := chess.A1; sq <= chess.H8; sq++ {
		piece := board.Piece(sq)
		if piece.Color() != color {
			continue
		}
		
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
	}
	
	// Calculate total centipawn value
	material.Total = material.Pawns*100 + material.Knights*320 + material.Bishops*330 + 
					 material.Rooks*500 + material.Queens*900
	
	return material
}

// IsPositionWinning checks if a position is overwhelmingly winning
func (mc *MoveCategorizer) IsPositionWinning(evaluation *models.EngineEvaluation, isWhiteToMove bool) bool {
	if evaluation == nil {
		return false
	}
	
	// Check for mate
	if evaluation.Mate != nil && *evaluation.Mate > 0 {
		return true
	}
	
	// Normalize evaluation for current player
	playerEval := evaluation.Score
	if !isWhiteToMove {
		playerEval = -playerEval
	}
	
	// Position is winning if advantage > +400 centipawns (roughly 4 pawns)
	return playerEval > 400
}

// CalculateAccuracyScores calculates overall accuracy for both players
func (mc *MoveCategorizer) CalculateAccuracyScores(moves []models.MoveAnalysis, excludeBookMoves bool) (whiteAccuracy, blackAccuracy float64) {
	var whiteMoves, blackMoves []models.MoveAnalysis
	
	// Separate moves by color
	for _, move := range moves {
		if excludeBookMoves && move.IsBookMove {
			continue
		}
		
		if move.MoveNumber%2 == 1 { // Odd move numbers are white
			whiteMoves = append(whiteMoves, move)
		} else { // Even move numbers are black
			blackMoves = append(blackMoves, move)
		}
	}
	
	// Calculate accuracy for each player
	whiteAccuracy = mc.calculatePlayerAccuracy(whiteMoves)
	blackAccuracy = mc.calculatePlayerAccuracy(blackMoves)
	
	return whiteAccuracy, blackAccuracy
}

// calculatePlayerAccuracy calculates accuracy for a single player's moves
func (mc *MoveCategorizer) calculatePlayerAccuracy(moves []models.MoveAnalysis) float64 {
	if len(moves) == 0 {
		return 100.0
	}
	
	totalAccuracy := 0.0
	for _, move := range moves {
		totalAccuracy += move.MoveAccuracy
	}
	
	return totalAccuracy / float64(len(moves))
}

// Enhanced move classification with context

type EnhancedMoveClassification struct {
	Classification models.MoveClassification `json:"classification"`
	Reason         string                    `json:"reason"`
	Confidence     float64                   `json:"confidence"`
	Context        ClassificationContext     `json:"context"`
}

type ClassificationContext struct {
	IsSacrifice       bool    `json:"isSacrifice"`
	MaterialChange    int     `json:"materialChange"`
	EPLoss            float64 `json:"epLoss"`
	PositionComplexity string  `json:"positionComplexity"`
	AlternativeCount  int     `json:"alternativeCount"`
}

// ClassifyMoveWithContext provides enhanced classification with reasoning
func (mc *MoveCategorizer) ClassifyMoveWithContext(data MoveCategoryData) EnhancedMoveClassification {
	classifier := &ChessComMoveClassifier{expectedPointsService: mc.expectedPointsService}
	classification := classifier.ClassifyMove(data)
	
	epLoss := classifier.calculateExpectedPointsLoss(data)
	
	context := ClassificationContext{
		IsSacrifice:       classifier.involvesPieceSacrifice(data),
		MaterialChange:    data.MaterialBefore.Total - data.MaterialAfter.Total,
		EPLoss:            epLoss,
		PositionComplexity: "medium", // Could be enhanced with position analysis
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

// generateClassificationReason provides human-readable reasoning
func (mc *MoveCategorizer) generateClassificationReason(classification models.MoveClassification, context ClassificationContext) string {
	switch classification {
	case models.Book:
		return "Well-known opening move from chess theory"
	case models.Brilliant:
		return "Excellent piece sacrifice that creates winning advantage"
	case models.Great:
		return "Critical move that significantly improves the position"
	case models.Best:
		return "Engine's top choice with no loss in winning chances"
	case models.Excellent:
		return "Very strong move with minimal loss in winning chances"
	case models.Good:
		return "Solid move with small loss in winning chances"
	case models.Inaccuracy:
		return "Slightly inaccurate move that worsens the position"
	case models.Mistake:
		return "Significant error that gives away advantage"
	case models.Blunder:
		return "Serious mistake that likely loses material or the game"
	default:
		return "Move classification unavailable"
	}
}

// calculateClassificationConfidence estimates confidence in the classification
func (mc *MoveCategorizer) calculateClassificationConfidence(classification models.MoveClassification, data MoveCategoryData) float64 {
	// Higher confidence for:
	// - Clear book moves
	// - Very large EP losses (clear blunders)
	// - Brilliant moves with clear sacrifices
	
		switch classification {
	case models.Book:
		if data.MoveNumber <= 8 && data.ECO != "" {
			return 0.95
		}
		return 0.80
	case models.Brilliant:
		if data.IsSacrifice && data.EPLoss < 0.01 {
			return 0.90
		}
		return 0.75
	case models.Blunder:
		if data.EPLoss > 0.30 {
			return 0.95
		}
				return 0.85
	case models.Best:
		if data.EPLoss == 0.0 {
			return 0.90
		}
		return 0.80
	default:
		return 0.75 // Default confidence
	}
}
