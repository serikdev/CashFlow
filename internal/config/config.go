package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DBConfig     DBConfig
	LoggerConfig LoggerConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Name     string
	Password string
	SllMode  string
}

type LoggerConfig struct {
	LogLevel string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Error("Failed to parse env file")
	}
	return &Config{
		DBConfig: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_port", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Name:     getEnv("DB_NAME", "cashflow"),
			Password: getEnv("DB_PASSWORD", ""),
			SllMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		LoggerConfig: LoggerConfig{
			LogLevel: getEnv("LOG_LEVEL", "debug"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
