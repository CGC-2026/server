package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"server/internal/app"
	"server/internal/log"
)

func main() {
	logger := log.DefaultLogger

	// Create server with default config
	config := app.DefaultServerConfig()
	server := app.NewServer(config)

	// Channel to listen for errors from server
	errChan := make(chan error)

	// Start server
	go func() {
		logger.Info("Starting server...")
		errChan <- server.Start()
	}()

	// Channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Block until either an error or interrupt is received
	select {
	case err := <-errChan:
		if err != nil {
			logger.Error("Server error: %v", err)
		}
	case <-sigChan:
		logger.Info("Received interrupt signal, shutting down...")
	}

	// Create a context with timeout to gracefully shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error: %v", err)
	}

	logger.Info("Server shutdown complete")
}
