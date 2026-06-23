## Testing k8s API and routeapi

### Testing access to api in a pod


https://kubernetes.io/docs/tasks/run-application/access-api-from-pod/

* Point to the internal API server hostname
APISERVER=https://kubernetes.default.svc

* Path to ServiceAccount token
SERVICEACCOUNT=/var/run/secrets/kubernetes.io/serviceaccount

* Read this Pod's namespace
NAMESPACE=$(cat ${SERVICEACCOUNT}/namespace)

* Read the ServiceAccount bearer token
TOKEN=$(cat ${SERVICEACCOUNT}/token)

* Reference the internal certificate authority (CA)
CACERT=${SERVICEACCOUNT}/ca.crt

```shell
curl --cacert ${CACERT} --header "Authorization: Bearer ${TOKEN}" -X GET ${APISERVER}/apis/gateway.networking.k8s.io/v1/namespaces/traefik/httproutes/weight-route
```
with label selector:

```shell
https://127.0.0.1:42051/apis/gateway.networking.k8s.io/v1/namespaces/traefik/httproutes?labelSelector=app%3Druta%2Ccomponent%3Dbla&limit=500
```

### Testing routeapi locally
If there is a kubernetes cluster in the local environment (i.e. kind or minikube)

1st Download Cert and token from an existing pod:
```shell
mkdir /var/run/secrets/kubernetes.io/serviceaccount/ -p

kubectl cp <pod>:/var/run/secrets/kubernetes.io/serviceaccount/..data/ca.crt ./ca.crt

kubectl cp <pod>:/var/run/secrets/kubernetes.io/serviceaccount/..data/token ./token

cp * /var/run/secrets/kubernetes.io/serviceaccount/
```
build the application:

```shell
go build cmd/server/main.go
```

run:

```shell
KUBERNETES_SERVICE_HOST=127.0.0.1  KUBERNETES_SERVICE_PORT=42051  ./main
```
