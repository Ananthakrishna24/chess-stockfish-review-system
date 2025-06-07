package services

import (
	"crypto/md5"
	"fmt"
	"sync"
	"time"

	"chess-backend/internal/models"

	"github.com/sirupsen/logrus"
)

// CacheService provides in-memory caching functionality
type CacheService struct {
	analysisCache map[string]*models.AnalysisCache
	positionCache map[string]*models.PositionCache
	mutex         sync.RWMutex
	stopCleanup   chan bool
}

// NewCacheService creates a new cache service
func NewCacheService() *CacheService {
	cache := &CacheService{
		analysisCache: make(map[string]*models.AnalysisCache),
		positionCache: make(map[string]*models.PositionCache),
		stopCleanup:   make(chan bool),
	}
	
	// Start cleanup goroutine
	go cache.startCleanupRoutine()
	
	return cache
}

// StoreAnalysis stores a game analysis result in cache
func (c *CacheService) StoreAnalysis(gameID string, result models.GameAnalysisResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	now := time.Now()
	cache := &models.AnalysisCache{
		GameID:      gameID,
		Result:      result,
		CreatedAt:   now,
		ExpiresAt:   now.Add(24 * time.Hour), // 24 hour expiration
		AccessedAt:  now,
		AccessCount: 1,
	}
	
	c.analysisCache[gameID] = cache
	
	logrus.Debugf("Stored analysis result for game %s in cache", gameID)
}

// GetAnalysis retrieves a game analysis result from cache
func (c *CacheService) GetAnalysis(gameID string) (*models.GameAnalysisResponse, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	cache, exists := c.analysisCache[gameID]
	if !exists {
		return nil, false
	}
	
	// Check if expired
	if cache.IsExpired() {
		delete(c.analysisCache, gameID)
		logrus.Debugf("Analysis cache entry for game %s expired and removed", gameID)
		return nil, false
	}
	
	// Update access tracking
	cache.UpdateAccess()
	
	logrus.Debugf("Retrieved analysis result for game %s from cache (access count: %d)", 
		gameID, cache.AccessCount)
	
	return &cache.Result, true
}

// GenerateGameID creates a unique game ID from PGN content
func (c *CacheService) GenerateGameID(pgn string) string {
	// Use MD5 hash of PGN content as game ID
	hash := md5.Sum([]byte(pgn))
	return fmt.Sprintf("%x", hash)
}

// GetStats returns cache statistics
func (c *CacheService) GetStats() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	var totalAccess int
	var expiredAnalysis int
	
	now := time.Now()
	
	for _, cache := range c.analysisCache {
		totalAccess += cache.AccessCount
		if cache.IsExpired() {
			expiredAnalysis++
		}
	}
	
	var expiredPositions int
	for _, cache := range c.positionCache {
		if cache.IsExpired() {
			expiredPositions++
		}
	}
	
	return map[string]interface{}{
		"analysis_cache_size":     len(c.analysisCache),
		"position_cache_size":     len(c.positionCache),
		"total_analysis_accesses": totalAccess,
		"expired_analysis":        expiredAnalysis,
		"expired_positions":       expiredPositions,
		"timestamp":              now,
	}
}

// startCleanupRoutine runs periodic cleanup of expired entries
func (c *CacheService) startCleanupRoutine() {
	ticker := time.NewTicker(30 * time.Minute) // Cleanup every 30 minutes
	defer ticker.Stop()
	
	logrus.Info("Started cache cleanup routine")
	
	for {
		select {
		case <-ticker.C:
			c.cleanupExpired()
		case <-c.stopCleanup:
			logrus.Info("Cache cleanup routine stopped")
			return
		}
	}
}

// cleanupExpired removes expired entries from cache
func (c *CacheService) cleanupExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	var expiredAnalysis []string
	
	// Find expired analysis entries
	for gameID, cache := range c.analysisCache {
		if cache.IsExpired() {
			expiredAnalysis = append(expiredAnalysis, gameID)
		}
	}
	
	// Remove expired entries
	for _, gameID := range expiredAnalysis {
		delete(c.analysisCache, gameID)
	}
	
	if len(expiredAnalysis) > 0 {
		logrus.Infof("Cache cleanup completed: removed %d analysis entries", len(expiredAnalysis))
	}
}

// Shutdown stops the cache service
func (c *CacheService) Shutdown() {
	logrus.Info("Shutting down cache service")
	close(c.stopCleanup)
	logrus.Info("Cache service shutdown complete")
} 