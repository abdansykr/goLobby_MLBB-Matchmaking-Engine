package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	AppEnv  string
	AppPort string
	
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	
	Redis struct {
		Host     string
		Port     string
		Password string
	}
	
	JWT struct {
		Secret string
	}
	
	Matchmaking struct {
		Timeout           int
		ReadyTimeout      int
		InitialRankRange  int
		ExtendedRankRange int
	}
	
	Reputation struct {
		GhostingPenalty int
	}
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file
	_ = godotenv.Load()
	
	cfg := &Config{}
	
	// App config
	cfg.AppEnv = getEnv("APP_ENV", "production")
	cfg.AppPort = getEnv("APP_PORT", "8080")
	
	// Database config
	cfg.DB.Host = getEnv("DB_HOST", "localhost")
	cfg.DB.Port = getEnv("DB_PORT", "5432")
	cfg.DB.User = getEnv("DB_USER", "postgres")
	cfg.DB.Password = getEnv("DB_PASSWORD", "postgres")
	cfg.DB.Name = getEnv("DB_NAME", "antigravity")
	
	// Redis config
	cfg.Redis.Host = getEnv("REDIS_HOST", "localhost")
	cfg.Redis.Port = getEnv("REDIS_PORT", "6379")
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")
	
	// JWT config
	cfg.JWT.Secret = getEnv("JWT_SECRET", "your-secret-key")
	
	// Matchmaking config
	cfg.Matchmaking.Timeout = getEnvAsInt("MATCHMAKING_TIMEOUT", 30)
	cfg.Matchmaking.ReadyTimeout = getEnvAsInt("READY_TIMEOUT", 60)
	cfg.Matchmaking.InitialRankRange = getEnvAsInt("INITIAL_RANK_RANGE", 2)
	cfg.Matchmaking.ExtendedRankRange = getEnvAsInt("EXTENDED_RANK_RANGE", 4)
	
	// Reputation config
	cfg.Reputation.GhostingPenalty = getEnvAsInt("GHOSTING_PENALTY", -10)
	
	return cfg, nil
}

// GetDatabaseDSN returns the PostgreSQL connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.Name,
	)
}

// GetRedisAddr returns the Redis address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	
	return value
}
