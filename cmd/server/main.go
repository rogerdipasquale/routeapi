package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"routeapi/internal/api"
	"routeapi/internal/k8s"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = ctx

	k8sClient, err := k8s.NewClient()
	if err != nil {
		slog.Error("Failed to create k8s client","ERROR", err)
		return 
	}

	router := api.NewRouter(k8sClient)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		cancel()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("Server shutdown error:", "ERROR", err)
			return 
		}
	}()

	slog.Info("Starting server on :8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("Server error:","ERROR", err)
		return 
	}
}
