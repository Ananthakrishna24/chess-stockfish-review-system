package services

import (
	"fmt"
	"strings"

	"chess-backend/internal/models"

	"github.com/notnil/chess"
	"github.com/sirupsen/logrus"
)

// ChessService handles chess game parsing and validation
type ChessService struct{}

// NewChessService creates a new chess service
func NewChessService() *ChessService {
	return &ChessService{}
}

// ParsePGN parses a PGN string and returns game information
func (s *ChessService) ParsePGN(pgnStr string) (*models.ParsedGame, error) {
	// Parse the PGN
	pgnFunc, err := chess.PGN(strings.NewReader(pgnStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse PGN: %v", err)
	}
	
	game := chess.NewGame(pgnFunc)
	
	// Extract game information (basic implementation)
	gameInfo := s.extractGameInfoFromGame(game)
	
	// Parse moves
	moves, err := s.extractMoves(game)
	if err != nil {
		return nil, fmt.Errorf("failed to extract moves: %v", err)
	}
	
	parsedGame := &models.ParsedGame{
		Game:       game,
		GameInfo:   gameInfo,
		Moves:      moves,
		TotalMoves: len(moves),
	}
	
	// Set starting FEN if different from standard
	if game.Position().String() != chess.StartingPosition().String() {
		parsedGame.StartingFEN = game.Position().String()
	}
	
	return parsedGame, nil
}

// extractGameInfoFromGame extracts basic game metadata
func (s *ChessService) extractGameInfoFromGame(game *chess.Game) models.GameInfo {
	gameInfo := models.GameInfo{
		White:  "Unknown",
		Black:  "Unknown", 
		Result: "*",
	}
	
	// Get game outcome
	if game.Outcome() != chess.NoOutcome {
		switch game.Outcome() {
		case chess.WhiteWon:
			gameInfo.Result = "1-0"
		case chess.BlackWon:
			gameInfo.Result = "0-1"
		case chess.Draw:
			gameInfo.Result = "1/2-1/2"
		}
	}
	
	// For now, we'll use basic defaults
	// In a full implementation, you'd parse the PGN headers properly
	gameInfo.Date = "Unknown"
	gameInfo.Event = "Unknown"
	gameInfo.Site = "Unknown"
	
	return gameInfo
}

// extractMoves extracts all moves from the game
func (s *ChessService) extractMoves(game *chess.Game) ([]models.ParsedMove, error) {
	var moves []models.ParsedMove
	
	// Create a new game to replay moves
	tempGame := chess.NewGame()
	
	for i, move := range game.Moves() {
		// Apply the move to get the position
		if err := tempGame.Move(move); err != nil {
			logrus.Errorf("Failed to apply move %d: %v", i, err)
			continue
		}
		
		// Determine if this is a white move
		isWhite := (i % 2) == 0
		
		// Calculate move number (1-based for display)
		moveNumber := (i / 2) + 1
		if !isWhite {
			moveNumber = (i + 1) / 2
		}
		
		parsedMove := models.ParsedMove{
			MoveNumber: moveNumber,
			Move:       move.String(),
			SAN:        s.moveToSAN(move),
			UCI:        s.moveToUCI(move),
			FEN:        tempGame.Position().String(),
			IsWhite:    isWhite,
		}
		
		moves = append(moves, parsedMove)
	}
	
	return moves, nil
}

// moveToSAN converts a move to Standard Algebraic Notation
func (s *ChessService) moveToSAN(move *chess.Move) string {
	return move.String()
}

// moveToUCI converts a move to UCI notation
func (s *ChessService) moveToUCI(move *chess.Move) string {
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

// parseRating converts a string rating to integer
func (s *ChessService) parseRating(ratingStr string) int {
	var rating int
	if _, err := fmt.Sscanf(ratingStr, "%d", &rating); err != nil {
		return 0
	}
	return rating
}

// ValidateFEN validates a FEN string
func (s *ChessService) ValidateFEN(fen string) error {
	_, err := chess.FEN(fen)
	if err != nil {
		return fmt.Errorf("invalid FEN: %v", err)
	}
	return nil
}

// GetPositionFromFEN creates a position from FEN
func (s *ChessService) GetPositionFromFEN(fen string) (*chess.Position, error) {
	fenFunc, err := chess.FEN(fen)
	if err != nil {
		return nil, fmt.Errorf("invalid FEN: %v", err)
	}
	
	game := chess.NewGame(fenFunc)
	return game.Position(), nil
}

// GetGamePhase determines the game phase based on material and moves
func (s *ChessService) GetGamePhase(position *chess.Position, moveCount int) models.GamePhase {
	// Count pieces for both sides
	whitePieces := s.countPieces(position, chess.White)
	blackPieces := s.countPieces(position, chess.Black)
	totalPieces := whitePieces + blackPieces
	
	// Opening phase: first 15 moves or if many pieces remain
	if moveCount <= 15 || totalPieces >= 28 {
		return models.Opening
	}
	
	// Endgame phase: few pieces remaining
	if totalPieces <= 12 {
		return models.Endgame
	}
	
	// Check for queen exchange (common endgame indicator)
	whiteQueens := s.countPieceType(position, chess.White, chess.Queen)
	blackQueens := s.countPieceType(position, chess.Black, chess.Queen)
	
	if whiteQueens == 0 && blackQueens == 0 && totalPieces <= 20 {
		return models.Endgame
	}
	
	// Otherwise, middlegame
	return models.Middlegame
}

// countPieces counts total pieces for a color
func (s *ChessService) countPieces(position *chess.Position, color chess.Color) int {
	count := 0
	board := position.Board()
	
	for square := chess.A1; square <= chess.H8; square++ {
		piece := board.Piece(square)
		if piece.Color() == color && piece.Type() != chess.NoPieceType {
			count++
		}
	}
	
	return count
}

// countPieceType counts specific piece types for a color
func (s *ChessService) countPieceType(position *chess.Position, color chess.Color, pieceType chess.PieceType) int {
	count := 0
	board := position.Board()
	
	for square := chess.A1; square <= chess.H8; square++ {
		piece := board.Piece(square)
		if piece.Color() == color && piece.Type() == pieceType {
			count++
		}
	}
	
	return count
}

// CalculateMaterialValue calculates material value for a position
func (s *ChessService) CalculateMaterialValue(position *chess.Position, color chess.Color) int {
	value := 0
	board := position.Board()
	
	// Standard piece values
	pieceValues := map[chess.PieceType]int{
		chess.Pawn:   1,
		chess.Knight: 3,
		chess.Bishop: 3,
		chess.Rook:   5,
		chess.Queen:  9,
		chess.King:   0, // King doesn't have material value
	}
	
	for square := chess.A1; square <= chess.H8; square++ {
		piece := board.Piece(square)
		if piece.Color() == color {
			if val, exists := pieceValues[piece.Type()]; exists {
				value += val
			}
		}
	}
	
	return value
}

// AssessKingSafety provides a basic king safety assessment
func (s *ChessService) AssessKingSafety(position *chess.Position, color chess.Color) string {
	// This is a simplified king safety assessment
	// In a full implementation, you would check for:
	// - Pawn shelter
	// - Open files near the king
	// - Enemy pieces attacking king zone
	// - King position relative to center
	
	board := position.Board()
	
	// Find the king
	var kingSquare chess.Square
	for square := chess.A1; square <= chess.H8; square++ {
		piece := board.Piece(square)
		if piece.Color() == color && piece.Type() == chess.King {
			kingSquare = square
			break
		}
	}
	
	// Simple assessment based on king file and rank
	file := kingSquare.File()
	rank := kingSquare.Rank()
	
	// Check if king is on back rank (safer for most of game)
	expectedBackRank := chess.Rank1
	if color == chess.Black {
		expectedBackRank = chess.Rank8
	}
	
	if rank == expectedBackRank {
		// King on back rank - check if castled
		if file == chess.FileG || file == chess.FileC {
			return "safe" // Likely castled
		} else if file == chess.FileE {
			return "exposed" // King still in center
		}
		return "safe"
	}
	
	// King advanced - potentially dangerous
	return "danger"
} 