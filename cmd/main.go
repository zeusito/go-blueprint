package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/zeusito/go-blueprint/internal/adapters/api"
	"github.com/zeusito/go-blueprint/internal/adapters/api/controllers"
	"github.com/zeusito/go-blueprint/internal/adapters/config"
	"go.uber.org/zap"
)

func main() {
	// Setup zap logger
	logger := config.NewLogger()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	// Load Configurations
	configurations, err := config.LoadConfigurations()
	if err != nil {
		logger.Sugar().Fatalf("failed to load configurations: %v", err)
	}

	// database
	//rdb := db.NewConnection(configurations.Database, logger.Sugar())

	// Http Server
	srv := api.NewHTTPServer(configurations.Server)

	// --- Controllers ---
	controllers.NewHealthRoutes(srv)

	// Start the server in the background
	go srv.Start(logger.Sugar())

	// --- Graceful shutdown ---
	gracefulShutdown(srv, logger)
}

func gracefulShutdown(srv *api.HTTPServer, logger *zap.Logger) {
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logger.Info("shutting down")
	os.Exit(0)
}
