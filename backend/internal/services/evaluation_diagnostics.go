package services

import (
	"chess-backend/internal/models"
	"fmt"
	"math"
)

// EvaluationDiagnostics helps debug issues with engine evaluation and expected points calculation
type EvaluationDiagnostics struct {
	epService *ExpectedPointsService
}

// NewEvaluationDiagnostics creates a new diagnostics service
func NewEvaluationDiagnostics() *EvaluationDiagnostics {
	return &EvaluationDiagnostics{
		epService: NewExpectedPointsService(),
	}
}

// DiagnoseEvaluationProblem analyzes potential issues with evaluation calculations
func (ed *EvaluationDiagnostics) DiagnoseEvaluationProblem() {
	fmt.Println("=== ENGINE EVALUATION DIAGNOSTICS ===")
	
	// Test 1: Check basic sigmoid function behavior
	fmt.Println("\n1. Testing basic sigmoid function behavior:")
	ed.testSigmoidFunction()
	
	// Test 2: Check rating factor calculations
	fmt.Println("\n2. Testing rating factor calculations:")
	ed.testRatingFactors()
	
	// Test 3: Check evaluation normalization
	fmt.Println("\n3. Testing evaluation normalization:")
	ed.testEvaluationNormalization()
	
	// Test 4: Check expected points for normal moves
	fmt.Println("\n4. Testing expected points for normal moves:")
	ed.testNormalMoveEvaluations()
	
	// Test 5: Check for potential overflow or extreme values
	fmt.Println("\n5. Testing for extreme values:")
	ed.testExtremeValues()
	
	// Test 6: Check Chess.com thresholds
	fmt.Println("\n6. Testing Chess.com threshold scaling:")
	ed.testThresholdScaling()
}

// testSigmoidFunction tests the sigmoid function for reasonable behavior
func (ed *EvaluationDiagnostics) testSigmoidFunction() {
	testCases := []struct {
		eval   int
		rating int
		desc   string
	}{
		{0, 1500, "Equal position"},
		{50, 1500, "Small advantage (0.5 pawns)"},
		{100, 1500, "Medium advantage (1 pawn)"},
		{200, 1500, "Large advantage (2 pawns)"},
		{300, 1500, "Very large advantage (3 pawns)"},
		{-50, 1500, "Small disadvantage"},
		{-100, 1500, "Medium disadvantage"},
	}
	
	for _, tc := range testCases {
		ep := ed.epService.CalculateExpectedPoints(tc.eval, tc.rating)
		ratingFactor := ed.epService.getRatingFactor(tc.rating)
		normalizedEval := float64(tc.eval) / 100.0
		adjustedEval := normalizedEval * ratingFactor
		
		fmt.Printf("  %s (%+d cp): EP=%.4f, RatingFactor=%.3f, AdjustedEval=%.3f\n", 
			tc.desc, tc.eval, ep, ratingFactor, adjustedEval)
		
		// Check for unreasonable values
		if ep < 0 || ep > 1 {
			fmt.Printf("    ⚠️  WARNING: EP value %.4f is outside [0,1] range!\n", ep)
		}
		if math.IsNaN(ep) || math.IsInf(ep, 0) {
			fmt.Printf("    ❌ ERROR: EP value is NaN or Inf!\n")
		}
	}
}

// testRatingFactors checks if rating factors are reasonable
func (ed *EvaluationDiagnostics) testRatingFactors() {
	ratings := []int{800, 1000, 1200, 1500, 1800, 2000, 2200, 2400, 2600}
	
	for _, rating := range ratings {
		factor := ed.epService.getRatingFactor(rating)
		fmt.Printf("  Rating %d: Factor %.3f\n", rating, factor)
		
		// Check for unreasonable factors
		if factor < 0.4 || factor > 2.5 {
			fmt.Printf("    ⚠️  WARNING: Rating factor %.3f seems extreme!\n", factor)
		}
	}
}

// testEvaluationNormalization checks evaluation normalization for both colors
func (ed *EvaluationDiagnostics) testEvaluationNormalization() {
	testEvals := []int{-200, -100, 0, 100, 200}
	
	for _, eval := range testEvals {
		whiteNorm := ed.epService.NormalizeEvaluationForPlayer(eval, true)
		blackNorm := ed.epService.NormalizeEvaluationForPlayer(eval, false)
		
		fmt.Printf("  Original: %+d cp → White: %+d cp, Black: %+d cp\n", 
			eval, whiteNorm, blackNorm)
		
		// Verify that normalization is correct
		if whiteNorm != eval {
			fmt.Printf("    ❌ ERROR: White normalization incorrect!\n")
		}
		if blackNorm != -eval {
			fmt.Printf("    ❌ ERROR: Black normalization incorrect!\n")
		}
	}
}

