package services

import (
	"fmt"
	"strings"

	"chess-backend/internal/models"

	"github.com/notnil/chess"
	"github.com/sirupsen/logrus"
)

// OpeningService handles opening database operations
type OpeningService struct {
	openings map[string]models.OpeningInfo
	ecoTree  map[string][]models.OpeningInfo
}

// NewOpeningService creates a new opening service with ECO database
func NewOpeningService() *OpeningService {
	service := &OpeningService{
		openings: make(map[string]models.OpeningInfo),
		ecoTree:  make(map[string][]models.OpeningInfo),
	}
	
	service.loadOpeningDatabase()
	return service
}

// SearchByECO finds opening information by ECO code
func (s *OpeningService) SearchByECO(ecoCode string) (*models.OpeningInfo, error) {
	ecoCode = strings.ToUpper(ecoCode)
	
	if opening, exists := s.openings[ecoCode]; exists {
		return &opening, nil
	}
	
	return nil, fmt.Errorf("opening with ECO code %s not found", ecoCode)
}

// SearchByPosition finds opening information by FEN position
func (s *OpeningService) SearchByPosition(fen string) (*models.OpeningInfo, error) {
	// Parse the position
	fenFunc, err := chess.FEN(fen)
	if err != nil {
		return nil, fmt.Errorf("invalid FEN: %v", err)
	}
	
	game := chess.NewGame(fenFunc)
	
	// Get the move history to reconstruct the opening
	moves := make([]string, 0)
	for _, move := range game.Moves() {
		moves = append(moves, move.String())
	}
	
	return s.SearchByMoves(moves)
}

// SearchByMoves finds opening information by move sequence
func (s *OpeningService) SearchByMoves(moves []string) (*models.OpeningInfo, error) {
	// Convert moves to standardized notation
	movesStr := strings.Join(moves, " ")
	
	// Search through openings for matching move sequences
	for _, opening := range s.openings {
		if len(opening.Moves) <= len(moves) {
			match := true
			for i, move := range opening.Moves {
				if i >= len(moves) || moves[i] != move {
					match = false
					break
				}
			}
			if match {
				return &opening, nil
			}
		}
	}
	
	return nil, fmt.Errorf("opening not found for move sequence: %s", movesStr)
}

// GetOpeningByName searches for opening by name (fuzzy)
func (s *OpeningService) GetOpeningByName(name string) ([]*models.OpeningInfo, error) {
	results := make([]*models.OpeningInfo, 0)
	name = strings.ToLower(name)
	
	for _, opening := range s.openings {
		if strings.Contains(strings.ToLower(opening.Name), name) ||
		   strings.Contains(strings.ToLower(opening.Variation), name) {
			openingCopy := opening
			results = append(results, &openingCopy)
		}
	}
	
	if len(results) == 0 {
		return nil, fmt.Errorf("no openings found matching: %s", name)
	}
	
	return results, nil
}

