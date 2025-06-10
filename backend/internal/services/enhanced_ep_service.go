package services

import (
	"chess-backend/internal/models"
	"fmt"

	"github.com/notnil/chess"
	"github.com/sirupsen/logrus"
)

// EnhancedEPService combines all EP algorithm components
type EnhancedEPService struct {
	expectedPointsService *ExpectedPointsService
	calibrationService    *CalibrationService
	moveCategorizer       *MoveCategorizer
	stockfishService      *StockfishService
	chessService          *ChessService
	openingBookService    *OpeningBookService
	logger                *logrus.Logger
	isInitialized         bool
}

// NewEnhancedEPService creates a new enhanced EP service with all dependencies
func NewEnhancedEPService(
	stockfish *StockfishService,
	chess *ChessService,
	opening *OpeningService,
	logger *logrus.Logger,
) *EnhancedEPService {
	// Create services in dependency order
	eps := NewExpectedPointsService()
	calibration := NewCalibrationService(stockfish, eps, chess, logger)
	categorizer := NewMoveCategorizer(eps, chess, opening)
	openingBook := NewOpeningBookService()
	
	// Connect services
	eps.SetCalibrationService(calibration)
	
	service := &EnhancedEPService{
		expectedPointsService: eps,
		calibrationService:    calibration,
		moveCategorizer:       categorizer,
		stockfishService:      stockfish,
		chessService:          chess,
		openingBookService:    openingBook,
		logger:                logger,
		isInitialized:         false,
	}
	
	// Initialize thresholds
	if err := service.Initialize(); err != nil {
		logger.Warnf("Failed to initialize enhanced EP service: %v", err)
	}
	
	return service
}

// Initialize loads thresholds and prepares the service for use
func (eep *EnhancedEPService) Initialize() error {
	eep.logger.Info("Initializing Enhanced EP Service...")
	
	// Load existing thresholds or use defaults
	thresholds, err := eep.calibrationService.LoadThresholds()
	if err != nil {
		eep.logger.Warnf("Using default thresholds: %v", err)
	} else {
		eep.logger.Infof("Loaded %d rating bucket thresholds", len(thresholds))
	}
	
	eep.isInitialized = true
	return nil
}

// ClassifyMove classifies a single move using the revised EP algorithm
func (eep *EnhancedEPService) ClassifyMove(ms models.MoveStat) models.MoveClassification {
	// Validate input
	if !eep.isInitialized {
		eep.Initialize()
	}
	
	// Get rating-specific thresholds
	T := eep.expectedPointsService.GetDynamicThresholds(ms.Rating)
	
	// Apply the REVISED algorithm specification exactly
	
	// 4.1: Book moves only in the opening phase (first 15 moves)
	if ms.MoveNumber <= 15 && eep.openingBookService.Contains(ms.MoveUCI) {
		return models.Book
	}
	
	// 4.2: Best engine move by UCI match (fixes "all → Excellent" problem)
	if ms.MoveUCI == ms.BestEngineUCI {
		return models.Best
	}
	
	// 4.3: Brilliant - best move + sacrifice + ultra-rare EP loss (P1)
	if ms.MoveUCI == ms.BestEngineUCI && ms.MaterialChange < 0 && ms.EPLoss <= T.P1 {
		return models.Brilliant
	}
	
	// 4.4: Great - best move + rare EP loss (P5)
	if ms.MoveUCI == ms.BestEngineUCI && ms.EPLoss <= T.P5 {
		return models.Great
	}
	
	// 4.5: Excellent - very low EP loss (P10)
	if ms.EPLoss <= T.P10 {
		return models.Excellent
	}
	
	// 4.6: Good - moderately low EP loss (P25)
	if ms.EPLoss <= T.P25 {
		return models.Good
	}
	
	// 4.7: Inaccuracy - small slip (P50)
	if ms.EPLoss <= T.P50 {
		return models.Inaccuracy
	}
	
	// 4.8: Mistake - noticeable error (P75)
	if ms.EPLoss <= T.P75 {
		return models.Mistake
	}
	
	// 4.9: Miss - big oversight (P90)
	if ms.EPLoss <= T.P90 {
		return models.Miss
	}
	
	// 4.10: Blunder - worst mistakes (>P90)
	return models.Blunder
}

