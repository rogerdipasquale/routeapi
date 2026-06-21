package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"routeapi/internal/k8s"
)

/*
returns a list of routes matching criteria:
  - namespace
  - label + value
    will return deployment data if requested by "include_deployment" parameter
*/
func (d Deps) HandleListRoutes(k8sClient *k8s.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		namespace := req.URL.Query().Get("namespace")
		queriesDeployment := req.URL.Query().Get("include_deployment")
		labelSelector := req.URL.Query().Get("labelSelector")
		labelSelectorParam := ""

		if d.ValidateLabelSelector(labelSelector) {
			labelSelectorParam = labelSelector
		}

		routes, err := k8sClient.ListHTTPRoutes(req.Context(), namespace, labelSelectorParam)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			slog.Error("Internal server error", "ERROR", err)
			return
		}

		slog.Info("::HandleListRoutes:: queryesDeployment", "queriesDeployment", queriesDeployment)
		if queriesDeployment == "true" {

			err = k8sClient.FillRoutesWithDeployments(req.Context(), routes)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"routes": routes})
	}
}
