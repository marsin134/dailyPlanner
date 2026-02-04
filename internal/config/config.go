package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type TokenConfig struct {
	JWTSecret            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type DB struct {
	DbHost    string
	DbPort    string
	DbUser    string
	DbPass    string
	DbName    string
	DbSSLMode string
}

type Config struct {
	ServerPort int
	DB         DB
	Token      TokenConfig
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if number, err := strconv.Atoi(value); err == nil {
			return number
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

func LoadDB() DB {
	return DB{
		DbHost:    getEnv("DB_HOST", "localhost"),
		DbPort:    getEnv("DB_PORT", "5432"),
		DbUser:    getEnv("DB_USER", "postgres"),
		DbPass:    getEnv("DB_PASSWORD", "123"),
		DbName:    getEnv("DB_NAME", "dailyPlanner"),
		DbSSLMode: getEnv("DB_SSL_MODE", "disable"),
	}
}

func LoadToken() TokenConfig {
	return TokenConfig{
		JWTSecret:            getEnv("JWT_SECRET_KEY", ""),
		AccessTokenDuration:  getEnvDuration("ACCESS_TOKEN_DURATION", time.Hour),
		RefreshTokenDuration: getEnvDuration("REFRESH_TOKEN_DURATION", 168*time.Hour),
	}
}

func LoadConfig() Config {
	return Config{
		ServerPort: getEnvInt("SERVER_PORT", 8080),
		DB:         LoadDB(),
		Token:      LoadToken(),
	}
}

func LoadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments and blank lines
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		// Separating the key and the value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Setting the environment variable
		// Only if it is not installed yet
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}
