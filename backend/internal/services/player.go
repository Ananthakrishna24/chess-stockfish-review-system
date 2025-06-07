package services

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"chess-backend/internal/models"

	"github.com/sirupsen/logrus"
)

// PlayerService handles player statistics and historical data
type PlayerService struct {
	playerStats map[string]*models.PlayerProfile
	gameHistory map[string][]models.GameRecord
	mutex       sync.RWMutex
}

// NewPlayerService creates a new player service
func NewPlayerService() *PlayerService {
	return &PlayerService{
		playerStats: make(map[string]*models.PlayerProfile),
		gameHistory: make(map[string][]models.GameRecord),
	}
}

// GetPlayerStatistics retrieves comprehensive player statistics
func (s *PlayerService) GetPlayerStatistics(playerName string) (*models.PlayerStatisticsResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	playerName = s.normalizePlayerName(playerName)
	
	profile, exists := s.playerStats[playerName]
	if !exists {
		return nil, fmt.Errorf("player not found: %s", playerName)
	}
	
	games := s.gameHistory[playerName]
	if games == nil {
		games = make([]models.GameRecord, 0)
	}
	
	// Get recent games (last 10)
	recentGames := s.getRecentGames(games, 10)
	
	// Calculate strengths and weaknesses
	strengths, weaknesses := s.analyzePlayerTendencies(profile)
	
	// Generate improvement suggestions
	suggestions := s.generateImprovementSuggestions(profile, games)
	
	response := &models.PlayerStatisticsResponse{
		PlayerName:            playerName,
		GamesAnalyzed:         profile.GamesAnalyzed,
		AverageAccuracy:       profile.AverageAccuracy,
		RatingRange:           profile.RatingRange,
		RecentGames:          recentGames,
		Strengths:            strengths,
		Weaknesses:           weaknesses,
		ImprovementSuggestions: suggestions,
		PhasePerformance:     profile.PhasePerformance,
		OpeningRepertoire:    profile.OpeningRepertoire,
		TacticalStats:        profile.TacticalStats,
		LastUpdated:          profile.LastUpdated,
	}
	
	return response, nil
}

// RecordGameAnalysis updates player statistics with new game analysis
func (s *PlayerService) RecordGameAnalysis(analysis *models.GameAnalysisResponse) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	// Process both players
	s.processPlayerFromGame(analysis.GameInfo.White, "white", analysis)
	s.processPlayerFromGame(analysis.GameInfo.Black, "black", analysis)
}

// processPlayerFromGame updates a single player's statistics
func (s *PlayerService) processPlayerFromGame(playerName, color string, analysis *models.GameAnalysisResponse) {
	playerName = s.normalizePlayerName(playerName)
	
	// Get or create player profile
	profile, exists := s.playerStats[playerName]
	if !exists {
		profile = &models.PlayerProfile{
			PlayerName:        playerName,
			GamesAnalyzed:    0,
			AverageAccuracy:  0.0,
			RatingRange:      models.RatingRange{Min: 9999, Max: 0, Current: 0},
			PhasePerformance: models.PhasePerformance{},
			OpeningRepertoire: make(map[string]models.OpeningPerformance),
			TacticalStats:    models.TacticalStats{},
			CreatedAt:        time.Now(),
			LastUpdated:      time.Now(),
		}
		s.playerStats[playerName] = profile
	}
	
	// Get player stats from analysis
	var playerStats models.PlayerStatistics
	var playerRating int
	var result string
	
	if color == "white" {
		playerStats = analysis.Analysis.WhiteStats
		playerRating = analysis.GameInfo.WhiteRating
		switch analysis.GameInfo.Result {
		case "1-0":
			result = "win"
		case "0-1":
			result = "loss"
		default:
			result = "draw"
		}
	} else {
		playerStats = analysis.Analysis.BlackStats
		playerRating = analysis.GameInfo.BlackRating
		switch analysis.GameInfo.Result {
		case "0-1":
			result = "win"
		case "1-0":
			result = "loss"
		default:
			result = "draw"
		}
	}
	
	// Update basic statistics
	profile.GamesAnalyzed++
	profile.AverageAccuracy = s.updateRunningAverage(
		profile.AverageAccuracy, 
		playerStats.Accuracy, 
		profile.GamesAnalyzed,
	)
	
	// Update rating range
	if playerRating > 0 {
		if playerRating < profile.RatingRange.Min {
			profile.RatingRange.Min = playerRating
		}
		if playerRating > profile.RatingRange.Max {
			profile.RatingRange.Max = playerRating
		}
		profile.RatingRange.Current = playerRating
	}
	
	// Update phase performance
	s.updatePhasePerformance(profile, &analysis.Analysis.PhaseAnalysis)
	
	// Update opening repertoire
	s.updateOpeningRepertoire(profile, &analysis.Analysis.OpeningAnalysis, result)
	
	// Update tactical statistics
	s.updateTacticalStats(profile, &playerStats)
	
	// Add game record
	gameRecord := models.GameRecord{
		GameID:   analysis.GameID,
		Opponent: s.getOpponentName(playerName, analysis.GameInfo),
		Color:    color,
		Result:   result,
		Accuracy: playerStats.Accuracy,
		Date:     analysis.GameInfo.Date,
		Rating:   playerRating,
		Opening:  analysis.Analysis.OpeningAnalysis.Name,
		ECO:      analysis.Analysis.OpeningAnalysis.ECO,
	}
	
	if s.gameHistory[playerName] == nil {
		s.gameHistory[playerName] = make([]models.GameRecord, 0)
	}
	s.gameHistory[playerName] = append(s.gameHistory[playerName], gameRecord)
	
	// Sort games by date (most recent first)
	sort.Slice(s.gameHistory[playerName], func(i, j int) bool {
		return s.gameHistory[playerName][i].Date > s.gameHistory[playerName][j].Date
	})
	
	profile.LastUpdated = time.Now()
	
	logrus.Debugf("Updated statistics for player %s (games: %d, accuracy: %.1f%%)", 
		playerName, profile.GamesAnalyzed, profile.AverageAccuracy)
}

