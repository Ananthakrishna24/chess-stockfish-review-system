package configs

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App       AppConfig
	Server    ServerConfig
	Engine    EngineConfig
	RateLimit RateLimitConfig
}

type AppConfig struct {
	Mode string
}

type ServerConfig struct {
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type EngineConfig struct {
	BinaryPath      string
	MaxWorkers      int
	DefaultDepth    int
	DefaultTimeMs   int
	MaxDepth        int
	MaxTimeMs       int
	Threads         int
	HashSizeMB      int
	Contempt        int
	AnalysisContempt string
}

type RateLimitConfig struct {
	GameAnalysisPerHour     int
	PositionAnalysisPerHour int
	OpeningLookupsPerHour   int
	PlayerStatsPerHour      int
}

func Load() *Config {
	viper.SetDefault("APP_MODE", "debug")
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("SERVER_READ_TIMEOUT", "30s")
	viper.SetDefault("SERVER_WRITE_TIMEOUT", "30s")
	viper.SetDefault("SERVER_SHUTDOWN_TIMEOUT", "30s")

	viper.SetDefault("ENGINE_BINARY_PATH", "stockfish")
	viper.SetDefault("ENGINE_MAX_WORKERS", 4)
	viper.SetDefault("ENGINE_DEFAULT_DEPTH", 15)
	viper.SetDefault("ENGINE_DEFAULT_TIME_MS", 1000)
	viper.SetDefault("ENGINE_MAX_DEPTH", 24)
	viper.SetDefault("ENGINE_MAX_TIME_MS", 30000)
	viper.SetDefault("ENGINE_THREADS", 1)
	viper.SetDefault("ENGINE_HASH_SIZE_MB", 128)
	viper.SetDefault("ENGINE_CONTEMPT", 0)
	viper.SetDefault("ENGINE_ANALYSIS_CONTEMPT", "off")

	viper.SetDefault("RATE_LIMIT_GAME_ANALYSIS_PER_HOUR", 10000)
	viper.SetDefault("RATE_LIMIT_POSITION_ANALYSIS_PER_HOUR", 100000)
	viper.SetDefault("RATE_LIMIT_OPENING_LOOKUPS_PER_HOUR", 1000000)
	viper.SetDefault("RATE_LIMIT_PLAYER_STATS_PER_HOUR", 500000)

	viper.AutomaticEnv()

	readTimeout, _ := time.ParseDuration(viper.GetString("SERVER_READ_TIMEOUT"))
	writeTimeout, _ := time.ParseDuration(viper.GetString("SERVER_WRITE_TIMEOUT"))
	shutdownTimeout, _ := time.ParseDuration(viper.GetString("SERVER_SHUTDOWN_TIMEOUT"))

	return &Config{
		App: AppConfig{
			Mode: viper.GetString("APP_MODE"),
		},
		Server: ServerConfig{
			Port:            viper.GetInt("SERVER_PORT"),
			ReadTimeout:     readTimeout,
			WriteTimeout:    writeTimeout,
			ShutdownTimeout: shutdownTimeout,
		},
		Engine: EngineConfig{
			BinaryPath:       viper.GetString("ENGINE_BINARY_PATH"),
			MaxWorkers:       viper.GetInt("ENGINE_MAX_WORKERS"),
			DefaultDepth:     viper.GetInt("ENGINE_DEFAULT_DEPTH"),
			DefaultTimeMs:    viper.GetInt("ENGINE_DEFAULT_TIME_MS"),
			MaxDepth:         viper.GetInt("ENGINE_MAX_DEPTH"),
			MaxTimeMs:        viper.GetInt("ENGINE_MAX_TIME_MS"),
			Threads:          viper.GetInt("ENGINE_THREADS"),
			HashSizeMB:       viper.GetInt("ENGINE_HASH_SIZE_MB"),
			Contempt:         viper.GetInt("ENGINE_CONTEMPT"),
			AnalysisContempt: viper.GetString("ENGINE_ANALYSIS_CONTEMPT"),
		},
		RateLimit: RateLimitConfig{
			GameAnalysisPerHour:     viper.GetInt("RATE_LIMIT_GAME_ANALYSIS_PER_HOUR"),
			PositionAnalysisPerHour: viper.GetInt("RATE_LIMIT_POSITION_ANALYSIS_PER_HOUR"),
			OpeningLookupsPerHour:   viper.GetInt("RATE_LIMIT_OPENING_LOOKUPS_PER_HOUR"),
			PlayerStatsPerHour:      viper.GetInt("RATE_LIMIT_PLAYER_STATS_PER_HOUR"),
		},
	}
} 