package services

import (
	"runtime"
	"time"

	"chess-backend/internal/models"

	"github.com/sirupsen/logrus"
)

// PerformanceOptimizer provides dynamic Stockfish optimization
type PerformanceOptimizer struct {
	systemCPUs      int
	systemMemoryMB  int
	analysisProfile string
}

// AnalysisProfile represents different optimization profiles
type AnalysisProfile struct {
	Name        string
	Description string
	Settings    OptimalSettings
}

// OptimalSettings contains optimized Stockfish settings
type OptimalSettings struct {
	Threads           int
	Hash              int
	DepthRecommended  int
	TimeRecommended   int
	WorkerCount       int
	ProfilePurpose    string
}

// NewPerformanceOptimizer creates a new performance optimizer
func NewPerformanceOptimizer() *PerformanceOptimizer {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	return &PerformanceOptimizer{
		systemCPUs:     runtime.NumCPU(),
		systemMemoryMB: int(memStats.Sys / 1024 / 1024),
	}
}

// GetOptimalSettings returns optimal settings for different use cases
func (po *PerformanceOptimizer) GetOptimalSettings(useCase string) OptimalSettings {
	switch useCase {
	case "fast_analysis":
		return po.getFastAnalysisSettings()
	case "deep_analysis":
		return po.getDeepAnalysisSettings()
	case "game_analysis":
		return po.getGameAnalysisSettings()
	case "bulk_analysis":
		return po.getBulkAnalysisSettings()
	default:
		return po.getBalancedSettings()
	}
}

// getFastAnalysisSettings optimizes for speed over accuracy
func (po *PerformanceOptimizer) getFastAnalysisSettings() OptimalSettings {
	threads := po.systemCPUs / 2 // Use half CPU cores for responsiveness
	if threads < 1 {
		threads = 1
	}
	if threads > 8 {
		threads = 8 // Fast analysis doesn't need more than 8 threads
	}
	
	hash := 64 // Small hash for quick lookup
	
	return OptimalSettings{
		Threads:          threads,
		Hash:             hash,
		DepthRecommended: 12,
		TimeRecommended:  500, // 0.5 seconds
		WorkerCount:      6,
		ProfilePurpose:   "Quick position analysis with good accuracy",
	}
}

// getDeepAnalysisSettings optimizes for maximum accuracy
func (po *PerformanceOptimizer) getDeepAnalysisSettings() OptimalSettings {
	threads := po.systemCPUs - 2 // Use nearly all CPU cores
	if threads < 1 {
		threads = 1
	}
	if threads > 32 {
		threads = 32 // Stockfish optimal max
	}
	
	// Use significant memory for deep analysis
	hash := po.systemMemoryMB / 4
	if hash < 512 {
		hash = 512
	}
	if hash > 8192 {
		hash = 8192
	}
	
	return OptimalSettings{
		Threads:          threads,
		Hash:             hash,
		DepthRecommended: 25,
		TimeRecommended:  5000, // 5 seconds
		WorkerCount:      3,     // Fewer workers for deep analysis
		ProfilePurpose:   "Maximum accuracy for critical positions",
	}
}

// getGameAnalysisSettings optimizes for analyzing full games
func (po *PerformanceOptimizer) getGameAnalysisSettings() OptimalSettings {
	threads := po.systemCPUs - 3 // Leave room for other processes
	if threads < 1 {
		threads = 1
	}
	if threads > 20 {
		threads = 20
	}
	
	hash := po.systemMemoryMB / 6 // Conservative memory usage for long analysis
	if hash < 256 {
		hash = 256
	}
	if hash > 4096 {
		hash = 4096
	}
	
	return OptimalSettings{
		Threads:          threads,
		Hash:             hash,
		DepthRecommended: 18,
		TimeRecommended:  1500, // 1.5 seconds per move
		WorkerCount:      4,
		ProfilePurpose:   "Balanced analysis for complete games with EP algorithm",
	}
}

// getBulkAnalysisSettings optimizes for analyzing many games/positions
func (po *PerformanceOptimizer) getBulkAnalysisSettings() OptimalSettings {
	threads := 4 // Lower thread count per engine to run more in parallel
	
	hash := 128 // Smaller hash to conserve memory across many workers
	
	workerCount := po.systemCPUs / 2 // More workers for parallel processing
	if workerCount < 4 {
		workerCount = 4
	}
	if workerCount > 12 {
		workerCount = 12
	}
	
	return OptimalSettings{
		Threads:          threads,
		Hash:             hash,
		DepthRecommended: 15,
		TimeRecommended:  1000, // 1 second
		WorkerCount:      workerCount,
		ProfilePurpose:   "High throughput for batch processing",
	}
}

