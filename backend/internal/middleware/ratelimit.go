package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"chess-backend/configs"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter manages rate limiting per IP address
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	config   configs.RateLimitConfig
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config configs.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		config:   config,
	}
}

// getLimiter returns or creates a rate limiter for an IP address
func (rl *RateLimiter) getLimiter(ip string, limit int) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		// Create new limiter with burst size of 5 and rate based on hourly limit
		// Convert hourly limit to per-second rate
		perSecondRate := rate.Limit(float64(limit) / 3600.0)
		limiter = rate.NewLimiter(perSecondRate, 5) // Allow burst of 5
		rl.limiters[ip] = limiter
	}

	return limiter
}

// Allow checks if a request should be allowed
func (rl *RateLimiter) Allow(ip string, limit int) bool {
	limiter := rl.getLimiter(ip, limit)
	return limiter.Allow()
}

// cleanupOldLimiters removes inactive limiters (run periodically)
func (rl *RateLimiter) cleanupOldLimiters() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Remove limiters that haven't been used for a while
	// This is a simplified cleanup - in production you'd track last access time
	if len(rl.limiters) > 1000 {
		// Clear half the limiters when we have too many
		for ip := range rl.limiters {
			delete(rl.limiters, ip)
			if len(rl.limiters) <= 500 {
				break
			}
		}
	}
}

// RateLimit returns a gin middleware for rate limiting
func RateLimit(config configs.RateLimitConfig) gin.HandlerFunc {
	limiter := NewRateLimiter(config)

	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			limiter.cleanupOldLimiters()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		path := c.FullPath()

		var limit int
		var limitType string

		// Determine rate limit based on endpoint
		switch {
		case path == "/api/games/analyze":
			limit = config.GameAnalysisPerHour
			limitType = "game_analysis"
		case path == "/api/positions/analyze":
			limit = config.PositionAnalysisPerHour
			limitType = "position_analysis"
		default:
			// Default limit for other endpoints
			limit = 1000
			limitType = "general"
		}

		// Check rate limit
		if !limiter.Allow(ip, limit) {
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Hour).Unix()))
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": fmt.Sprintf("Too many %s requests. Limit: %d per hour", limitType, limit),
				"retry_after": 3600,
			})
			c.Abort()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Header("X-RateLimit-Type", limitType)

		c.Next()
	}
} 