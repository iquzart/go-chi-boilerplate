package server

import (
	"context"
	"fmt"
	"go-chi-boilerplate/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

// serverConfigs contains configuration options for the HTTP server.
type serverConfigs struct {
	port             string // The port number on which to start the server.
	gracefulShutdown bool   // Whether to use graceful shutdown when stopping the server.
	serviceName      string // The service name to configure on logs and metrics
}

// Run starts the HTTP server.
func Run() {
	// Get the server configuration options from environment variables.
	serverConfigs := getConfigs()

	// Initialize the router with the application's routes.
	router := routes.SetupRouter(serverConfigs.serviceName)

	// Create an HTTP server with the specified address and router.
	server := &http.Server{
		Addr:    serverConfigs.port,
		Handler: router,
	}

	// Start the server with or without graceful shutdown.
	if serverConfigs.gracefulShutdown {
		startWithGracefulShutdown(serverConfigs.serviceName, server)
	} else {
		start(serverConfigs.serviceName, server)
	}
}

// getConfigs gets the server configuration options from environment variables.
func getConfigs() *serverConfigs {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get the serviceName from the environment variable or use a default value
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "go-chi-boilerplate"
	}

	gracefulShutdown, err := strconv.ParseBool(os.Getenv("ENABLE_GRACEFUL_SHUTDOWN"))
	if err != nil {
		gracefulShutdown = true
	}
	return &serverConfigs{
		port:             fmt.Sprintf(":%s", port),
		gracefulShutdown: gracefulShutdown,
		serviceName:      serviceName,
	}
}

// start starts the HTTP server without graceful shutdown.
func start(serviceName string, server *http.Server) {
	log.Printf("Started %s service on port %s", serviceName, server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %s\n", err)
	}
}

// startWithGracefulShutdown starts the HTTP server with graceful shutdown.
func startWithGracefulShutdown(serviceName string, server *http.Server) {
	log.Printf("Started %s service on port %s with graceful shutdown ", serviceName, server.Addr)

	// Create a channel to signal when all idle connections are closed.
	idleConnsClosed := make(chan struct{})

	// Start a goroutine to listen for interrupts and shut down the server gracefully when one is received.
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Error shutting down server: %s\n", err)
		}

		close(idleConnsClosed)
	}()

	// Start the server and wait for it to return.
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %s\n", err)
	}

	<-idleConnsClosed
	log.Println("Server stopped")
}
