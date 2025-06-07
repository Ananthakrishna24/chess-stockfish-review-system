package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"chess-backend/configs"
	"chess-backend/internal/services"
)

func main() {
	fmt.Println("=== Stockfish Performance Testing ===")
	fmt.Printf("System: %d CPU cores, %s architecture\n", runtime.NumCPU(), runtime.GOARCH)
	
	// Load configuration
	cfg := configs.Load()
	
	// Test different optimization profiles
	testProfiles := []string{
		"fast_analysis",
		"balanced", 
		"game_analysis",
		"deep_analysis",
	}
	
	// Sample positions for testing
	testPositions := []string{
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", // Opening
		"r1bqkb1r/pppp1ppp/2n2n2/4p3/2B1P3/3P1N2/PPP2PPP/RNBQK2R w KQkq - 4 4", // Middlegame
		"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1", // Endgame
	}
	
	optimizer := services.NewPerformanceOptimizer()
	
	fmt.Println("\n=== Available Optimization Profiles ===")
	profiles := optimizer.GetAllProfiles()
	for _, profile := range profiles {
		fmt.Printf("Profile: %s\n", profile.Name)
		fmt.Printf("  Description: %s\n", profile.Description)
		fmt.Printf("  Threads: %d, Hash: %dMB, Depth: %d\n", 
			profile.Settings.Threads, profile.Settings.Hash, profile.Settings.DepthRecommended)
		fmt.Printf("  Time per move: %dms, Workers: %d\n", 
			profile.Settings.TimeRecommended, profile.Settings.WorkerCount)
		fmt.Printf("  Purpose: %s\n\n", profile.Settings.ProfilePurpose)
	}
	
	// Performance comparison
	fmt.Println("=== Performance Comparison ===")
	
	for _, profileName := range testProfiles {
		fmt.Printf("\nTesting profile: %s\n", profileName)
		settings := optimizer.GetOptimalSettings(profileName)
		
		// Create service with this profile
		stockfishService := services.NewStockfishService(settings.WorkerCount, cfg.Engine.BinaryPath)
		engineOptions := optimizer.ConvertToEngineOptions(settings)
		
		if err := stockfishService.Initialize(); err != nil {
			log.Printf("Failed to initialize Stockfish for profile %s: %v", profileName, err)
			continue
		}
		
		if err := stockfishService.UpdateConfig(engineOptions); err != nil {
			log.Printf("Failed to update config for profile %s: %v", profileName, err)
		}
		
		// Test analysis performance
		totalTime := time.Duration(0)
		successfulTests := 0
		
		for i, position := range testPositions {
			fmt.Printf("  Testing position %d...", i+1)
			
			start := time.Now()
			_, _, err := stockfishService.AnalyzePosition(position, settings.DepthRecommended, settings.TimeRecommended, 1)
			elapsed := time.Since(start)
			
			if err != nil {
				fmt.Printf(" ERROR: %v\n", err)
				continue
			}
			
			fmt.Printf(" %v\n", elapsed.Round(time.Millisecond))
			totalTime += elapsed
			successfulTests++
		}
		
		if successfulTests > 0 {
			avgTime := totalTime / time.Duration(successfulTests)
			fmt.Printf("  Average time: %v\n", avgTime.Round(time.Millisecond))
			
			// Estimate game analysis time
			gameTime := optimizer.EstimateAnalysisTime(40, settings) // 40 moves
			fmt.Printf("  Estimated game analysis time: %v\n", gameTime.Round(time.Second))
		}
		
		stockfishService.Shutdown()
	}
	
	// Performance metrics
	fmt.Println("\n=== System Performance Metrics ===")
	metrics := optimizer.GetPerformanceMetrics()
	for key, value := range metrics {
		fmt.Printf("%s: %v\n", key, value)
	}
	
	// Recommendations
	fmt.Println("\n=== Performance Recommendations ===")
	fmt.Println("1. Use 'fast_analysis' for real-time position evaluation")
	fmt.Println("2. Use 'game_analysis' for complete game analysis with EP algorithm")
	fmt.Println("3. Use 'deep_analysis' for critical position analysis")
	fmt.Println("4. Use 'bulk_analysis' for processing many games in parallel")
	fmt.Println("5. Download optimized Stockfish binaries from https://stockfishchess.org/download/")
	fmt.Printf("6. For your system (%d cores), consider x86-64-bmi2 or x86-64-avx2 variants\n", runtime.NumCPU())
	
	fmt.Println("\n=== Performance Testing Complete ===")
} 