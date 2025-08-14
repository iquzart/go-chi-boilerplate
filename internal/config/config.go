package config

import (
	"fmt"
	"os"
)

type ServerConfigs struct {
	Port         string
	ServiceName  string
	LogLevel     string
	OTLPEndpoint string
}

func GetServerConfigs() *ServerConfigs {
	return &ServerConfigs{
		Port:         fmt.Sprintf(":%s", getEnvOrDefault("PORT", "8080")),
		ServiceName:  getEnvOrDefault("SERVICE_NAME", "go-chi-boilerplate"),
		LogLevel:     getEnvOrDefault("LOG_LEVEL", "info"),
		OTLPEndpoint: getEnvOrDefault("OTLP_ENDPOINT", "otelcollector:4317"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
