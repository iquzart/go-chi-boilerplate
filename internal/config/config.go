package config

import (
	"errors"
	"fmt"
	"os"
)

type ServerConfigs struct {
	Port         string
	ServiceName  string
	LogLevel     string
	OTLPEndpoint string
}

type DatabaseConfigs struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// AppConfigs holds all configs for the service
type AppConfigs struct {
	Server   *ServerConfigs
	Database *DatabaseConfigs
}

// GetAppConfigs loads all configs (server + db) and validates them
func GetAppConfigs() (*AppConfigs, error) {
	serverCfg := &ServerConfigs{
		Port:         fmt.Sprintf(":%s", getEnvOrDefault("PORT", "8080")),
		ServiceName:  getEnvOrDefault("SERVICE_NAME", "go-chi-boilerplate"),
		LogLevel:     getEnvOrDefault("LOG_LEVEL", "info"),
		OTLPEndpoint: getEnvOrDefault("OTLP_ENDPOINT", "otelcollector:4317"),
	}

	dbCfg := &DatabaseConfigs{
		Host:     getEnvOrDefault("DB_HOST", ""),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", ""),
		Password: getEnvOrDefault("DB_PASSWORD", ""),
		DBName:   getEnvOrDefault("DB_NAME", ""),
		SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
	}

	if err := dbCfg.Validate(); err != nil {
		return nil, err
	}

	return &AppConfigs{
		Server:   serverCfg,
		Database: dbCfg,
	}, nil
}

// Validate checks if required DB configs are present
func (d *DatabaseConfigs) Validate() error {
	if d.Host == "" || d.User == "" || d.Password == "" || d.DBName == "" {
		return errors.New("database configuration is incomplete: DB_HOST, DB_USER, DB_PASSWORD, and DB_NAME are required")
	}
	return nil
}

// getEnvOrDefault returns the value of an environment variable or a default
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
