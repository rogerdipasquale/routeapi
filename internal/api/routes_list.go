package api

import (
	"encoding/json"
	"net/http"
	"log/slog"
	"routeapi/internal/k8s"
)

func HandleListRoutes(k8sClient *k8s.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		namespace := req.URL.Query().Get("namespace")
		queriesDeployment := req.URL.Query().Get("include_deployment")

		routes, err := k8sClient.ListHTTPRoutes(req.Context(), namespace)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			slog.Error("Internal server error","ERROR", err)
			return
		}

		slog.Info("::HandleListRoutes:: queryesDeployment", "queriesDeployment", queriesDeployment)
		if (queriesDeployment == "true") {
			
			err = k8sClient.FillRoutesWithDeployments(req.Context(), routes)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"routes": routes})
	}
}
