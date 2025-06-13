package main

import (
	"fmt"
	"chess-backend/internal/services"
)

func main() {
	// Test smoothing with sample evaluation sequence that should be smooth
	rawEvaluations := []int{0, 15, 25, 20, 30, 35, 25, 40, 45, 35, 50}
	isWhiteToMove := []bool{true, false, true, false, true, false, true, false, true, false, true}
	
	fmt.Println("=== LICHESS SMOOTHING TEST ===")
	fmt.Println("Raw evaluations:", rawEvaluations)
	
	lichessService := services.NewLichessEvaluationService()
	displayEvals := lichessService.ProcessEvaluationHistory(rawEvaluations, isWhiteToMove)
	
	fmt.Println("\nProcessed evaluations:")
	for i, eval := range displayEvals {
		fmt.Printf("Move %d: %d cp -> %d cp (%.1f%% win, %.3f bar)\n", 
			i+1, rawEvaluations[i], eval.DisplayScore, eval.WinProbability*100, eval.EvaluationBar)
	}
	
	// Test edge case: equal position
	fmt.Println("\n=== EQUAL POSITION TEST ===")
	equalEval := lichessService.CreateDisplayEvaluation(0, true, nil)
	fmt.Printf("0 cp -> %.1f%% win, %.3f bar, %s\n", 
		equalEval.WinProbability*100, equalEval.EvaluationBar, equalEval.PositionAssessment)
		
	// Test small changes
	fmt.Println("\n=== SMALL CHANGES TEST ===")
	smallChanges := []int{0, 5, 10, 8, 12, 7, 15}
	isWhite := []bool{true, false, true, false, true, false, true}
	
	smallDisplayEvals := lichessService.ProcessEvaluationHistory(smallChanges, isWhite)
	for i, eval := range smallDisplayEvals {
		fmt.Printf("%d cp -> %d cp (Δ=%d)\n", 
			smallChanges[i], eval.DisplayScore, eval.DisplayScore - smallChanges[i])
	}
} 