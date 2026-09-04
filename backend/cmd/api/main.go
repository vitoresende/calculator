package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vitoresende/calculator/backend/internal/calculator"
	transportHttp "vitoresende/calculator/backend/internal/transport/http"
)

func buildServer(port, allowedOrigins string) *http.Server {
	if port == "" {
		port = "8080"
	}
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	calcService := calculator.NewService()
	handler := transportHttp.NewHandler(calcService, allowedOrigins)

	return &http.Server{
		Addr:              ":" + port,
		Handler:           handler.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	port := os.Getenv("PORT")
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	server := buildServer(port, allowedOrigins)

	// Server run context for graceful shutdown
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		// Shutdown signal with grace period of 15 seconds
		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 15*time.Second)
		defer shutdownCancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				slog.Error("graceful shutdown timed out.. forcing exit.")
				os.Exit(1)
			}
		}()

		slog.Info("shutting down HTTP server...")
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("error during server shutdown", "error", err)
		}
		serverStopCtx()
	}()

	slog.Info("calculator backend service starting", "port", port, "allowed_origins", allowedOrigins)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server startup failed", "error", err)
		os.Exit(1)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	slog.Info("calculator backend service gracefully stopped")
}
