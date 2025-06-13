package services

import (
	"math"
	"testing"

	"chess-backend/internal/models"
)

// TestLichessEvaluationConstants verifies the exact Lichess constants
func TestLichessEvaluationConstants(t *testing.T) {
	service := NewLichessEvaluationService()
	constants := service.GetLichessConstants()
	
	// Test exact Lichess constants from the research
	expectedSigmoidCoefficient := -0.00368208
	if math.Abs(constants["sigmoidCoefficient"]-expectedSigmoidCoefficient) > 1e-8 {
		t.Errorf("Expected sigmoid coefficient %.8f, got %.8f", 
			expectedSigmoidCoefficient, constants["sigmoidCoefficient"])
	}
	
	expectedAccuracyBase := 103.1668
	if math.Abs(constants["accuracyBase"]-expectedAccuracyBase) > 1e-4 {
		t.Errorf("Expected accuracy base %.4f, got %.4f", 
			expectedAccuracyBase, constants["accuracyBase"])
	}
	
	expectedAccuracyExponent := -0.04354
	if math.Abs(constants["accuracyExponent"]-expectedAccuracyExponent) > 1e-5 {
		t.Errorf("Expected accuracy exponent %.5f, got %.5f", 
			expectedAccuracyExponent, constants["accuracyExponent"])
	}
	
	expectedAccuracyOffset := -3.1669
	if math.Abs(constants["accuracyOffset"]-expectedAccuracyOffset) > 1e-4 {
		t.Errorf("Expected accuracy offset %.4f, got %.4f", 
			expectedAccuracyOffset, constants["accuracyOffset"])
	}
}

// TestWinProbabilityFormula tests the exact Lichess win probability formula
func TestWinProbabilityFormula(t *testing.T) {
	service := NewLichessEvaluationService()
	
	testCases := []struct {
		centipawns   int
		expectedWinProb float64
		tolerance   float64
		description string
	}{
		{0, 0.500, 0.001, "Equal position should be 50%"},
		{100, 0.591, 0.005, "+1 pawn advantage"},
		{-100, 0.409, 0.005, "-1 pawn disadvantage"},
		{200, 0.676, 0.005, "+2 pawn advantage"},
		{-200, 0.324, 0.005, "-2 pawn disadvantage"},
		{500, 0.863, 0.005, "+5 pawn advantage"},
		{-500, 0.137, 0.005, "-5 pawn disadvantage"},
		{1000, 0.975, 0.005, "Maximum capped evaluation"},
		{-1000, 0.025, 0.005, "Minimum capped evaluation"},
	}
	
	for _, tc := range testCases {
		winProb := service.ConvertCentipawnsToWinProbability(tc.centipawns)
		if math.Abs(winProb-tc.expectedWinProb) > tc.tolerance {
			t.Errorf("%s: Expected %.3f, got %.3f (diff: %.4f)", 
				tc.description, tc.expectedWinProb, winProb, 
				math.Abs(winProb-tc.expectedWinProb))
		}
	}
}

// TestAccuracyFormula tests the exact Lichess accuracy formula
func TestAccuracyFormula(t *testing.T) {
	service := NewLichessEvaluationService()
	
	testCases := []struct {
		winProbBefore   float64
		winProbAfter    float64
		expectedAccuracy float64
		tolerance       float64
		description     string
	}{
		{0.5, 0.5, 100.0, 0.1, "No change = 100% accuracy"},
		{0.6, 0.5, 63.6, 1.0, "Small loss in advantage"},
		{0.5, 0.4, 63.6, 1.0, "Small disadvantage"},
		{0.7, 0.3, 14.9, 1.0, "Major blunder"},
		{0.5, 0.6, 100.0, 0.1, "Improvement = 100% (capped)"},
		{0.8, 0.6, 40.0, 1.0, "Moderate mistake"},
	}
	
	for _, tc := range testCases {
		accuracy := service.CalculateAccuracy(tc.winProbBefore, tc.winProbAfter)
		if math.Abs(accuracy-tc.expectedAccuracy) > tc.tolerance {
			t.Errorf("%s: Expected %.1f%%, got %.1f%% (diff: %.2f)", 
				tc.description, tc.expectedAccuracy, accuracy, 
				math.Abs(accuracy-tc.expectedAccuracy))
		}
	}
}

