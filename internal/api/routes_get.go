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

func HandleGetRoute(k8sClient *k8s.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		routeName, svcName, ok := parseRoutePath(path)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid path format")
			return
		}

		namespace := req.URL.Query().Get("namespace")
		if namespace == "" {
			namespace = "default"
		}

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

		route, err := k8sClient.GetHTTPRoute(req.Context(), routeName, namespace)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(route)
	}
}

func parseRoutePath(path string) (routeName, svcName string, ok bool) {
	gvr := gatewayGroup + "/" + gatewayVersion + "/" + gatewayResource

	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	if len(path) > len(gvr) && path[:len(gvr)] == gvr {
		path = path[len(gvr)+1:]
	}

	parts := strings.Split(path, "/")

	if len(parts) >= 2 && parts[0] == "namespaces" {
		if len(parts) >= 4 && parts[2] == "httproutes" {
			routeName = parts[3]
			if len(parts) >= 6 && parts[4] == "services" {
				svcName = parts[5]
			}
			return routeName, svcName, true
		}
	}

	if len(parts) >= 1 {
		routeName = parts[0]
		if len(parts) >= 3 && parts[1] == "services" {
			svcName = parts[2]
		}
		return routeName, svcName, true
	}

	return "", "", false
}
