package k8s

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	GatewayAPIGroup   = "gateway.networking.k8s.io"
	GatewayAPIVersion = "v1"
	HTTPRouteResource = "httproutes"
)

type RouteInfo struct {
	Name       string          `json:"name"`
	Namespace  string          `json:"namespace"`
	Hostnames  []string        `json:"hostnames"`
	ParentRefs []ParentRefInfo `json:"parentRefs"`
	Rules      []RouteRuleInfo `json:"rules"`
}

type ParentRefInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type RouteRuleInfo struct {
	BackendRefs []BackendRefInfo `json:"backendRefs"`
}

type BackendRefInfo struct {
	ServiceName string `json:"serviceName"`
	ServicePort int32  `json:"servicePort"`
	Weight      *int32 `json:"weight,omitempty"`
	Service ServiceInfo `json:"serviceInfo,omitempty"`
	Deployments	[]DeploymentInfo `json:"deploymentInfo,omitempty"`
}

type ServiceInfo struct {
	Name             string            `json:"name"`
	Namespace        string            `json:"namespace"`
	ClusterIP        string            `json:"clusterIP"`
	Selector         map[string]string `json:"selector"`
	AssociatedRoutes []string          `json:"associatedRoutes,omitempty"`
}

type DeploymentInfo struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Replicas      int32             `json:"replicas"`
	Image         string            `json:"image"`
	LabelSelector map[string]string `json:"labelSelector"`
}

type Client struct {
	httpClient *http.Client
	baseURL    string
	namespace  string
}

func NewClient() (*Client, error) {
	ns := os.Getenv("POD_NAMESPACE")
	if ns == "" {
		ns = "default"
	}

	baseURL := os.Getenv("KUBERNETES_SERVICE_HOST")
	if baseURL != "" {
		port := os.Getenv("KUBERNETES_SERVICE_PORT")
		if port == "" {
			port = "443"
		}
		baseURL = fmt.Sprintf("https://%s:%s", baseURL, port)
	} else {
		baseURL = "https://kubernetes.default.svc"
	}

	caCertPool := x509.NewCertPool()
	if caCert, ok := os.LookupEnv("CA_CERT"); ok {
		caCertPool.AppendCertsFromPEM([]byte(caCert))
	} else {
		if data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"); err == nil {
			caCertPool.AppendCertsFromPEM(data)
		}
	}

	token := ""
	if data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token"); err == nil {
		token = strings.TrimSpace(string(data))
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: caCertPool,
		},
	}

	client := &Client{
		httpClient: &http.Client{Transport: transport},
		baseURL:    baseURL,
		namespace:  ns,
	}

	if token != "" {
		client.httpClient.Transport = &tokenRoundTripper{
			token: token,
			base:  transport,
		}
	}

	return client, nil
}

type tokenRoundTripper struct {
	token string
	base  http.RoundTripper
}

func (t *tokenRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(req)
}

func (c *Client) doRequest(method, url string, body io.Reader) ([]byte, int, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response: %w", err)
	}

	return data, resp.StatusCode, nil
}

type HTTPRoute struct {
	APIVersion string        `json:"apiVersion"`
	Kind       string        `json:"kind"`
	Metadata   HTTPMetadata  `json:"metadata"`
	Spec       HTTPRouteSpec `json:"spec"`
}

