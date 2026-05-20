package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"routeapi/internal/api"
	"routeapi/internal/config"
	"routeapi/internal/k8s"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = ctx

	k8sClient, err := k8s.NewClient()
	if err != nil {
		slog.Error("Failed to create k8s client", "ERROR", err)
		return
	}

	cfg := config.Load()

	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{Level: lvl}
	log := slog.New(slog.NewTextHandler(os.Stdout, opts))

	deps := api.Deps{
		Config:    cfg,
		Log:       log,
		K8sClient: k8sClient,
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(api.CORSMiddleware(cfg))
	/*	r.Get("/swagger", func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, "/swagger/index.html", http.StatusMovedPermanently)
		})
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		))*/
	router.Route("/api", func(apiRouter chi.Router) {
		api.Register(apiRouter, deps)
	})

	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("listening", "addr", addr, "api", "http://127.0.0.1"+addr+"/api")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server", "err", err)
			os.Exit(1)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)

}
