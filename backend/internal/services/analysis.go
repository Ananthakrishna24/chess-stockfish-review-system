package services

import (
	"fmt"
	"sync"
	"time"

	"chess-backend/internal/models"

	"github.com/notnil/chess"
	"github.com/sirupsen/logrus"
)

// AnalysisService orchestrates game and position analysis
type AnalysisService struct {
	stockfishService      *StockfishService
	chessService          *ChessService
	cacheService          *CacheService
	playerService         *PlayerService
	openingService        *OpeningService
	enhancedAnalysisService *EnhancedAnalysisService
	activeJobs            map[string]*models.AnalysisJob
	jobsMutex             sync.RWMutex
}

// NewAnalysisService creates a new analysis service
func NewAnalysisService(stockfish *StockfishService, chess *ChessService, cache *CacheService, player *PlayerService, opening *OpeningService) *AnalysisService {
	// Initialize enhanced analysis service
	enhancedService := NewEnhancedAnalysisService(stockfish, chess, cache, player, opening)
	
	return &AnalysisService{
		stockfishService:        stockfish,
		chessService:            chess,
		cacheService:            cache,
		playerService:           player,
		openingService:          opening,
		enhancedAnalysisService: enhancedService,
		activeJobs:              make(map[string]*models.AnalysisJob),
	}
}

// StartGameAnalysis starts asynchronous game analysis
func (s *AnalysisService) StartGameAnalysis(pgn string, options models.AnalysisOptions) string {
	// Generate game ID
	gameID := s.cacheService.GenerateGameID(pgn)
	
	// Check if already in cache
	if _, found := s.cacheService.GetAnalysis(gameID); found {
		logrus.Debugf("Analysis for game %s found in cache", gameID)
		return gameID
	}
	
	// Check if already being processed
	s.jobsMutex.RLock()
	if _, exists := s.activeJobs[gameID]; exists {
		s.jobsMutex.RUnlock()
		logrus.Debugf("Analysis for game %s already in progress", gameID)
		return gameID
	}
	s.jobsMutex.RUnlock()
	
	// Create new analysis job
	job := &models.AnalysisJob{
		ID:        gameID,
		PGN:       pgn,
		Options:   options,
		Status:    models.StatusQueued,
		CreatedAt: time.Now(),
	}
	
	// Store job
	s.jobsMutex.Lock()
	s.activeJobs[gameID] = job
	s.jobsMutex.Unlock()
	
	// Start analysis in goroutine
	go s.processGameAnalysis(job)
	
	logrus.Infof("Started analysis for game %s", gameID)
	return gameID
}

// processGameAnalysis performs the actual game analysis
func (s *AnalysisService) processGameAnalysis(job *models.AnalysisJob) {
	defer func() {
		// Remove from active jobs when done
		s.jobsMutex.Lock()
		delete(s.activeJobs, job.ID)
		s.jobsMutex.Unlock()
	}()
	
	job.SetStatus(models.StatusAnalyzing)
	
	// Choose analysis method based on options
	var response *models.GameAnalysisResponse
	var err error
	
	// Use Enhanced EP-based analysis if player ratings are provided or specifically requested
	if (job.Options.PlayerRatings.White > 0 || job.Options.PlayerRatings.Black > 0) && s.enhancedAnalysisService != nil {
		logrus.Infof("Using Enhanced EP-based analysis for game %s", job.ID)
		response, err = s.enhancedAnalysisService.AnalyzeGameWithEP(job.PGN, job.Options, func(current, total int) {
			job.UpdateProgress(current, total)
		})
	} else {
		// Fall back to standard analysis
		logrus.Infof("Using standard analysis for game %s", job.ID)
		response, err = s.processStandardAnalysis(job)
	}
	
	if err != nil {
		logrus.Errorf("Failed to analyze game %s: %v", job.ID, err)
		job.SetError(fmt.Sprintf("Analysis failed: %v", err))
		return
	}
	
	// Store in cache
	s.cacheService.StoreAnalysis(job.ID, *response)
	
	// Update player statistics
	if s.playerService != nil {
		s.playerService.RecordGameAnalysis(response)
		logrus.Debugf("Updated player statistics for game %s", job.ID)
	}
	
	// Update job with result
	job.SetResult(response)
	
	logrus.Infof("Completed analysis for game %s in %.2f seconds", 
		job.ID, response.ProcessingTime)
}

