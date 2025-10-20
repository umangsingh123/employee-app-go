package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerAddr      string
	DatabaseDSN     string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func Load() *Config {
	// Defaults
	serverAddr := getEnv("SERVER_ADDR", ":8080")
	// Use sqlite file by default. For production switch to Postgres DSN.
	dsn := getEnv("DATABASE_DSN", "file:employees.db?_busy_timeout=5000&_foreign_keys=1")
	maxOpen := mustAtoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	maxIdle := mustAtoi(getEnv("DB_MAX_IDLE_CONNS", "25"))
	connLifeS := mustAtoi(getEnv("DB_CONN_MAX_LIFETIME_SECONDS", "300"))

	return &Config{
		ServerAddr:      serverAddr,
		DatabaseDSN:     dsn,
		MaxOpenConns:    maxOpen,
		MaxIdleConns:    maxIdle,
		ConnMaxLifetime: time.Duration(connLifeS) * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("invalid number %s", s))
	}
	return i
}
