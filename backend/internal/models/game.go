package models

import (
	"sync"
	"time"
)

// AnalysisJob - Represents an analysis job with progress tracking
type AnalysisJob struct {
	ID              string                `json:"id"`
	PGN             string                `json:"pgn"`
	Options         AnalysisOptions       `json:"options"`
	Status          AnalysisStatus        `json:"status"`
	Progress        ProgressDetails       `json:"progress"`
	Result          *GameAnalysisResponse `json:"result,omitempty"`
	Error           string                `json:"error,omitempty"`
	CreatedAt       time.Time             `json:"createdAt"`
	CompletedAt     *time.Time            `json:"completedAt,omitempty"`
	ProcessingTime  float64               `json:"processingTime"`
	CurrentMove     int                   `json:"currentMove"`
	TotalMoves      int                   `json:"totalMoves"`
	mutex           sync.RWMutex          `json:"-"`
}

// AnalysisStatus - Status of analysis job
type AnalysisStatus string

const (
	StatusQueued     AnalysisStatus = "queued"
	StatusAnalyzing  AnalysisStatus = "analyzing"
	StatusCompleted  AnalysisStatus = "completed"
	StatusFailed     AnalysisStatus = "failed"
	StatusCancelled  AnalysisStatus = "cancelled"
)

// UpdateProgress - Thread-safe progress update
func (j *AnalysisJob) UpdateProgress(currentMove, totalMoves int) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	
	j.CurrentMove = currentMove
	j.TotalMoves = totalMoves
	j.Progress.CurrentMove = currentMove
	j.Progress.TotalMoves = totalMoves
	
	if totalMoves > 0 {
		j.Progress.Percentage = float64(currentMove) / float64(totalMoves) * 100
	}
	
	// Estimate remaining time based on current progress
	if currentMove > 0 && j.ProcessingTime > 0 {
		avgTimePerMove := j.ProcessingTime / float64(currentMove)
		remainingMoves := totalMoves - currentMove
		j.Progress.EstimatedTimeRemaining = int(avgTimePerMove * float64(remainingMoves))
	}
}

// SetStatus - Thread-safe status update
func (j *AnalysisJob) SetStatus(status AnalysisStatus) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.Status = status
	
	if status == StatusCompleted || status == StatusFailed {
		now := time.Now()
		j.CompletedAt = &now
		j.ProcessingTime = now.Sub(j.CreatedAt).Seconds()
	}
}

// SetError - Thread-safe error setting
func (j *AnalysisJob) SetError(err string) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.Error = err
	j.Status = StatusFailed
	now := time.Now()
	j.CompletedAt = &now
}

// SetResult - Thread-safe result setting
func (j *AnalysisJob) SetResult(result *GameAnalysisResponse) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.Result = result
	j.Status = StatusCompleted
	now := time.Now()
	j.CompletedAt = &now
	j.ProcessingTime = now.Sub(j.CreatedAt).Seconds()
}

// GetProgress - Thread-safe progress retrieval
func (j *AnalysisJob) GetProgress() ProgressResponse {
	j.mutex.RLock()
	defer j.mutex.RUnlock()
	
	return ProgressResponse{
		GameID:   j.ID,
		Status:   string(j.Status),
		Progress: j.Progress,
		Error:    j.Error,
	}
}

// ParsedGame - Represents a parsed chess game
type ParsedGame struct {
	Game         interface{}        `json:"game"` // chess.Game object
	GameInfo     GameInfo           `json:"gameInfo"`
	Moves        []ParsedMove       `json:"moves"`
	TotalMoves   int                `json:"totalMoves"`
	StartingFEN  string             `json:"startingFen,omitempty"`
}

// ParsedMove - Represents a single parsed move
type ParsedMove struct {
	MoveNumber int    `json:"moveNumber"`
	Move       string `json:"move"`
	SAN        string `json:"san"`
	UCI        string `json:"uci"`
	FEN        string `json:"fen"`
	IsWhite    bool   `json:"isWhite"`
}

// Note: MoveClassification is now defined in analysis.go

// GamePhase - Enumeration of game phases
type GamePhase string

const (
	Opening    GamePhase = "opening"
	Middlegame GamePhase = "middlegame"
	Endgame    GamePhase = "endgame"
)

// String - Convert phase to string
func (gp GamePhase) String() string {
	return string(gp)
}

// AnalysisCache - Cache entry for analysis results
type AnalysisCache struct {
	GameID     string                `json:"gameId"`
	Result     GameAnalysisResponse  `json:"result"`
	CreatedAt  time.Time             `json:"createdAt"`
	ExpiresAt  time.Time             `json:"expiresAt"`
	AccessedAt time.Time             `json:"accessedAt"`
	AccessCount int                  `json:"accessCount"`
}

// IsExpired - Check if cache entry is expired
func (ac *AnalysisCache) IsExpired() bool {
	return time.Now().After(ac.ExpiresAt)
}

// UpdateAccess - Update access tracking
func (ac *AnalysisCache) UpdateAccess() {
	ac.AccessedAt = time.Now()
	ac.AccessCount++
}

// PositionCache - Cache entry for position analysis
type PositionCache struct {
	FEN        string                   `json:"fen"`
	Result     PositionAnalysisResponse `json:"result"`
	CreatedAt  time.Time                `json:"createdAt"`
	ExpiresAt  time.Time                `json:"expiresAt"`
}

// IsExpired - Check if position cache entry is expired
func (pc *PositionCache) IsExpired() bool {
	return time.Now().After(pc.ExpiresAt)
} 