// processStandardAnalysis performs the legacy standard analysis
func (s *AnalysisService) processStandardAnalysis(job *models.AnalysisJob) (*models.GameAnalysisResponse, error) {
	// Parse PGN
	parsedGame, err := s.chessService.ParsePGN(job.PGN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PGN: %v", err)
	}
	
	job.UpdateProgress(0, parsedGame.TotalMoves)
	
	// Perform standard analysis
	analysis, err := s.stockfishService.AnalyzeGame(parsedGame, job.Options, func(current, total int) {
		job.UpdateProgress(current, total)
	})
	if err != nil {
		return nil, fmt.Errorf("analysis failed: %v", err)
	}
	
	// Create response
	response := &models.GameAnalysisResponse{
		GameID:         job.ID,
		GameInfo:       parsedGame.GameInfo,
		Analysis:       *analysis,
		ProcessingTime: time.Since(job.CreatedAt).Seconds(),
		Timestamp:      time.Now(),
	}
	
	return response, nil
}

// GetAnalysisProgress returns the progress of an analysis job
func (s *AnalysisService) GetAnalysisProgress(gameID string) (*models.ProgressResponse, error) {
	// Check cache first
	if _, found := s.cacheService.GetAnalysis(gameID); found {
		return &models.ProgressResponse{
			GameID: gameID,
			Status: string(models.StatusCompleted),
			Progress: models.ProgressDetails{
				CurrentMove: 0,
				TotalMoves:  0,
				Percentage:  100.0,
			},
		}, nil
	}
	
	// Check active jobs
	s.jobsMutex.RLock()
	job, exists := s.activeJobs[gameID]
	s.jobsMutex.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("analysis job not found")
	}
	
	progress := job.GetProgress()
	return &progress, nil
}

// GetAnalysisResult returns the completed analysis result
func (s *AnalysisService) GetAnalysisResult(gameID string) (*models.GameAnalysisResponse, error) {
	// Check cache first
	if result, found := s.cacheService.GetAnalysis(gameID); found {
		return result, nil
	}
	
	// Check if job is completed
	s.jobsMutex.RLock()
	job, exists := s.activeJobs[gameID]
	s.jobsMutex.RUnlock()
	
	if exists && job.Status == models.StatusCompleted && job.Result != nil {
		return job.Result, nil
	}
	
	return nil, fmt.Errorf("analysis not found or not completed")
}

// AnalyzePosition analyzes a single chess position
func (s *AnalysisService) AnalyzePosition(request models.AnalyzePositionRequest) (*models.PositionAnalysisResponse, error) {
	// Validate FEN
	if err := s.chessService.ValidateFEN(request.FEN); err != nil {
		return nil, fmt.Errorf("invalid FEN: %v", err)
	}
	
	// Set defaults
	depth := request.Depth
	if depth == 0 {
		depth = 15
	}
	
	timeLimit := request.TimeLimit
	if timeLimit == 0 {
		timeLimit = 5000
	}
	
	multiPV := request.MultiPV
	if multiPV == 0 {
		multiPV = 1
	}
	
	// Analyze position
	evaluation, alternatives, err := s.stockfishService.AnalyzePosition(
		request.FEN, depth, timeLimit, multiPV)
	if err != nil {
		return nil, fmt.Errorf("analysis failed: %v", err)
	}
	
	// Get position information
	position, err := s.chessService.GetPositionFromFEN(request.FEN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse position: %v", err)
	}
	
	// Calculate additional position info
	phase := s.chessService.GetGamePhase(position, 20) // Assume middlegame for single positions
	whiteMaterial := s.chessService.CalculateMaterialValue(position, chess.White)
	blackMaterial := s.chessService.CalculateMaterialValue(position, chess.Black)
	whiteKingSafety := s.chessService.AssessKingSafety(position, chess.White)
	blackKingSafety := s.chessService.AssessKingSafety(position, chess.Black)
	
	response := &models.PositionAnalysisResponse{
		FEN:              request.FEN,
		Evaluation:       *evaluation,
		AlternativeMoves: alternatives,
		PositionInfo: models.PositionInfo{
			Phase: phase.String(),
			Material: models.MaterialInfo{
				White: whiteMaterial,
				Black: blackMaterial,
			},
			Safety: models.SafetyInfo{
				WhiteKing: whiteKingSafety,
				BlackKing: blackKingSafety,
			},
		},
	}
	
	return response, nil
}

// GetEngineConfig returns current engine configuration
func (s *AnalysisService) GetEngineConfig() *models.EngineConfigResponse {
	config := s.stockfishService.GetConfig()
	
	return &models.EngineConfigResponse{
		Version:  "Stockfish 16", // This would be detected from engine
		Features: []string{"UCI", "MultiPV", "Hash", "Threads"},
		Limits: models.EngineLimits{
			MaxDepth: 24,
			MaxTime:  30000,
			MaxNodes: 10000000,
		},
		CurrentConfig: config,
	}
}

