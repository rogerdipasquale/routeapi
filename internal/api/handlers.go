package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"routeapi/internal/config"
	"routeapi/internal/k8s"
	"routeapi/internal/version"
)

// processStart matches Node process.uptime() reference point.
var processStart = time.Now()

// Deps bundles server dependencies.
type Deps struct {
	Config    config.Config
	Log       *slog.Logger
	K8sClient *k8s.Client
}

func (d Deps) hello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Welcome to the Router API!"))
	}
}

func (d Deps) health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":    "ok",
			"version":   version.Version,
			"uptime":    time.Since(processStart).Seconds(),
			"timestamp": time.Now().UnixMilli(),
		})
	}
}

// Register mounts all routes on mux (caller wraps /api prefix outside or inside).
func Register(mux chi.Router, d Deps) {
	
	mux.Get("/health", d.health())

	mux.Route("/route", func(r chi.Router) {
		r.Get("/", d.HandleListRoutes(d.K8sClient))
	})

	mux.Route("/getRoute", func(r chi.Router) {
		r.Get("/", d.HandleGetRoute(d.K8sClient))
	})

}
