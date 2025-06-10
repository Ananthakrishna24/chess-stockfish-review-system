package services

import (
	"chess-backend/internal/models"
	"fmt"
)

// Example demonstrates the Chess.com move classification algorithm
func ExampleChessComMoveClassification() {
	// Initialize services
	epService := NewExpectedPointsService()
	classifier := NewChessComMoveClassifier(epService)
	
	// Example 1: Book move (opening theory)
	bookMove := MoveCategoryData{
		MoveNumber:       4,
		Move:            "Nf3",
		UCI:             "g1f3",
		IsWhiteToMove:   true,
		BeforeEvaluation: &models.EngineEvaluation{Score: 20},
		AfterEvaluation:  &models.EngineEvaluation{Score: 25},
		ECO:             "C20", // King's Pawn Game
		OpeningName:     "King's Pawn Game",
		PlayerRating:    1500,
		MaterialBefore:  models.MaterialValue{Total: 3900},
		MaterialAfter:   models.MaterialValue{Total: 3900},
	}
	
	classification := classifier.ClassifyMove(bookMove)
	fmt.Printf("Book move example: %s -> %s\n", bookMove.Move, classification)
	
	// Example 2: Brilliant move (piece sacrifice)
	brilliantMove := MoveCategoryData{
		MoveNumber:       15,
		Move:            "Nxf7",
		UCI:             "g5f7",
		IsWhiteToMove:   true,
		BeforeEvaluation: &models.EngineEvaluation{Score: 50},   // Slightly better
		AfterEvaluation:  &models.EngineEvaluation{Score: 180},  // Much better after sacrifice
		PlayerRating:    1600,
		MaterialBefore:  models.MaterialValue{Total: 3900, Knights: 2},
		MaterialAfter:   models.MaterialValue{Total: 3580, Knights: 1}, // Lost a knight
	}
	
	classification = classifier.ClassifyMove(brilliantMove)
	fmt.Printf("Brilliant sacrifice example: %s -> %s\n", brilliantMove.Move, classification)
	
	// Example 3: Excellent move (minimal EP loss)
	excellentMove := MoveCategoryData{
		MoveNumber:       20,
		Move:            "Qd2",
		UCI:             "d1d2",
		IsWhiteToMove:   true,
		BeforeEvaluation: &models.EngineEvaluation{Score: 120},
		AfterEvaluation:  &models.EngineEvaluation{Score: 115}, // Tiny loss
		PlayerRating:    1800,
		MaterialBefore:  models.MaterialValue{Total: 3900},
		MaterialAfter:   models.MaterialValue{Total: 3900},
	}
	
	classification = classifier.ClassifyMove(excellentMove)
	fmt.Printf("Excellent move example: %s -> %s\n", excellentMove.Move, classification)
	
	// Example 4: Blunder (large EP loss)
	blunderMove := MoveCategoryData{
		MoveNumber:       25,
		Move:            "Qxh7??",
		UCI:             "d4h7",
		IsWhiteToMove:   true,
		BeforeEvaluation: &models.EngineEvaluation{Score: 80},   // Slightly better
		AfterEvaluation:  &models.EngineEvaluation{Score: -250}, // Much worse
		PlayerRating:    1400,
		MaterialBefore:  models.MaterialValue{Total: 3900},
		MaterialAfter:   models.MaterialValue{Total: 3900},
	}
	
	classification = classifier.ClassifyMove(blunderMove)
	fmt.Printf("Blunder example: %s -> %s\n", blunderMove.Move, classification)
	
	// Example 5: Rating-dependent thresholds
	fmt.Println("\nRating-dependent thresholds:")
	
	ratings := []int{800, 1200, 1600, 2000, 2400}
	for _, rating := range ratings {
		thresholds := classifier.getChessComThresholds(rating)
		fmt.Printf("Rating %d: Excellent ≤%.3f, Good ≤%.3f, Inaccuracy ≤%.3f, Mistake ≤%.3f\n",
			rating, thresholds.Excellent, thresholds.Good, thresholds.Inaccuracy, thresholds.Mistake)
	}
}