// Helper functions

func (s *PlayerService) normalizePlayerName(name string) string {
	return strings.TrimSpace(strings.ToLower(name))
}

func (s *PlayerService) updateRunningAverage(currentAvg, newValue float64, count int) float64 {
	if count <= 1 {
		return newValue
	}
	return ((currentAvg * float64(count-1)) + newValue) / float64(count)
}

func (s *PlayerService) updatePhasePerformance(profile *models.PlayerProfile, phaseAnalysis *models.PhaseAnalysis) {
	profile.PhasePerformance.OpeningAccuracy = s.updateRunningAverage(
		profile.PhasePerformance.OpeningAccuracy,
		phaseAnalysis.OpeningAccuracy,
		profile.GamesAnalyzed,
	)
	profile.PhasePerformance.MiddlegameAccuracy = s.updateRunningAverage(
		profile.PhasePerformance.MiddlegameAccuracy,
		phaseAnalysis.MiddlegameAccuracy,
		profile.GamesAnalyzed,
	)
	profile.PhasePerformance.EndgameAccuracy = s.updateRunningAverage(
		profile.PhasePerformance.EndgameAccuracy,
		phaseAnalysis.EndgameAccuracy,
		profile.GamesAnalyzed,
	)
}

func (s *PlayerService) updateOpeningRepertoire(profile *models.PlayerProfile, openingAnalysis *models.OpeningAnalysis, result string) {
	opening := openingAnalysis.ECO
	if opening == "" {
		opening = "Unknown"
	}
	
	perf, exists := profile.OpeningRepertoire[opening]
	if !exists {
		perf = models.OpeningPerformance{
			ECO:      opening,
			Name:     openingAnalysis.Name,
			Games:    0,
			Accuracy: 0.0,
			Results:  models.OpeningResults{},
		}
	}
	
	perf.Games++
	perf.Accuracy = s.updateRunningAverage(perf.Accuracy, openingAnalysis.Accuracy, perf.Games)
	
	switch result {
	case "win":
		perf.Results.Wins++
	case "loss":
		perf.Results.Losses++
	case "draw":
		perf.Results.Draws++
	}
	
	profile.OpeningRepertoire[opening] = perf
}

func (s *PlayerService) updateTacticalStats(profile *models.PlayerProfile, playerStats *models.PlayerStatistics) {
	if playerStats.TacticalMoves > 0 {
		profile.TacticalStats.TotalTacticalMoves += playerStats.TacticalMoves
	}
	if playerStats.ForcingMoves > 0 {
		profile.TacticalStats.TotalForcingMoves += playerStats.ForcingMoves
	}
	if playerStats.CriticalMoments > 0 {
		profile.TacticalStats.TotalCriticalMoments += playerStats.CriticalMoments
	}
	
	// Update move quality distribution
	profile.TacticalStats.BrilliantMoves += playerStats.MoveCounts.Brilliant
	profile.TacticalStats.BlunderRate = s.updateRunningAverage(
		profile.TacticalStats.BlunderRate,
		float64(playerStats.MoveCounts.Blunder),
		profile.GamesAnalyzed,
	)
}

func (s *PlayerService) getOpponentName(playerName string, gameInfo models.GameInfo) string {
	normalizedPlayer := s.normalizePlayerName(playerName)
	normalizedWhite := s.normalizePlayerName(gameInfo.White)
	
	if normalizedPlayer == normalizedWhite {
		return gameInfo.Black
	}
	return gameInfo.White
}