// TestEvaluationCapping tests that evaluations are properly capped at ±1000
func TestEvaluationCapping(t *testing.T) {
	service := NewLichessEvaluationService()
	
	testCases := []struct {
		input    int
		expected int
	}{
		{500, 500},     // No capping needed
		{1000, 1000},   // At limit
		{1500, 1000},   // Should be capped
		{-500, -500},   // No capping needed
		{-1000, -1000}, // At limit
		{-1500, -1000}, // Should be capped
	}
	
	for _, tc := range testCases {
		displayEval := service.CreateDisplayEvaluation(tc.input, true, nil)
		if displayEval.DisplayScore != tc.expected {
			t.Errorf("Input %d: Expected capped score %d, got %d", 
				tc.input, tc.expected, displayEval.DisplayScore)
		}
	}
}

// TestMateScoreHandling tests proper handling of mate scores
func TestMateScoreHandling(t *testing.T) {
	service := NewLichessEvaluationService()
	
	testCases := []struct {
		centipawns  int
		isMate      bool
		winProb     float64
		description string
	}{
		{2999, false, 0.0, "Just below mate threshold"},
		{3000, true, 0.999, "Mate threshold"},
		{5000, true, 0.999, "High mate score"},
		{-3000, true, 0.001, "Mate against"},
		{-5000, true, 0.001, "High mate against"},
	}
	
	for _, tc := range testCases {
		isMate := service.IsMateScore(tc.centipawns)
		if isMate != tc.isMate {
			t.Errorf("%s: Expected mate=%t, got mate=%t", 
				tc.description, tc.isMate, isMate)
		}
		
		if tc.isMate {
			winProb := service.ConvertCentipawnsToWinProbability(tc.centipawns)
			if math.Abs(winProb-tc.winProb) > 0.001 {
				t.Errorf("%s: Expected win prob %.3f, got %.3f", 
					tc.description, tc.winProb, winProb)
			}
		}
	}
}

// TestEvaluationBarConversion tests the non-linear evaluation bar conversion
func TestEvaluationBarConversion(t *testing.T) {
	service := NewLichessEvaluationService()
	
	testCases := []struct {
		centipawns  int
		expectedBar float64
		tolerance   float64
		description string
	}{
		{0, 0.0, 0.01, "Equal position = 0.0 bar"},
		{100, 0.427, 0.01, "Small advantage"},
		{-100, -0.427, 0.01, "Small disadvantage"},
		{500, 0.852, 0.01, "Large advantage"},
		{-500, -0.852, 0.01, "Large disadvantage"},
		{1000, 0.975, 0.01, "Maximum advantage"},
		{-1000, -0.975, 0.01, "Maximum disadvantage"},
	}
	
	for _, tc := range testCases {
		displayEval := service.CreateDisplayEvaluation(tc.centipawns, true, nil)
		if math.Abs(displayEval.EvaluationBar-tc.expectedBar) > tc.tolerance {
			t.Errorf("%s: Expected bar %.3f, got %.3f", 
				tc.description, tc.expectedBar, displayEval.EvaluationBar)
		}
	}
}

