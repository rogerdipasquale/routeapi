package k8s

import (
	"testing"
	"os"
	"fmt"
	"encoding/json"
)

func testFillRouteWithDeployment(t *testing.T) {

	var routeList HTTPRouteList
	data, _ := os.ReadFile("../test/httpRoutes.json")

	json.Unmarshal(data, &routeList)

	for _, httpRoute := range routeList.Items {
		for _, rule := range httpRoute.Spec.Rules {
			for _, backend := range rule.BackendRefs {
				fmt.Printf("%s service \n", backend.Name)
				if (*backend.Kind != "Service") {
					t.Errorf("test found a backend with kind different than 'Service': %s-%s", *backend.Kind, backend.Name)
				}
			}
		}
	}	

}
