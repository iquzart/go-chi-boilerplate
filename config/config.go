package config

import (
	"fmt"
	"os"
	"strconv"
)

// ServerConfigs contains configuration options for the HTTP server.
type ServerConfigs struct {
	Port             string // The port number on which to start the server.
	GracefulShutdown bool   // Whether to use graceful shutdown when stopping the server.
	ServiceName      string // The service name to configure on logs and metrics.
}

// GetServerConfigs gets the server configuration options from environment variables.
func GetServerConfigs() *ServerConfigs {

	// Get the PORT from the environment variable or use a default value.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get the SERVICE_NAME from the environment variable or use a default value.
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "go-chi-boilerplate"
	}

	// Parse the ENABLE_GRACEFUL_SHUTDOWN environment variable to determine if graceful shutdown is enabled.
	gracefulShutdown, err := strconv.ParseBool(os.Getenv("ENABLE_GRACEFUL_SHUTDOWN"))
	if err != nil {
		gracefulShutdown = true
	}

	return &ServerConfigs{
		Port:             fmt.Sprintf(":%s", port),
		GracefulShutdown: gracefulShutdown,
		ServiceName:      serviceName,
	}
}
