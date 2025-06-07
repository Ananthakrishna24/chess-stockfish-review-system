package models

import (
	"time"
)

// AnalyzeGameRequest - Request structure for game analysis
type AnalyzeGameRequest struct {
	PGN     string           `json:"pgn" binding:"required"`
	Options AnalysisOptions  `json:"options"`
}

// AnalysisOptions - Configuration options for analysis
type AnalysisOptions struct {
	Depth                     int           `json:"depth,omitempty"`
	TimePerMove               int           `json:"timePerMove,omitempty"`
	IncludeBookMoves          bool          `json:"includeBookMoves,omitempty"`
	IncludeTacticalAnalysis   bool          `json:"includeTacticalAnalysis,omitempty"`
	PlayerRatings             PlayerRatings `json:"playerRatings,omitempty"`
}

// PlayerRatings - Player rating information for accuracy calculation
type PlayerRatings struct {
	White int `json:"white,omitempty"`
	Black int `json:"black,omitempty"`
}

// AnalyzePositionRequest - Request structure for position analysis
type AnalyzePositionRequest struct {
	FEN       string `json:"fen" binding:"required"`
	Depth     int    `json:"depth,omitempty"`
	MultiPV   int    `json:"multiPv,omitempty"`
	TimeLimit int    `json:"timeLimit,omitempty"`
}

// GameAnalysisResponse - Complete game analysis response
type GameAnalysisResponse struct {
	GameID         string          `json:"gameId"`
	GameInfo       GameInfo        `json:"gameInfo"`
	Analysis       GameAnalysis    `json:"analysis"`
	ProcessingTime float64         `json:"processingTime"`
	Timestamp      time.Time       `json:"timestamp"`
}

// GameInfo - Basic game metadata
type GameInfo struct {
	White       string `json:"white"`
	Black       string `json:"black"`
	WhiteRating int    `json:"whiteRating,omitempty"`
	BlackRating int    `json:"blackRating,omitempty"`
	Result      string `json:"result"`
	Date        string `json:"date,omitempty"`
	Event       string `json:"event,omitempty"`
	Site        string `json:"site,omitempty"`
	Opening     string `json:"opening,omitempty"`
	ECO         string `json:"eco,omitempty"`
}

// GameAnalysis - Complete analysis of a chess game
type GameAnalysis struct {
	Moves              []MoveAnalysis     `json:"moves"`
	WhiteStats         PlayerStatistics   `json:"whiteStats"`
	BlackStats         PlayerStatistics   `json:"blackStats"`
	OpeningAnalysis    OpeningAnalysis    `json:"openingAnalysis"`
	GamePhases         GamePhases         `json:"gamePhases"`
	PhaseAnalysis      PhaseAnalysis      `json:"phaseAnalysis"`
	CriticalMoments    []CriticalMoment   `json:"criticalMoments"`
	EvaluationHistory  []EngineEvaluation `json:"evaluationHistory"`
}

// MoveAnalysis - Analysis of a single move
type MoveAnalysis struct {
	MoveNumber        int                  `json:"moveNumber"`
	Move              string               `json:"move"`
	SAN               string               `json:"san"`
	FEN               string               `json:"fen"`
	Evaluation        EngineEvaluation     `json:"evaluation"`
	Classification    string               `json:"classification"`
	AlternativeMoves  []AlternativeMove    `json:"alternativeMoves,omitempty"`
	TacticalAnalysis  *TacticalAnalysis    `json:"tacticalAnalysis,omitempty"`
	Comment           string               `json:"comment,omitempty"`
}

// EngineEvaluation - Stockfish evaluation of a position
type EngineEvaluation struct {
	Score              int      `json:"score"`
	Depth              int      `json:"depth"`
	BestMove           string   `json:"bestMove"`
	PrincipalVariation []string `json:"principalVariation"`
	Nodes              int64    `json:"nodes"`
	Time               int      `json:"time"`
	Mate               *int     `json:"mate,omitempty"`
}

// AlternativeMove - Alternative move with evaluation
type AlternativeMove struct {
	Move       string           `json:"move"`
	SAN        string           `json:"san"`
	Evaluation EngineEvaluation `json:"evaluation"`
}

// TacticalAnalysis - Tactical pattern analysis
type TacticalAnalysis struct {
	Patterns     []string `json:"patterns"`
	IsForcing    bool     `json:"isForcing"`
	IsTactical   bool     `json:"isTactical"`
	ThreatLevel  string   `json:"threatLevel"`
	Description  string   `json:"description"`
}

