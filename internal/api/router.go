package api

import (
	"net/http"
	"path/filepath"
	"routeapi/internal/k8s"
)

type Router struct {
	k8s *k8s.Client
}

func NewRouter(k8sClient *k8s.Client) *Router {
	return &Router{k8s: k8sClient}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case req.Method == http.MethodGet && req.URL.Path == "/api/health":
		HandleHealth(w, req)
	case req.Method == http.MethodGet && req.URL.Path == "/api/routes":
		HandleListRoutes(r.k8s)(w, req)
	case req.Method == http.MethodGet && req.URL.Path == "/api/getRoute":
		HandleGetRoute(r.k8s)(w, req)
	/* Serving web site */ 
	case req.Method == http.MethodGet && req.URL.Path == "/":
		http.ServeFile(w, req, filepath.Join("web", "index.html"))
	case req.Method == http.MethodGet:
		if filepath.Ext(req.URL.Path) == "" {
			http.ServeFile(w, req, filepath.Join("web", "index.html"))
		} else {
			http.ServeFile(w, req, filepath.Join("web", req.URL.Path))
		}
	default:
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
	}
}
