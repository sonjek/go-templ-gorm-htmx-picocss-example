package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sonjek/go-full-stack-example/internal/service"
	"github.com/sonjek/go-full-stack-example/internal/storage"
	"github.com/sonjek/go-full-stack-example/internal/web"
	"github.com/sonjek/go-full-stack-example/internal/web/handlers"
)

// Number of notes to load per page for lazy loading
const defaultPageSize = 4

func main() {
	db, err := storage.NewDbStorage()
	if err != nil {
		exitOnError("Failed to initialize database", err)
	}
	if err := storage.DBMigrate(db); err != nil {
		exitOnError("Failed to run database migrations", err)
	}
	if err := storage.SeedData(db); err != nil {
		exitOnError("Failed to seed database", err)
	}

	noteService := service.NewNoteService(db, defaultPageSize)
	appHandlers := handlers.NewHandler(noteService)
	webServer := web.NewServer(appHandlers)
	webServer.SetupMiddleware()
	if err := webServer.SetupRoutes(); err != nil {
		exitOnError("Failed to setup routes", err)
	}

	// Run web server in a goroutine for graceful shutdown
	go func() {
		if err := webServer.Start(); err != nil {
			slog.Error("Server error", "error", err)
		}
	}()

	// Wait for shutdown signal
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	stop()

	slog.Info("Shutting down...")

	// Create new context with timeout to let the web server finish its work
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := webServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server shutdown error", "error", err)
	}

	if sqlDB, err := db.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			slog.Error("Failed to close database connection", "error", err)
		}
	}
}

func exitOnError(msg string, err error) {
	slog.Error(msg, "error", err)
	os.Exit(1) //nolint:revive // Helper function for main package
}