// testNormalMoveEvaluations tests EP differences for typical normal moves
func (ed *EvaluationDiagnostics) testNormalMoveEvaluations() {
	rating := 1500
	
	testCases := []struct {
		before, after int
		desc          string
		expectedEPLoss float64 // Rough expectation
	}{
		{100, 95, "Tiny evaluation change", 0.005},
		{100, 85, "Small evaluation change", 0.02},
		{100, 70, "Medium evaluation change", 0.05},
		{100, 50, "Larger evaluation change", 0.10},
		{0, -50, "From equal to disadvantage", 0.10},
		{200, 180, "From winning to slightly less winning", 0.02},
	}
	
	for _, tc := range testCases {
		epBefore := ed.epService.CalculateExpectedPoints(tc.before, rating)
		epAfter := ed.epService.CalculateExpectedPoints(tc.after, rating)
		epLoss := epBefore - epAfter
		
		fmt.Printf("  %s: %+d → %+d cp\n", tc.desc, tc.before, tc.after)
		fmt.Printf("    EP: %.4f → %.4f (loss: %.4f)\n", epBefore, epAfter, epLoss)
		
		// Check if the EP loss is reasonable
		if epLoss > 0.5 {
			fmt.Printf("    ⚠️  WARNING: EP loss %.4f is very large for this change!\n", epLoss)
		}
		if epLoss < 0 {
			fmt.Printf("    ⚠️  WARNING: Negative EP loss %.4f (this should be a gain!)\n", epLoss)
		}
		if math.Abs(epLoss - tc.expectedEPLoss) > 0.05 {
			fmt.Printf("    ⚠️  WARNING: EP loss %.4f differs significantly from expected %.4f\n", 
				epLoss, tc.expectedEPLoss)
		}
	}
}

// testExtremeValues checks for potential overflow or extreme value issues
func (ed *EvaluationDiagnostics) testExtremeValues() {
	extremeCases := []struct {
		eval   int
		rating int
		desc   string
	}{
		{5000, 1500, "Mate score (high)"},
		{-5000, 1500, "Mate score (low)"},
		{1000, 1500, "Very high advantage"},
		{-1000, 1500, "Very high disadvantage"},
		{100, 3000, "Normal eval, extreme rating (high)"},
		{100, 500, "Normal eval, extreme rating (low)"},
	}
	
	for _, tc := range extremeCases {
		ep := ed.epService.CalculateExpectedPoints(tc.eval, tc.rating)
		fmt.Printf("  %s: EP=%.6f\n", tc.desc, ep)
		
		// Check for problematic values
		if math.IsNaN(ep) || math.IsInf(ep, 0) {
			fmt.Printf("    ❌ ERROR: EP is NaN or Inf!\n")
		}
		if ep < 0 || ep > 1 {
			fmt.Printf("    ❌ ERROR: EP %.6f is outside valid range [0,1]!\n", ep)
		}
	}
}

// testThresholdScaling checks Chess.com threshold scaling behavior
func (ed *EvaluationDiagnostics) testThresholdScaling() {
	ratings := []int{800, 1200, 1600, 2000, 2400}
	
	for _, rating := range ratings {
		thresholds := ed.epService.GetChessComThresholds(rating)
		scalingFactor := ed.epService.getChessComRatingFactor(rating)
		
		fmt.Printf("  Rating %d (factor %.3f):\n", rating, scalingFactor)
		fmt.Printf("    Excellent: ≤%.4f\n", thresholds["excellent"])
		fmt.Printf("    Good: ≤%.4f\n", thresholds["good"])
		fmt.Printf("    Inaccuracy: ≤%.4f\n", thresholds["inaccuracy"])
		fmt.Printf("    Mistake: ≤%.4f\n", thresholds["mistake"])
	}
}