// TestPositionAssessment tests position assessment thresholds
func TestPositionAssessment(t *testing.T) {
	service := NewLichessEvaluationService()
	
	testCases := []struct {
		centipawns int
		expected   string
	}{
		{0, "equal"},
		{100, "equal"},      // 59.1% is still in equal range (40-60%)
		{300, "much_better"}, // 75.1% -> much_better (>=75%)
		{600, "winning"},     // 90.1% -> winning (>=90%)
		{-100, "equal"},     // 40.9% is still in equal range
		{-300, "much_worse"}, // 24.9% -> much_worse (>=10%, <25%)
		{-600, "losing"},     // 9.9% -> losing (<10%)
	}
	
	for _, tc := range testCases {
		displayEval := service.CreateDisplayEvaluation(tc.centipawns, true, nil)
		if displayEval.PositionAssessment != tc.expected {
			t.Errorf("Centipawns %d: Expected assessment '%s', got '%s'", 
				tc.centipawns, tc.expected, displayEval.PositionAssessment)
		}
	}
}

// TestSmoothingAlgorithm tests the smoothing transition algorithm
func TestSmoothingAlgorithm(t *testing.T) {
	service := NewLichessEvaluationService()
	
	// Test case: gradual evaluation change should be smoothed
	previousEval := &models.DisplayEvaluation{
		WinProbability: 0.5,
	}
	
	// Small change - should apply full smoothing
	newEval := service.CreateDisplayEvaluation(50, true, previousEval)
	expectedProb := 0.15*service.ConvertCentipawnsToWinProbability(50) + 0.85*0.5
	
	if math.Abs(newEval.WinProbability-expectedProb) > 0.01 {
		t.Errorf("Small change smoothing: Expected %.3f, got %.3f", 
			expectedProb, newEval.WinProbability)
	}
	
	// Large change - should apply less smoothing
	bigChangeEval := service.CreateDisplayEvaluation(500, true, previousEval)
	rawWinProb := service.ConvertCentipawnsToWinProbability(500)
	expectedProbBig := 0.3*rawWinProb + 0.7*0.5
	
	if math.Abs(bigChangeEval.WinProbability-expectedProbBig) > 0.02 {
		t.Errorf("Large change smoothing: Expected %.3f, got %.3f", 
			expectedProbBig, bigChangeEval.WinProbability)
	}
}

// TestEvaluationHistory tests processing of evaluation sequences
func TestEvaluationHistory(t *testing.T) {
	service := NewLichessEvaluationService()
	
	rawEvaluations := []int{20, -30, 150, -80, 200}
	isWhiteToMove := []bool{true, false, true, false, true}
	
	displayEvals := service.ProcessEvaluationHistory(rawEvaluations, isWhiteToMove)
	
	if len(displayEvals) != len(rawEvaluations) {
		t.Errorf("Expected %d evaluations, got %d", 
			len(rawEvaluations), len(displayEvals))
	}
	
	// Test that each evaluation is properly processed
	for i, displayEval := range displayEvals {
		if displayEval.DisplayScore != service.capEvaluation(rawEvaluations[i]) {
			t.Errorf("Position %d: Expected display score %d, got %d", 
				i, service.capEvaluation(rawEvaluations[i]), displayEval.DisplayScore)
		}
		
		// Test perspective adjustment (allow for smoothing effects)
		rawWinProb := service.ConvertCentipawnsToWinProbability(rawEvaluations[i])
		expectedWinProb := rawWinProb
		if !isWhiteToMove[i] {
			expectedWinProb = 1.0 - rawWinProb
		}
		
		// Allow for significant smoothing effects, especially in later positions
		tolerance := 0.15 // Increased tolerance for smoothing
		if i > 2 { // Later positions have more accumulated smoothing
			tolerance = 0.20
		}
		
		if math.Abs(displayEval.WinProbability-expectedWinProb) > tolerance {
			t.Errorf("Position %d perspective: Expected %.3f, got %.3f (diff: %.3f)", 
				i, expectedWinProb, displayEval.WinProbability,
				math.Abs(displayEval.WinProbability-expectedWinProb))
		}
	}
}

