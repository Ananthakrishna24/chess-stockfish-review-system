package main

import (
	"chess-backend/internal/services"
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	// Parse command line arguments
	var (
		pgnPath    = flag.String("pgn", "", "Path to PGN file for calibration")
		outputPath = flag.String("output", "data/thresholds.json", "Output path for thresholds")
		verbose    = flag.Bool("v", false, "Verbose logging")
	)
	flag.Parse()

	if *pgnPath == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -pgn <path_to_pgn_file> [-output <output_path>] [-v]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Setup logging
	logger := logrus.New()
	if *verbose {
		logger.SetLevel(logrus.DebugLevel)
	}

	// Check if PGN file exists
	if _, err := os.Stat(*pgnPath); os.IsNotExist(err) {
		logger.Fatalf("PGN file does not exist: %s", *pgnPath)
	}

	logger.Infof("Starting calibration with PGN: %s", *pgnPath)

	// Initialize services
	// Note: For calibration, we need minimal services
	stockfishService := services.NewStockfishService(4, "stockfish") // 4 workers, default binary path
	chessService := services.NewChessService()
	openingService := services.NewOpeningService()

	// Initialize the enhanced EP service
	enhancedEP := services.NewEnhancedEPService(
		stockfishService,
		chessService,
		openingService,
		logger,
	)

	// Run calibration
	logger.Info("Running calibration phase...")
	if err := enhancedEP.RunCalibrationFromPGN(*pgnPath); err != nil {
		logger.Fatalf("Calibration failed: %v", err)
	}

	// Display results
	thresholds := enhancedEP.GetThresholds()
	logger.Info("Calibration completed successfully!")
	
	fmt.Println("\n=== Calibration Results ===")
	for bucket, threshold := range thresholds {
		fmt.Printf("\nRating Bucket: %s\n", bucket)
		fmt.Printf("  P1  (Brilliant): %.4f\n", threshold.P1)
		fmt.Printf("  P5  (Great):     %.4f\n", threshold.P5)
		fmt.Printf("  P10 (Excellent): %.4f\n", threshold.P10)
		fmt.Printf("  P25 (Good):      %.4f\n", threshold.P25)
		fmt.Printf("  P50 (Inaccuracy): %.4f\n", threshold.P50)
		fmt.Printf("  P75 (Mistake):   %.4f\n", threshold.P75)
		fmt.Printf("  P90 (Miss):      %.4f\n", threshold.P90)
	}

	fmt.Printf("\nThresholds saved to: %s\n", *outputPath)
	logger.Info("Calibration process complete")
} 