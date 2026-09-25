package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port       string
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

	if cfg.Port == "" {
		return nil, fmt.Errorf("Failed to get PORT env variable")
	}
	if cfg.DBHost == "" {
		return nil, fmt.Errorf("Failed to get DB_HOST env variable")
	}
	if cfg.DBPort == "" {
		return nil, fmt.Errorf("Failed to get DB_PORT env variable")
	}
	if cfg.DBUser == "" {
		return nil, fmt.Errorf("Failed to get DB_USER env variable")
	}
	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("Failed to get DB_PASSWORD env variable")
	}
	if cfg.DBName == "" {
		return nil, fmt.Errorf("Failed to get DB_NAME env variable")
	}

	return cfg, nil
}
