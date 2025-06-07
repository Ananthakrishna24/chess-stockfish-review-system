package services

import (
	"chess-backend/internal/models"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// EnhancedAnalysisService provides the main interface for EP-based chess analysis
type EnhancedAnalysisService struct {
	stockfishService      *StockfishService
	chessService          *ChessService
	cacheService          *CacheService
	playerService         *PlayerService
	openingService        *OpeningService
	expectedPointsService *ExpectedPointsService
	moveCategorizer       *MoveCategorizer
}

// NewEnhancedAnalysisService creates a new enhanced analysis service with all components
func NewEnhancedAnalysisService(
	stockfish *StockfishService,
	chess *ChessService,
	cache *CacheService,
	player *PlayerService,
	opening *OpeningService,
) *EnhancedAnalysisService {
	
	// Initialize EP service and categorizer
	epsService := NewExpectedPointsService()
	categorizer := NewMoveCategorizer(epsService, chess, opening)
	
	return &EnhancedAnalysisService{
		stockfishService:      stockfish,
		chessService:          chess,
		cacheService:          cache,
		playerService:         player,
		openingService:        opening,
		expectedPointsService: epsService,
		moveCategorizer:       categorizer,
	}
}

// AnalyzeGameWithEP performs comprehensive EP-based game analysis
func (eas *EnhancedAnalysisService) AnalyzeGameWithEP(
	pgn string,
	options models.AnalysisOptions,
	progressCallback func(int, int),
) (*models.GameAnalysisResponse, error) {
	
	startTime := time.Now()
	logrus.Infof("Starting enhanced EP-based analysis")
	
	// Parse PGN
	parsedGame, err := eas.chessService.ParsePGN(pgn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PGN: %w", err)
	}
	
	// Validate player ratings
	if options.PlayerRatings.White == 0 {
		options.PlayerRatings.White = 1500
		logrus.Warn("White rating not provided, using default 1500")
	}
	if options.PlayerRatings.Black == 0 {
		options.PlayerRatings.Black = 1500
		logrus.Warn("Black rating not provided, using default 1500")
	}
	
	// Set analysis depth for EP calculations
	if options.Depth == 0 {
		options.Depth = 18 // Deeper analysis for better EP calculations
	}
	
	// Perform the enhanced analysis
	analysis, err := eas.stockfishService.AnalyzeGameEnhanced(parsedGame, options, progressCallback)
	if err != nil {
		return nil, fmt.Errorf("enhanced analysis failed: %w", err)
	}
	
	// Post-process with advanced categorization if needed
	eas.enhanceMoveCategorization(analysis, options)
	
	// Calculate additional EP-based metrics
	eas.calculateEPMetrics(analysis, options)
	
	// Create response
	response := &models.GameAnalysisResponse{
		GameID:         eas.cacheService.GenerateGameID(pgn),
		GameInfo:       parsedGame.GameInfo,
		Analysis:       *analysis,
		ProcessingTime: time.Since(startTime).Seconds(),
		Timestamp:      time.Now(),
	}
	
	logrus.Infof("Enhanced EP analysis completed in %.2f seconds", response.ProcessingTime)
	return response, nil
}

// enhanceMoveCategorization applies advanced categorization to moves
func (eas *EnhancedAnalysisService) enhanceMoveCategorization(analysis *models.GameAnalysis, options models.AnalysisOptions) {
	logrus.Debug("Applying enhanced move categorization")
	
	for i := range analysis.Moves {
		move := &analysis.Moves[i]
		
		// Create categorization data
		categoryData := MoveCategoryData{
			MoveNumber:       move.MoveNumber,
			Move:             move.Move,
			SAN:              move.SAN,
			UCI:              move.Move,
			IsWhiteToMove:    move.MoveNumber%2 == 1,
			BeforeEvaluation: move.BeforeEvaluation,
			AfterEvaluation:  &move.Evaluation,
			EPLoss:           move.ExpectedPoints.Loss,
			MoveAccuracy:     move.ExpectedPoints.Accuracy,
			AlternativeMoves: move.AlternativeMoves,
		}
		
		// Set best move
		if move.BeforeEvaluation != nil {
			categoryData.BestMove = move.BeforeEvaluation.BestMove
		}
		
		// Apply enhanced categorization
		enhanced := eas.moveCategorizer.ClassifyMoveWithContext(categoryData)
		
		// Update the move with enhanced classification
		move.Classification = enhanced.Classification.String()
		
		// Add enhanced comment if available
		if enhanced.Reason != "" && move.Comment == "" {
			move.Comment = enhanced.Reason
		}
	}
}

// calculateEPMetrics calculates additional Expected Points metrics
func (eas *EnhancedAnalysisService) calculateEPMetrics(analysis *models.GameAnalysis, options models.AnalysisOptions) {
	logrus.Debug("Calculating additional EP metrics")
	
	// Calculate accuracy excluding book moves
	whiteAccuracy, blackAccuracy := eas.moveCategorizer.CalculateAccuracyScores(analysis.Moves, true)
	
	// Update player statistics with book-move-excluded accuracy
	analysis.WhiteStats.Accuracy = whiteAccuracy
	analysis.BlackStats.Accuracy = blackAccuracy
	
	// Calculate EP-based critical moments (already done in main analysis)
	eas.identifyEPCriticalMoments(analysis)
	
	// Add phase-based accuracy analysis
	eas.enhancePhaseAnalysis(analysis)
}

// identifyEPCriticalMoments identifies critical moments based on EP loss
func (eas *EnhancedAnalysisService) identifyEPCriticalMoments(analysis *models.GameAnalysis) {
	logrus.Debug("Identifying EP-based critical moments")
	
	for _, move := range analysis.Moves {
		// Large EP loss indicates critical moment
		if move.ExpectedPoints.Loss > 0.15 {
			// Check if this critical moment is already recorded
			found := false
			for _, existing := range analysis.CriticalMoments {
				if existing.MoveNumber == move.MoveNumber {
					found = true
					break
				}
			}
			
			if !found {
				critical := models.CriticalMoment{
					MoveNumber:  move.MoveNumber,
					Description: fmt.Sprintf("Large EP loss: %.3f", move.ExpectedPoints.Loss),
				}
				
				if move.BeforeEvaluation != nil {
					critical.BeforeEval = move.BeforeEvaluation.Score
					critical.AfterEval = move.Evaluation.Score
					
					if move.Evaluation.Score > move.BeforeEvaluation.Score {
						critical.Advantage = "white"
					} else {
						critical.Advantage = "black"
					}
				}
				
				analysis.CriticalMoments = append(analysis.CriticalMoments, critical)
			}
		}
	}
}

// enhancePhaseAnalysis calculates phase-specific EP metrics
func (eas *EnhancedAnalysisService) enhancePhaseAnalysis(analysis *models.GameAnalysis) {
	logrus.Debug("Enhancing phase analysis with EP metrics")
	
	var openingMoves, middlegameMoves, endgameMoves []models.MoveAnalysis
	
	// Separate moves by phase
	for _, move := range analysis.Moves {
		switch {
		case move.MoveNumber <= analysis.GamePhases.Opening:
			openingMoves = append(openingMoves, move)
		case move.MoveNumber <= analysis.GamePhases.Middlegame:
			middlegameMoves = append(middlegameMoves, move)
		default:
			endgameMoves = append(endgameMoves, move)
		}
	}
	
	// Calculate phase accuracies using EP data
	analysis.PhaseAnalysis.OpeningAccuracy = eas.calculatePhaseAccuracy(openingMoves)
	analysis.PhaseAnalysis.MiddlegameAccuracy = eas.calculatePhaseAccuracy(middlegameMoves)
	analysis.PhaseAnalysis.EndgameAccuracy = eas.calculatePhaseAccuracy(endgameMoves)
}

// calculatePhaseAccuracy calculates average accuracy for a phase
func (eas *EnhancedAnalysisService) calculatePhaseAccuracy(moves []models.MoveAnalysis) float64 {
	if len(moves) == 0 {
		return 0.0
	}
	
	total := 0.0
	count := 0
	
	for _, move := range moves {
		// Skip book moves for accuracy calculation
		if !move.IsBookMove {
			total += move.ExpectedPoints.Accuracy
			count++
		}
	}
	
	if count == 0 {
		return 100.0 // All book moves
	}
	
	return total / float64(count)
}

// GetExpectedPointsService returns the EP service for external use
func (eas *EnhancedAnalysisService) GetExpectedPointsService() *ExpectedPointsService {
	return eas.expectedPointsService
}

// GetMoveCategorizer returns the move categorizer for external use
func (eas *EnhancedAnalysisService) GetMoveCategorizer() *MoveCategorizer {
	return eas.moveCategorizer
}

// CalculatePositionEP calculates Expected Points for a single position
func (eas *EnhancedAnalysisService) CalculatePositionEP(fen string, playerRating int, isWhiteToMove bool) (float64, error) {
	// Analyze the position
	eval, _, err := eas.stockfishService.AnalyzePosition(fen, 15, 2000, 1)
	if err != nil {
		return 0, fmt.Errorf("failed to analyze position: %w", err)
	}
	
	// Normalize evaluation for player
	normalizedEval := eas.expectedPointsService.NormalizeEvaluationForPlayer(eval.Score, isWhiteToMove)
	
	// Calculate EP
	ep := eas.expectedPointsService.CalculateExpectedPoints(normalizedEval, playerRating)
	
	return ep, nil
}

// ValidateEPAnalysisOptions validates and sets defaults for EP analysis options
func (eas *EnhancedAnalysisService) ValidateEPAnalysisOptions(options *models.AnalysisOptions) {
	// Use performance optimizer to get optimal settings
	optimizer := NewPerformanceOptimizer()
	optimalSettings := optimizer.GetOptimalSettings("game_analysis")
	
	// Set minimum depth for EP calculations based on optimization
	if options.Depth < optimalSettings.DepthRecommended {
		options.Depth = optimalSettings.DepthRecommended
		logrus.Debugf("Increased analysis depth to %d for better EP calculations", options.Depth)
	}
	
	// Set default ratings if not provided
	if options.PlayerRatings.White == 0 {
		options.PlayerRatings.White = 1500
	}
	if options.PlayerRatings.Black == 0 {
		options.PlayerRatings.Black = 1500
	}
	
	// Enable book move analysis by default
	options.IncludeBookMoves = true
	
	// Set optimal time per move based on system capabilities
	if options.TimePerMove == 0 {
		options.TimePerMove = optimalSettings.TimeRecommended
		logrus.Debugf("Set optimal time per move to %dms based on system capabilities", options.TimePerMove)
	}
	
	// Update Stockfish configuration with optimal settings
	optimalConfig := optimizer.ConvertToEngineOptions(optimalSettings)
	if err := eas.stockfishService.UpdateConfig(optimalConfig); err != nil {
		logrus.Warnf("Could not update Stockfish configuration: %v", err)
	} else {
		logrus.Debugf("Updated Stockfish with optimal settings: %d threads, %dMB hash", 
			optimalConfig.Threads, optimalConfig.Hash)
	}
	
	// Log optimization report for transparency
	optimizer.LogOptimizationReport("game_analysis")
} 