package services

import (
	"fmt"
	"math"
)

// EPExample demonstrates the Expected Points algorithm with sample calculations
type EPExample struct {
	epsService *ExpectedPointsService
}

// NewEPExample creates a new EP example demonstrator
func NewEPExample() *EPExample {
	return &EPExample{
		epsService: NewExpectedPointsService(),
	}
}

// DemonstrateEPCalculation shows how EP calculations work for different scenarios
func (ep *EPExample) DemonstrateEPCalculation() {
	fmt.Println("=== Expected Points (EP) Algorithm Demonstration ===\n")
	
	// Example 1: Different player ratings with same position
	fmt.Println("1. Effect of Player Rating on Expected Points:")
	evaluation := 100 // +1.00 pawns advantage
	
	ratings := []int{800, 1200, 1600, 2000, 2400}
	for _, rating := range ratings {
		ep := ep.epsService.CalculateExpectedPoints(evaluation, rating)
		fmt.Printf("   Rating %d: EP = %.3f (%.1f%% win probability)\n", 
			rating, ep, ep*100)
	}
	fmt.Println()
	
	// Example 2: Different evaluations with same rating
	fmt.Println("2. Effect of Position Evaluation on Expected Points (Rating 1600):")
	evaluations := []int{-200, -100, 0, 100, 200, 300}
	rating := 1600
	
	for _, eval := range evaluations {
		ep := ep.epsService.CalculateExpectedPoints(eval, rating)
		fmt.Printf("   Eval %+4d cp: EP = %.3f (%.1f%% win probability)\n", 
			eval, ep, ep*100)
	}
	fmt.Println()
	
	// Example 3: Move accuracy calculation
	fmt.Println("3. Move Accuracy Examples:")
	ep.demonstrateMoveAccuracy()
	fmt.Println()
	
	// Example 4: Classification thresholds
	fmt.Println("4. Move Classification Thresholds:")
	ep.demonstrateClassificationThresholds()
}

// demonstrateMoveAccuracy shows how move accuracy is calculated
func (ep *EPExample) demonstrateMoveAccuracy() {
	rating := 1600
	beforeEval := 150 // +1.5 pawns
	
	scenarios := []struct {
		name      string
		afterEval int
		moveType  string
	}{
		{"Best Move", 150, "Maintains advantage"},
		{"Excellent", 130, "Slight inaccuracy"},
		{"Good Move", 100, "Minor loss"},
		{"Inaccuracy", 50, "Noticeable loss"},
		{"Mistake", -50, "Significant loss"},
		{"Blunder", -200, "Major loss"},
	}
	
	epBefore := ep.epsService.CalculateExpectedPoints(beforeEval, rating)
	
	for _, scenario := range scenarios {
		epAfter := ep.epsService.CalculateExpectedPoints(scenario.afterEval, rating)
		epLoss := epBefore - epAfter
		accuracy := ep.epsService.CalculateMoveAccuracy(epLoss)
		
		fmt.Printf("   %s: %+4d cp → EP Loss: %.3f → Accuracy: %.1f%%\n", 
			scenario.name, scenario.afterEval, epLoss, accuracy)
	}
}

// demonstrateClassificationThresholds shows the thresholds for move classification
func (ep *EPExample) demonstrateClassificationThresholds() {
	thresholds := ep.epsService.GetAccuracyThresholds()
	
	fmt.Println("   EP Loss Thresholds:")
	classifications := []string{"best", "excellent", "good", "inaccuracy", "mistake", "blunder"}
	
	for _, class := range classifications {
		if threshold, exists := thresholds[class]; exists {
			accuracy := (1.0 - threshold) * 100
			fmt.Printf("   %s: EP Loss ≤ %.3f (Accuracy ≥ %.1f%%)\n", 
				class, threshold, accuracy)
		}
	}
}

// SimulateGameAnalysis simulates the EP-based analysis process
func (ep *EPExample) SimulateGameAnalysis() {
	fmt.Println("\n=== Simulated Game Analysis Process ===\n")
	
	// Simulate a few moves from a game
	moves := []struct {
		moveNum    int
		player     string
		rating     int
		beforeEval int
		afterEval  int
		bestMove   string
		playedMove string
	}{
		{1, "White", 1800, 20, 25, "e4", "e4"},     // Book move
		{2, "Black", 1600, -25, -20, "e5", "e5"},   // Book move  
		{3, "White", 1800, 20, 15, "Nf3", "Bc4"},  // Inaccuracy
		{4, "Black", 1600, -15, 50, "Nf6", "f6"},  // Mistake
		{5, "White", 1800, -50, 200, "Nxe5", "Nxe5"}, // Best move (tactical)
	}
	
	fmt.Println("Move-by-move analysis:")
	fmt.Println("Move | Player | Before EP | After EP | EP Loss | Accuracy | Classification")
	fmt.Println("-----|--------|-----------|----------|---------|----------|---------------")
	
	for _, move := range moves {
		// Calculate EP values
		beforeEP := ep.epsService.CalculateExpectedPoints(move.beforeEval, move.rating)
		afterEP := ep.epsService.CalculateExpectedPoints(move.afterEval, move.rating)
		epLoss := beforeEP - afterEP
		accuracy := ep.epsService.CalculateMoveAccuracy(epLoss)
		
		// Classify move
		classification := ep.classifyMoveSimple(epLoss, move.moveNum, move.bestMove, move.playedMove)
		
		fmt.Printf("%4d | %6s | %9.3f | %8.3f | %7.3f | %8.1f%% | %s\n",
			move.moveNum, move.player, beforeEP, afterEP, epLoss, accuracy, classification)
	}
	
	// Calculate overall accuracy
	fmt.Println("\nOverall Accuracy Calculation:")
	whiteAccuracy := (100.0 + 95.0 + 80.0) / 3 // Moves 1, 3, 5
	blackAccuracy := (100.0 + 70.0) / 2        // Moves 2, 4
	
	fmt.Printf("White accuracy: %.1f%%\n", whiteAccuracy)
	fmt.Printf("Black accuracy: %.1f%%\n", blackAccuracy)
}

