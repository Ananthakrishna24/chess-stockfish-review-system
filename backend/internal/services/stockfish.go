package services

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"chess-backend/internal/models"
	"chess-backend/pkg/uci"

	"github.com/sirupsen/logrus"
)

// StockfishService manages a pool of Stockfish engines
type StockfishService struct {
	engines    []*uci.Engine
	available  chan *uci.Engine
	maxWorkers int
	binaryPath string
	mutex      sync.RWMutex
	config     models.EngineOptions
	engineInfo *uci.EngineInfo
	optimized  bool
}

// NewStockfishService creates a new Stockfish service
func NewStockfishService(maxWorkers int, binaryPath string) *StockfishService {
	// Auto-detect optimal configuration based on system resources
	optimalConfig := getOptimalEngineConfig()
	
	return &StockfishService{
		engines:    make([]*uci.Engine, 0, maxWorkers),
		available:  make(chan *uci.Engine, maxWorkers),
		maxWorkers: maxWorkers,
		binaryPath: binaryPath,
		config:     optimalConfig,
		optimized:  true,
	}
}

// getOptimalEngineConfig calculates optimal Stockfish settings based on system resources
func getOptimalEngineConfig() models.EngineOptions {
	cpuCount := runtime.NumCPU()
	
	// Use max cores - 2 for optimal performance (leave some for OS/other tasks)
	optimalThreads := cpuCount - 2
	if optimalThreads < 1 {
		optimalThreads = 1
	}
	if optimalThreads > 32 { // Stockfish performs best with <= 32 threads
		optimalThreads = 32
	}
	
	// Calculate optimal hash size based on available memory
	// Rule of thumb: Use up to 1/4 of available RAM, but cap at 4GB for stability
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	// Convert to MB and calculate 25% of total system memory
	totalMemMB := int(memStats.Sys / 1024 / 1024)
	optimalHash := totalMemMB / 4
	
	// Apply constraints based on Stockfish documentation
	if optimalHash < 64 {
		optimalHash = 64 // Minimum for good performance
	}
	if optimalHash > 4096 {
		optimalHash = 4096 // Maximum for stability
	}
	
	logrus.Infof("Auto-optimized Stockfish config: %d threads, %dMB hash (detected %d CPU cores)", 
		optimalThreads, optimalHash, cpuCount)
	
	return models.EngineOptions{
		Threads:          optimalThreads,
		Hash:             optimalHash,
		Contempt:         0,
		AnalysisContempt: "off",
	}
}

// Initialize creates and initializes the engine pool
func (s *StockfishService) Initialize() error {
	logrus.Infof("Initializing optimized Stockfish service with %d workers", s.maxWorkers)
	
	// First, detect and validate the Stockfish binary
	if err := s.detectStockfishCapabilities(); err != nil {
		logrus.Warnf("Could not detect Stockfish capabilities: %v", err)
	}
	
	for i := 0; i < s.maxWorkers; i++ {
		engine, err := uci.NewEngine(s.binaryPath)
		if err != nil {
			logrus.Errorf("Failed to create engine %d: %v", i, err)
			return fmt.Errorf("failed to create engine %d: %v", i, err)
		}
		
		if err := engine.Initialize(); err != nil {
			logrus.Errorf("Failed to initialize engine %d: %v", i, err)
			return fmt.Errorf("failed to initialize engine %d: %v", i, err)
		}
		
		// Configure engine with optimized settings
		if err := s.configureEngineOptimized(engine); err != nil {
			logrus.Errorf("Failed to configure engine %d: %v", i, err)
			return fmt.Errorf("failed to configure engine %d: %v", i, err)
		}
		
		s.engines = append(s.engines, engine)
		s.available <- engine
		
		logrus.Debugf("Optimized engine %d initialized successfully", i)
	}
	
	logrus.Infof("Stockfish service initialized with %d optimized engines (Threads: %d, Hash: %dMB)", 
		len(s.engines), s.config.Threads, s.config.Hash)
	return nil
}

// detectStockfishCapabilities detects Stockfish version and optimal binary
func (s *StockfishService) detectStockfishCapabilities() error {
	// Try to detect if we're using an optimal Stockfish binary
	engine, err := uci.NewEngine(s.binaryPath)
	if err != nil {
		return err
	}
	defer engine.Close()
	
	if err := engine.Initialize(); err != nil {
		return err
	}
	
	info, err := engine.GetEngineInfo()
	if err != nil {
		return err
	}
	
	s.engineInfo = info
	logrus.Infof("Detected Stockfish: %s by %s", info.Name, info.Author)
	
	// Log performance recommendations based on detected engine
	s.logPerformanceRecommendations()
	
	return nil
}

