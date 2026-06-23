package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"routeapi/internal/k8s"
)

/*
Returns a list of Gateway API HttpRoutes matching criteria:
  - namespace
  - label_selector
  - include_deployment
*/
func (d Deps) HandleListRoutes(k8sClient *k8s.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		namespace := req.URL.Query().Get("namespace")
		queriesDeployment := req.URL.Query().Get("include_deployment")
		labelSelector := req.URL.Query().Get("label_selector")
		labelSelectorParam := ""

		if d.ValidateLabelSelector(labelSelector) {
			labelSelectorParam = labelSelector
		}
		d.Log.Debug("::HandleListRoutes", "labelSelectorParam", labelSelectorParam)
		routes, err := k8sClient.ListHTTPRoutes(req.Context(), namespace, labelSelectorParam)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			slog.Error("Internal server error", "ERROR", err)
			return
		}

		d.Log.Debug("::HandleListRoutes:: queryesDeployment", "queriesDeployment", queriesDeployment)
		if queriesDeployment == "true" {

			err = k8sClient.FillRoutesWithDeployments(req.Context(), routes)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"routes": routes})
	}
}
