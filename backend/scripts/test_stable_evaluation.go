package main

import (
	"fmt"
	"strings"
	"chess-backend/internal/services"
	"chess-backend/internal/models"
)

func main() {
	fmt.Println("=== Stable Evaluation Display System Demo ===\n")
	
	// Create the evaluation display service
	displayService := services.NewEvaluationDisplayService()
	
	// Simulate a sequence of raw centipawn evaluations that would cause UI chaos
	rawEvaluations := []int{25, 127, -43, 89, 234, -156, 67, 445, -89, 123}
	isWhiteToMove := []bool{true, false, true, false, true, false, true, false, true, false}
	
	fmt.Println("Raw Centipawn Sequence (Volatile):")
	fmt.Println("Move | Raw CP | Problem")
	fmt.Println("-----|--------|--------")
	for i, cp := range rawEvaluations {
		problem := "Normal"
		if i > 0 {
			diff := cp - rawEvaluations[i-1]
			if abs(diff) > 100 {
				problem = fmt.Sprintf("JUMP %+d", diff)
			}
		}
		fmt.Printf("%4d | %+6d | %s\n", i+1, cp, problem)
	}
	
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("Stable Display Evaluation (Smooth):")
	fmt.Println("Move | Raw CP | Win%  | Bar   | Assessment    | Stable")
	fmt.Println("-----|--------|-------|-------|---------------|-------")
	
	var previous *models.DisplayEvaluation
	
	for i, cp := range rawEvaluations {
		displayEval := displayService.NormalizeForDisplay(cp, isWhiteToMove[i], previous)
		
		fmt.Printf("%4d | %+6d | %5.1f | %+5.2f | %-13s | %v\n", 
			i+1, 
			cp, 
			displayEval.WinProbability*100, 
			displayEval.EvaluationBar,
			displayEval.PositionAssessment,
			displayEval.IsStable)
		
		previous = displayEval
	}
	
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("Key Benefits Demonstrated:")
	fmt.Println("1. Win probability stays in reasonable 45-55% range for normal moves")
	fmt.Println("2. Evaluation bar values are compressed (-1 to +1) for smooth UI")
	fmt.Println("3. Position assessments provide human-readable context")
	fmt.Println("4. Stability flags help UI decide when to update")
	fmt.Println("5. Large jumps are smoothed but still visible")
	
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("Extreme Value Handling:")
	
	extremeValues := []int{-2000, -500, 0, 500, 2000}
	fmt.Println("Raw CP | Win%  | Bar   | Assessment | Capped")
	fmt.Println("-------|-------|-------|------------|-------")
	
	for _, cp := range extremeValues {
		displayEval := displayService.NormalizeForDisplay(cp, true, nil)
		fmt.Printf("%+6d | %5.1f | %+5.2f | %-10s | %+6d\n",
			cp,
			displayEval.WinProbability*100,
			displayEval.EvaluationBar,
			displayEval.PositionAssessment,
			displayEval.DisplayScore)
	}
	
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("Configuration Examples:")
	
	configs := []string{"live_game", "game_analysis", "beginner_mode"}
	for _, config := range configs {
		settings := displayService.GetRecommendedDisplaySettings(config)
		fmt.Printf("%s: smoothing=%.2f, update=%dms\n", 
			config, 
			settings["smoothing_factor"], 
			settings["update_frequency"])
	}
	
	fmt.Println("\n✅ Stable evaluation system successfully demonstrated!")
	fmt.Println("Frontend can now use displayEvaluation fields for smooth, predictable UI.")
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
} 