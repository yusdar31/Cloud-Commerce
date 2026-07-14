package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds common configuration values shared across services.
type Config struct {
	Port        string
	DatabaseURL string
	RedisURL    string
	NATSURL     string
	JWTSecret   string
	JWTExpiry   string
	Environment string
	ServiceName string
}

// Load reads configuration from environment variables with sensible defaults.
func Load(serviceName string) *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		NATSURL:     getEnv("NATS_URL", "nats://localhost:4222"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-me"),
		JWTExpiry:   getEnv("JWT_EXPIRY", "24h"),
		Environment: getEnv("APP_ENV", "development"),
		ServiceName: serviceName,
	}
}

// IsDevelopment returns true if the environment is development.
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if the environment is production.
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// DBMaxConns returns the max DB connections from env or default 20.
func (c *Config) DBMaxConns() int {
	return getEnvInt("DB_MAX_CONNS", 20)
}

// DBMinConns returns the min DB connections from env or default 5.
func (c *Config) DBMinConns() int {
	return getEnvInt("DB_MIN_CONNS", 5)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

// DSN returns a formatted PostgreSQL DSN from individual env vars,
// used when DATABASE_URL is not provided.
func DSN() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	pass := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "cloudcommerce")
	sslmode := getEnv("DB_SSLMODE", "disable")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbname, sslmode)
}
