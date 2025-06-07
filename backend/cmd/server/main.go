package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chess-backend/configs"
	"chess-backend/internal/handlers"
	"chess-backend/internal/middleware"
	"chess-backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize configuration
	cfg := configs.Load()

	// Setup logging
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// Create services
	cacheService := services.NewCacheService()
	stockfishService := services.NewStockfishService(cfg.Engine.MaxWorkers, cfg.Engine.BinaryPath)
	chessService := services.NewChessService()
	analysisService := services.NewAnalysisService(stockfishService, chessService, cacheService)

	// Initialize Stockfish pool
	if err := stockfishService.Initialize(); err != nil {
		logrus.Fatalf("Failed to initialize Stockfish service: %v", err)
	}
	defer stockfishService.Shutdown()

	// Setup Gin
	if cfg.App.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:          12 * time.Hour,
	}))

	// Rate limiting middleware
	router.Use(middleware.RateLimit(cfg.RateLimit))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
		})
	})

	// Initialize handlers
	analysisHandler := handlers.NewAnalysisHandler(analysisService)
	healthHandler := handlers.NewHealthHandler()

	// API routes
	api := router.Group("/api")
	{
		// Game analysis endpoints
		games := api.Group("/games")
		{
			games.POST("/analyze", analysisHandler.AnalyzeGame)
			games.GET("/analyze/:gameId", analysisHandler.GetAnalysis)
			games.GET("/analyze/:gameId/progress", analysisHandler.GetProgress)
		}

		// Position analysis endpoints
		positions := api.Group("/positions")
		{
			positions.POST("/analyze", analysisHandler.AnalyzePosition)
		}

		// Engine configuration endpoints
		engine := api.Group("/engine")
		{
			engine.GET("/config", analysisHandler.GetEngineConfig)
			engine.POST("/config", analysisHandler.UpdateEngineConfig)
		}

		// Health and stats
		api.GET("/health", healthHandler.Health)
		api.GET("/stats", healthHandler.Stats)
	}

	// Create server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		logrus.Infof("Starting server on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %v", err)
	}

	logrus.Info("Server exited")
} 