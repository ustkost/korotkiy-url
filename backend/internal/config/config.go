package config

import (
	"os"
	"fmt"
)

type Config struct {
	Port string
}

func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("Failed to get PORT env variable")
	}
	return &Config{Port: port}, nil
}
