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
// HandleListRoutes returns a list of Gateway API HttpRoutes matching criteria.
//	@Summary		List HTTPRoute objects
//	@Description	Returns a list of Gateway API HttpRoutes matching criteria
//	@Tags			HTTPRoute
//	@Produce		json
//	@Success		200					{array}	string
//	@Router			/routes				[get]
//	@param			namespace			query	string	false	"Namespace to search (default)"
//	@param			include_deployment	query	boolean	false	"decides if it adds Deployment data related to the route"
//	@param			label_selector		query	string	false	"Applies a label selector to k8s API query in format label_name=value,label_name2,value2"
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
