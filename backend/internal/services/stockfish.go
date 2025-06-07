package services

import (
	"fmt"
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
}

// NewStockfishService creates a new Stockfish service
func NewStockfishService(maxWorkers int, binaryPath string) *StockfishService {
	return &StockfishService{
		engines:    make([]*uci.Engine, 0, maxWorkers),
		available:  make(chan *uci.Engine, maxWorkers),
		maxWorkers: maxWorkers,
		binaryPath: binaryPath,
		config: models.EngineOptions{
			Threads:          1,
			Hash:             128,
			Contempt:         0,
			AnalysisContempt: "off",
		},
	}
}

// Initialize creates and initializes the engine pool
func (s *StockfishService) Initialize() error {
	logrus.Infof("Initializing Stockfish service with %d workers", s.maxWorkers)
	
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
		
		// Configure engine with default settings
		if err := s.configureEngine(engine); err != nil {
			logrus.Errorf("Failed to configure engine %d: %v", i, err)
			return fmt.Errorf("failed to configure engine %d: %v", i, err)
		}
		
		s.engines = append(s.engines, engine)
		s.available <- engine
		
		logrus.Debugf("Engine %d initialized successfully", i)
	}
	
	logrus.Infof("Stockfish service initialized with %d engines", len(s.engines))
	return nil
}

// configureEngine applies configuration to an engine
func (s *StockfishService) configureEngine(engine *uci.Engine) error {
	if err := engine.SetOption("Threads", fmt.Sprintf("%d", s.config.Threads)); err != nil {
		return err
	}
	if err := engine.SetOption("Hash", fmt.Sprintf("%d", s.config.Hash)); err != nil {
		return err
	}
	if err := engine.SetOption("Contempt", fmt.Sprintf("%d", s.config.Contempt)); err != nil {
		return err
	}
	if err := engine.SetOption("Analysis Contempt", s.config.AnalysisContempt); err != nil {
		return err
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

// AnalyzeGame analyzes a complete game
func (s *StockfishService) AnalyzeGame(game *models.ParsedGame, options models.AnalysisOptions, progressCallback func(int, int)) (*models.GameAnalysis, error) {
	// Set default options
	depth := options.Depth
	if depth == 0 {
		depth = 15
	}
	
	timePerMove := options.TimePerMove
	if timePerMove == 0 {
		timePerMove = 1000
	}
	
	analysis := &models.GameAnalysis{
		Moves:             make([]models.MoveAnalysis, 0, len(game.Moves)),
		EvaluationHistory: make([]models.EngineEvaluation, 0, len(game.Moves)),
		CriticalMoments:   make([]models.CriticalMoment, 0),
	}
	
	var previousEval *models.EngineEvaluation
	
	// Analyze each position
	for i, move := range game.Moves {
		if progressCallback != nil {
			progressCallback(i, len(game.Moves))
		}
		
		// Analyze the position after this move
		eval, _, err := s.AnalyzePosition(move.FEN, depth, timePerMove, 1)
		if err != nil {
			logrus.Errorf("Failed to analyze position %d: %v", i, err)
			continue
		}
		
		// Classify the move
		classification := s.classifyMove(previousEval, eval, move.IsWhite)
		
		moveAnalysis := models.MoveAnalysis{
			MoveNumber:     move.MoveNumber,
			Move:           move.UCI,
			SAN:            move.SAN,
			FEN:            move.FEN,
			Evaluation:     *eval,
			Classification: classification.String(),
		}
		
		analysis.Moves = append(analysis.Moves, moveAnalysis)
		analysis.EvaluationHistory = append(analysis.EvaluationHistory, *eval)
		
		// Check for critical moments
		if previousEval != nil {
			if s.isCriticalMoment(previousEval, eval) {
				critical := models.CriticalMoment{
					MoveNumber: move.MoveNumber,
					BeforeEval: previousEval.Score,
					AfterEval:  eval.Score,
				}
				
				if eval.Score > previousEval.Score {
					critical.Advantage = "white"
				} else {
					critical.Advantage = "black"
				}
				
				analysis.CriticalMoments = append(analysis.CriticalMoments, critical)
			}
		}
		
		previousEval = eval
	}
	
	// Calculate statistics
	analysis.WhiteStats = s.calculatePlayerStats(analysis.Moves, true)
	analysis.BlackStats = s.calculatePlayerStats(analysis.Moves, false)
	
	// Determine game phases (simplified)
	analysis.GamePhases = models.GamePhases{
		Opening:    min(15, len(game.Moves)/3),
		Middlegame: min(40, len(game.Moves)*2/3),
		Endgame:    len(game.Moves),
	}
	
	// Calculate phase accuracy
	analysis.PhaseAnalysis = s.calculatePhaseAnalysis(analysis.Moves, analysis.GamePhases)
	
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
		if err := s.configureEngine(engine); err != nil {
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