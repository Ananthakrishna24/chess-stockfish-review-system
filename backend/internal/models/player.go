package models

import (
	"time"
)

// Player Statistics Models

// PlayerStatisticsResponse - Complete player statistics response
type PlayerStatisticsResponse struct {
	PlayerName            string                       `json:"playerName"`
	GamesAnalyzed         int                          `json:"gamesAnalyzed"`
	AverageAccuracy       float64                      `json:"averageAccuracy"`
	RatingRange           RatingRange                  `json:"ratingRange"`
	RecentGames          []RecentGame                  `json:"recentGames"`
	Strengths            []string                      `json:"strengths"`
	Weaknesses           []string                      `json:"weaknesses"`
	ImprovementSuggestions []string                    `json:"improvementSuggestions"`
	PhasePerformance     PhasePerformance             `json:"phasePerformance"`
	OpeningRepertoire    map[string]OpeningPerformance `json:"openingRepertoire"`
	TacticalStats        TacticalStats                `json:"tacticalStats"`
	LastUpdated          time.Time                    `json:"lastUpdated"`
}

// PlayerProfile - Internal player profile for statistics tracking
type PlayerProfile struct {
	PlayerName        string                       `json:"playerName"`
	GamesAnalyzed     int                          `json:"gamesAnalyzed"`
	AverageAccuracy   float64                      `json:"averageAccuracy"`
	RatingRange       RatingRange                  `json:"ratingRange"`
	PhasePerformance  PhasePerformance             `json:"phasePerformance"`
	OpeningRepertoire map[string]OpeningPerformance `json:"openingRepertoire"`
	TacticalStats     TacticalStats                `json:"tacticalStats"`
	CreatedAt         time.Time                    `json:"createdAt"`
	LastUpdated       time.Time                    `json:"lastUpdated"`
}

// RatingRange - Player rating information
type RatingRange struct {
	Min     int `json:"min"`
	Max     int `json:"max"`
	Current int `json:"current"`
}

// RecentGame - Recent game information
type RecentGame struct {
	GameID   string  `json:"gameId"`
	Opponent string  `json:"opponent"`
	Result   string  `json:"result"`
	Accuracy float64 `json:"accuracy"`
	Date     string  `json:"date"`
	Opening  string  `json:"opening,omitempty"`
	ECO      string  `json:"eco,omitempty"`
}

// GameRecord - Internal game record for tracking
type GameRecord struct {
	GameID   string  `json:"gameId"`
	Opponent string  `json:"opponent"`
	Color    string  `json:"color"`
	Result   string  `json:"result"`
	Accuracy float64 `json:"accuracy"`
	Date     string  `json:"date"`
	Rating   int     `json:"rating"`
	Opening  string  `json:"opening"`
	ECO      string  `json:"eco"`
}

// PhasePerformance - Performance statistics by game phase
type PhasePerformance struct {
	OpeningAccuracy    float64 `json:"openingAccuracy"`
	MiddlegameAccuracy float64 `json:"middlegameAccuracy"`
	EndgameAccuracy    float64 `json:"endgameAccuracy"`
}

// OpeningPerformance - Performance with specific opening
type OpeningPerformance struct {
	ECO      string         `json:"eco"`
	Name     string         `json:"name"`
	Games    int            `json:"games"`
	Accuracy float64        `json:"accuracy"`
	Results  OpeningResults `json:"results"`
}

// OpeningResults - Win/draw/loss results for opening
type OpeningResults struct {
	Wins   int `json:"wins"`
	Draws  int `json:"draws"`
	Losses int `json:"losses"`
}

// TacticalStats - Tactical performance statistics
type TacticalStats struct {
	TotalTacticalMoves   int     `json:"totalTacticalMoves"`
	TotalForcingMoves    int     `json:"totalForcingMoves"`
	TotalCriticalMoments int     `json:"totalCriticalMoments"`
	BrilliantMoves       int     `json:"brilliantMoves"`
	BlunderRate          float64 `json:"blunderRate"`
}

// PlayerRanking - Player ranking information
type PlayerRanking struct {
	PlayerName      string  `json:"playerName"`
	GamesAnalyzed   int     `json:"gamesAnalyzed"`
	AverageAccuracy float64 `json:"averageAccuracy"`
	CurrentRating   int     `json:"currentRating"`
}

// Opening Database Models

// OpeningInfo - Complete opening information
type OpeningInfo struct {
	ECO        string             `json:"eco"`
	Name       string             `json:"name"`
	Variation  string             `json:"variation"`
	Moves      []string           `json:"moves"`
	Popularity float64            `json:"popularity"`
	Statistics OpeningStatistics  `json:"statistics"`
	Theory     string             `json:"theory"`
	KeyIdeas   []string           `json:"keyIdeas"`
}

// OpeningStatistics - Statistical performance of opening
type OpeningStatistics struct {
	White float64 `json:"white"`
	Draw  float64 `json:"draw"`
	Black float64 `json:"black"`
}

// OpeningSearchRequest - Request for opening search
type OpeningSearchRequest struct {
	ECO   string   `json:"eco,omitempty"`
	FEN   string   `json:"fen,omitempty"`
	Moves []string `json:"moves,omitempty"`
	Name  string   `json:"name,omitempty"`
}

// OpeningSearchResponse - Response for opening search
type OpeningSearchResponse struct {
	Results []OpeningInfo `json:"results"`
	Count   int           `json:"count"`
} 