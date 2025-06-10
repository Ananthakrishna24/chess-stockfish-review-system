package services

import (
	"chess-backend/internal/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/notnil/chess"
	"github.com/sirupsen/logrus"
)

// CalibrationService handles the calibration phase of the EP algorithm
type CalibrationService struct {
	stockfishService     *StockfishService
	expectedPointsService *ExpectedPointsService
	chessService         *ChessService
	logger               *logrus.Logger
	thresholdsPath       string
	calibrationData      *models.CalibrationData
}

// NewCalibrationService creates a new calibration service
func NewCalibrationService(stockfish *StockfishService, eps *ExpectedPointsService, chess *ChessService, logger *logrus.Logger) *CalibrationService {
	return &CalibrationService{
		stockfishService:     stockfish,
		expectedPointsService: eps,
		chessService:         chess,
		logger:               logger,
		thresholdsPath:       "data/thresholds.json",
		calibrationData:      &models.CalibrationData{
			RatingBuckets: make(map[models.RatingBucket][]models.MoveStat),
			Thresholds:    make(map[models.RatingBucket]models.EPThresholds),
			Version:       "1.0",
		},
	}
}

// RunCalibration processes a PGN archive to collect move statistics
func (cs *CalibrationService) RunCalibration(pgnArchivePath string) error {
	cs.logger.Info("Starting calibration phase...")
	
	// Read PGN file
	pgnContent, err := ioutil.ReadFile(pgnArchivePath)
	if err != nil {
		return fmt.Errorf("failed to read PGN file: %w", err)
	}
	
	// Parse games from PGN
	games, err := cs.parseMultiplePGNGames(string(pgnContent))
	if err != nil {
		return fmt.Errorf("failed to parse PGN games: %w", err)
	}
	
	cs.logger.Infof("Processing %d games for calibration...", len(games))
	
	totalMoves := 0
	for i, game := range games {
		if i%100 == 0 {
			cs.logger.Infof("Processing game %d/%d...", i+1, len(games))
		}
		
		moves, err := cs.processGameForCalibration(game)
		if err != nil {
			cs.logger.Warnf("Failed to process game %d: %v", i+1, err)
			continue
		}
		
		totalMoves += len(moves)
		cs.calibrationData.GamesParsed++
	}
	
	cs.calibrationData.TotalMoves = totalMoves
	cs.calibrationData.LastUpdated = time.Now()
	
	cs.logger.Infof("Calibration complete: %d moves from %d games", totalMoves, cs.calibrationData.GamesParsed)
	
	// Save raw calibration data
	if err := cs.saveCalibrationData(); err != nil {
		return fmt.Errorf("failed to save calibration data: %w", err)
	}
	
	return nil
}

// ComputePercentiles calculates percentile thresholds for each rating bucket
func (cs *CalibrationService) ComputePercentiles() error {
	cs.logger.Info("Computing percentiles from calibration data...")
	
	for bucket, moves := range cs.calibrationData.RatingBuckets {
		if len(moves) < 100 {
			cs.logger.Warnf("Insufficient data for bucket %s: %d moves", bucket, len(moves))
			continue
		}
		
		// Extract EP losses and sort them
		epLosses := make([]float64, len(moves))
		for i, move := range moves {
			epLosses[i] = move.EPLoss
		}
		sort.Float64s(epLosses)
		
		// Calculate percentiles
		thresholds := models.EPThresholds{
			P1:  cs.calculatePercentile(epLosses, 1),
			P5:  cs.calculatePercentile(epLosses, 5),
			P10: cs.calculatePercentile(epLosses, 10),
			P25: cs.calculatePercentile(epLosses, 25),
			P50: cs.calculatePercentile(epLosses, 50),
			P75: cs.calculatePercentile(epLosses, 75),
			P90: cs.calculatePercentile(epLosses, 90),
		}
		
		cs.calibrationData.Thresholds[bucket] = thresholds
		
		cs.logger.Infof("Bucket %s thresholds: P1=%.4f, P5=%.4f, P10=%.4f, P25=%.4f, P50=%.4f, P75=%.4f, P90=%.4f",
			bucket, thresholds.P1, thresholds.P5, thresholds.P10, thresholds.P25, thresholds.P50, thresholds.P75, thresholds.P90)
	}
	
	// Save thresholds to JSON file
	if err := cs.saveThresholds(); err != nil {
		return fmt.Errorf("failed to save thresholds: %w", err)
	}
	
	cs.logger.Info("Percentile computation complete")
	return nil
}

