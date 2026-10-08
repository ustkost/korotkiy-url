package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	Port       string
	LogLevel   slog.Level
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:       os.Getenv("PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}

	if v := os.Getenv("LOG_LEVEL"); v != "" {
		if err := cfg.LogLevel.UnmarshalText([]byte(v)); err != nil {
			return nil, fmt.Errorf("LOG_LEVEL: %w", err)
		}
	}

	if cfg.Port == "" {
		return nil, errors.New("PORT is required but not set")
	}
	if cfg.DBHost == "" {
		return nil, errors.New("DB_HOST is required but not set")
	}
	if cfg.DBPort == "" {
		return nil, errors.New("DB_PORT is required but not set")
	}
	if cfg.DBUser == "" {
		return nil, errors.New("DB_USER is required but not set")
	}
	if cfg.DBPassword == "" {
		return nil, errors.New("DB_PASSWORD is required but not set")
	}
	if cfg.DBName == "" {
		return nil, errors.New("DB_NAME is required but not set")
	}

	return cfg, nil
}