// UpdateEngineConfig updates engine configuration
func (s *AnalysisService) UpdateEngineConfig(request models.UpdateEngineConfigRequest) error {
	currentConfig := s.stockfishService.GetConfig()
	
	// Update only provided fields
	if request.Threads != nil {
		currentConfig.Threads = *request.Threads
	}
	if request.Hash != nil {
		currentConfig.Hash = *request.Hash
	}
	if request.Contempt != nil {
		currentConfig.Contempt = *request.Contempt
	}
	if request.AnalysisContempt != nil {
		currentConfig.AnalysisContempt = *request.AnalysisContempt
	}
	
	// Validate configuration
	if currentConfig.Threads < 1 || currentConfig.Threads > 8 {
		return fmt.Errorf("threads must be between 1 and 8")
	}
	if currentConfig.Hash < 1 || currentConfig.Hash > 2048 {
		return fmt.Errorf("hash must be between 1 and 2048 MB")
	}
	if currentConfig.Contempt < -100 || currentConfig.Contempt > 100 {
		return fmt.Errorf("contempt must be between -100 and 100")
	}
	
	// Apply configuration
	return s.stockfishService.UpdateConfig(currentConfig)
}

// GetStats returns analysis service statistics
func (s *AnalysisService) GetStats() map[string]interface{} {
	s.jobsMutex.RLock()
	activeJobCount := len(s.activeJobs)
	
	jobStatuses := make(map[string]int)
	for _, job := range s.activeJobs {
		jobStatuses[string(job.Status)]++
	}
	s.jobsMutex.RUnlock()
	
	cacheStats := s.cacheService.GetStats()
	
	stats := map[string]interface{}{
		"active_jobs":    activeJobCount,
		"job_statuses":   jobStatuses,
		"cache_stats":    cacheStats,
		"timestamp":      time.Now(),
	}
	
	// Add EP analysis availability
	if s.enhancedAnalysisService != nil {
		stats["enhanced_analysis_available"] = true
	}
	
	return stats
}

// StartEnhancedGameAnalysis starts EP-based game analysis directly
func (s *AnalysisService) StartEnhancedGameAnalysis(pgn string, options models.AnalysisOptions) string {
	if s.enhancedAnalysisService == nil {
		logrus.Warn("Enhanced analysis service not available, falling back to standard analysis")
		return s.StartGameAnalysis(pgn, options)
	}
	
	// Validate and set EP-specific options
	s.enhancedAnalysisService.ValidateEPAnalysisOptions(&options)
	
	// Use standard job creation but force enhanced analysis
	gameID := s.cacheService.GenerateGameID(pgn)
	
	// Check cache first
	if _, found := s.cacheService.GetAnalysis(gameID); found {
		logrus.Debugf("Enhanced analysis for game %s found in cache", gameID)
		return gameID
	}
	
	// Create job with enhanced flag
	job := &models.AnalysisJob{
		ID:        gameID,
		PGN:       pgn,
		Options:   options,
		Status:    models.StatusQueued,
		CreatedAt: time.Now(),
	}
	
	// Store job
	s.jobsMutex.Lock()
	s.activeJobs[gameID] = job
	s.jobsMutex.Unlock()
	
	// Start enhanced analysis
	go s.processEnhancedGameAnalysis(job)
	
	logrus.Infof("Started enhanced EP analysis for game %s", gameID)
	return gameID
}

// processEnhancedGameAnalysis processes EP-based analysis
func (s *AnalysisService) processEnhancedGameAnalysis(job *models.AnalysisJob) {
	defer func() {
		s.jobsMutex.Lock()
		delete(s.activeJobs, job.ID)
		s.jobsMutex.Unlock()
	}()
	
	job.SetStatus(models.StatusAnalyzing)
	
	// Use enhanced analysis
	response, err := s.enhancedAnalysisService.AnalyzeGameWithEP(job.PGN, job.Options, func(current, total int) {
		job.UpdateProgress(current, total)
	})
	
	if err != nil {
		logrus.Errorf("Enhanced analysis failed for game %s: %v", job.ID, err)
		job.SetError(fmt.Sprintf("Enhanced analysis failed: %v", err))
		return
	}
	
	// Store results
	s.cacheService.StoreAnalysis(job.ID, *response)
	
	if s.playerService != nil {
		s.playerService.RecordGameAnalysis(response)
	}
	
	job.SetResult(response)
	logrus.Infof("Enhanced analysis completed for game %s", job.ID)
}

// GetEnhancedAnalysisService returns the enhanced analysis service for direct access
func (s *AnalysisService) GetEnhancedAnalysisService() *EnhancedAnalysisService {
	return s.enhancedAnalysisService
}

// AnalyzePositionWithEP analyzes a position using Expected Points
func (s *AnalysisService) AnalyzePositionWithEP(fen string, playerRating int, isWhiteToMove bool) (float64, error) {
	if s.enhancedAnalysisService == nil {
		return 0, fmt.Errorf("enhanced analysis service not available")
	}
	
	return s.enhancedAnalysisService.CalculatePositionEP(fen, playerRating, isWhiteToMove)
} 