package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// Config represents the application configuration.
type Config struct {
	KafkaBrokers   []string `json:"kafka_brokers"`
	KafkaTopic     string   `json:"kafka_topic"`
	DatabaseURL    string   `json:"database_url"`
	BatchSize      int      `json:"batch_size"`
	BatchTimeoutMs int      `json:"batch_timeout_ms"`
	LogLevel       string   `json:"log_level"`
}

// LoadConfig loads the configuration from a JSON file.
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return &config, nil
}

// SetupLogger configures the logger based on the log level.
func SetupLogger(logLevel string) *log.Logger {
	level := log.LstdFlags | log.Lshortfile
	switch logLevel {
	case "debug":
		level = level | log.Lmicroseconds
	case "info":
		// Use default level
	case "warn":
		// Use default level
	case "error":
		// Use default level
	default:
		log.Printf("Invalid log level '%s', defaulting to 'info'", logLevel)
	}

	return log.New(os.Stdout, "[analytics-worker] ", level)
}

// Retry executes a function with retries.
func Retry(attempts int, sleep time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}

		log.Printf("attempt %d failed: %v; retrying in %s", i+1, err, sleep)
		time.Sleep(sleep)
	}
	return fmt.Errorf("after %d attempts, the last error was: %w", attempts, err)
}