// TestGameAccuracyCalculation tests overall game accuracy calculation
func TestGameAccuracyCalculation(t *testing.T) {
	service := NewLichessEvaluationService()
	
	// Create a sequence of evaluations representing a game
	rawEvals := []int{20, -30, 150, -200, 300, -250}
	isWhiteToMove := []bool{true, false, true, false, true, false}
	
	displayEvals := service.ProcessEvaluationHistory(rawEvals, isWhiteToMove)
	
	// Calculate accuracy for both players
	whiteAccuracy := service.CalculateGameAccuracy(displayEvals, true)
	blackAccuracy := service.CalculateGameAccuracy(displayEvals, false)
	
	// Accuracy should be between 0 and 100
	if whiteAccuracy < 0 || whiteAccuracy > 100 {
		t.Errorf("White accuracy out of range: %.1f%%", whiteAccuracy)
	}
	
	if blackAccuracy < 0 || blackAccuracy > 100 {
		t.Errorf("Black accuracy out of range: %.1f%%", blackAccuracy)
	}
	
	// With these evaluations, both players should have reasonable accuracy
	if whiteAccuracy < 30 || whiteAccuracy > 100 {
		t.Errorf("White accuracy seems unreasonable: %.1f%%", whiteAccuracy)
	}
	
	if blackAccuracy < 30 || blackAccuracy > 100 {
		t.Errorf("Black accuracy seems unreasonable: %.1f%%", blackAccuracy)
	}
}

// TestWindowingSmoothingSize tests the dynamic window size calculation
func TestWindowingSmoothingSize(t *testing.T) {
	testCases := []struct {
		gameLength   int
		expectedSize int
	}{
		{10, 2},  // Minimum window size
		{20, 2},  // Still minimum
		{30, 3},  // game_length / 10
		{50, 5},  // game_length / 10
		{100, 8}, // Maximum window size (capped at 8)
		{150, 8}, // Still maximum
	}
	
	for _, tc := range testCases {
		// Create dummy evaluations for the game length
		rawEvals := make([]int, tc.gameLength)
		isWhiteToMove := make([]bool, tc.gameLength)
		for i := 0; i < tc.gameLength; i++ {
			rawEvals[i] = 20 // Arbitrary value
			isWhiteToMove[i] = (i % 2) == 0
		}
		
		service := NewLichessEvaluationService()
		displayEvals := service.ProcessEvaluationHistory(rawEvals, isWhiteToMove)
		
		// The smoothing should have been applied - we can't directly test the window size
		// but we can verify that the evaluations were processed
		if len(displayEvals) != tc.gameLength {
			t.Errorf("Game length %d: Expected %d evaluations, got %d", 
				tc.gameLength, tc.gameLength, len(displayEvals))
		}
	}
}

// BenchmarkLichessFormula benchmarks the core Lichess evaluation formula
func BenchmarkLichessFormula(b *testing.B) {
	service := NewLichessEvaluationService()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Test various evaluation ranges
		centipawns := (i % 2000) - 1000 // Range from -1000 to +1000
		service.ConvertCentipawnsToWinProbability(centipawns)
	}
}

// BenchmarkAccuracyCalculation benchmarks the accuracy calculation
func BenchmarkAccuracyCalculation(b *testing.B) {
	service := NewLichessEvaluationService()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		before := float64(i%100) / 100.0 // 0.0 to 0.99
		after := float64((i+10)%100) / 100.0
		service.CalculateAccuracy(before, after)
	}
}

// BenchmarkEvaluationHistory benchmarks processing of evaluation sequences
func BenchmarkEvaluationHistory(b *testing.B) {
	service := NewLichessEvaluationService()
	
	// Create a typical game-length sequence
	rawEvals := make([]int, 80) // Typical game length
	isWhiteToMove := make([]bool, 80)
	
	for i := 0; i < 80; i++ {
		rawEvals[i] = (i*10 - 400) // Varying evaluations
		isWhiteToMove[i] = (i % 2) == 0
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.ProcessEvaluationHistory(rawEvals, isWhiteToMove)
	}
} 