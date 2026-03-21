# Route API

This project intends to provide an API to query Gateway API Objects and derivates in order to offer a way from outside the cluster to know about routing and applications. 

## First steps

Initially we are implementing querying HTTPRoutes and deployment serving them. As we are querying Gateway API, the API is agnostic from your implementation gateway.

## Installing / Using

This is a Go project to be run as a pod inside the cluster, so there is a todo list to provide such deliverables:

- GitHub action to build the image
- FluxCD objects to deploy using FluxCD
- More readmes and functionalities