func (s *PlayerService) getRecentGames(games []models.GameRecord, limit int) []models.RecentGame {
	recentGames := make([]models.RecentGame, 0)
	
	count := len(games)
	if count > limit {
		count = limit
	}
	
	for i := 0; i < count; i++ {
		game := games[i]
		recentGames = append(recentGames, models.RecentGame{
			GameID:   game.GameID,
			Opponent: game.Opponent,
			Result:   game.Result,
			Accuracy: game.Accuracy,
			Date:     game.Date,
			Opening:  game.Opening,
			ECO:      game.ECO,
		})
	}
	
	return recentGames
}

func (s *PlayerService) analyzePlayerTendencies(profile *models.PlayerProfile) ([]string, []string) {
	strengths := make([]string, 0)
	weaknesses := make([]string, 0)
	
	// Analyze phase performance
	if profile.PhasePerformance.OpeningAccuracy > 85 {
		strengths = append(strengths, "opening preparation")
	} else if profile.PhasePerformance.OpeningAccuracy < 70 {
		weaknesses = append(weaknesses, "opening knowledge")
	}
	
	if profile.PhasePerformance.MiddlegameAccuracy > 85 {
		strengths = append(strengths, "middlegame tactics")
	} else if profile.PhasePerformance.MiddlegameAccuracy < 70 {
		weaknesses = append(weaknesses, "middlegame planning")
	}
	
	if profile.PhasePerformance.EndgameAccuracy > 85 {
		strengths = append(strengths, "endgame technique")
	} else if profile.PhasePerformance.EndgameAccuracy < 70 {
		weaknesses = append(weaknesses, "endgame fundamentals")
	}
	
	// Analyze tactical stats
	avgTacticalPerGame := float64(profile.TacticalStats.TotalTacticalMoves) / float64(profile.GamesAnalyzed)
	if avgTacticalPerGame > 3 {
		strengths = append(strengths, "tactical awareness")
	} else if avgTacticalPerGame < 1 {
		weaknesses = append(weaknesses, "tactical vision")
	}
	
	if profile.TacticalStats.BlunderRate > 2 {
		weaknesses = append(weaknesses, "time management")
	}
	
	if profile.TacticalStats.BrilliantMoves > profile.GamesAnalyzed/10 {
		strengths = append(strengths, "creative play")
	}
	
	return strengths, weaknesses
}

func (s *PlayerService) generateImprovementSuggestions(profile *models.PlayerProfile, games []models.GameRecord) []string {
	suggestions := make([]string, 0)
	
	// Opening suggestions
	if profile.PhasePerformance.OpeningAccuracy < 75 {
		suggestions = append(suggestions, "Study opening principles and common opening traps")
		suggestions = append(suggestions, "Analyze your most played openings more deeply")
	}
	
	// Middlegame suggestions
	if profile.PhasePerformance.MiddlegameAccuracy < 75 {
		suggestions = append(suggestions, "Practice tactical puzzles daily")
		suggestions = append(suggestions, "Study pawn structures and strategic concepts")
	}
	
	// Endgame suggestions
	if profile.PhasePerformance.EndgameAccuracy < 75 {
		suggestions = append(suggestions, "Learn basic endgame patterns (K+Q vs K, K+R vs K)")
		suggestions = append(suggestions, "Practice endgame technique with fewer pieces")
	}
	
	// Blunder prevention
	if profile.TacticalStats.BlunderRate > 1.5 {
		suggestions = append(suggestions, "Take more time for critical moves")
		suggestions = append(suggestions, "Always check for tactical threats before moving")
	}
	
	// Opening repertoire diversity
	if len(profile.OpeningRepertoire) < 3 {
		suggestions = append(suggestions, "Expand your opening repertoire with both white and black")
	}
	
	return suggestions
}

// GetAllPlayers returns a list of all tracked players
func (s *PlayerService) GetAllPlayers() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	players := make([]string, 0, len(s.playerStats))
	for playerName := range s.playerStats {
		players = append(players, playerName)
	}
	
	sort.Strings(players)
	return players
}

// GetTopPlayers returns players ranked by average accuracy
func (s *PlayerService) GetTopPlayers(limit int) []models.PlayerRanking {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	rankings := make([]models.PlayerRanking, 0, len(s.playerStats))
	
	for _, profile := range s.playerStats {
		rankings = append(rankings, models.PlayerRanking{
			PlayerName:      profile.PlayerName,
			GamesAnalyzed:   profile.GamesAnalyzed,
			AverageAccuracy: profile.AverageAccuracy,
			CurrentRating:   profile.RatingRange.Current,
		})
	}
	
	// Sort by accuracy (descending)
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].AverageAccuracy > rankings[j].AverageAccuracy
	})
	
	if limit > 0 && len(rankings) > limit {
		rankings = rankings[:limit]
	}
	
	return rankings
} 