package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"routeapi/internal/k8s"
)

const (
	gatewayGroup    = "gateway.networking.k8s.io"
	gatewayVersion  = "v1"
	gatewayResource = "httproutes"
)

// HandleGetRoute godoc
// @Summary Retrieves details of a specific route
// @Description Queries API Gateway HTTP Route object
// @Accept json
// @Produce json
// @Param namespace path string true "Namespace to search for params"
// @Param routeName path string false "Label string to match selector: comma separated values of type key=value"
// @Success 201 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /route/{namespace}/{routeName} [get]
func (d Deps) HandleGetRoute(k8sClient *k8s.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		routeName := req.PathValue("routeName")
		namespace := req.PathValue("namespace")
		labelSelector := req.URL.Query().Get("labelSelector")
		labelSelectorParam := ""
		// includeDeployment := req.URL.Query().Get("include_deployment")

		if namespace == "" {
			namespace = "default"
		}

		// verifies that labelSelector has the format needed; will be omited if not, to avoid injection
		if d.ValidateLabelSelector(labelSelector) {
			labelSelectorParam = labelSelector
		}
		/*
			Will analyze querying routes based on service in a different method
			svcNamespace := req.URL.Query().Get("serviceNamespace")
					if svcNamespace == "" {
						svcNamespace = namespace
					}

					if svcName != "" {
						route, svc, deploy, err := k8sClient.GetRouteWithService(req.Context(), routeName, namespace, svcName, svcNamespace)
						if err != nil {
							writeError(w, http.StatusInternalServerError, err.Error())
							return
						}
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(map[string]interface{}{
							"route":      route,
							"service":    svc,
							"deployment": deploy,
						})
						return
					}
		*/

		/* labels to select: ?labelSelector=app%3Druta%2Ccomponent%3Dbla */
		route, err := k8sClient.GetHTTPRoute(req.Context(), routeName, namespace, labelSelectorParam)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(route)
	}
}

func (d Deps) ValidateLabelSelector(labelSelector string) bool {
	validLabelSelector := true
	d.Log.Debug("::ValidateLabelSelector::", "input", labelSelector)
	if labelSelector != "" {
		lowerLS := strings.ToLower(labelSelector)
		selectorArr := strings.Split(lowerLS, ",")
		if len(selectorArr) > 0 {
			for i := 0; i < len(selectorArr); i++ {
				keyVal := strings.Split(selectorArr[i], "=")
				validLabelSelector = validLabelSelector && len(keyVal) == 2
				d.Log.Debug("label selector: %s", keyVal[0])
			}
		} else {
			d.Log.Warn("Label selector is wrong, ommitting it", "error", labelSelector)
		}
	}
	return validLabelSelector
}