// LoadThresholds loads thresholds from JSON file
func (cs *CalibrationService) LoadThresholds() (map[models.RatingBucket]models.EPThresholds, error) {
	if _, err := os.Stat(cs.thresholdsPath); os.IsNotExist(err) {
		// Return default thresholds if file doesn't exist
		return cs.getDefaultThresholds(), nil
	}
	
	data, err := ioutil.ReadFile(cs.thresholdsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read thresholds file: %w", err)
	}
	
	var thresholds map[models.RatingBucket]models.EPThresholds
	if err := json.Unmarshal(data, &thresholds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal thresholds: %w", err)
	}
	
	return thresholds, nil
}

// processGameForCalibration processes a single game to extract move statistics
func (cs *CalibrationService) processGameForCalibration(game *chess.Game) ([]models.MoveStat, error) {
	var moves []models.MoveStat
	
	// Extract player ratings from headers
	whiteRating := cs.extractRatingFromHeaders(game, "WhiteElo")
	blackRating := cs.extractRatingFromHeaders(game, "BlackElo")
	
	if whiteRating == 0 || blackRating == 0 {
		return nil, fmt.Errorf("missing player ratings")
	}
	
	// Replay the game and analyze each move
	pos := chess.StartingPosition()
	gamePos := chess.NewGame()
	
	for i, move := range game.Moves() {
		isWhiteToMove := pos.Turn() == chess.White
		playerRating := whiteRating
		if !isWhiteToMove {
			playerRating = blackRating
		}
		
		// Get evaluation before the move
		evalBefore, _, err := cs.stockfishService.AnalyzePosition(pos.String(), 12, 1000, 1)
		if err != nil {
			continue // Skip moves that can't be evaluated
		}
		
		// Make the move
		if err := gamePos.Move(move); err != nil {
			continue
		}
		pos = gamePos.Position()
		
		// Get evaluation after the move
		evalAfter, _, err := cs.stockfishService.AnalyzePosition(pos.String(), 12, 1000, 1)
		if err != nil {
			continue
		}
		
		// Calculate EP loss
		epBefore := cs.expectedPointsService.CalculateExpectedPoints(evalBefore.Score, playerRating)
		epAfter := cs.expectedPointsService.CalculateExpectedPoints(evalAfter.Score, playerRating)
		epLoss := epBefore - epAfter
		
		// Get UCI moves
		moveUCI := cs.moveToUCI(move)
		bestEngineUCI := evalBefore.BestMove
		
		// Determine if this was the best move by UCI comparison
		isBestMove := moveUCI == bestEngineUCI
		
		// Calculate material change (simplified - would need proper material counting)
		materialChange := 0 // This would require proper material balance calculation
		
		// Determine if it's a sacrifice
		isSacrifice := materialChange < -100 // Lost more than a pawn's worth
		
		// Create move stat
		moveStat := models.MoveStat{
			Rating:         playerRating,
			EPLoss:         epLoss,
			MaterialChange: materialChange,
			MoveUCI:        moveUCI,
			BestEngineUCI:  bestEngineUCI,
			IsBestMove:     isBestMove,
			IsSacrifice:    isSacrifice,
			MoveNumber:     i + 1,
			GamePhase:      cs.determineGamePhase(i + 1),
		}
		
		// Add to appropriate bucket
		bucket := cs.getRatingBucket(playerRating)
		cs.calibrationData.RatingBuckets[bucket] = append(cs.calibrationData.RatingBuckets[bucket], moveStat)
		
		moves = append(moves, moveStat)
	}
	
	return moves, nil
}