// logPerformanceRecommendations provides performance optimization tips
func (s *StockfishService) logPerformanceRecommendations() {
	if s.engineInfo == nil {
		return
	}
	
	logrus.Info("=== Stockfish Performance Recommendations ===")
	
	// Check if using optimal binary
	engineName := strings.ToLower(s.engineInfo.Name)
	if strings.Contains(engineName, "stockfish") {
		logrus.Info("✓ Using Stockfish engine")
		
		// Recommend optimal binary based on CPU features
		cpuInfo := runtime.GOARCH
		if cpuInfo == "amd64" {
			logrus.Info("💡 For best performance, use optimized binaries:")
			logrus.Info("   - Download from: https://stockfishchess.org/download/")
			logrus.Info("   - Prefer: x86-64-bmi2 or x86-64-avx2 variants")
		}
	}
	
	logrus.Infof("✓ Threads: %d (optimal for %d CPU cores)", s.config.Threads, runtime.NumCPU())
	logrus.Infof("✓ Hash: %dMB (optimal for available memory)", s.config.Hash)
	logrus.Info("================================================")
}

// configureEngineOptimized applies optimized configuration to an engine
func (s *StockfishService) configureEngineOptimized(engine *uci.Engine) error {
	// Apply thread configuration
	if err := engine.SetOption("Threads", fmt.Sprintf("%d", s.config.Threads)); err != nil {
		return fmt.Errorf("failed to set Threads: %v", err)
	}
	
	// Apply hash configuration
	if err := engine.SetOption("Hash", fmt.Sprintf("%d", s.config.Hash)); err != nil {
		return fmt.Errorf("failed to set Hash: %v", err)
	}
	
	// Apply contempt settings
	if err := engine.SetOption("Contempt", fmt.Sprintf("%d", s.config.Contempt)); err != nil {
		// Some Stockfish versions might not support this option, log but continue
		logrus.Debugf("Could not set Contempt option: %v", err)
	}
	
	if err := engine.SetOption("Analysis Contempt", s.config.AnalysisContempt); err != nil {
		// Some Stockfish versions might not support this option, log but continue
		logrus.Debugf("Could not set Analysis Contempt option: %v", err)
	}
	
	// Additional performance optimizations
	if err := engine.SetOption("MultiPV", "1"); err != nil {
		logrus.Debugf("Could not set MultiPV option: %v", err)
	}
	
	return nil
}

// GetEngine acquires an engine from the pool
func (s *StockfishService) GetEngine() (*uci.Engine, error) {
	select {
	case engine := <-s.available:
		return engine, nil
	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("timeout waiting for available engine")
	}
}

// ReturnEngine returns an engine to the pool
func (s *StockfishService) ReturnEngine(engine *uci.Engine) {
	// Prepare engine for next use
	if err := engine.NewGame(); err != nil {
		logrus.Errorf("Failed to reset engine for new game: %v", err)
	}
	
	select {
	case s.available <- engine:
		// Successfully returned
	default:
		logrus.Warn("Engine pool is full, this shouldn't happen")
	}
}

// AnalyzePosition analyzes a single position
func (s *StockfishService) AnalyzePosition(fen string, depth int, timeMs int, multiPV int) (*models.EngineEvaluation, []models.AlternativeMove, error) {
	engine, err := s.GetEngine()
	if err != nil {
		return nil, nil, err
	}
	defer s.ReturnEngine(engine)
	
	// Set the position
	if err := engine.SetPosition(fen, nil); err != nil {
		return nil, nil, fmt.Errorf("failed to set position: %v", err)
	}
	
	// Perform the search
	result, err := engine.Search(depth, timeMs, multiPV)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search position: %v", err)
	}
	
	// Convert result to our model
	evaluation := &models.EngineEvaluation{
		Score:              result.Score,
		Depth:              result.Depth,
		BestMove:           result.BestMove,
		PrincipalVariation: result.PrincipalVariation,
		Nodes:              result.Nodes,
		Time:               result.Time,
	}
	
	if result.ScoreType == "mate" {
		evaluation.Mate = &result.Score
	}
	
	// For now, return empty alternative moves - this would need MultiPV support
	var alternatives []models.AlternativeMove
	
	return evaluation, alternatives, nil
}

