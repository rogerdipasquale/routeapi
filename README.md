# Route API

This project intends to provide an API to query Gateway API Objects and derivates in order to offer a way from outside the cluster to know about routing and applications. 

## First steps

Initially we are implementing querying HTTPRoutes and deployment serving them. As we are querying Gateway API, the API is agnostic from your implementation gateway.

## Installing / Using

This is a Go project to be run as a pod inside the cluster, so there is a todo list to provide such deliverables:

- FluxCD objects to deploy using FluxCD
- More readmes and functionalities

## Examples

In this case API returns an HTTPRoute object with two rules; one of them has redirection to an existing service/deployment:

```
{
  "routes": [
    {
      "name": "weight-route",
      "namespace": "traefik",
      "hostnames": [],
      "parentRefs": [
        {
          "name": "weight-gateway",
          "namespace": "traefik"
        }
      ],
      "rules": [
        {
          "backendRefs": [
            {
              "namespace": "default",
              "serviceName": "whoami",
              "servicePort": 80,
              "weight": 1,
              "serviceInfo": {
                "name": "",
                "namespace": "",
                "clusterIP": "",
                "selector": null
              },
              "deploymentInfo": [
                {
                  "name": "whoami",
                  "namespace": "default",
                  "replicas": 2,
                  "image": "traefik/whoami",
                  "labelSelector": {
                    "app": "whoami"
                  }
                }
              ]
            },
            {
              "namespace": "default",
              "serviceName": "whoami-api",
              "servicePort": 80,
              "weight": 3,
              "serviceInfo": {
                "name": "",
                "namespace": "",
                "clusterIP": "",
                "selector": null
              }
            }
          ]
        }
      ]
    }
  ]
}
```

Why did not whoami-api returned deployment information? Selector labels are not available as labels in the metadata section of the deployment.
