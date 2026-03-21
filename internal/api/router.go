package api

import (
	"net/http"
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
	case req.Method == http.MethodGet && req.URL.Path == "/health":
		HandleHealth(w, req)
	case req.Method == http.MethodGet && req.URL.Path == "/routes":
		HandleListRoutes(r.k8s)(w, req)
	case req.Method == http.MethodGet:
		HandleGetRoute(r.k8s)(w, req)
	default:
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
	}
}