// AnalyzeGame analyzes a complete game using the Enhanced Expected Points algorithm
func (s *StockfishService) AnalyzeGame(game *models.ParsedGame, options models.AnalysisOptions, progressCallback func(int, int)) (*models.GameAnalysis, error) {
	return s.AnalyzeGameEnhanced(game, options, progressCallback)
}

// AnalyzeGameEnhanced implements the full EP-based analysis algorithm
func (s *StockfishService) AnalyzeGameEnhanced(game *models.ParsedGame, options models.AnalysisOptions, progressCallback func(int, int)) (*models.GameAnalysis, error) {
	// Set default options
	depth := options.Depth
	if depth == 0 {
		depth = 18 // Higher depth for better analysis
	}
	
	timePerMove := options.TimePerMove
	if timePerMove == 0 {
		timePerMove = 1000
	}
	
	// Initialize services for enhanced analysis
	epsService := NewExpectedPointsService()
	displayService := NewEvaluationDisplayService()
	
	// Initialize analysis structure
	analysis := &models.GameAnalysis{
		Moves:             make([]models.MoveAnalysis, 0, len(game.Moves)),
		EvaluationHistory: make([]models.EngineEvaluation, 0, len(game.Moves)),
		CriticalMoments:   make([]models.CriticalMoment, 0),
	}
	
	// Track accuracy scores
	var whiteAccuracyTotal, blackAccuracyTotal float64
	var whiteMoveCount, blackMoveCount int
	
	var currentPosition *models.ParsedMove
	var previousDisplayEval *models.DisplayEvaluation
	
	logrus.Infof("Starting enhanced EP-based analysis for %d moves", len(game.Moves))
	
	// Implementation of the core algorithm loop
	for i, move := range game.Moves {
		if progressCallback != nil {
			progressCallback(i, len(game.Moves))
		}
		
		// A. Evaluate position BEFORE the move (if not first move)
		var beforeEval *models.EngineEvaluation
		if i > 0 && currentPosition != nil {
			// Analyze the position before this move was made
			eval, _, err := s.AnalyzePosition(currentPosition.FEN, depth, timePerMove, 1)
			if err != nil {
				logrus.Errorf("Failed to analyze position before move %d: %v", i, err)
				continue
			}
			beforeEval = eval
		}
		
		// B. Calculate pre-move Expected Points
		var epBefore float64
		if beforeEval != nil {
			playerRating := s.getPlayerRating(options.PlayerRatings, move.IsWhite)
			normalizedEval := epsService.NormalizeEvaluationForPlayer(beforeEval.Score, move.IsWhite)
			epBefore = epsService.CalculateExpectedPoints(normalizedEval, playerRating)
		}
		
		// C. Apply the player's actual move (position after move)
		currentPosition = &move
		
		// D. Evaluate position AFTER the move
		afterEval, alternatives, err := s.AnalyzePosition(move.FEN, depth, timePerMove, 3) // Multi-PV for alternatives
		if err != nil {
			logrus.Errorf("Failed to analyze position after move %d: %v", i, err)
			continue
		}
		
		// E. Calculate post-move Expected Points
		playerRating := s.getPlayerRating(options.PlayerRatings, move.IsWhite)
		normalizedAfterEval := epsService.NormalizeEvaluationForPlayer(afterEval.Score, move.IsWhite)
		epAfter := epsService.CalculateExpectedPoints(normalizedAfterEval, playerRating)
		
		// F. Determine Expected Points Loss and Move Accuracy
		var epLoss float64
		var moveAccuracy float64
		if beforeEval != nil {
			epLoss = epBefore - epAfter
			moveAccuracy = epsService.CalculateMoveAccuracy(epLoss)
		} else {
			// First move - assume perfect accuracy
			epLoss = 0.0
			moveAccuracy = 100.0
		}
		
		// Add to player's total accuracy
		if move.IsWhite {
			whiteAccuracyTotal += moveAccuracy
			whiteMoveCount++
		} else {
			blackAccuracyTotal += moveAccuracy
			blackMoveCount++
		}
		
		// G. Categorize the Move
		classification := s.classifyMoveEnhanced(beforeEval, afterEval, move, alternatives, epLoss, options)
		
		// H. Create stable display evaluation
		displayEval := displayService.NormalizeForDisplay(afterEval.Score, move.IsWhite, previousDisplayEval)
		
		// Create enhanced move analysis
		moveAnalysis := models.MoveAnalysis{
			MoveNumber:       move.MoveNumber,
			Move:             move.UCI,
			SAN:              move.SAN,
			FEN:              move.FEN,
			Evaluation:       *afterEval,
			DisplayEvaluation: displayEval,
			BeforeEvaluation: beforeEval,
			Classification:   classification.String(),
			AlternativeMoves: alternatives,
			MoveAccuracy:     moveAccuracy,
			ExpectedPoints: models.ExpectedPointsData{
				Before:   epBefore,
				After:    epAfter,
				Loss:     epLoss,
				Accuracy: moveAccuracy,
			},
		}
		
		// Update for next iteration
		previousDisplayEval = displayEval
		
		// Set book move flag if in opening
		if move.MoveNumber <= 15 && classification == models.Book {
			moveAnalysis.IsBookMove = true
		}
		
		analysis.Moves = append(analysis.Moves, moveAnalysis)
		analysis.EvaluationHistory = append(analysis.EvaluationHistory, *afterEval)
		
		// Check for critical moments with enhanced detection
		if beforeEval != nil && s.isCriticalMomentEnhanced(beforeEval, afterEval, epLoss) {
			critical := models.CriticalMoment{
				MoveNumber: move.MoveNumber,
				BeforeEval: beforeEval.Score,
				AfterEval:  afterEval.Score,
				Description: s.describeCriticalMoment(epLoss, classification),
			}
			
			if afterEval.Score > beforeEval.Score {
				critical.Advantage = "white"
			} else {
				critical.Advantage = "black"
			}
			
			analysis.CriticalMoments = append(analysis.CriticalMoments, critical)
		}
		
		// Move to next iteration
	}
	
	// 3. Finalization - Calculate final accuracies
	finalWhiteAccuracy := 0.0
	finalBlackAccuracy := 0.0
	
	if whiteMoveCount > 0 {
		finalWhiteAccuracy = whiteAccuracyTotal / float64(whiteMoveCount)
	}
	if blackMoveCount > 0 {
		finalBlackAccuracy = blackAccuracyTotal / float64(blackMoveCount)
	}
	
	// Calculate enhanced statistics
	analysis.WhiteStats = s.calculateEnhancedPlayerStats(analysis.Moves, true, finalWhiteAccuracy)
	analysis.BlackStats = s.calculateEnhancedPlayerStats(analysis.Moves, false, finalBlackAccuracy)
	
	// Determine game phases with better detection
	analysis.GamePhases = s.determineGamePhases(analysis.Moves)
	analysis.PhaseAnalysis = s.calculatePhaseAnalysis(analysis.Moves, analysis.GamePhases)
	
	logrus.Infof("Enhanced analysis complete: White accuracy: %.1f%%, Black accuracy: %.1f%%", 
		finalWhiteAccuracy, finalBlackAccuracy)
	
	return analysis, nil
}

