package main

import (
	"fmt"
	"chess-backend/internal/services"
	"chess-backend/internal/models"
)

func main() {
	// Simulate typical opening sequence evaluations
	openingEvals := []int{0, 15, -10, 20, -15, 25, -20, 30}
	isWhiteToMove := []bool{true, false, true, false, true, false, true, false}
	
	fmt.Println("=== OPENING MOVES ANALYSIS ===")
	
	lichessService := services.NewLichessEvaluationService()
	
	// Test individual moves
	fmt.Println("\nIndividual move processing:")
	var previous *models.DisplayEvaluation
	for i, eval := range openingEvals {
		display := lichessService.CreateDisplayEvaluation(eval, true, previous)
		fmt.Printf("Move %d: %d cp -> Display: %d cp, %.1f%% win, %.3f bar\n", 
			i+1, eval, display.DisplayScore, display.WinProbability*100, display.EvaluationBar)
		previous = display
	}
	
	// Test batch processing (what the API uses)
	fmt.Println("\nBatch processing (API response):")
	displayEvals := lichessService.ProcessEvaluationHistory(openingEvals, isWhiteToMove)
	
	for i, eval := range displayEvals {
		fmt.Printf("Move %d: %d cp -> Display: %d cp, %.1f%% win, %.3f bar\n", 
			i+1, openingEvals[i], eval.DisplayScore, eval.WinProbability*100, eval.EvaluationBar)
	}
	
	// Check difference between moves 1 and 2
	if len(displayEvals) >= 2 {
		move1 := displayEvals[0]
		move2 := displayEvals[1]
		
		fmt.Printf("\n=== MOVES 1->2 COMPARISON ===\n")
		fmt.Printf("Move 1: %d cp -> %.3f bar\n", openingEvals[0], move1.EvaluationBar)
		fmt.Printf("Move 2: %d cp -> %.3f bar\n", openingEvals[1], move2.EvaluationBar)
		fmt.Printf("Bar change: %.3f (should be minimal for opening)\n", move2.EvaluationBar - move1.EvaluationBar)
	}
} 