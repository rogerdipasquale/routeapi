package api

import (
	"encoding/json"
	"net/http"

	"routeapi/internal/k8s"
)

func HandleListRoutes(k8sClient *k8s.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		namespace := req.URL.Query().Get("namespace")

		routes, err := k8sClient.ListHTTPRoutes(req.Context(), namespace)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"routes": routes})
	}
}
