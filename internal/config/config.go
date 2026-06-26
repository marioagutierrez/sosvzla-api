package config

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
)

// Config holds the configuration for the application.
type Config struct {
    DBConnectionString string
}

// Load loads the configuration from environment variables.
func Load() (*Config, error) {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using environment variables")
    }

    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        return nil, fmt.Errorf("missing required environment variable: DATABASE_URL")
    }

    return &Config{
        DBConnectionString: connStr,
    }, nil
}