// Helper functions

func (cs *CalibrationService) parseMultiplePGNGames(pgnContent string) ([]*chess.Game, error) {
	var games []*chess.Game
	
	// Split PGN content by games (games are typically separated by blank lines)
	gameStrings := strings.Split(pgnContent, "\n\n\n")
	
	for _, gameString := range gameStrings {
		gameString = strings.TrimSpace(gameString)
		if gameString == "" {
			continue
		}
		
		pgn, err := chess.PGN(strings.NewReader(gameString))
		if err != nil {
			cs.logger.Warnf("Failed to parse game: %v", err)
			continue
		}
		
		game := chess.NewGame(pgn)
		games = append(games, game)
	}
	
	return games, nil
}

func (cs *CalibrationService) extractRatingFromHeaders(game *chess.Game, header string) int {
	// For now, return a default rating since we can't access tag pairs directly
	// In a real implementation, you would parse the PGN headers before creating the game
	// or use a different chess library that provides access to headers
	return 1500 // Default rating
}

func (cs *CalibrationService) getRatingBucket(rating int) models.RatingBucket {
	switch {
	case rating >= 2001:
		return models.RatingBucket2001Plus
	case rating >= 1601:
		return models.RatingBucket1601to2000
	case rating >= 1201:
		return models.RatingBucket1201to1600
	default:
		return models.RatingBucket800to1200
	}
}

func (cs *CalibrationService) determineGamePhase(moveNumber int) string {
	switch {
	case moveNumber <= 15:
		return "opening"
	case moveNumber <= 40:
		return "middlegame"
	default:
		return "endgame"
	}
}

func (cs *CalibrationService) moveToUCI(move *chess.Move) string {
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

func (cs *CalibrationService) calculatePercentile(sortedData []float64, percentile int) float64 {
	if len(sortedData) == 0 {
		return 0
	}
	
	index := float64(percentile) / 100.0 * float64(len(sortedData)-1)
	lower := int(math.Floor(index))
	upper := int(math.Ceil(index))
	
	if lower == upper {
		return sortedData[lower]
	}
	
	weight := index - float64(lower)
	return sortedData[lower]*(1-weight) + sortedData[upper]*weight
}

func (cs *CalibrationService) saveCalibrationData() error {
	// Ensure data directory exists
	if err := os.MkdirAll(filepath.Dir("data/calibration.json"), 0755); err != nil {
		return err
	}
	
	data, err := json.MarshalIndent(cs.calibrationData, "", "  ")
	if err != nil {
		return err
	}
	
	return ioutil.WriteFile("data/calibration.json", data, 0644)
}

func (cs *CalibrationService) saveThresholds() error {
	// Ensure data directory exists
	if err := os.MkdirAll(filepath.Dir(cs.thresholdsPath), 0755); err != nil {
		return err
	}
	
	data, err := json.MarshalIndent(cs.calibrationData.Thresholds, "", "  ")
	if err != nil {
		return err
	}
	
	return ioutil.WriteFile(cs.thresholdsPath, data, 0644)
}

func (cs *CalibrationService) getDefaultThresholds() map[models.RatingBucket]models.EPThresholds {
	// Default thresholds based on typical chess performance
	return map[models.RatingBucket]models.EPThresholds{
		models.RatingBucket800to1200: {
			P1: 0.002, P5: 0.008, P10: 0.015, P25: 0.040, P50: 0.080, P75: 0.150, P90: 0.250,
		},
		models.RatingBucket1201to1600: {
			P1: 0.001, P5: 0.005, P10: 0.012, P25: 0.030, P50: 0.060, P75: 0.120, P90: 0.200,
		},
		models.RatingBucket1601to2000: {
			P1: 0.001, P5: 0.003, P10: 0.008, P25: 0.020, P50: 0.045, P75: 0.090, P90: 0.150,
		},
		models.RatingBucket2001Plus: {
			P1: 0.000, P5: 0.002, P10: 0.005, P25: 0.015, P50: 0.035, P75: 0.070, P90: 0.120,
		},
	}
} 