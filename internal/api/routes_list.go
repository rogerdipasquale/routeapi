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
// HandleListRoutes godoc
// @Summary Retrieves a list of routes matching criteria
// @Description Queries API Gateway HTTP Route objects based on criteria
// @Accept json
// @Produce json
// @Param namespace query string true "Namespace to search for params"
// @Param label_selector query string false "Label string to match selector: comma separated values of type key=value"
// @param include_deployment query boolean false "Flag to indicate if response should include deployment data for each route"
// @Success 201 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /routes/ [get]
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