type HTTPMetadata struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	UID       string            `json:"uid,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}

type HTTPRouteSpec struct {
	Hostnames  []HostnameWrapper `json:"hostnames"`
	ParentRefs []ParentRef       `json:"parentRefs"`
	Rules      []RouteRule       `json:"rules"`
}

type HostnameWrapper struct {
	value string
}

func (h *HostnameWrapper) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	h.value = s
	return nil
}

func (h HostnameWrapper) String() string {
	return h.value
}

type ParentRef struct {
	Group     *string `json:"group,omitempty"`
	Kind      *string `json:"kind,omitempty"`
	Name      string  `json:"name"`
	Namespace *string `json:"namespace,omitempty"`
	UID       *string `json:"uid,omitempty"`
}

type RouteRule struct {
	BackendRefs []BackendRef `json:"backendRefs"`
}

type BackendRef struct {
	BackendObjectReference
	Weight *int32 `json:"weight,omitempty"`
}

type BackendObjectReference struct {
	Group     *string `json:"group,omitempty"`
	Kind      *string `json:"kind,omitempty"`
	Name      string  `json:"name"`
	Namespace *string `json:"namespace,omitempty"`
	Port      *int32  `json:"port,omitempty"`
	Deployment      string  `json:"deployment,omitempty"`
	Image      string  `json:"image,omitempty"`
}

type HTTPRouteList struct {
	APIVersion string       `json:"apiVersion"`
	Kind       string       `json:"kind"`
	Metadata   ListMetadata `json:"metadata"`
	Items      []HTTPRoute  `json:"items"`
}

type ListMetadata struct {
	Continue           string `json:"continue,omitempty"`
	RemainingItemCount *int64 `json:"remainingItemCount,omitempty"`
}

func (c *Client) GetHTTPRoute(ctx context.Context, name, namespace string) (*RouteInfo, error) {
	url := fmt.Sprintf("%s/apis/%s/%s/namespaces/%s/%s/%s",
		c.baseURL, GatewayAPIGroup, GatewayAPIVersion, namespace, HTTPRouteResource, name)

	data, statusCode, err := c.doRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if statusCode == http.StatusNotFound {
		return nil, fmt.Errorf("httproute %s/%s not found", namespace, name)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", statusCode, string(data))
	}

	var route HTTPRoute
	if err := json.Unmarshal(data, &route); err != nil {
		return nil, fmt.Errorf("failed to unmarshal HTTPRoute: %w", err)
	}

	return toRouteInfo(&route), nil
}

func (c *Client) ListHTTPRoutes(ctx context.Context, namespace string) ([]RouteInfo, error) {
	if namespace == "" {
		namespace = c.namespace
	}

	url := fmt.Sprintf("%s/apis/%s/%s/namespaces/%s/%s",
		c.baseURL, GatewayAPIGroup, GatewayAPIVersion, namespace, HTTPRouteResource)

	data, statusCode, err := c.doRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", statusCode, string(data))
	}

	var routeList HTTPRouteList
	if err := json.Unmarshal(data, &routeList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal HTTPRouteList: %w", err)
	}

	result := make([]RouteInfo, len(routeList.Items))
	for i, r := range routeList.Items {
		result[i] = *toRouteInfo(&r)
	}

	return result, nil
}

func (c *Client) GetRouteWithService(ctx context.Context, routeName, routeNamespace, svcName, svcNamespace string) (*RouteInfo, *ServiceInfo, *DeploymentInfo, error) {
	route, err := c.GetHTTPRoute(ctx, routeName, routeNamespace)
	if err != nil {
		return nil, nil, nil, err
	}

	svc, err := c.GetService(ctx, svcName, svcNamespace)
	if err != nil {
		return nil, nil, nil, err
	}

	var associatedRoutes []string
	for _, rule := range route.Rules {
		for _, ref := range rule.BackendRefs {
			if ref.ServiceName == svcName {
				associatedRoutes = append(associatedRoutes, ref.ServiceName)
			}
		}
	}
	svc.AssociatedRoutes = associatedRoutes

	deploy, err := c.GetDeploymentBySelector(ctx, svc.Namespace, svc.Selector)
	if err != nil {
		return route, svc, nil, nil
	}

	return route, svc, deploy, nil
}

func toRouteInfo(route *HTTPRoute) *RouteInfo {
	info := &RouteInfo{
		Name:       route.Metadata.Name,
		Namespace:  route.Metadata.Namespace,
		Hostnames:  make([]string, len(route.Spec.Hostnames)),
		ParentRefs: make([]ParentRefInfo, len(route.Spec.ParentRefs)),
		Rules:      make([]RouteRuleInfo, len(route.Spec.Rules)),
	}

	for i, h := range route.Spec.Hostnames {
		info.Hostnames[i] = h.String()
	}

	for i, p := range route.Spec.ParentRefs {
		ns := route.Metadata.Namespace
		if p.Namespace != nil {
			ns = *p.Namespace
		}
		info.ParentRefs[i] = ParentRefInfo{
			Name:      p.Name,
			Namespace: ns,
		}
	}

	for i, r := range route.Spec.Rules {
		rule := RouteRuleInfo{
			BackendRefs: make([]BackendRefInfo, len(r.BackendRefs)),
		}
		for j, ref := range r.BackendRefs {
			ns := route.Metadata.Namespace
			if ref.Namespace != nil {
				ns = *ref.Namespace
			}
			rule.BackendRefs[j] = BackendRefInfo{
				ServiceName: ref.Name,
				ServicePort: func() int32 {
					if ref.Port != nil {
						return *ref.Port
					}
					return 0
				}(),
				Weight: ref.Weight,
			}
			_ = ns
		}
		info.Rules[i] = rule
	}

	return info
}


func (c *Client) FillRoutesWithDeployments(ctx context.Context, routes []RouteInfo) (error) {
	var namespace string
	
	for _, route := range routes {
		for _, rule := range route.Rules {
			for _, backendRef := range rule.BackendRefs {
				namespace = route.Namespace
				fmt.Printf("setvice %s", backendRef.ServiceName)

				svc, err := c.GetService(ctx, backendRef.ServiceName, namespace)
				if err != nil {
					return err
				}
				deployment, err := c.GetDeploymentBySelector(ctx, namespace, svc.Selector)
				if err != nil {
					return  err
				}
				
				backendRef.Deployments[0] = *deployment

			}
		}
	}
/*
	for _, httpRoute := range routes.Items {
		for _, rule := range httpRoute.Spec.Rules {
			for _, backend := range rule.BackendRefs {
				fmt.Printf("%s service \n", backend.Name)
				if (*backend.Namespace != "") {
					namespace = *backend.Namespace
				}else {
					namespace = "default"
				} 
				//var svc ServiceInfo
				//var deployment DeploymentInfo
				
				svc, err := c.GetService(ctx, backend.Name, namespace)
				if err != nil {
					return err
				}
				deployment, err := c.GetDeploymentBySelector(ctx, namespace, svc.Selector)
				if err != nil {
					return  err
				}
				backend.Deployments[0] = toDeploymentInfo(deployment)
			}
		}
	}
*/
	return nil	
}