// AnalyzeMoveSequence analyzes a sequence of moves to identify evaluation issues
func (ed *EvaluationDiagnostics) AnalyzeMoveSequence(moves []models.MoveAnalysis, playerRating int) {
	fmt.Println("\n=== MOVE SEQUENCE ANALYSIS ===")
	
	for i, move := range moves {
		if i == 0 || move.BeforeEvaluation == nil {
			continue // Skip first move or moves without before evaluation
		}
		
		beforeEval := move.BeforeEvaluation.Score
		afterEval := move.Evaluation.Score
		
		// Normalize for current player
		normalizedBefore := ed.epService.NormalizeEvaluationForPlayer(beforeEval, move.MoveNumber%2 == 1)
		normalizedAfter := ed.epService.NormalizeEvaluationForPlayer(afterEval, move.MoveNumber%2 == 1)
		
		epLoss := ed.epService.CalculateExpectedPointsLoss(normalizedBefore, normalizedAfter, playerRating)
		
		fmt.Printf("Move %d (%s): %+d → %+d cp (normalized: %+d → %+d) → EP loss: %.4f\n",
			move.MoveNumber, move.SAN, beforeEval, afterEval, normalizedBefore, normalizedAfter, epLoss)
		
		// Flag potential issues
		if epLoss > 0.3 {
			fmt.Printf("  ⚠️  LARGE EP LOSS: %.4f - may indicate evaluation issue\n", epLoss)
		}
		if epLoss < -0.1 {
			fmt.Printf("  ⚠️  NEGATIVE EP LOSS: %.4f - this is a significant gain\n", epLoss)
		}
		if math.Abs(float64(beforeEval-afterEval)) < 10 && epLoss > 0.1 {
			fmt.Printf("  ⚠️  SMALL CENTIPAWN CHANGE but LARGE EP LOSS: eval change %d cp, EP loss %.4f\n", 
				afterEval-beforeEval, epLoss)
		}
	}
}

// RecommendFixes suggests potential fixes for identified issues
func (ed *EvaluationDiagnostics) RecommendFixes() {
	fmt.Println("\n=== RECOMMENDED FIXES ===")
	
	fmt.Println("1. SIGMOID FUNCTION TUNING:")
	fmt.Println("   - Current: EP = 1/(1 + exp(-eval * rating_factor))")
	fmt.Println("   - Consider gentler curve: EP = 1/(1 + exp(-eval * rating_factor * 0.5))")
	fmt.Println("   - Or Chess.com-style: EP = 1/(1 + 10^(-eval/400))")
	
	fmt.Println("\n2. RATING FACTOR ADJUSTMENT:")
	fmt.Println("   - Current range: [0.5, 2.0] may be too extreme")
	fmt.Println("   - Consider narrower range: [0.8, 1.3]")
	
	fmt.Println("\n3. EVALUATION NORMALIZATION:")
	fmt.Println("   - Consider centipawn scaling: eval/200 instead of eval/100")
	fmt.Println("   - Add position-based factors (opening/endgame)")
	
	fmt.Println("\n4. THRESHOLD VALIDATION:")
	fmt.Println("   - Verify thresholds against actual Chess.com data")
	fmt.Println("   - Consider dynamic thresholds based on position type")
	
	fmt.Println("\n5. ENGINE EVALUATION:")
	fmt.Println("   - Ensure consistent depth across positions")
	fmt.Println("   - Consider position-dependent time allocation")
	fmt.Println("   - Validate mate score handling")
}

// DiagnoseTuningParameters helps find optimal parameters for the sigmoid function
func (ed *EvaluationDiagnostics) DiagnoseTuningParameters() {
	fmt.Println("\n=== SIGMOID FUNCTION TUNING ===")
	
	// Test different scaling factors
	scalingFactors := []float64{0.3, 0.5, 0.7, 1.0, 1.5, 2.0}
	rating := 1500
	
	for _, scale := range scalingFactors {
		ratingFactor := ed.epService.getRatingFactor(rating)
		
		// Calculate EP with different scaling
		normalizedEval1 := float64(100) / 100.0 * scale
		normalizedEval2 := float64(0) / 100.0 * scale
		
		ep1 := 1.0 / (1.0 + math.Exp(-normalizedEval1*ratingFactor))
		ep2 := 1.0 / (1.0 + math.Exp(-normalizedEval2*ratingFactor))
		epLoss := ep1 - ep2
		
		fmt.Printf("Scale %.1f: 100cp advantage → EP %.4f, Equal → EP %.4f, Loss: %.4f\n",
			scale, ep1, ep2, epLoss)
	}
	
	fmt.Println("\nOptimal scaling should give:")
	fmt.Println("- 1 pawn (100cp) ≈ 0.05-0.10 EP loss")
	fmt.Println("- 2 pawns (200cp) ≈ 0.15-0.25 EP loss")
	fmt.Println("- 3+ pawns (300cp+) ≈ 0.30+ EP loss")
} 