// loadOpeningDatabase initializes the ECO opening database
func (s *OpeningService) loadOpeningDatabase() {
	logrus.Info("Loading opening database...")
	
	// This is a simplified opening database
	// In a production system, this would be loaded from a file or database
	openings := []models.OpeningInfo{
		// King's Pawn Openings (B00-B99)
		{
			ECO:        "B00",
			Name:       "King's Pawn Game",
			Variation:  "",
			Moves:      []string{"e4"},
			Popularity: 45.2,
			Statistics: models.OpeningStatistics{
				White: 52.1,
				Draw:  23.4,
				Black: 24.5,
			},
			Theory:   "The most popular opening move, controlling the center and developing pieces",
			KeyIdeas: []string{"Central control", "Quick development", "King safety"},
		},
		{
			ECO:        "B01",
			Name:       "Scandinavian Defense",
			Variation:  "",
			Moves:      []string{"e4", "d5"},
			Popularity: 2.1,
			Statistics: models.OpeningStatistics{
				White: 56.8,
				Draw:  21.2,
				Black: 22.0,
			},
			Theory:   "Immediate central challenge, often leads to early queen development",
			KeyIdeas: []string{"Central challenge", "Queen activity", "Rapid development"},
		},
		{
			ECO:        "B02",
			Name:       "Alekhine's Defense",
			Variation:  "",
			Moves:      []string{"e4", "Nf6"},
			Popularity: 1.8,
			Statistics: models.OpeningStatistics{
				White: 55.4,
				Draw:  22.8,
				Black: 21.8,
			},
			Theory: "Hypermodern opening that invites white to build a pawn center",
			KeyIdeas: []string{"Hypermodern play", "Knight mobility", "Center provocation"},
		},
		{
			ECO:        "B10",
			Name:       "Caro-Kann Defense",
			Variation:  "",
			Moves:      []string{"e4", "c6"},
			Popularity: 4.2,
			Statistics: models.OpeningStatistics{
				White: 52.8,
				Draw:  27.1,
				Black: 20.1,
			},
			Theory: "Solid defense preparing d5, leads to solid pawn structures",
			KeyIdeas: []string{"Solid structure", "d5 preparation", "Positional play"},
		},
		{
			ECO:        "B20",
			Name:       "Sicilian Defense",
			Variation:  "",
			Moves:      []string{"e4", "c5"},
			Popularity: 16.8,
			Statistics: models.OpeningStatistics{
				White: 52.3,
				Draw:  23.1,
				Black: 24.6,
			},
			Theory: "Most popular response to e4, creates imbalanced positions",
			KeyIdeas: []string{"Asymmetrical structure", "Counterplay", "Sharp positions"},
		},
		
		// Queen's Pawn Openings (D00-D99)
		{
			ECO:        "D00",
			Name:       "Queen's Pawn Game",
			Variation:  "",
			Moves:      []string{"d4"},
			Popularity: 35.6,
			Statistics: models.OpeningStatistics{
				White: 54.2,
				Draw:  28.3,
				Black: 17.5,
			},
			Theory: "Classical opening focusing on central control and development",
			KeyIdeas: []string{"Central control", "Positional play", "Strategic complexity"},
		},
		{
			ECO:        "D02",
			Name:       "Queen's Pawn Game",
			Variation:  "London System",
			Moves:      []string{"d4", "d5", "Bf4"},
			Popularity: 8.4,
			Statistics: models.OpeningStatistics{
				White: 53.7,
				Draw:  29.1,
				Black: 17.2,
			},
			Theory: "Solid system setup with bishop on f4",
			KeyIdeas: []string{"System play", "Flexible development", "Solid structure"},
		},
		{
			ECO:        "D06",
			Name:       "Queen's Gambit",
			Variation:  "",
			Moves:      []string{"d4", "d5", "c4"},
			Popularity: 12.3,
			Statistics: models.OpeningStatistics{
				White: 55.1,
				Draw:  31.2,
				Black: 13.7,
			},
			Theory: "Classical gambit offering a pawn for central control",
			KeyIdeas: []string{"Central control", "Piece activity", "Initiative"},
		},
		
		// King's Indian and Other Defenses
		{
			ECO:        "E60",
			Name:       "King's Indian Defense",
			Variation:  "",
			Moves:      []string{"d4", "Nf6", "c4", "g6"},
			Popularity: 3.2,
			Statistics: models.OpeningStatistics{
				White: 54.8,
				Draw:  24.7,
				Black: 20.5,
			},
			Theory: "Hypermodern defense with fianchetto setup",
			KeyIdeas: []string{"Fianchetto", "Kingside attack", "Dynamic play"},
		},
		
		// Flank Openings (A00-A99)
		{
			ECO:        "A00",
			Name:       "Uncommon Opening",
			Variation:  "",
			Moves:      []string{},
			Popularity: 0.1,
			Statistics: models.OpeningStatistics{
				White: 50.0,
				Draw:  25.0,
				Black: 25.0,
			},
			Theory: "Various uncommon first moves",
			KeyIdeas: []string{"Surprise value", "Transpositional possibilities"},
		},
		{
			ECO:        "A04",
			Name:       "Reti Opening",
			Variation:  "",
			Moves:      []string{"Nf3"},
			Popularity: 5.7,
			Statistics: models.OpeningStatistics{
				White: 53.4,
				Draw:  26.8,
				Black: 19.8,
			},
			Theory: "Hypermodern opening delaying central pawn moves",
			KeyIdeas: []string{"Flexible development", "Hypermodern principles", "Transpositional"},
		},
		{
			ECO:        "A10",
			Name:       "English Opening",
			Variation:  "",
			Moves:      []string{"c4"},
			Popularity: 7.2,
			Statistics: models.OpeningStatistics{
				White: 54.1,
				Draw:  27.4,
				Black: 18.5,
			},
			Theory: "Flank opening controlling d5 and preparing development",
			KeyIdeas: []string{"Flank control", "Flexible structure", "Positional play"},
		},
	}
	
	// Build the opening database
	for _, opening := range openings {
		s.openings[opening.ECO] = opening
		
		// Group by ECO category (first letter)
		category := string(opening.ECO[0])
		if s.ecoTree[category] == nil {
			s.ecoTree[category] = make([]models.OpeningInfo, 0)
		}
		s.ecoTree[category] = append(s.ecoTree[category], opening)
	}
	
	logrus.Infof("Loaded %d openings into database", len(s.openings))
}

// GetAllOpenings returns all openings in the database
func (s *OpeningService) GetAllOpenings() []models.OpeningInfo {
	openings := make([]models.OpeningInfo, 0, len(s.openings))
	for _, opening := range s.openings {
		openings = append(openings, opening)
	}
	return openings
}

// GetECOCategories returns openings grouped by ECO categories
func (s *OpeningService) GetECOCategories() map[string][]models.OpeningInfo {
	return s.ecoTree
} 