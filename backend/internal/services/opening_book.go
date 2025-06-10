package services

import (
	"strings"
)

// OpeningBookService provides opening move validation
type OpeningBookService struct {
	// Common opening moves in UCI format
	openingMoves map[string]bool
}

// NewOpeningBookService creates a new opening book service
func NewOpeningBookService() *OpeningBookService {
	return &OpeningBookService{
		openingMoves: initializeOpeningMoves(),
	}
}

// Contains checks if a move is in the opening book
func (obs *OpeningBookService) Contains(moveUCI string) bool {
	return obs.openingMoves[strings.ToLower(moveUCI)]
}

// IsOpeningPhase checks if we're still in the opening phase
func (obs *OpeningBookService) IsOpeningPhase(moveNumber int) bool {
	return moveNumber <= 15
}

// initializeOpeningMoves creates a map of common opening moves
func initializeOpeningMoves() map[string]bool {
	// Common opening moves from major openings
	moves := []string{
		// King's pawn openings
		"e2e4", "e7e5", "g1f3", "b8c6", "f1b5", "a7a6", "b5a4", "g8f6",
		"e1g1", "f8e7", "f1e1", "b7b5", "a4b3", "d7d6", "c2c3", "e8g8",
		"h2h3", "c8b7", "d2d3", "f6d7", "b1d2", "c6b8", "d2f1", "b8d7",
		
		// Queen's pawn openings
		"d2d4", "d7d5", "c2c4", "e7e6", "b1c3", "g8f6", "c4d5", "e6d5",
		"c1g5", "c7c6", "e2e3", "f8e7", "f1d3", "e8g8", "g1e2", "b8d7",
		"e1g1", "f8d6", "f2f4", "c8g4", "h2h3", "g4h5", "g2g4", "h5g6",
		
		// English opening
		"c2c4", "e7e5", "b1c3", "g8f6", "g1f3", "b8c6", "g2g3", "f8b4",
		"f1g2", "e8g8", "e1g1", "e5e4", "f3h4", "b4c3", "d2c3", "h7h6",
		
		// Reti opening
		"g1f3", "d7d5", "c2c4", "c7c6", "b2b3", "c8f5", "c1b2", "e7e6",
		"e2e3", "g8f6", "f1e2", "b8d7", "e1g1", "f8d6", "d2d3", "e8g8",
		
		// Sicilian Defense
		"c7c5", "g1f3", "d7d6", "d2d4", "c5d4", "f3d4", "g8f6", "b1c3",
		"g7g6", "c1e3", "f8g7", "f2f3", "e8g8", "d1d2", "b8c6", "e1c1",
		
		// French Defense
		"e7e6", "d2d4", "d7d5", "b1c3", "f8b4", "e2e5", "c7c5", "a2a3",
		"b4c3", "b2c3", "g8e7", "g1f3", "b8c6", "f1e2", "d8a5", "c1d2",
		
		// Caro-Kann Defense
		"c7c6", "d2d4", "d7d5", "b1c3", "d5c4", "e2e4", "c8f5", "g1f3",
		"e7e6", "f1c4", "f8b4", "e1g1", "g8f6", "d1e2", "e8g8", "f1d1",
		
		// Italian Game
		"f1c4", "f8c5", "c2c3", "f6e4", "d2d4", "e5d4", "c4f7", "e8f7",
		"d1b3", "e7f6", "c3d4", "c5b6", "f3e4", "d5e4", "b3f7", "f8e8",
		
		// Spanish Opening (Ruy Lopez)
		"b5c6", "d7c6", "d2d3", "f6d7", "c1e3", "f7f6", "f3d2", "a6a5",
		"f2f4", "e5f4", "e3f4", "d7c5", "d1f3", "c5e6", "f4e3", "c8d7",
		
		// Common pawn moves
		"a2a3", "a7a6", "a2a4", "a7a5", "b2b3", "b7b6", "b2b4", "b7b5",
		"c2c3", "c7c6", "f2f3", "f7f6", "g2g3", "g7g6", "h2h3", "h7h6",
		"g2g4", "g7g5", "h2h4", "h7h5",
	}
	
	openingMap := make(map[string]bool)
	for _, move := range moves {
		openingMap[strings.ToLower(move)] = true
	}
	
	return openingMap
} 