package k8s

import (
	"context"
//	"crypto/tls"
//	"crypto/x509"
	"encoding/json"
	"fmt"
	"log/slog"
//	"io"
	"net/http"
//	"os"
	"strings"
)

type DeploymentList struct {
	APIVersion string       `json:"apiVersion"`
	Kind       string       `json:"kind"`
	Metadata   ListMetadata `json:"metadata"`
	Items      []Deployment `json:"items"`
}

type Deployment struct {
	APIVersion string             `json:"apiVersion"`
	Kind       string             `json:"kind"`
	Metadata   DeploymentMetadata `json:"metadata"`
	Spec       DeploymentSpec     `json:"spec"`
}

type DeploymentMetadata struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	UID       string            `json:"uid,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}

type DeploymentSpec struct {
	Replicas *int32             `json:"replicas,omitempty"`
	Selector DeploymentSelector `json:"selector"`
	Template PodTemplateSpec    `json:"template"`
}

type DeploymentSelector struct {
	MatchLabels map[string]string `json:"matchLabels,omitempty"`
}

type PodTemplateSpec struct {
	Spec PodSpec `json:"spec"`
}

type PodSpec struct {
	Containers []Container `json:"containers"`
}

type Container struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Service struct {
	APIVersion string          `json:"apiVersion"`
	Kind       string          `json:"kind"`
	Metadata   ServiceMetadata `json:"metadata"`
	Spec       ServiceSpec     `json:"spec"`
}

type ServiceMetadata struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	UID       string            `json:"uid,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}

type ServiceSpec struct {
	ClusterIP string            `json:"clusterIP"`
	Selector  map[string]string `json:"selector"`
	Ports     []ServicePort     `json:"ports"`
}

type ServicePort struct {
	Name     string `json:"name"`
	Port     int32  `json:"port"`
	Protocol string `json:"protocol"`
}

func (c *Client) GetService(ctx context.Context, name, namespace string) (*ServiceInfo, error) {
	url := fmt.Sprintf("%s/api/v1/namespaces/%s/services/%s",
		c.baseURL, namespace, name)

	data, statusCode, err := c.doRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if statusCode == http.StatusNotFound {
		return nil, fmt.Errorf("service %s/%s not found", namespace, name)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", statusCode, string(data))
	}

	var svc Service
	if err := json.Unmarshal(data, &svc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Service: %w", err)
	}

	return toServiceInfo(&svc), nil
}

func (c *Client) GetDeploymentBySelector(ctx context.Context, namespace string, selector map[string]string) (*DeploymentInfo, error) {
	if len(selector) == 0 {
		return nil, fmt.Errorf("no selector provided")
	}

	var selParts []string
	for k, v := range selector {
		selParts = append(selParts, fmt.Sprintf("%s=%s", k, v))
	}
	labelSelector := strings.Join(selParts, ",")

	url := fmt.Sprintf("%s/apis/apps/v1/namespaces/%s/deployments?labelSelector=%s",
		c.baseURL, namespace, labelSelector)

	slog.Info("::GetDeploymentBySelector:: querying", "url", url)
	data, statusCode, err := c.doRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", statusCode, string(data))
	}
	slog.Info("::GetDeploymentBySelector::", "response", data)
	var deployList DeploymentList
	if err := json.Unmarshal(data, &deployList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal DeploymentList: %w", err)
	}

	if len(deployList.Items) == 0 {
		return nil, fmt.Errorf("no deployment found with selector")
	}

	return toDeploymentInfo(&deployList.Items[0]), nil
}

func toServiceInfo(svc *Service) *ServiceInfo {
	return &ServiceInfo{
		Name:      svc.Metadata.Name,
		Namespace: svc.Metadata.Namespace,
		ClusterIP: svc.Spec.ClusterIP,
		Selector:  svc.Spec.Selector,
	}
}

func toDeploymentInfo(deploy *Deployment) *DeploymentInfo {
	var image string
	if len(deploy.Spec.Template.Spec.Containers) > 0 {
		image = deploy.Spec.Template.Spec.Containers[0].Image
	}
	replicas := int32(0)
	if deploy.Spec.Replicas != nil {
		replicas = *deploy.Spec.Replicas
	}
	return &DeploymentInfo{
		Name:          deploy.Metadata.Name,
		Namespace:     deploy.Metadata.Namespace,
		Replicas:      replicas,
		Image:         image,
		LabelSelector: deploy.Spec.Selector.MatchLabels,
	}
}