// PlayerStatistics - Player performance statistics
type PlayerStatistics struct {
	Accuracy        float64    `json:"accuracy"`
	MoveCounts      MoveCounts `json:"moveCounts"`
	TacticalMoves   int        `json:"tacticalMoves,omitempty"`
	ForcingMoves    int        `json:"forcingMoves,omitempty"`
	CriticalMoments int        `json:"criticalMoments,omitempty"`
}

// MoveCounts - Count of different move types
type MoveCounts struct {
	Brilliant   int `json:"brilliant"`
	Great       int `json:"great"`
	Best        int `json:"best"`
	Excellent   int `json:"excellent"`
	Good        int `json:"good"`
	Book        int `json:"book"`
	Inaccuracy  int `json:"inaccuracy"`
	Mistake     int `json:"mistake"`
	Blunder     int `json:"blunder"`
	Miss        int `json:"miss"`
}

// OpeningAnalysis - Opening phase analysis
type OpeningAnalysis struct {
	Name          string  `json:"name"`
	ECO           string  `json:"eco"`
	Accuracy      float64 `json:"accuracy"`
	Theory        string  `json:"theory"`
	DeviationMove int     `json:"deviationMove"`
}

// GamePhases - Move numbers where game phases end/begin
type GamePhases struct {
	Opening    int `json:"opening"`
	Middlegame int `json:"middlegame"`
	Endgame    int `json:"endgame"`
}

// PhaseAnalysis - Accuracy statistics by game phase
type PhaseAnalysis struct {
	OpeningAccuracy    float64 `json:"openingAccuracy"`
	MiddlegameAccuracy float64 `json:"middlegameAccuracy"`
	EndgameAccuracy    float64 `json:"endgameAccuracy"`
}

// CriticalMoment - Significant evaluation swing
type CriticalMoment struct {
	MoveNumber  int     `json:"moveNumber"`
	BeforeEval  int     `json:"beforeEval"`
	AfterEval   int     `json:"afterEval"`
	Advantage   string  `json:"advantage"`
	Description string  `json:"description"`
}

// PositionAnalysisResponse - Response for position analysis
type PositionAnalysisResponse struct {
	FEN              string             `json:"fen"`
	Evaluation       EngineEvaluation   `json:"evaluation"`
	AlternativeMoves []AlternativeMove  `json:"alternativeMoves"`
	PositionInfo     PositionInfo       `json:"positionInfo"`
}

// PositionInfo - Additional position information
type PositionInfo struct {
	Phase    string       `json:"phase"`
	Material MaterialInfo `json:"material"`
	Safety   SafetyInfo   `json:"safety"`
}

// MaterialInfo - Material balance information
type MaterialInfo struct {
	White int `json:"white"`
	Black int `json:"black"`
}

// SafetyInfo - King safety assessment
type SafetyInfo struct {
	WhiteKing string `json:"whiteKing"`
	BlackKing string `json:"blackKing"`
}

// ProgressResponse - Analysis progress information
type ProgressResponse struct {
	GameID   string          `json:"gameId"`
	Status   string          `json:"status"`
	Progress ProgressDetails `json:"progress"`
	Error    string          `json:"error,omitempty"`
}

// ProgressDetails - Detailed progress information
type ProgressDetails struct {
	CurrentMove              int     `json:"currentMove"`
	TotalMoves               int     `json:"totalMoves"`
	Percentage               float64 `json:"percentage"`
	EstimatedTimeRemaining   int     `json:"estimatedTimeRemaining"`
}

// EngineConfigResponse - Engine configuration response
type EngineConfigResponse struct {
	Version       string        `json:"version"`
	Features      []string      `json:"features"`
	Limits        EngineLimits  `json:"limits"`
	CurrentConfig EngineOptions `json:"currentConfig"`
}

// EngineLimits - Engine capability limits
type EngineLimits struct {
	MaxDepth int `json:"maxDepth"`
	MaxTime  int `json:"maxTime"`
	MaxNodes int `json:"maxNodes"`
}

// EngineOptions - Engine configuration options
type EngineOptions struct {
	Threads          int    `json:"threads"`
	Hash             int    `json:"hash"`
	Contempt         int    `json:"contempt"`
	AnalysisContempt string `json:"analysisContempt"`
}

// UpdateEngineConfigRequest - Engine configuration update request
type UpdateEngineConfigRequest struct {
	Threads          *int    `json:"threads,omitempty"`
	Hash             *int    `json:"hash,omitempty"`
	Contempt         *int    `json:"contempt,omitempty"`
	AnalysisContempt *string `json:"analysisContempt,omitempty"`
} 