// classifyMove classifies a move based on evaluation change
func (s *StockfishService) classifyMove(prevEval, currEval *models.EngineEvaluation, isWhite bool) models.MoveClassification {
	if prevEval == nil {
		return models.Book
	}
	
	// Calculate centipawn difference
	var evalDiff int
	if isWhite {
		evalDiff = currEval.Score - prevEval.Score
	} else {
		evalDiff = prevEval.Score - currEval.Score
	}
	
	// Classify based on centipawn loss
	switch {
	case evalDiff >= 100:
		return models.Brilliant
	case evalDiff >= 50:
		return models.Great
	case evalDiff >= 0:
		return models.Best
	case evalDiff >= -20:
		return models.Excellent
	case evalDiff >= -50:
		return models.Good
	case evalDiff >= -100:
		return models.Inaccuracy
	case evalDiff >= -200:
		return models.Mistake
	default:
		return models.Blunder
	}
}

// isCriticalMoment determines if there's a significant evaluation swing
func (s *StockfishService) isCriticalMoment(prevEval, currEval *models.EngineEvaluation) bool {
	diff := abs(currEval.Score - prevEval.Score)
	return diff >= 150 // 1.5 pawn evaluation swing
}

// calculatePlayerStats calculates statistics for a player
func (s *StockfishService) calculatePlayerStats(moves []models.MoveAnalysis, isWhite bool) models.PlayerStatistics {
	stats := models.PlayerStatistics{
		MoveCounts: models.MoveCounts{},
	}
	
	var totalMoves, accurateSum int
	
	for _, move := range moves {
		// Check if this move belongs to the player
		moveIsWhite := (move.MoveNumber % 2) == 1
		if moveIsWhite != isWhite {
			continue
		}
		
		totalMoves++
		
		// Count move types
		switch move.Classification {
		case "brilliant":
			stats.MoveCounts.Brilliant++
			accurateSum += 100
		case "great":
			stats.MoveCounts.Great++
			accurateSum += 95
		case "best":
			stats.MoveCounts.Best++
			accurateSum += 90
		case "excellent":
			stats.MoveCounts.Excellent++
			accurateSum += 85
		case "good":
			stats.MoveCounts.Good++
			accurateSum += 80
		case "book":
			stats.MoveCounts.Book++
			accurateSum += 85
		case "inaccuracy":
			stats.MoveCounts.Inaccuracy++
			accurateSum += 70
		case "mistake":
			stats.MoveCounts.Mistake++
			accurateSum += 50
		case "blunder":
			stats.MoveCounts.Blunder++
			accurateSum += 30
		case "miss":
			stats.MoveCounts.Miss++
			accurateSum += 20
		}
	}
	
	if totalMoves > 0 {
		stats.Accuracy = float64(accurateSum) / float64(totalMoves)
	}
	
	return stats
}