// AnalyzeMoveComplete performs complete move analysis with EP calculation
func (eep *EnhancedEPService) AnalyzeMoveComplete(
	beforeFEN, afterFEN string,
	move *chess.Move,
	playerRating int,
	moveNumber int,
) (*models.MoveAnalysis, error) {
	
	// Get evaluations for before and after positions
	evalBefore, _, err := eep.stockfishService.AnalyzePosition(beforeFEN, 12, 1000, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate before position: %w", err)
	}
	
	evalAfter, _, err := eep.stockfishService.AnalyzePosition(afterFEN, 12, 1000, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate after position: %w", err)
	}
	
	// Calculate Expected Points
	epBefore := eep.expectedPointsService.CalculateExpectedPoints(
		evalBefore.Score, playerRating)
	epAfter := eep.expectedPointsService.CalculateExpectedPoints(
		evalAfter.Score, playerRating)
	epLoss := epBefore - epAfter
	
	// Get UCI notation for moves
	moveUCI := eep.moveToUCI(move)
	bestEngineUCI := evalBefore.BestMove
	
	// Determine if this was the best move by UCI comparison
	isBestMove := moveUCI == bestEngineUCI
	
	// Calculate material change (simplified implementation)
	materialChange := eep.calculateMaterialChange(beforeFEN, afterFEN, move)
	isSacrifice := materialChange < -100 // Lost more than a pawn's worth
	
	// Create MoveStat for classification
	moveStat := models.MoveStat{
		Rating:         playerRating,
		EPLoss:         epLoss,
		MaterialChange: materialChange,
		MoveUCI:        moveUCI,
		BestEngineUCI:  bestEngineUCI,
		IsBestMove:     isBestMove,
		IsSacrifice:    isSacrifice,
		MoveNumber:     moveNumber,
		GamePhase:      eep.determineGamePhase(moveNumber),
	}
	
	// Classify the move
	classification := eep.ClassifyMove(moveStat)
	
	// Build complete analysis
	moveAnalysis := &models.MoveAnalysis{
		MoveNumber:     moveNumber,
		Move:           move.String(),
		SAN:            move.String(), // Would need proper SAN conversion
		FEN:            afterFEN,
		Evaluation:     *evalAfter,
		Classification: string(classification),
		BeforeEvaluation: evalBefore,
		ExpectedPoints: models.ExpectedPointsData{
			Before:   epBefore,
			After:    epAfter,
			Loss:     epLoss,
			Accuracy: (1.0 - epLoss) * 100.0,
		},
		MoveAccuracy: (1.0 - epLoss) * 100.0,
	}
	
	return moveAnalysis, nil
}

// RunCalibrationFromPGN runs the calibration phase on a PGN file
func (eep *EnhancedEPService) RunCalibrationFromPGN(pgnPath string) error {
	eep.logger.Infof("Starting calibration from PGN: %s", pgnPath)
	
	// Step 1: Run calibration to collect move statistics
	if err := eep.calibrationService.RunCalibration(pgnPath); err != nil {
		return fmt.Errorf("calibration failed: %w", err)
	}
	
	// Step 2: Compute percentiles from collected data
	if err := eep.calibrationService.ComputePercentiles(); err != nil {
		return fmt.Errorf("percentile computation failed: %w", err)
	}
	
	// Step 3: Reload thresholds
	if err := eep.Initialize(); err != nil {
		return fmt.Errorf("failed to reload thresholds: %w", err)
	}
	
	eep.logger.Info("Calibration completed successfully")
	return nil
}

// GetThresholds returns the current thresholds for all rating buckets
func (eep *EnhancedEPService) GetThresholds() map[models.RatingBucket]models.EPThresholds {
	thresholds, _ := eep.calibrationService.LoadThresholds()
	return thresholds
}

// GetThresholdsForRating returns thresholds for a specific rating
func (eep *EnhancedEPService) GetThresholdsForRating(rating int) models.EPThresholds {
	return eep.expectedPointsService.GetDynamicThresholds(rating)
}

// Helper functions

func (eep *EnhancedEPService) calculateMaterialChange(beforeFEN, afterFEN string, move *chess.Move) int {
	// Simplified material calculation - in production this would need proper implementation
	// For now, detect obvious captures and promotions
	
	if move.HasTag(chess.Capture) {
		// Estimate captured piece value
		// This is a simplification - real implementation would parse FENs
		return -100 // Assume pawn capture
	}
	
	if move.Promo() != chess.NoPieceType {
		// Promotion
		switch move.Promo() {
		case chess.Queen:
			return 800 // Queen value minus pawn
		case chess.Rook:
			return 400
		case chess.Bishop, chess.Knight:
			return 200
		}
	}
	
	return 0
}

func (eep *EnhancedEPService) determineGamePhase(moveNumber int) string {
	switch {
	case moveNumber <= 15:
		return "opening"
	case moveNumber <= 40:
		return "middlegame"
	default:
		return "endgame"
	}
}

func (eep *EnhancedEPService) moveToUCI(move *chess.Move) string {
	uci := move.S1().String() + move.S2().String()
	
	// Add promotion piece if applicable
	if move.Promo() != chess.NoPieceType {
		switch move.Promo() {
		case chess.Queen:
			uci += "q"
		case chess.Rook:
			uci += "r"
		case chess.Bishop:
			uci += "b"
		case chess.Knight:
			uci += "n"
		}
	}
	
	return uci
}

func convertToModelEvaluation(eval *models.EngineEvaluation) models.EngineEvaluation {
	if eval == nil {
		return models.EngineEvaluation{}
	}
	return *eval
} 