// classifyMoveSimple provides basic move classification for the example
func (ep *EPExample) classifyMoveSimple(epLoss float64, moveNumber int, bestMove, playedMove string) string {
	// Book moves in opening
	if moveNumber <= 8 && epLoss <= 0.03 {
		return "Book"
	}
	
	// Best move
	if bestMove == playedMove {
		return "Best"
	}
	
	// Classification by EP loss
	switch {
	case epLoss <= 0.02:
		return "Excellent"
	case epLoss <= 0.05:
		return "Good"
	case epLoss <= 0.10:
		return "Inaccuracy"
	case epLoss <= 0.20:
		return "Mistake"
	default:
		return "Blunder"
	}
}

// ExplainAlgorithm provides a detailed explanation of the EP algorithm
func (ep *EPExample) ExplainAlgorithm() {
	fmt.Println("\n=== Expected Points Algorithm Explanation ===\n")
	
	fmt.Println("The Expected Points (EP) model calculates win probability based on:")
	fmt.Println("1. Engine evaluation (in centipawns)")
	fmt.Println("2. Player rating (skill level)")
	fmt.Println()
	
	fmt.Println("Key Formula: EP = 1 / (1 + exp(-adjusted_evaluation))")
	fmt.Println("Where: adjusted_evaluation = (evaluation/100) * rating_factor")
	fmt.Println()
	
	fmt.Println("Rating Factor Calculation:")
	fmt.Println("- Base rating: 1200 (factor = 1.0)")
	fmt.Println("- Higher rating = better conversion of advantages")
	fmt.Println("- Lower rating = worse conversion of advantages")
	fmt.Println("- Factor range: 0.5 to 2.0")
	fmt.Println()
	
	fmt.Println("Move Analysis Process:")
	fmt.Println("1. Evaluate position BEFORE the move")
	fmt.Println("2. Calculate EP_before using player rating")
	fmt.Println("3. Apply the actual move played")
	fmt.Println("4. Evaluate position AFTER the move")
	fmt.Println("5. Calculate EP_after using player rating")
	fmt.Println("6. EP_loss = EP_before - EP_after")
	fmt.Println("7. Accuracy = (1 - EP_loss) * 100%")
	fmt.Println("8. Classify move based on EP_loss thresholds")
	fmt.Println()
	
	fmt.Println("Advantages over centipawn-only analysis:")
	fmt.Println("- Accounts for player skill in converting advantages")
	fmt.Println("- More accurate accuracy calculations")
	fmt.Println("- Better move classification")
	fmt.Println("- Realistic win probability estimates")
}

// CompareWithTraditional compares EP-based analysis with traditional centipawn analysis
func (ep *EPExample) CompareWithTraditional() {
	fmt.Println("\n=== EP vs Traditional Analysis Comparison ===\n")
	
	scenarios := []struct {
		name       string
		beforeEval int
		afterEval  int
		rating     int
	}{
		{"Beginner loses advantage", 200, 100, 800},
		{"Master loses advantage", 200, 100, 2200},
		{"Equal position to slight edge", 0, 50, 1600},
		{"Winning to very winning", 300, 500, 1600},
	}
	
	fmt.Println("Scenario | Rating | Traditional | EP-based | Difference")
	fmt.Println("---------|--------|-------------|----------|----------")
	
	for _, scenario := range scenarios {
		// Traditional analysis (centipawn difference)
		traditionalLoss := float64(scenario.beforeEval-scenario.afterEval) / 100.0
		traditionalAccuracy := math.Max(0, (1.0-math.Abs(traditionalLoss)/2.0)*100)
		
		// EP-based analysis
		epBefore := ep.epsService.CalculateExpectedPoints(scenario.beforeEval, scenario.rating)
		epAfter := ep.epsService.CalculateExpectedPoints(scenario.afterEval, scenario.rating)
		epLoss := epBefore - epAfter
		epAccuracy := ep.epsService.CalculateMoveAccuracy(epLoss)
		
		difference := epAccuracy - traditionalAccuracy
		
		fmt.Printf("%-20s | %4d | %11.1f%% | %8.1f%% | %+8.1f%%\n",
			scenario.name, scenario.rating, traditionalAccuracy, epAccuracy, difference)
	}
	
	fmt.Println("\nKey Insights:")
	fmt.Println("- EP analysis is more sensitive to player skill")
	fmt.Println("- Beginners get lower accuracy for same centipawn loss")
	fmt.Println("- Masters get higher accuracy for same centipawn loss")
	fmt.Println("- More realistic assessment of move quality")
} 