// getBalancedSettings provides a good balance for general use
func (po *PerformanceOptimizer) getBalancedSettings() OptimalSettings {
	threads := po.systemCPUs - 2
	if threads < 1 {
		threads = 1
	}
	if threads > 16 {
		threads = 16
	}
	
	hash := po.systemMemoryMB / 8
	if hash < 128 {
		hash = 128
	}
	if hash > 2048 {
		hash = 2048
	}
	
	return OptimalSettings{
		Threads:          threads,
		Hash:             hash,
		DepthRecommended: 16,
		TimeRecommended:  1200,
		WorkerCount:      5,
		ProfilePurpose:   "General purpose analysis with good performance",
	}
}

// GetAllProfiles returns all available optimization profiles
func (po *PerformanceOptimizer) GetAllProfiles() []AnalysisProfile {
	return []AnalysisProfile{
		{
			Name:        "fast_analysis",
			Description: "Quick analysis for real-time position evaluation",
			Settings:    po.getFastAnalysisSettings(),
		},
		{
			Name:        "balanced",
			Description: "Balanced performance for general use",
			Settings:    po.getBalancedSettings(),
		},
		{
			Name:        "game_analysis",
			Description: "Optimal for analyzing complete games with EP algorithm",
			Settings:    po.getGameAnalysisSettings(),
		},
		{
			Name:        "deep_analysis",
			Description: "Maximum accuracy for critical positions",
			Settings:    po.getDeepAnalysisSettings(),
		},
		{
			Name:        "bulk_analysis",
			Description: "High throughput for batch processing",
			Settings:    po.getBulkAnalysisSettings(),
		},
	}
}

// EstimateAnalysisTime estimates how long an analysis will take
func (po *PerformanceOptimizer) EstimateAnalysisTime(moveCount int, settings OptimalSettings) time.Duration {
	// Base time per move in milliseconds
	baseTimeMs := settings.TimeRecommended
	
	// Factor in depth complexity (exponential)
	depthFactor := float64(settings.DepthRecommended) / 15.0
	
	// Factor in thread efficiency (diminishing returns)
	threadFactor := 1.0 / (1.0 + float64(settings.Threads)/10.0)
	
	// Calculate total estimated time
	totalMs := float64(moveCount) * float64(baseTimeMs) * depthFactor * threadFactor
	
	return time.Duration(totalMs) * time.Millisecond
}

// LogOptimizationReport logs detailed optimization information
func (po *PerformanceOptimizer) LogOptimizationReport(profile string) {
	settings := po.GetOptimalSettings(profile)
	
	logrus.Info("=== Stockfish Performance Optimization Report ===")
	logrus.Infof("System Resources: %d CPU cores, ~%dMB memory", po.systemCPUs, po.systemMemoryMB)
	logrus.Infof("Selected Profile: %s", profile)
	logrus.Infof("Purpose: %s", settings.ProfilePurpose)
	logrus.Info("--- Optimized Settings ---")
	logrus.Infof("Threads per Engine: %d", settings.Threads)
	logrus.Infof("Hash Table Size: %dMB", settings.Hash)
	logrus.Infof("Recommended Depth: %d", settings.DepthRecommended)
	logrus.Infof("Time per Move: %dms", settings.TimeRecommended)
	logrus.Infof("Engine Workers: %d", settings.WorkerCount)
	
	// Performance estimates
	singleGameTime := po.EstimateAnalysisTime(40, settings) // Average 40 moves per game
	logrus.Infof("Estimated Time per Game: %v", singleGameTime.Round(time.Second))
	
	logrus.Info("=== Performance Tips ===")
	logrus.Info("• For better performance, ensure Stockfish binary matches your CPU architecture")
	logrus.Info("• Download optimized binaries from: https://stockfishchess.org/download/")
	logrus.Info("• Monitor system load - reduce workers if system becomes unresponsive")
	
	if settings.Hash > po.systemMemoryMB/2 {
		logrus.Warn("⚠ Warning: Hash size may be too large for available memory")
	}
	
	if settings.Threads*settings.WorkerCount > po.systemCPUs*2 {
		logrus.Warn("⚠ Warning: Total thread count exceeds system capabilities")
	}
	
	logrus.Info("==================================================")
}

// ConvertToEngineOptions converts optimal settings to engine options
func (po *PerformanceOptimizer) ConvertToEngineOptions(settings OptimalSettings) models.EngineOptions {
	return models.EngineOptions{
		Threads:          settings.Threads,
		Hash:             settings.Hash,
		Contempt:         0,
		AnalysisContempt: "off",
	}
}

// GetPerformanceMetrics returns current performance metrics
func (po *PerformanceOptimizer) GetPerformanceMetrics() map[string]interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	return map[string]interface{}{
		"system_cpu_cores":     po.systemCPUs,
		"system_memory_mb":     po.systemMemoryMB,
		"go_routines":          runtime.NumGoroutine(),
		"heap_alloc_mb":        memStats.HeapAlloc / 1024 / 1024,
		"heap_sys_mb":          memStats.HeapSys / 1024 / 1024,
		"gc_cycles":           memStats.NumGC,
		"available_profiles":   len(po.GetAllProfiles()),
	}
} 