// calculatePhaseAnalysis calculates accuracy by game phase
func (s *StockfishService) calculatePhaseAnalysis(moves []models.MoveAnalysis, phases models.GamePhases) models.PhaseAnalysis {
	var openingAccuracy, middlegameAccuracy, endgameAccuracy float64
	var openingCount, middlegameCount, endgameCount int
	
	for _, move := range moves {
		var accuracy float64
		switch move.Classification {
		case "brilliant", "great":
			accuracy = 100
		case "best":
			accuracy = 90
		case "excellent":
			accuracy = 85
		case "good", "book":
			accuracy = 80
		case "inaccuracy":
			accuracy = 70
		case "mistake":
			accuracy = 50
		case "blunder":
			accuracy = 30
		case "miss":
			accuracy = 20
		}
		
		if move.MoveNumber <= phases.Opening {
			openingAccuracy += accuracy
			openingCount++
		} else if move.MoveNumber <= phases.Middlegame {
			middlegameAccuracy += accuracy
			middlegameCount++
		} else {
			endgameAccuracy += accuracy
			endgameCount++
		}
	}
	
	return models.PhaseAnalysis{
		OpeningAccuracy:    safeDiv(openingAccuracy, float64(openingCount)),
		MiddlegameAccuracy: safeDiv(middlegameAccuracy, float64(middlegameCount)),
		EndgameAccuracy:    safeDiv(endgameAccuracy, float64(endgameCount)),
	}
}

// UpdateConfig updates the engine configuration
func (s *StockfishService) UpdateConfig(config models.EngineOptions) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	s.config = config
	
	// Update all engines
	for _, engine := range s.engines {
		if err := s.configureEngineOptimized(engine); err != nil {
			return fmt.Errorf("failed to update engine configuration: %v", err)
		}
	}
	
	return nil
}

// GetConfig returns the current engine configuration
func (s *StockfishService) GetConfig() models.EngineOptions {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.config
}

// Shutdown closes all engines
func (s *StockfishService) Shutdown() {
	logrus.Info("Shutting down Stockfish service")
	
	// Close all engines
	for _, engine := range s.engines {
		if err := engine.Quit(); err != nil {
			logrus.Errorf("Failed to quit engine: %v", err)
		}
	}
	
	// Clear the pool
	s.engines = nil
	close(s.available)
	
	logrus.Info("Stockfish service shutdown complete")
}