// TestMoveClassificationAccuracy tests the accuracy of classifications
func TestMoveClassificationAccuracy() {
	epService := NewExpectedPointsService()
	classifier := NewChessComMoveClassifier(epService)
	
	testCases := []struct {
		name        string
		moveData    MoveCategoryData
		expected    models.MoveClassification
		description string
	}{
		{
			name: "Perfect move",
			moveData: MoveCategoryData{
				MoveNumber:       10,
				IsWhiteToMove:   true,
				BeforeEvaluation: &models.EngineEvaluation{Score: 100},
				AfterEvaluation:  &models.EngineEvaluation{Score: 100}, // No change
				PlayerRating:    1500,
				MaterialBefore:  models.MaterialValue{Total: 3900},
				MaterialAfter:   models.MaterialValue{Total: 3900},
			},
			expected:    models.Best,
			description: "No EP loss should be classified as Best",
		},
		{
			name: "Small mistake",
			moveData: MoveCategoryData{
				MoveNumber:       15,
				IsWhiteToMove:   true,
				BeforeEvaluation: &models.EngineEvaluation{Score: 100},
				AfterEvaluation:  &models.EngineEvaluation{Score: 50}, // ~0.05 EP loss
				PlayerRating:    1500,
				MaterialBefore:  models.MaterialValue{Total: 3900},
				MaterialAfter:   models.MaterialValue{Total: 3900},
			},
			expected:    models.Good, // or Inaccuracy depending on exact EP calculation
			description: "Small loss should be Good or Inaccuracy",
		},
		{
			name: "Queen blunder",
			moveData: MoveCategoryData{
				MoveNumber:       20,
				IsWhiteToMove:   true,
				BeforeEvaluation: &models.EngineEvaluation{Score: 0},
				AfterEvaluation:  &models.EngineEvaluation{Score: -900}, // Lost a queen
				PlayerRating:    1500,
				MaterialBefore:  models.MaterialValue{Total: 3900, Queens: 1},
				MaterialAfter:   models.MaterialValue{Total: 3000, Queens: 0},
			},
			expected:    models.Blunder,
			description: "Losing a queen should be a Blunder",
		},
	}
	
	fmt.Println("Testing move classification accuracy:")
	for _, tc := range testCases {
		result := classifier.ClassifyMove(tc.moveData)
		status := "✓"
		if result != tc.expected {
			status = "✗"
		}
		fmt.Printf("%s %s: Expected %s, Got %s - %s\n",
			status, tc.name, tc.expected, result, tc.description)
	}
}

// DemonstrateEPCalculation shows how expected points are calculated
func DemonstrateEPCalculation() {
	epService := NewExpectedPointsService()
	
	fmt.Println("Expected Points Calculation Examples:")
	fmt.Println("(Evaluation in centipawns -> Win Probability)")
	
	evaluations := []int{-300, -100, 0, 100, 300, 500}
	ratings := []int{1200, 1600, 2000}
	
	for _, rating := range ratings {
		fmt.Printf("\nRating %d:\n", rating)
		for _, eval := range evaluations {
			ep := epService.CalculateExpectedPoints(eval, rating)
			fmt.Printf("  %+4d cp -> %.3f EP (%.1f%% win chance)\n", eval, ep, ep*100)
		}
	}
	
	// Show EP loss calculation
	fmt.Println("\nEP Loss Examples (1500 rated player):")
	examples := []struct {
		before, after int
		description   string
	}{
		{100, 95, "Tiny loss (excellent move)"},
		{100, 80, "Small loss (good move)"},
		{100, 50, "Medium loss (inaccuracy)"},
		{100, 0, "Large loss (mistake)"},
		{100, -200, "Huge loss (blunder)"},
	}
	
	for _, ex := range examples {
		loss := epService.CalculateExpectedPointsLoss(ex.before, ex.after, 1500)
		fmt.Printf("  %+3d -> %+3d cp: %.3f EP loss (%s)\n",
			ex.before, ex.after, loss, ex.description)
	}
} 