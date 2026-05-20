package api

import (
	"net/http"

	"routeapi/internal/config"
)

// CORSMiddleware applies Nest-like allowed origins (split ALLOWED_ORIGINS by comma).
func CORSMiddleware(cfg config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if len(cfg.AllowedOrigins) > 0 && origin != "" {
				for _, o := range cfg.AllowedOrigins {
					if o == origin {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						w.Header().Set("Vary", "Origin")
						break
					}
				}
			}
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Authorization")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