// Helper functions
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func safeDiv(numerator, denominator float64) float64 {
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

// Enhanced analysis helper methods for EP-based algorithm

// getPlayerRating gets the rating for the current player
func (s *StockfishService) getPlayerRating(ratings models.PlayerRatings, isWhite bool) int {
	if isWhite {
		if ratings.White > 0 {
			return ratings.White
		}
	} else {
		if ratings.Black > 0 {
			return ratings.Black
		}
	}
	return 1500 // Default rating
}

// classifyMoveEnhanced uses the sophisticated EP-based categorization
func (s *StockfishService) classifyMoveEnhanced(beforeEval, afterEval *models.EngineEvaluation, move models.ParsedMove, alternatives []models.AlternativeMove, epLoss float64, options models.AnalysisOptions) models.MoveClassification {
	// Extract best move
	bestMove := ""
	if beforeEval != nil {
		bestMove = beforeEval.BestMove
	}
	
	// Use simple classification for now - can be enhanced with MoveCategorizer
	return s.classifyMoveSimple(epLoss, move.MoveNumber, bestMove, move.UCI)
}

// classifyMoveSimple provides basic classification based on EP loss
func (s *StockfishService) classifyMoveSimple(epLoss float64, moveNumber int, bestMove, playedMove string) models.MoveClassification {
	// Book moves in opening
	if moveNumber <= 12 && epLoss <= 0.03 {
		return models.Book
	}
	
	// Best move
	if strings.EqualFold(bestMove, playedMove) {
		return models.Best
	}
	
	// Classification by EP loss
	switch {
	case epLoss <= 0.02:
		return models.Excellent
	case epLoss <= 0.05:
		return models.Good
	case epLoss <= 0.10:
		return models.Inaccuracy
	case epLoss <= 0.20:
		return models.Mistake
	default:
		return models.Blunder
	}
}

// isCriticalMomentEnhanced detects critical moments using EP loss
func (s *StockfishService) isCriticalMomentEnhanced(beforeEval, afterEval *models.EngineEvaluation, epLoss float64) bool {
	// Significant EP loss indicates a critical moment
	if epLoss > 0.15 {
		return true
	}
	
	// Large evaluation swings
	evalDiff := abs(afterEval.Score - beforeEval.Score)
	return evalDiff > 150 // 150 centipawn swing
}

// describeCriticalMoment provides description for critical moments
func (s *StockfishService) describeCriticalMoment(epLoss float64, classification models.MoveClassification) string {
	switch classification {
	case models.Blunder:
		return "Major blunder changes game outcome"
	case models.Mistake:
		return "Significant mistake in crucial position"
	case models.Brilliant:
		return "Brilliant move in critical position"
	default:
		if epLoss > 0.2 {
			return "Game-changing moment"
		}
		return "Critical decision point"
	}
}

// calculateEnhancedPlayerStats calculates player statistics with accuracy
func (s *StockfishService) calculateEnhancedPlayerStats(moves []models.MoveAnalysis, isWhite bool, accuracy float64) models.PlayerStatistics {
	stats := models.PlayerStatistics{
		Accuracy:   accuracy,
		MoveCounts: models.MoveCounts{},
	}
	
	for _, move := range moves {
		// Check if this move belongs to the player
		isMoveForPlayer := (move.MoveNumber%2 == 1) == isWhite
		if !isMoveForPlayer {
			continue
		}
		
		// Count move types
		switch move.Classification {
		case "brilliant":
			stats.MoveCounts.Brilliant++
		case "great":
			stats.MoveCounts.Great++
		case "best":
			stats.MoveCounts.Best++
		case "excellent":
			stats.MoveCounts.Excellent++
		case "good":
			stats.MoveCounts.Good++
		case "book":
			stats.MoveCounts.Book++
		case "inaccuracy":
			stats.MoveCounts.Inaccuracy++
		case "mistake":
			stats.MoveCounts.Mistake++
		case "blunder":
			stats.MoveCounts.Blunder++
		}
	}
	
	return stats
}

// determineGamePhases determines game phases based on move analysis
func (s *StockfishService) determineGamePhases(moves []models.MoveAnalysis) models.GamePhases {
	totalMoves := len(moves)
	
	// Default phase boundaries
	opening := min(15, totalMoves/4)
	middlegame := min(35, totalMoves*3/4)
	
	// Could be enhanced with piece count analysis, etc.
	return models.GamePhases{
		Opening:    opening,
		Middlegame: middlegame,
		Endgame:    totalMoves,
	}
}

// isPositionWinning determines if a position is overwhelmingly winning
func (s *StockfishService) isPositionWinning(eval *models.EngineEvaluation, isWhiteToMove bool) bool {
	if eval.Mate != nil {
		mate := *eval.Mate
		if isWhiteToMove {
			return mate > 0
		} else {
			return mate < 0
		}
	}
	
	// Consider winning if advantage > 300 centipawns
	threshold := 300
	if isWhiteToMove {
		return eval.Score > threshold
	} else {
		return eval.Score < -threshold
	}
} 